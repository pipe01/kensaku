package query

type Query []Operator

type Operator interface {
	operator()
}

type NumberComparison byte

const (
	CompareEquals NumberComparison = iota
	CompareLessThan
	CompareLessOrEqual
	CompareGreaterThan
	CompareGreaterOrEqual
)

type TextOperator struct {
	Field string
	Text  string
	Exact bool
}

func (t *TextOperator) operator() {}

type NumberOperator struct {
	Field      string
	Value      float64
	Comparison NumberComparison
}

func (t *NumberOperator) operator() {}
