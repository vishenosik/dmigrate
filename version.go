package migrate

import (
	"iter"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/vishenosik/dmigrate/pkg/collections"
)

type Version struct {
	Timestamp      int64    `json:"version_timestamp,omitempty"`
	CurrentVersion int64    `json:"version_current,omitempty"`
	DType          []string `json:"dgraph.type,omitempty"`
}

type migration struct {
	version  int64
	filename string
}

type migrations = []migration

func migrationsToApply(
	filenamesIter iter.Seq[string],
	currentVersion int64,
	toVersion int64,
) iter.Seq[migration] {

	filtered := collections.Filter(
		collectMigrations(filenamesIter),
		func(m migration) bool {
			return m.version > currentVersion &&
				((toVersion <= 0 || m.version <= toVersion) ||
					(toVersion > 0 && m.version <= toVersion))
		},
	)

	sorted := slices.SortedFunc(
		filtered,
		func(m1, m2 migration) int {
			return int(m1.version - m2.version)
		},
	)

	return collections.Iter(sorted)
}

func collectMigrations(filenames iter.Seq[string]) iter.Seq[migration] {
	return func(yield func(migration) bool) {
		for filename := range filenames {
			version, ok := parseVersion(filename)
			if !ok {
				continue
			}

			_migration := migration{
				version:  version,
				filename: filename,
			}

			if !yield(_migration) {
				return
			}
		}
	}

}

func parseVersion(filename string) (int64, bool) {

	base := filepath.Base(filename)
	if ext := filepath.Ext(base); ext != gqlExt {
		return 0, false
	}

	idx := strings.Index(base, "_")
	if idx < 0 {
		return 0, false
	}

	n, err := strconv.ParseInt(base[:idx], 10, 64)
	if err != nil {
		return 0, false
	}

	if n < 1 {
		return 0, false
	}

	return n, true
}
