module gorpheus-cli

go 1.23.3

require (
	github.com/alecthomas/kong v1.8.1
	github.com/lmittmann/tint v1.0.7
	github.com/toudi/gorpheus/v2 v2.0.0-00010101000000-000000000000
	github.com/toudi/gorpheus/v2/dialect/firebird v0.0.0-00010101000000-000000000000
	github.com/toudi/gorpheus/v2/dialect/postgresql v0.0.0-00010101000000-000000000000
	github.com/toudi/gorpheus/v2/dialect/sqlite3 v0.0.0-00010101000000-000000000000
	github.com/xo/dburl v0.23.4
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.2 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/mattn/go-sqlite3 v1.14.24 // indirect
	github.com/nakagami/chacha20 v0.1.0 // indirect
	github.com/nakagami/firebirdsql v0.9.14 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/samber/lo v1.49.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	gitlab.com/nyarla/go-crypt v0.0.0-20160106005555-d9a5dc2b789b // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/exp v0.0.0-20250228200357-dead58393ab7 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/mathutil v1.7.1 // indirect
)

replace github.com/toudi/gorpheus/v2/dialect/sqlite3 => ../../dialect/sqlite3

replace github.com/toudi/gorpheus/v2/dialect/postgresql => ../../dialect/postgresql

replace github.com/toudi/gorpheus/v2/dialect/firebird => ../../dialect/firebird

replace github.com/toudi/gorpheus/v2 => ../../
