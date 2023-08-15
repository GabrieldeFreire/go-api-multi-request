package log

import (
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

var lock = &sync.Mutex{}

var logger *slog.Logger

func GetInstance() *slog.Logger {
	if logger == nil {
		lock.Lock()
		defer lock.Unlock()
		if logger == nil {
			logger = getLogger()
		}
	}

	return logger
}

func getLogger() *slog.Logger {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		// Remove the directory from the source's filename.
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
			source.Function = filepath.Base(source.Function)
		}
		return a
	}

	opts := &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: replace,
	}

	var handler slog.Handler = slog.NewJSONHandler(os.Stdout, opts)

	logger := slog.New(handler)
	return logger
}
