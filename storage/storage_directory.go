package storage

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/toudi/gorpheus"
	"github.com/toudi/gorpheus/migration"
)

func ScanDirectory(dir string, target *gorpheus.Collection) error {
	var namespace string = "default"
	var ext string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if info.Name() != dir {
				namespace = info.Name()
			}
		} else {

			log.Debugf("Found file: %s", path)
			ext = strings.ToLower(filepath.Ext(path))

			log.Debugf("Extension: %v", ext)

			m := &migration.FileMigration{
				Path: path,
				Migration: migration.Migration{
					Namespace: namespace,
					Version:   strings.Replace(info.Name(), ext, "", 1),
				},
			}

			target.Register(m)
		}

		return nil
	})

	return err
}
