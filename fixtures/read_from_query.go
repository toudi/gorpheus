package fixtures

import "database/sql"

func (ff *FixtureFile) ReadDataFromQuery(rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// https://kylewbanks.com/blog/query-result-to-map-in-golang
	var columnValues []any = make([]any, len(columns))
	var columnPointers []any = make([]any, len(columns))

	for rows.Next() {
		var row map[string]any = make(map[string]any)

		for colIdx := range columns {
			columnValues[colIdx] = nil
			columnPointers[colIdx] = &columnValues[colIdx]
		}

		if err = rows.Scan(columnPointers...); err != nil {
			return err
		}

		for colIdx, colName := range columns {
			row[colName] = columnValues[colIdx]
		}
		ff.Data = append(ff.Data, row)
	}

	return nil
}
