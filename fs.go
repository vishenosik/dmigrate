package migrate

import (
	"io"
	"io/fs"
	"iter"
	"path"

	"github.com/pkg/errors"
	"github.com/vishenosik/dmigrate/pkg/collections"
)

func collectFilenames(fsys fs.FS, dirpath string) (iter.Seq[string], error) {
	if _, err := fs.Stat(fsys, dirpath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errors.Wrap(err, "migrations directory not found")
		}

		return nil, errors.Wrap(err, "migrations directory unexpected "+dirpath)
	}

	filenames, err := fs.Glob(fsys, path.Join(dirpath, "*"+gqlExt))
	if err != nil {
		return nil, errors.Wrap(err, "migrations not found")
	}

	return collections.Iter(filenames), nil
}

func readUpMigration(fsys fs.FS, filepath string) ([]byte, error) {
	schemaFile, err := fsys.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer schemaFile.Close()

	schema, err := io.ReadAll(schemaFile)
	if err != nil {
		return nil, err
	}

	return schema, nil
}
