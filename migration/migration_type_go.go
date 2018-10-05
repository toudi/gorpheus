package migration

type GoMigration struct {
	Migration
}

func (g *GoMigration) UpScript() (string, uint8, error) {
	return "", TypeGo, nil
}

func (g *GoMigration) DownScript() (string, uint8, error) {
	return "", TypeGo, nil
}
