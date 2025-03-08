package operation

type DropConstraint struct {
	Table      string `yaml:"table"`
	Constraint string `yaml:"constraint"`
}
