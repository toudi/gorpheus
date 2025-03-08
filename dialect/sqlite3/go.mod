module github.com/toudi/gorpheus/v2/dialect/sqlite3

go 1.23.3

replace github.com/toudi/gorpheus/v2 => ../../

require (
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/samber/lo v1.49.1
	github.com/stretchr/testify v1.10.0
	github.com/toudi/gorpheus/v2 v2.0.0-00010101000000-000000000000
	github.com/xo/dburl v0.23.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
