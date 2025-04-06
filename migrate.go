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
	cancel         cancelFunc
}

func NewDgraphMigrator(
	config Config,
	fsys fs.FS,
	opts ...MigratorOption,
) (*dgraphMigrator, error) {
	return NewDgraphMigratorContext(context.Background(), config, fsys, opts...)
}

func NewDgraphMigratorContext(
	ctx context.Context,
	config Config,
	fsys fs.FS,
	opts ...MigratorOption,
) (*dgraphMigrator, error) {

	client, cancel := mustConnect(ctx, config)

	if err := applySchema(ctx, client); err != nil {
		return nil, err
	}

	version, err := fetchVersion(ctx, client)
	if err != nil && !errors.Is(err, ErrVersionFetch) {
		return nil, err
	}

	dmr := defaultMigrator(client, fsys, version.CurrentVersion, cancel)

	for _, opt := range opts {
		opt(dmr)
	}

	return dmr, nil
}

func defaultMigrator(
	client *dgo.Dgraph,
	fsys fs.FS,
	currentVersion int64,
	cancel cancelFunc,
) *dgraphMigrator {
	return &dgraphMigrator{
		client:         client,
		fsys:           fsys,
		currentVersion: currentVersion,
		log:            &stdLogger{},
		cancel:         cancel,
	}
}

// Up applies all available migrations from the specified path.
// It uses a background context and applies all migrations regardless of version.
// Returns an error if any migration fails.
func (dmr *dgraphMigrator) Up(path string) error {
	return dmr.UpToContext(context.Background(), path, 0)
}

// UpContext applies all available migrations from the specified path using the provided context.
// It applies all migrations regardless of version.
// Returns an error if any migration fails.
func (dmr *dgraphMigrator) UpContext(ctx context.Context, path string) error {
	return dmr.UpToContext(ctx, path, 0)
}

// UpTo applies migrations from the specified path up to the given version (inclusive).
// It uses a background context.
// Returns an error if any migration fails or if the specified version doesn't exist.
func (dmr *dgraphMigrator) UpTo(path string, toVersion int64) error {
	return dmr.UpToContext(context.Background(), path, toVersion)
}

// UpToContext applies migrations from the specified path up to the given version (inclusive) using the provided context.
// It handles the migration process including:
//   - Collecting migration filenames
//   - Determining which migrations need to be applied
//   - Executing each migration in order
//   - Updating version tracking
//
// Parameters:
//
//	ctx: Context for cancellation and timeouts
//	path: Directory path where migration files are located
//	toVersion: Target version to migrate to (0 means apply all)
//
// Returns:
//
//	error: Any error that occurred during migration, or nil if successful
func (dmr *dgraphMigrator) UpToContext(ctx context.Context, path string, toVersion int64) error {
	defer dmr.cancel()

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
