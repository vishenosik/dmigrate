package migrate

import (
	"embed"
	"path"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	//go:embed test/migrations
	test_migrations embed.FS
)

type fsSuite struct {
	dir            string
	validFilenames []string
	filenames      []string
}

func newFsSuite() fsSuite {

	testdir := path.Join("test", "migrations")

	return fsSuite{
		dir: testdir,
		validFilenames: []string{
			path.Join(testdir, "0001_create_initial_schema.gql"),
			path.Join(testdir, "0002_add_user_table.gql"),
			path.Join(testdir, "0003_add_post_table.gql"),
		},
		filenames: []string{
			path.Join(testdir, "0001_create_initial_schema.gql"),
			path.Join(testdir, "0002_add_user_table.gql"),
			path.Join(testdir, "0003_add_post_table.gql"),
			path.Join(testdir, "000_invalid.gql"),
			path.Join(testdir, "invalid.gql"),
		},
	}

}

func Test_collectFilenames(t *testing.T) {

	suite := newFsSuite()

	actualIter, err := collectFilenames(test_migrations, suite.dir)
	actual := slices.Collect(actualIter)
	assert.NoError(t, err)
	assert.Len(t, actual, 5)
}
