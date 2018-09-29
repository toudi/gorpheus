package storage

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/toudi/gorpheus"
	"github.com/toudi/gorpheus/dialect"
)

type FileStorage struct {
	Path    string
	Factory dialect.DialectFactory
	Id      gorpheus.Id
}

func (f *FileStorage) ID() gorpheus.Id {
	return f.Id
}

func (f *FileStorage) Execute(tx *sqlx.Tx) error {
	return nil
}

func (f *FileStorage) Depends() *gorpheus.Id {
	buffer := f.extractData("depends")
	if buffer != nil {
		log.Debugf("buff: %s", buffer.String())
		return gorpheus.Id_from_string(strings.TrimSpace(buffer.String()))
	}
	return nil
}

func (f *FileStorage) Up() gorpheus.MigrationI {
	// log.Debugf("Section parsed: %v", buffer.String())
	// return f.Factory(buffer)
	return nil
}

func (f *FileStorage) extractData(sequence string) *bytes.Buffer {
	var copyToBuffer bool = false
	file, err := os.Open(f.Path)
	if err != nil {
		return nil
	}
	defer file.Close()
	buffer := bytes.NewBuffer([]byte{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		copyToBuffer = copyToBuffer
		// log.Debugf("Copy to buffer: %v", copyToBuffer)
		line := scanner.Text()
		if strings.Contains(line, "-- gorph "+strings.ToUpper(sequence)) {
			copyToBuffer = true
		} else if strings.Contains(line, "-- end --") {
			copyToBuffer = false
		} else if copyToBuffer == true {
			// log.Debugf("Read a single line: %s", line)
			buffer.WriteString(line)
		}
	}

	if buffer.Len() == 0 {
		return nil
	}

	return buffer

}

func (f *FileStorage) Down() gorpheus.MigrationI {
	return nil
}
