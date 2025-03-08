package interfaces

import "strings"

type DDL interface {
	DDL(buffer *strings.Builder) error
}

type DDLFactory func() DDL
