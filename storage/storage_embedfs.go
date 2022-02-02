package storage

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/toudi/gorpheus"
	"github.com/toudi/gorpheus/migration"
)

type EmbedFsMigration struct {
	migration.Migration
}

func RegisterEmbedFS(target *gorpheus.Collection, f *embed.FS, namespace string) error {
	var ext string

	fs.WalkDir(f, ".", func(path string, info fs.DirEntry, err error) error {
		if !info.IsDir() {
			//fmt.Printf("Found file: %s, %+v\n", path, info)

			ext = strings.ToLower(filepath.Ext(path))

			// because reading from embed.FS is cheaper than reading from a real file we can
			// parse the dependencies as we go which will speed up the indexing process.
			reader, err := f.Open(path)
			if err != nil {
				return fmt.Errorf("could not open embed.FS file for reading: %v", err)
			}

			target.Register(&EmbedFsMigration{
				Migration: migration.Migration{
					Version: strings.Replace(namespace+"/"+info.Name(), ext, "", 1),
					Type:    migration.TypeMappings[ext],
					Reader:  reader.(io.ReadSeekCloser),
				},
			})
		}
		return nil
	})

	return nil
}
