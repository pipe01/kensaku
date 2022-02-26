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
	var textop *query.TextOperator

	for tk := range tokens {
		switch tk.Type {
		case TokenQuoted:
			ops = append(ops, &query.TextOperator{Text: tk.Content, Exact: true})

		case TokenOpenParen:
			op, ok := takeOperator(tokens)
			if !ok {
				return nil, false
			}

			if _, ok = tokens.Take(TokenCloseParen); !ok {
				return nil, false
			}

			ops = append(ops, op)

		default:
			if textop == nil {
				textop = &query.TextOperator{Text: tk.Content, Exact: false}
				ops = append(ops, textop)
			} else {
				textop.Text = strings.Join([]string{textop.Text, tk.Content}, " ")
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

	optk, ok := tokench.TakeAny(TokenColon, TokenEquals, TokenGreater, TokenGreaterEquals, TokenLess, TokenLessEquals)
	if !ok {
		return nil, false
	}

	var op query.Operator

	if optk.Type == TokenColon {
		op, ok = takeTextOperator(tokench, field.Content)
	} else {
		op, ok = takeNumberOperator(tokench, optk, field.Content)
	}

	if !ok {
		return nil, false
	}

	return op, true
}

func takeNumberOperator(tokench TokenStream, op Token, field string) (*query.NumberOperator, bool) {
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
	}

	return &query.NumberOperator{Field: field, Value: n, Comparison: comp}, true
}

func takeTextOperator(tokench TokenStream, field string) (*query.TextOperator, bool) {
	value, ok := tokench.TakeAny(TokenText, TokenQuoted)
	if !ok {
		return nil, false
	}
	exact := value.Type == TokenQuoted

	return &query.TextOperator{Field: field, Text: value.Content, Exact: exact}, true
}
