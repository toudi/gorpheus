package storage

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/toudi/gorpheus"
	"github.com/toudi/gorpheus/migration"
)

type FileMigration struct {
	migration.Migration
}

func _register(target *gorpheus.Collection, root string, namespace string, recurse bool) error {
	if namespace == "" {
		namespace = "default"
	}
	var ext string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path != root {
				if !recurse {
					return filepath.SkipDir
				}
				namespace = d.Name()
			}
		} else {
			ext = strings.ToLower(filepath.Ext(path))

			reader, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("could not create os.File reader: %v", err)
			}

			target.Register(&FileMigration{
				Migration: migration.Migration{
					Version: strings.Replace(namespace+"/"+d.Name(), ext, "", 1),
					Reader:  reader,
					Type:    migration.TypeMappings[ext],
				},
			})
		}
		return nil
	})

	return err
}

func RegisterFSRecurse(target *gorpheus.Collection, root string) error {
	return _register(target, root, "", true)
}

func RegisterFS(target *gorpheus.Collection, dir string, namespace string) error {
	return _register(target, dir, namespace, false)
}
