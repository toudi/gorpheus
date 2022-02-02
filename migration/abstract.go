package migration

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	TypeSQL  = iota
	TypeFizz = iota
	TypeGo   = iota
)

const tokenDepends = "DEPENDS"
const tokenUp = "UP"
const tokenDown = "DOWN"
const sectionBegin = "-- gorph %s"
const sectionEnd = "-- end --"

var TypeMappings = map[string]uint8{
	".sql":  TypeSQL,
	".fizz": TypeFizz,
}

type Migration struct {
	Type    uint8
	Version string
	Depends []string
	Reader  io.ReadSeekCloser
}

func (m Migration) GetVersion() string {
	return m.Version
}

func (m *Migration) GetVersionNumber() (int, error) {
	start, err := m.GetNamespace()
	if err != nil {
		return -1, fmt.Errorf("could not extract namespace: %v", err)
	}
	if underscoreIdx := strings.Index(m.Version, "_"); underscoreIdx != -1 {
		versionNumber, err := strconv.Atoi(m.Version[len(start)+1 : underscoreIdx])
		if err != nil {
			return -1, fmt.Errorf("could not parse version number")
		}
		return versionNumber, nil
	}
	return -1, errors.New("could not match version number")
}

func (m Migration) GetNamespace() (string, error) {
	if slashIndex := strings.Index(m.Version, "/"); slashIndex != -1 {
		return m.Version[:slashIndex], nil
	}
	return "", errors.New("could not find namespace in version")
}

func (m *Migration) Dependencies() []string {
	return m.Depends
}

func (m *Migration) Up(tx *sqlx.Tx) error {
	return nil
}

func (m *Migration) Down(tx *sqlx.Tx) error {
	return nil
}

func (m *Migration) extractSection(section string) (string, error) {
	m.Reader.Seek(0, io.SeekStart)

	var line string
	var buffer string
	var parsing = false

	token := fmt.Sprintf(sectionBegin, section)

	scanner := bufio.NewScanner(m.Reader)
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

func (m *Migration) GetType() uint8 {
	return m.Type
}

func (m *Migration) UpScript() (string, uint8, error) {
	// fileMigration has the type set based on file extension.
	script, err := m.extractSection(tokenUp)
	return script, m.Type, err
}

func (m *Migration) DownScript() (string, uint8, error) {
	// fileMigration has the type set based on file extension.
	script, err := m.extractSection(tokenDown)
	return script, m.Type, err
}

func (m *Migration) SetDependencies() error {
	if len(m.Depends) == 0 {
		depends, err := m.extractSection(tokenDepends)
		if err != nil {
			return fmt.Errorf("unable to parse dependencies: %v", err)
		}

		if depends != "" {
			for _, d := range strings.Split(depends, ",") {
				m.Depends = append(m.Depends, strings.Replace(d, " ", "", -1))
			}
		}
	}
	return nil
}

func (m *Migration) Close() error {
	// if this is an in-memory migration then it doesn't rely on external reader
	// therefore the reader can be nil
	if m.Reader != nil {
		return m.Reader.Close()
	}
	return nil
}

type MigrationI interface {
	UpScript() (string, uint8, error)
	DownScript() (string, uint8, error)
	Up(tx *sqlx.Tx) error
	Down(tx *sqlx.Tx) error
	GetVersion() string
	GetVersionNumber() (int, error)
	GetNamespace() (string, error)
	Dependencies() []string
	SetDependencies() error
	GetType() uint8
	Close() error
}
