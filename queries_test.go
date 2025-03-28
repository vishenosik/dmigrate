package migrate

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_queries(t *testing.T) {

	suite := newClientSuite(t)
	defer suite.cancel()

	client := suite.client
	ctx := suite.ctx

	err := applySchema(ctx, client)
	assert.NoError(t, err)

	schema := getSchema()

	const expectedVersion int64 = 3

	err = upVersion(ctx, client, expectedVersion, schema)
	assert.NoError(t, err)

	version, err := fetchVersion(ctx, client)
	assert.NoError(t, err)
	assert.Equal(t, expectedVersion, version.CurrentVersion)

	userID := "321"
	userNickname := "nickname"
	userEmail := "email@email.dev"

	if err = saveUser(ctx, client, userID, userNickname, userEmail, []byte(userEmail)); err != nil {
		t.Fatal(err)
	}

	user, err := userByEmail(ctx, client, userEmail)
	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, userNickname, user.Nickname)
	assert.Equal(t, userEmail, user.Email)

}

type User struct {
	UID          string   `json:"uid,omitempty"`
	Nickname     string   `json:"nickname,omitempty"`
	Email        string   `json:"email,omitempty"`
	ID           string   `json:"uuid,omitempty"`
	PasswordHash []byte   `json:"pass_hash,omitempty"`
	DType        []string `json:"dgraph.type,omitempty"`
}

func getSchema() []byte {
	return []byte(`
	uuid: string @index(exact) .
	nickname: string @index(exact) .
	email: string @index(exact) .
	pass_hash: string .

	type User {
		uuid: string
		nickname: string
		email: string
		pass_hash: string
	}`)
}

// SaveUser saves user to db.
func saveUser(ctx context.Context, client *dgo.Dgraph, id, nickname, email string, passHash []byte) error {
	const op = "Store.Dgraph.SaveUser"

	variables := map[string]string{
		"$email": email,
	}

	q := `query UserByEmail($email: string){
		users(func: eq(email, $email)) {
			uid
		}
	}`

	txn := client.NewTxn()
	defer txn.Discard(ctx)

	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return errors.Wrap(err, op)
	}

	type Root struct {
		Users []User `json:"users"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return errors.Wrap(err, op)
	}

	if len(r.Users) != 0 {
		return errors.New("already exist")
	}

	user := &User{
		Nickname:     nickname,
		Email:        email,
		ID:           id,
		PasswordHash: passHash,
	}

	userPB, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, op)
	}

	mu := &api.Mutation{
		CommitNow: true,
		SetJson:   userPB,
	}

	// TODO: apply metrics here
	_, err = txn.Mutate(ctx, mu)
	if err != nil {
		// TODO: handle error
		return errors.Wrap(err, op)
	}

	return nil
}

func userByEmail(ctx context.Context, client *dgo.Dgraph, email string) (User, error) {
	const op = "Store.Dgraph.UserByEmail"

	variables := map[string]string{
		"$email": email,
	}

	q := `query UserByEmail($email: string){
		users(func: eq(email, $email)) {
			uuid
			nickname
			email
			pass_hash
		}
	}`

	txn := client.NewTxn()
	defer txn.Discard(ctx)

	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return User{}, errors.Wrap(err, op)
	}

	type Root struct {
		Users []User `json:"users"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return User{}, errors.Wrap(err, op)
	}

	if len(r.Users) == 0 {
		return User{}, errors.New("not found")
	}

	user := r.Users[0]

	return User{
		Nickname:     user.Nickname,
		Email:        user.Email,
		ID:           user.ID,
		PasswordHash: user.PasswordHash,
	}, nil
}
