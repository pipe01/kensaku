package internal

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/pipe01/kensaku/query"
)

func TestTakeOperator(t *testing.T) {
	data := []struct {
		tks []Token
		op  Operator
	}{
		{[]Token{{TokenText, "field"}, {TokenColon, ":"}, {TokenText, "value"}}, &TextOperator{Field: "field", Text: "value", Exact: false}},
		{[]Token{{TokenText, "field"}, {TokenColon, ":"}, {TokenQuoted, "value"}}, &TextOperator{Field: "field", Text: "value", Exact: true}},
		{[]Token{{TokenText, "field"}, {TokenEquals, "="}, {TokenText, "123.45"}}, &NumberOperator{Field: "field", Value: 123.45, Comparison: CompareEquals}},
		{[]Token{{TokenText, "field"}, {TokenGreaterEquals, ">="}, {TokenText, "123.45"}}, &NumberOperator{Field: "field", Value: 123.45, Comparison: CompareGreaterOrEqual}},
		{[]Token{{TokenText, "field"}, {TokenGreater, ">"}, {TokenText, "123.45"}}, &NumberOperator{Field: "field", Value: 123.45, Comparison: CompareGreaterThan}},
		{[]Token{{TokenText, "field"}, {TokenLessEquals, "<="}, {TokenText, "123.45"}}, &NumberOperator{Field: "field", Value: 123.45, Comparison: CompareLessOrEqual}},
		{[]Token{{TokenText, "field"}, {TokenLess, "<"}, {TokenText, "123.45"}}, &NumberOperator{Field: "field", Value: 123.45, Comparison: CompareLessThan}},
		{[]Token{{TokenText, "field"}, {TokenOpenParen, "("}, {TokenText, "123.45"}}, nil},
		{[]Token{{TokenText, ""}, {TokenColon, ":"}, {TokenText, "123.45"}}, nil},
	}

	for i, d := range data {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tks := tokenChannel(d.tks...)

			op, ok := takeOperator(TokenStream(tks))
			if !ok {
				if d.op == nil {
					return // We expected a failure
				}

				t.Fatal("parsing failed")
			}
			if d.op == nil {
				t.Fatal("expected a parsing failure")
			}

			if !reflect.DeepEqual(op, d.op) {
				t.Fatalf("expected operator %#v, got %#v", d.op, op)
			}
		})
	}
}

func TestParseOperators(t *testing.T) {
	data := []struct {
		name string
		tks  []Token
		ops  []Operator
	}{
		{
			"single operator",
			[]Token{
				{TokenOpenParen, "("},
				{TokenText, "field"},
				{TokenEquals, "="},
				{TokenText, "123.45"},
				{TokenCloseParen, ")"},
			},
			[]Operator{&NumberOperator{Field: "field", Value: 123.45, Comparison: CompareEquals}},
		},
		{
			"mixed quotes",
			[]Token{
				{TokenText, "non quoted"},
				{TokenQuoted, "quoted text"},
				{TokenText, "again not"},
			},
			[]Operator{
				&TextOperator{Field: "", Text: "non quoted again not"},
				&TextOperator{Field: "", Text: "quoted text", Exact: true},
			},
		},
		{
			"malformed operator, missing close",
			[]Token{
				{TokenOpenParen, "("},
				{TokenText, "field"},
				{TokenEquals, "="},
				{TokenText, "123.45"},
			},
			nil,
		},
		{
			"malformed operator, missing field",
			[]Token{
				{TokenOpenParen, "("},
				{TokenEquals, "="},
				{TokenText, "123.45"},
				{TokenCloseParen, ")"},
			},
			nil,
		},
		{
			"malformed operator, missing comparator",
			[]Token{
				{TokenOpenParen, "("},
				{TokenText, "field"},
				{TokenText, "123.45"},
				{TokenCloseParen, ")"},
			},
			nil,
		},
		{
			"malformed text operator, missing value",
			[]Token{
				{TokenOpenParen, "("},
				{TokenText, "field"},
				{TokenColon, ":"},
				{TokenCloseParen, ")"},
			},
			nil,
		},
		{
			"malformed number operator, missing value",
			[]Token{
				{TokenOpenParen, "("},
				{TokenText, "field"},
				{TokenEquals, "="},
				{TokenCloseParen, ")"},
			},
			nil,
		},
		{
			"malformed number operator, invalid number",
			[]Token{
				{TokenOpenParen, "("},
				{TokenText, "field"},
				{TokenEquals, "="},
				{TokenText, "asd"},
				{TokenCloseParen, ")"},
			},
			nil,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			tks := tokenChannel(d.tks...)

			ops, ok := ParseOperators(TokenStream(tks))
			if !ok {
				if d.ops == nil {
					return // We expected a failure
				}

				t.Fatal("parsing failed")
			}
			if d.ops == nil {
				t.Fatal("expected a parsing failure")
			}

			if len(ops) != len(d.ops) {
				t.Fatalf("expected %d operators, got %d", len(d.ops), len(ops))
			}

			for i, op := range ops {
				if !reflect.DeepEqual(op, d.ops[i]) {
					t.Fatalf("expected operator %#v, got %#v", d.ops[i], op)
				}
			}
		})
	}
}

func tokenChannel(tks ...Token) chan Token {
	ch := make(chan Token, len(tks))

	for _, tk := range tks {
		ch <- tk
	}

	close(ch)
	return ch
}
