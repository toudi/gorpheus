package utils

func MigrationID(namespace string, name string) string {
	if name != "" {
		return namespace + "/" + name
	}
	return namespace
}
