package internal

import "github.com/pipe01/kensaku/query"

type TokenStream <-chan Token

func (ts TokenStream) Take(typ TokenType) (Token, bool) {
	tk := <-ts

	if tk.Type != typ {
		return Token{}, false
	}

	return tk, true
}

func (ts TokenStream) TakeEither(typa, typb TokenType) (Token, bool) {
	tk := <-ts

	if tk.Type != typa && tk.Type != typb {
		return Token{}, false
	}

	return tk, true
}

func ParseOperators(tokens TokenStream) ([]query.Operator, bool) {
	ops := make([]query.Operator, 0)

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
			ops = append(ops, &textOperator{text: tk.Content, exact: false})
		}
	}

	return ops, true
}

func takeOperator(tokench TokenStream) (query.Operator, bool) {
	field, ok := tokench.Take(TokenText)
	if !ok {
		return nil, false
	}

	if _, ok := tokench.Take(TokenColon); !ok {
		return nil, false
	}

	value, ok := tokench.TakeEither(TokenText, TokenQuoted)
	if !ok {
		return nil, false
	}
	exact := value.Type == TokenQuoted

	if _, ok := tokench.Take(TokenCloseParen); !ok {
		return nil, false
	}

	return &textOperator{field: field.Content, text: value.Content, exact: exact}, true
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
