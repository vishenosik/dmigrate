package migrate

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UpTo(t *testing.T) {

	suite := newClientSuite(t)
	defer suite.cancel()

	migrator, err := NewDgraphMigrator(
		Config{
			User:     user,
			Password: password,
			Host:     host,
			Port:     port,
		},
		test_migrations,
	)
	require.NoError(t, err)

	testdir := path.Join("test", "migrations")

	err = migrator.UpTo(testdir, 1)
	require.NoError(t, err)

	const expectedVersion int64 = 1

	version, err := fetchVersion(suite.ctx, suite.client)
	require.NoError(t, err)
	assert.Equal(t, expectedVersion, version.CurrentVersion)

	userID := "321"
	userNickname := "nickname"
	userEmail := "email@email.dev"

	if err = saveUser(suite.ctx, suite.client, userID, userNickname, userEmail, []byte(userEmail)); err != nil {
		t.Fatal(err)
	}

	user, err := userByEmail(suite.ctx, suite.client, userEmail)
	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, userNickname, user.Nickname)
	assert.Equal(t, userEmail, user.Email)

}
