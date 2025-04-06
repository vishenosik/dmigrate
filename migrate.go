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

type MigratorOption func(*dgraphMigrator)

type dgraphMigrator struct {
	client         *dgo.Dgraph
	fsys           fs.FS
	currentVersion int64
	log            Logger
}

func NewDgraphMigrator(client *dgo.Dgraph, fsys fs.FS, opts ...MigratorOption) (*dgraphMigrator, error) {
	return NewDgraphMigratorContext(context.Background(), client, fsys, opts...)
}

func NewDgraphMigratorContext(
	ctx context.Context,
	client *dgo.Dgraph,
	fsys fs.FS,
	opts ...MigratorOption,
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

	dmr := defaultMigrator(client, fsys, version.CurrentVersion)

	for _, opt := range opts {
		opt(dmr)
	}

	return dmr, nil
}

func defaultMigrator(
	client *dgo.Dgraph,
	fsys fs.FS,
	currentVersion int64,
) *dgraphMigrator {
	return &dgraphMigrator{
		client:         client,
		fsys:           fsys,
		currentVersion: currentVersion,
		log:            &stdLogger{},
	}
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
		dmr.log.Fatalf("collect filenames failed", err)
		return err
	}

	migrations := migrationsToApply(filenamesIter, dmr.currentVersion, toVersion)

	lastVersion := int64(0)
	for migration := range migrations {

		schemaUp, err := readUpMigration(dmr.fsys, migration.filename)
		if err != nil {
			dmr.log.Fatalf("failed to read migration file", err)
			return err
		}

		if err := upVersion(ctx, dmr.client, migration.version, schemaUp); err != nil {
			dmr.log.Fatalf("failed to update version", err)
			return err
		}
		lastVersion = migration.version
	}

	dmr.log.Printf("migrated successfully", "current version", lastVersion)

	return nil
}
