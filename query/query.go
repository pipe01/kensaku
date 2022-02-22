package query

type Query []Operator

type Operator interface {
	Field() string
}

type baseOperator struct {
	field string
}

func (b baseOperator) Field() string {
	return b.field
}

type TextOperator struct {
	baseOperator
	text  string
	exact bool
}

func (t *TextOperator) Text() string {
	return t.text
}

func (t *TextOperator) Exact() bool {
	return t.exact
}
