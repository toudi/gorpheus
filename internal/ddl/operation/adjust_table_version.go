package operation

type AdjustTableVersion struct {
	TableName string // which table version to adjust
	Delta     int    // by how much (either +1 or -1)
}
