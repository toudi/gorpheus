package operation

type CreateView struct {
	Name         string `yaml:"name"`
	Definition   string `yaml:"definition"`
	Materialized bool   `yaml:"materialized"`
}
