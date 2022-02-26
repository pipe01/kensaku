package internal

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pipe01/kensaku/query"
)

func TestTakeOperator(t *testing.T) {
	data := []struct {
		field string
		tks   []Token
		op    query.Operator
	}{
		{"field", []Token{{TokenText, "field"}, {TokenColon, ":"}, {TokenText, "value"}}, &textOperator{field: "field", text: "value", exact: false}},
		{"field", []Token{{TokenText, "field"}, {TokenColon, ":"}, {TokenQuoted, "value"}}, &textOperator{field: "field", text: "value", exact: true}},
		{"field", []Token{{TokenText, "field"}, {TokenEquals, "="}, {TokenText, "123.45"}}, &numberOperator{field: "field", value: 123.45, comp: query.CompareEquals}},
	}

	for i, d := range data {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tks := tokenChannel(d.tks...)

			op, ok := takeOperator(TokenStream(tks))
			if !ok {
				t.Fatal("parsing failed")
			}

			if op.Field() != d.op.Field() {
				t.Fatalf(`expected field "%s", got "%s"`, d.op.Field(), op.Field())
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
		ops  []query.Operator
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
			[]query.Operator{&numberOperator{field: "field", value: 123.45, comp: query.CompareEquals}},
		},
		{
			"mixed quotes",
			[]Token{
				{TokenText, "non quoted"},
				{TokenQuoted, "quoted text"},
				{TokenText, "again not"},
			},
			[]query.Operator{
				&textOperator{field: "", text: "non quoted again not"},
				&textOperator{field: "", text: "quoted text", exact: true},
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

			if len(ops) != len(d.ops) {
				t.Fatalf("expected %d operators, got %d", len(d.ops), len(ops))
			}

			for i, op := range ops {
				if op.Field() != d.ops[i].Field() {
					t.Fatalf(`expected field "%s", got "%s"`, d.ops[i].Field(), op.Field())
				}

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
