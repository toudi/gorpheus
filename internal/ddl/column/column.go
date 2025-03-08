package column

import (
	"errors"

	"github.com/toudi/gorpheus/v2/internal/ddl/constraint"
)

var (
	ErrPrimaryKeyCannotBeNull     = errors.New("primary key cannot be null")
	ErrInvalidColumnType          = errors.New("invalid column type")
	ErrInvalidMaxLength           = errors.New("invalid max length")
	ErrMissingGeneratedColumnType = errors.New("missing generated column type")
)

type GeneratedColumn struct {
	Definition string `yaml:"definition"`
	Stored     bool   `yaml:"stored"`
}

type Column struct {
	Name string `yaml:"name"`

	// these properties describe things that can be applied to
	// any column
	Null          *bool `yaml:"null"`
	Default       any   `yaml:"default"`
	Index         any   `yaml:"index"`  // either `true` or `unique`
	Unique        *bool `yaml:"unique"` // unique constraint
	AutoIncrement *bool `yaml:"auto-increment"`
	PrimaryKey    *bool `yaml:"primary-key"`

	// column type that will determine rest of the DDL
	TypeAlias string `yaml:"type"` // string-specified alias
	Type      int    `yaml:"-"`
	// type, as detected from `TypeAlias`. Due to zero values,
	// it will default to `TypeInvalid``

	// now let's continue with properties that can be defined per-`Type`
	// applicable for date/time/timestamp
	WithTimeZone bool  `yaml:"tzinfo"`
	Precision    uint8 `yaml:"precision"`

	// applicable for decimal
	MaxDigits     int `yaml:"max-digits"`
	DecimalPlaces int `yaml:"decimal-places"`

	// applicable for varchar
	MaxLength int `yaml:"max-length"`

	// applicable for generated columns
	Generated *GeneratedColumn `yaml:"generated"`

	// applicable for foreign keys
	ForeignKey *constraint.ForeignKey `yaml:"foreign-key"`
}

func (c *Column) Validate() error {
	if c.PrimaryKey != nil && *c.PrimaryKey && c.Null != nil && *c.Null {
		return ErrPrimaryKeyCannotBeNull
	}

	if c.Type == TypeUnknown {
		var exists bool
		if c.Type, exists = aliasToColumnType[c.TypeAlias]; !exists {
			return errors.Join(ErrInvalidColumnType, errors.New(c.TypeAlias))
		}
	}

	if c.Type == TypeString && c.MaxLength <= 0 {
		return ErrInvalidMaxLength
	}

	if c.ForeignKey != nil {
		if c.ForeignKey.Columns == nil && c.ForeignKey.Column != "" {
			c.ForeignKey.Columns = append(c.ForeignKey.Columns, c.ForeignKey.Column)
		}
	}

	return nil
}
