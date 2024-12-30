package globals

type ColumnRange struct {
	From int
	To   int
}

type Loc struct {
	Line  int
	Range ColumnRange
}

type Node interface {
	Reference() string
	GetLocs() []Loc
}
