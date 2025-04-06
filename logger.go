package migrate

import "log"

type Logger interface {
	Fatalf(format string, v ...any)
	Printf(format string, v ...any)
}

type stdLogger struct{}

func (*stdLogger) Fatalf(format string, v ...any) {
	log.Fatalf(format, v...)
}

func (*stdLogger) Printf(format string, v ...any) {
	log.Printf(format, v...)
}

func WithLogger(logger Logger) MigratorOption {
	return func(dmr *dgraphMigrator) {
		if logger == nil {
			return
		}
		dmr.log = logger
	}
}
