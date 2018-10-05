package migration

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

const tokenDepends = "DEPENDS"
const tokenUp = "UP"
const tokenDown = "DOWN"
const sectionBegin = "-- gorph %s"
const sectionEnd = "-- end --"

type FileMigration struct {
	Migration
	Path string
}

func (f *FileMigration) getType() uint8 {
	ext := strings.ToLower(filepath.Ext(f.Path))
	return TypeMappings[ext]
}

func (f *FileMigration) extractSection(section string) (string, error) {
	var line string
	var buffer string
	var parsing = false

	token := fmt.Sprintf(sectionBegin, section)

	file, err := os.Open(f.Path)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		if strings.HasPrefix(line, token) {
			line = strings.Replace(line, token, "", 1)
			parsing = true
		}
		if parsing {
			buffer += line
			if strings.Contains(buffer, sectionEnd) {
				buffer = strings.Replace(buffer, sectionEnd, "", 1)
				break
			}
		}
	}
	return buffer, nil
}

func (f *FileMigration) SetDependencies() error {

	depends, err := f.extractSection(tokenDepends)
	if err != nil {
		log.WithError(err).Error("Unable to parse dependencies")
		return err
	}

	if depends != "" {
		for _, d := range strings.Split(depends, ",") {
			f.Migration.Depends = append(f.Migration.Depends, strings.Replace(d, " ", "", -1))
		}
		log.Debugf("Parsed dependencies as : %v", f.Migration.Depends)
	}

	return nil
}

func (f *FileMigration) UpScript() (string, uint8, error) {
	// fileMigration has the type set based on file extension.
	script, err := f.extractSection(tokenUp)
	return script, f.getType(), err
}

func (f *FileMigration) DownScript() (string, uint8, error) {
	// fileMigration has the type set based on file extension.
	script, err := f.extractSection(tokenDown)
	return script, f.getType(), err
}

func (f *FileMigration) Up(tx *sqlx.Tx) error {
	return nil
}

func (f *FileMigration) Down(tx *sqlx.Tx) error {
	return nil
}
