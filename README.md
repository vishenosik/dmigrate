# dmigrate
golang dgraph migrator package

## Usage

### Files

Files should be named as *.gql - will add more extentions later. Best way to pass fsys to migrator is to embed.FS.

### Code

See basic usage with up migration in [tests](migrate_test.go)

```go
    // migrator init
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

    // migrations path 
	testdir := path.Join("test", "migrations")

    // running migration
	err = migrator.UpTo(testdir, 1)
	require.NoError(t, err)
```

## Docs

* [CHANGELOG](docs/CHANGELOG.md)
* [CONTRIBUTING](docs/CONTRIBUTING.md)
* [RELEASING](docs/RELEASING.md)
* [LICENSE](docs/LICENSE)

## Tools

* [Taskfile](https://taskfile.dev/)
* Linter:

```bash
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest
```