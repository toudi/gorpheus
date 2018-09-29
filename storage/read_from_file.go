package storage

import (
	"errors"
	"strings"

	"github.com/toudi/gorpheus"
)

const sectionUp = "-- gorph UP"
const sectionDown = "-- gorph DOWN"
const sectionDepends = "-- gorph DEPENDS"
const sectionEnd = "-- end --"

type SectionIndexes struct {
	Up      [2]int
	Down    [2]int
	Depends [2]int
}

type FileMigration struct {
	Depends *gorpheus.Id
	Up      []byte
	Down    []byte
}

var UpSectionNotFoundErr = errors.New("Unable to parse UP migration")
var DownSectionNotFoundErr = errors.New("Unable to parse DOWN migration")

func ParseContent(content []byte, dest *FileMigration) error {
	indexes := SectionIndexes{}

	var parsing *[2]int

	for pos := 0; pos <= len(content); pos++ {
		if pos >= 9 {
			parsed := string(content[pos-9 : pos])

			if parsed == "-- gorph " {
				// we're parsing either UP, DOWN or DEPENDS.
				// since DEPENDS is the longest one, try to parse this and
				// take it from there.
				if string(content[pos:pos+2]) == "UP" {
					indexes.Up[0] = pos + 2 + 1
					parsing = &indexes.Up
				} else if string(content[pos:pos+4]) == "DOWN" {
					indexes.Down[0] = pos + 4 + 1
					parsing = &indexes.Down
				} else if string(content[pos:pos+7]) == "DEPENDS" {
					indexes.Depends[0] = pos + 7 + 1
					parsing = &indexes.Depends
				}
			} else if parsed == sectionEnd {
				if parsing != nil {
					parsing[1] = pos - 9
				}
			}
		}
	}

	if indexes.Up[1] == 0 {
		return UpSectionNotFoundErr
	}
	if indexes.Down[1] == 0 {
		return DownSectionNotFoundErr
	}
	dest.Up = content[indexes.Up[0]:indexes.Up[1]]
	dest.Down = content[indexes.Down[0]:indexes.Down[1]]

	if indexes.Depends[1] != 0 {
		dest.Depends = &gorpheus.Id{}
		depends := string(content[indexes.Depends[0]:indexes.Depends[1]])
		// let's parse the DEPENDS section into Namespace and Version
		parts := strings.Split(depends, "/")
		if len(parts) == 2 {
			dest.Depends.Namespace = strings.Trim(parts[0], " ")
		}
		dest.Depends.Version = strings.Trim(parts[len(parts)-1], " \n")
	}

	return nil
}
