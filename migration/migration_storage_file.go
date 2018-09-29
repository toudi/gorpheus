package migration

import (
	"bufio"
	"database/sql"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

const sectionDepends = "-- gorph DEPENDS"
const sectionEnd = "-- end --"

type FileMigration struct {
	Migration
	Path string
}

func (f *FileMigration) Parse() error {
	log.Debugf("Opening %v\n", f.Path)

	file, err := os.Open(f.Path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()
	var line string
	var buffer string
	var parsing bool = false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		if strings.HasPrefix(line, sectionDepends) {
			log.Debugf("Has depends section")
			line = strings.Replace(line, sectionDepends, "", 1)
			parsing = true
		}
		if parsing {
			buffer += line
			if strings.Contains(buffer, sectionEnd) {
				buffer = strings.Replace(buffer, sectionEnd, "", 1)
				parsing = false
			}
		}
	}

	depends := strings.Split(buffer, ",")
	for _, d := range depends {
		f.Migration.Depends = append(f.Migration.Depends, strings.Replace(d, " ", "", -1))
	}

	log.Debugf("PArsed dependencies as : %v", f.Migration.Depends)
	return nil
}

func (f *FileMigration) UpFn(tx *sql.Tx) error {
	return nil
}

func (f *FileMigration) DownFn(tx *sql.Tx) error {
	return nil
}
