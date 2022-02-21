package kensaku

type Query []Operator

type Operator interface {
	Field() string
}

// func Parse(str string) {
// 	input := antlr.NewInputStream(str)
// 	lexer := parser.NewQueryLexer(input)
// 	stream := antlr.NewCommonTokenStream(lexer, 0)

// 	p := parser.NewQueryParser(stream)
// 	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
// 	p.BuildParseTrees = true

// 	q := p.Query()
// 	_ = q
// }

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
