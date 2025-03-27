package migrate

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"
)

const node = "dmigrate_version_node"

func applySchema(ctx context.Context, client *dgo.Dgraph) error {
	return client.Alter(ctx, &api.Operation{
		Schema: `
		version_index_name: string @index(exact) @upsert .
		version_timestamp: int @index(int) @upsert .
		version_current: int @index(int) @upsert .
		type SchemaVersion {
			version_index_name: string
			version_timestamp: int
			version_current: int
		}`,
	})
}

func fetchVersion(ctx context.Context, client *dgo.Dgraph) (Version, error) {

	vars := map[string]string{"$node": node}
	q := `query node_name($node: string){
		current_version(func: eq(version_index_name, $node)) {
			version_timestamp	
			version_current
		}
	}`

	txn := client.NewTxn()
	defer txn.Discard(ctx)

	resp, err := txn.QueryWithVars(ctx, q, vars)
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

	vars := map[string]string{"$node": node}
	q := `
		query node_name($node: string){
			version_node as var(func: eq(version_index_name, $node)) {}
		}`

	mu := &api.Mutation{
		SetNquads: []byte(fmt.Sprintf(`
		    uid(version_node) <version_index_name> "%s" .
			uid(version_node) <version_current> "%d" .
			uid(version_node) <version_timestamp> "%d" .`,
			node,
			version,
			time.Now().Unix(),
		)),
	}

	req := &api.Request{
		Query:     q,
		Vars:      vars,
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
