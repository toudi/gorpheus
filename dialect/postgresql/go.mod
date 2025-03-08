module github.com/toudi/gorpheus/v2/dialect/postgresql

go 1.23.3

replace github.com/toudi/gorpheus/v2 => ../../

require (
	github.com/jackc/pgx/v5 v5.7.2
	github.com/samber/lo v1.49.1
	github.com/toudi/gorpheus/v2 v2.0.0-00010101000000-000000000000
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
