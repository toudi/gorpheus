package storage

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/toudi/gorpheus"
	"github.com/toudi/gorpheus/migration"
)

func RegisterEmbedFS(namespace string, f *embed.FS, target *gorpheus.Collection) error {
	var ext string

	fs.WalkDir(f, ".", func(path string, info fs.DirEntry, err error) error {
		if !info.IsDir() {
			//fmt.Printf("Found file: %s, %+v\n", path, info)

			ext = strings.ToLower(filepath.Ext(path))

			m := &migration.FileMigration{
				Path: path,
				Migration: migration.Migration{
					Namespace: namespace,
					Version:   strings.Replace(info.Name(), ext, "", 1),
				},
				EmbeddedFs: f,
			}
			target.Register(m)
		}
		return nil
	})

	return nil
}
