package migrate

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vishenosik/dmigrate/pkg/collections"
)

func Test_collectMigrations(t *testing.T) {

	suite := NewSuite()

	expected := migrations{
		{version: 1, filename: suite.validFilenames[0]},
		{version: 2, filename: suite.validFilenames[1]},
		{version: 3, filename: suite.validFilenames[2]},
	}

	filenamesIter := collections.Iter(suite.filenames)

	actual := slices.Collect(collectMigrations(filenamesIter))
	assert.Equal(t, expected, actual)

}

func Test_migrationsToApply(t *testing.T) {

	suite := NewSuite()

	filenamesIter := collections.Iter(suite.validFilenames)

	expected := migrations{
		{version: 1, filename: suite.validFilenames[0]},
		{version: 2, filename: suite.validFilenames[1]},
		{version: 3, filename: suite.validFilenames[2]},
	}

	t.Run("first full migration", func(t *testing.T) {
		actual := slices.Collect(migrationsToApply(filenamesIter, 0, 0))
		assert.Equal(t, expected, actual)
	})

	t.Run("first migration to version 2", func(t *testing.T) {
		actual := slices.Collect(migrationsToApply(filenamesIter, 0, 2))
		assert.Equal(t, expected[:2], actual)
	})

	t.Run("full migration from 1", func(t *testing.T) {
		actual := slices.Collect(migrationsToApply(filenamesIter, 1, 0))
		assert.Equal(t, expected[1:], actual)
	})

	t.Run("no migration (current >= toVersion)", func(t *testing.T) {
		actual := slices.Collect(migrationsToApply(filenamesIter, 2, 1))
		assert.Len(t, actual, 0)
	})

}

func Test_parseVersion(t *testing.T) {

	testingTable := []struct {
		name            string
		filename        string
		expectedVersion int64
		expectedOk      bool
	}{
		{
			name:            "parse_version_success_1",
			filename:        "0002_add_user_table.gql",
			expectedVersion: 2,
			expectedOk:      true,
		},
		{
			name:            "parse_version_success_2",
			filename:        "0001_create_initial_schema.gql",
			expectedVersion: 1,
			expectedOk:      true,
		},
		{
			name:            "parse_version_error_1",
			filename:        "0000_invalid_version.gql",
			expectedVersion: 0,
			expectedOk:      false,
		},
		{
			name:            "parse_version_error_invalid_filename_1",
			filename:        "invalid_filename.gql",
			expectedVersion: 0,
			expectedOk:      false,
		},
		{
			name:            "parse_version_error_invalid_filename_2",
			filename:        "invalidFilename.gql",
			expectedVersion: 0,
			expectedOk:      false,
		},
		{
			name:            "parse_version_error_invalid_extension",
			filename:        "0002_add_user_table.sql",
			expectedVersion: 0,
			expectedOk:      false,
		},
	}

	for _, tt := range testingTable {

		t.Run(tt.name, func(t *testing.T) {
			actualVersion, ok := parseVersion(tt.filename)
			assert.Equal(t, ok, tt.expectedOk)
			assert.Equal(t, actualVersion, tt.expectedVersion)
		})

	}

}
