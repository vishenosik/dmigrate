package migrate

import (
	"context"
	"encoding/json"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
)

func applySchema(client *dgo.Dgraph, ctx context.Context) error {
	return client.Alter(ctx, &api.Operation{
		Schema: `
		version_index_name: string @index(exact) .
		version_timestamp: datetime .
		version_current: int .
		type SchemaVersion {
			version_index_name: string
			version_timestamp: datetime	
			version_current: int
		}`,
	})
}

func fetchVersion(client *dgo.Dgraph, ctx context.Context) (Version, error) {

	q := `query {
		current_version(func: eq(version_index_name, "current schema version")) {
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

func (dmr *dgraphMigrator) upsertVersion(ctx context.Context) error {

	txn := dmr.client.NewTxn()
	defer txn.Discard(ctx)

	/*
	   query = `
	   	query {
	   		user as var(func: eq(email, "wrong_email@dgraph.io"))
	   	}`
	     mu := &api.Mutation{
	   	SetNquads: []byte(`uid(user) <email> "correct_email@dgraph.io" .`),
	     }
	     req := &api.Request{
	   	Query: query,
	   	Mutations: []*api.Mutation{mu},
	   	CommitNow:true,
	     }

	     // Update email only if matching uid found.
	     _, err := dg.NewTxn().Do(ctx, req)
	     // Check error

	*/

	_ = `
	upsert {
  		query {
  		  	q(func: eq(email, "user@company1.io")) {
  		  	  	v as uid
  		  	  	name
  		  	}
  		}

  		mutation {
  		  	set {
  		  	  	uid(v) <name> "first last" .
  		  	  	uid(v) <email> "user@company1.io" .
  		  	}
  		}
	}`

	return nil
}
