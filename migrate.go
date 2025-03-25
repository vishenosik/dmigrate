package migrate

import (
	"context"
	"io/fs"

	"github.com/pkg/errors"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
)

const (
	gqlExt = ".gql"
)

var (
	ErrVersionFetch = errors.New("no version fetched")
)

type Logger interface {
	Fatalf(format string, v ...any)
	Printf(format string, v ...any)
}

type dgraphMigrator struct {
	client         *dgo.Dgraph
	fsys           fs.FS
	currentVersion int64
	log            Logger
}

func NewDgraphMigrator(client *dgo.Dgraph, fsys fs.FS) (*dgraphMigrator, error) {
	return NewDgraphMigratorContext(context.Background(), client, fsys)
}

func NewDgraphMigratorContext(
	ctx context.Context,
	client *dgo.Dgraph,
	fsys fs.FS,
) (*dgraphMigrator, error) {

	if client == nil {
		return nil, errors.New("dgraph client not initialized")
	}

	if err := applySchema(ctx, client); err != nil {
		return nil, err
	}

	version, err := fetchVersion(ctx, client)
	if err != nil && !errors.Is(err, ErrVersionFetch) {
		return nil, err
	}

	return &dgraphMigrator{
		client:         client,
		fsys:           fsys,
		currentVersion: version.CurrentVersion,
	}, nil
}

func (dmr *dgraphMigrator) Up(path string) error {
	return dmr.UpContext(context.Background(), path)
}

// 1. Get current version
// 2. Collect .gql files
// 3. Sort provided versions
// 4. Alter versions greater than current
// 5. If success - update version in dgraph
// 6. If failure - Rollback
func (dmr *dgraphMigrator) UpContext(ctx context.Context, path string) error {

	filenamesIter, err := collectFilenames(dmr.fsys, path)
	if err != nil {
		return err
	}

	migrations := migrationsToApply(filenamesIter, dmr.currentVersion)

	for migration := range migrations {

		schemaUp, err := readUpMigration(dmr.fsys, migration.filename)
		if err != nil {
			return err
		}

		op := &api.Operation{
			Schema: string(schemaUp),
		}

		if err := dmr.client.Alter(ctx, op); err != nil {
			return err
		}

		if err := upsertVersion(ctx, dmr.client, migration.version); err != nil {
			return err
		}
	}

	return nil
}
