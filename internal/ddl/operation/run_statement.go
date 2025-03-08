package operation

type RunStatement struct {
	DDL          string
	ReverseDDL   string
	ForwardsOnly bool
}
