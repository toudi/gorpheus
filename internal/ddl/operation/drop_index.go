package operation

type DropIndex struct {
	Table     string `yaml:"table"`
	IndexName string `yaml:"index"`
}
