package internal

import (
	"strconv"
	"strings"

	"github.com/pipe01/kensaku/query"
)

type TokenStream <-chan Token

func (ts TokenStream) Take(typ TokenType) (Token, bool) {
	tk := <-ts

	if tk.Type != typ {
		return Token{}, false
	}

	return tk, true
}

func (ts TokenStream) TakeAny(typ ...TokenType) (Token, bool) {
	tk := <-ts

	for _, t := range typ {
		if tk.Type == t {
			return tk, true
		}
	}

	return Token{}, false
}

func ParseOperators(tokens TokenStream) ([]query.Operator, bool) {
	ops := make([]query.Operator, 0)
	var textop *textOperator

	for tk := range tokens {
		switch tk.Type {
		case TokenQuoted:
			ops = append(ops, &textOperator{text: tk.Content, exact: true})

		case TokenOpenParen:
			op, ok := takeOperator(tokens)
			if !ok {
				return nil, false
			}
			ops = append(ops, op)

		default:
			if textop == nil {
				textop = &textOperator{text: tk.Content, exact: false}
				ops = append(ops, textop)
			} else {
				textop.text = strings.Join([]string{textop.text, tk.Content}, " ")
			}
		}
	}

	return ops, true
}

func takeOperator(tokench TokenStream) (query.Operator, bool) {
	field, ok := tokench.Take(TokenText)
	if !ok {
		return nil, false
	}

	op, ok := tokench.TakeAny(TokenColon, TokenEquals, TokenGreater, TokenGreaterEquals, TokenLess, TokenLessEquals)
	if !ok {
		return nil, false
	}

	if op.Type == TokenColon {
		return takeTextOperator(tokench, field.Content)
	}

	return takeNumberOperator(tokench, op, field.Content)
}

func takeNumberOperator(tokench TokenStream, op Token, field string) (query.Operator, bool) {
	valuetk, ok := tokench.Take(TokenText)
	if !ok {
		return nil, false
	}

	n, err := strconv.ParseFloat(valuetk.Content, 64)
	if err != nil {
		return nil, false
	}

	var comp query.NumberComparison

	switch op.Type {
	case TokenEquals:
		comp = query.CompareEquals
	case TokenGreater:
		comp = query.CompareGreaterThan
	case TokenGreaterEquals:
		comp = query.CompareGreaterOrEqual
	case TokenLess:
		comp = query.CompareLessThan
	case TokenLessEquals:
		comp = query.CompareLessOrEqual
	default:
		return nil, false
	}

	return &numberOperator{field: field, value: n, comp: comp}, true
}

func takeTextOperator(tokench TokenStream, field string) (query.Operator, bool) {
	value, ok := tokench.TakeAny(TokenText, TokenQuoted)
	if !ok {
		return nil, false
	}
	exact := value.Type == TokenQuoted

	if _, ok := tokench.Take(TokenCloseParen); !ok {
		return nil, false
	}

	return &textOperator{field: field, text: value.Content, exact: exact}, true
}

type textOperator struct {
	field string
	text  string
	exact bool
}

func (t *textOperator) Field() string {
	return t.field
}

func (t *textOperator) Text() string {
	return t.text
}

func (t *textOperator) Exact() bool {
	return t.exact
}

type numberOperator struct {
	field string
	value float64
	comp  query.NumberComparison
}

func (n *numberOperator) Field() string {
	return n.field
}

func (n *numberOperator) Value() float64 {
	return n.value
}

func (n *numberOperator) Comparison() query.NumberComparison {
	return n.comp
}
