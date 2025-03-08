package dialect

import "github.com/toudi/gorpheus/v2/interfaces"

var registry map[string]interfaces.DialectConstructor

func Register(name string, constructor interfaces.DialectConstructor) {
	registry[name] = constructor
}

func GetConstructor(name string) (interfaces.DialectConstructor, bool) {
	constructor, exists := registry[name]
	return constructor, exists
}

func init() {
	registry = make(map[string]interfaces.DialectConstructor)
}
