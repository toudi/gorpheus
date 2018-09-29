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
			// ext = ext[1:len(ext)]
			log.Debugf("Extension: %v", ext)

			m := &migration.FileMigration{
				Path: path,
				Migration: migration.Migration{
					Namespace: namespace,
					Version:   strings.Replace(info.Name(), ext, "", 1),
					Type:      migration.TypeSQL,
					Storage:   migration.StorageFile,
				},
			}

			log.Debugf("Constructed FileMigration: %+v", m)

			target.Register(m)
			// migration := &FileStorage{
			// 	Path: path,
			// 	Id: gorpheus.Id{
			// 		Namespace: namespace,
			// 		Version:   info.Name(),
			// 	},
			// }

			// if ext == "sql" {
			// 	migration.Factory = dialect.SQLMigration_init
			// }

			// target.Register(migration)
			// if upMigration := migration.Up(); upMigration != nil {
			// 	target.Register(upMigration, gorpheus.UP)
			// } else {
			// 	log.Debugf("UP section missing")
			// }
			// if downMigration := migration.Down(); downMigration != nil {
			// 	target.Register(downMigration, gorpheus.DOWN)
			// } else {
			// 	log.Debugf("DOWN section missing")
			// }
		}

		return nil
	})
	if err != nil {
		return err
	}
	// err = SortMigrations()
	return err
}
