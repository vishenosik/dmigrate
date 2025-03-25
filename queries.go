package migrate

import (
	"context"
	"encoding/json"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
)

func applySchema(ctx context.Context, client *dgo.Dgraph) error {
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

func fetchVersion(ctx context.Context, client *dgo.Dgraph) (Version, error) {

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

func upsertVersion(ctx context.Context, client *dgo.Dgraph, version int64) error {

	txn := client.NewTxn()
	defer txn.Discard(ctx)

	/*
		q := `query all($a: string) {
		    all(func: eq(name, $a)) {
		      name
		    }
		  }`

		res, err := txn.QueryWithVars(ctx, q, map[string]string{"$a": "Alice"})
		fmt.Printf("%s\n", res.Json)

		req := &api.Request{
		  Query: q,
		  Vars: map[string]string{"$a": "Alice"},
		}
		res, err := txn.Do(ctx, req)
		// Check error
		fmt.Printf("%s\n", res.Json)


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

		query UserByEmail($email: string){
			user as var(func: eq(email, $email))
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
