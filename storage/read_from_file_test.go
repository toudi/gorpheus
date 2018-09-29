package storage_test

import (
	"testing"

	"github.com/toudi/gorpheus"
)

func Test_file_parse(t *testing.T) {
	content := `
-- gorph DEPENDS users/0001_something -- end --
-- gorph UP
this is UP migration
-- end --

-- gorph DOWN
this is DOWN migration
-- end --
`
	sections, err := gorpheus.ParseContent([]byte(content))
	if err != nil {
		t.Errorf("Error parsing content: %v\n", err)
	}

	if string(sections.Up) != "this is UP migration\n" {
		t.Errorf("Unexpected UP migration: %s\n", string(sections.Up))
	}

	if string(sections.Down) != "this is DOWN migration\n" {
		t.Errorf("Unexpected DOWN migration: %s\n", string(sections.Down))
	}

	if sections.Depends.Namespace != "users" {
		t.Errorf("Invalid namespace parsed from DEPENDS section: %s\n", sections.Depends.Namespace)
	}

	if sections.Depends.Version != "0001_something" {
		t.Errorf("Invalid version parsed from DEPENDS section: %s\n", sections.Depends.Version)
	}
}
