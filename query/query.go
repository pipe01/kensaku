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
