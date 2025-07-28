package operation

type DropView struct {
	Name         string `yaml:"name"`
	Materialized bool   `yaml:"materialized"`
}
