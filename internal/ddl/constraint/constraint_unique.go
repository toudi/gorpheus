package constraint

type Unique struct {
	Columns  []string `yaml:"columns"`
	Null     *bool    `yaml:"nulls"`
	Distinct *bool    `yaml:"distinct"`
}
