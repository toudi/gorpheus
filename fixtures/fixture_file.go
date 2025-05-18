package fixtures

type FixtureFile struct {
	Table string           `yaml:"table" json:"table"`
	Data  []map[string]any `yaml:"data" json:"data"`
}
