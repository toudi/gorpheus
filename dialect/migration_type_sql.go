package dialect

import (
	"bytes"

	"github.com/jmoiron/sqlx"
	"github.com/toudi/gorpheus"
)

type SQLMigration struct {
	gorpheus.BaseMigration
	// sourceFile string
	// fileInfo   *storage.FileMigration
}

func SQLMigration_init(body bytes.Buffer) gorpheus.MigrationI {
	return nil
}
func (s *SQLMigration) Depends() *gorpheus.Id {
	return nil
}

func (s *SQLMigration) Execute(tx *sqlx.Tx) error {
	return nil
}

func (s *SQLMigration) ID() gorpheus.Id {
	return s.Id
}

// func SQLMigration_init(path string, info os.FileInfo, namespace string) (*SQLMigration, error) {
// 	var out *SQLMigration
// 	var err error

// 	baseName := info.Name()

// 	out = &SQLMigration{
// 		Id: gorpheus.Id{
// 			Version:   strings.TrimSuffix(baseName, filepath.Ext(baseName)),
// 			Namespace: namespace,
// 		},
// 		sourceFile: path,
// 		// fileInfo:   &FileMigration{},
// 	}

// 	err = out.Parse()

// 	return out, err
// }

// func (s *SQLMigration) ID() gorpheus.Id {
// 	return s.Id
// }

// func (s *SQLMigration) Depends() *gorpheus.Id {
// 	return nil
// 	// s.fileInfo.Depends
// }

// func (s *SQLMigration) Parse() error {
// 	log.Debugf("Parsing source file: %s", s.sourceFile)
// 	// TODO: use NewReader?
// 	_, err := ioutil.ReadFile(s.sourceFile)
// 	if err != nil {
// 		log.Error("Cannot open source file")
// 		return err
// 	}
// 	// err = ParseContent(data, s.fileInfo)
// 	if err != nil {
// 		log.Error("Error parsing source file")
// 		return err
// 	}
// 	return nil
// }
