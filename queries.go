package migrate

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"
)

func applySchema(ctx context.Context, client *dgo.Dgraph) error {
	return client.Alter(ctx, &api.Operation{
		Schema: `
		version_timestamp: int @index(int) @upsert .
		version_current: int @index(int) @upsert .
		type SchemaVersion {
			version_timestamp: int
			version_current: int
		}`,
	})
}

func fetchVersion(ctx context.Context, client *dgo.Dgraph) (Version, error) {

	q := `query {
		current_version(func: eq(dgraph.type, "SchemaVersion")) {
			version_timestamp	
			version_current
		}
	}`

	txn := client.NewTxn()
	defer txn.Discard(ctx)

	resp, err := txn.Query(ctx, q)
	if err != nil {
		return Version{}, err
	}

	type Root struct {
		Version []Version `json:"current_version"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return Version{}, err
	}

	if len(r.Version) == 0 {
		return Version{}, ErrVersionFetch
	}

	return r.Version[0], nil
}

func upVersion(
	ctx context.Context,
	client *dgo.Dgraph,
	version int64,
	schemaUp []byte,
) error {

	op := &api.Operation{
		Schema: string(schemaUp),
	}

	if err := client.Alter(ctx, op); err != nil {
		return err
	}

	q := `
		query {
			version_node as var(func: eq(dgraph.type, "SchemaVersion")) {}
		}`

	mu := &api.Mutation{
		SetNquads: fmt.Appendf(
			[]byte{}, `
		    uid(version_node) <dgraph.type> "SchemaVersion" .
			uid(version_node) <version_current> "%d" .
			uid(version_node) <version_timestamp> "%d" .`,
			version,
			time.Now().Unix(),
		),
	}

	req := &api.Request{
		Query:     q,
		Mutations: []*api.Mutation{mu},
		CommitNow: true,
	}

	txn := client.NewTxn()
	defer txn.Discard(ctx)

	if _, err := txn.Do(ctx, req); err != nil {
		return err
	}

	return nil
}
