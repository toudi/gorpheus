package operation

type DropField struct {
	TableName string `yaml:"table"`
	Field     string
}
