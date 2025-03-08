package column

const (
	TypeUnknown = iota
	TypeBoolean
	TypeDate
	TypeTime
	TypeTimestamp
	TypeSmallInt
	TypeInt
	TypeBigInt
	TypeInt128
	TypeDecimal
	TypeFloat
	TypeDouble
	TypeString
	TypeText
	TypeBlob
	// postgreSQL-specific types
	TypeSmallSerial
	TypeSerial
	TypeBigSerial
)

var columnTypeAliases = map[int][]string{
	TypeBoolean:     {"bool", "boolean"},
	TypeDate:        {"date"},
	TypeTime:        {"time"},
	TypeTimestamp:   {"datetime", "timestamp"},
	TypeSmallInt:    {"smallint"},
	TypeInt:         {"int", "integer"},
	TypeBigInt:      {"bigint"},
	TypeInt128:      {"int128"},
	TypeSmallSerial: {"smallserial"},
	TypeSerial:      {"serial"},
	TypeBigSerial:   {"bigserial"},
	TypeDecimal:     {"decimal", "numeric"},
	TypeFloat:       {"float", "float32", "real"},
	TypeDouble:      {"double", "float64"},
	TypeString:      {"string", "varchar"},
	TypeText:        {"text"},
}

var aliasToColumnType = map[string]int{}

func init() {
	// generate reverse mapping, i.e. the one from
	// string alias to the corresponding type so that we could
	// quickly deduce the internal type
	for columnType, aliases := range columnTypeAliases {
		for _, alias := range aliases {
			aliasToColumnType[alias] = columnType
		}
	}
}
