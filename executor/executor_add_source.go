package executor

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/toudi/gorpheus/v2/interfaces"
)

func (e *Executor) AddMigrationsFromPath(directory string) error {
	// let's scan the directories recursively and find yaml files.
	// for each of them, let's add migration and set the namespace as
	// the directory name that contains said yaml files.
	var visitedNamespaces = make(map[string]bool)

	return filepath.WalkDir(directory, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// that's a directory, we're only interested in files.
		if info.IsDir() {
			return nil
		}

		var namespace = filepath.Base(filepath.Dir(path))
		var reader *os.File

		// that's what we're looking for!
		if strings.ToLower(filepath.Ext(info.Name())) == ".yaml" {
			if !visitedNamespaces[namespace] {
				e.logger.Log(interfaces.LogMessage{
					Level:   slog.LevelDebug,
					Message: "discovered new namespace",
					Payload: []any{slog.String("namespace", namespace)},
				})
				visitedNamespaces[namespace] = true
			}
			if reader, err = os.Open(path); err != nil {
				// unable to open file for reading
				return err
			}
			defer reader.Close()
			err = e.AddSourceToNamespace(namespace, filepath.Base(path), reader)
		}
		return err
	})
}
