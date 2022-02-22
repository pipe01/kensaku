package query

type Query []Operator

type Operator interface {
	Field() string
}

type TextOperator interface {
	Operator
	Text() string
	Exact() bool
}

type NumberComparison byte

const (
	CompareEquals NumberComparison = iota
	CompareLessThan
	CompareLessOrEqual
	CompareGreaterThan
	CompareGreaterOrEqual
)

type NumberOperator interface {
	Operator
	Value() float64
	Comparison() NumberComparison
}
