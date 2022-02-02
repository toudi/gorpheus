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

func (g *GoMigration) GetType() uint8 {
	return TypeGo
}

func (m *GoMigration) Close() error {
	return nil
}
