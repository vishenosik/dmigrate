package migrate

import (
	"context"
	"io/fs"

	"github.com/pkg/errors"

	"github.com/dgraph-io/dgo/v240"
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
	return dmr.UpToContext(context.Background(), path, 0)
}

func (dmr *dgraphMigrator) UpContext(ctx context.Context, path string) error {
	return dmr.UpToContext(ctx, path, 0)
}

func (dmr *dgraphMigrator) UpTo(path string, toVersion int64) error {
	return dmr.UpToContext(context.Background(), path, toVersion)
}

func (dmr *dgraphMigrator) UpToContext(ctx context.Context, path string, toVersion int64) error {

	filenamesIter, err := collectFilenames(dmr.fsys, path)
	if err != nil {
		return err
	}

	migrations := migrationsToApply(filenamesIter, dmr.currentVersion, toVersion)

	for migration := range migrations {

		schemaUp, err := readUpMigration(dmr.fsys, migration.filename)
		if err != nil {
			return err
		}

		if err := upVersion(ctx, dmr.client, migration.version, schemaUp); err != nil {
			return err
		}
	}

	return nil
}
