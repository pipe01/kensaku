package kensaku

import (
	"github.com/pipe01/kensaku/internal"
	"github.com/pipe01/kensaku/query"
)

type Query = query.Query
type Operator = query.Operator

func Parse(str string) (q query.Query, ok bool) {
	tks := make(chan internal.Token)
	defer close(tks)

	l := internal.NewLexer(str, tks)

	go l.Lex()

	ops, ok := internal.ParseOperators(internal.TokenStream(tks))
	if !ok {
		return nil, false
	}

	return query.Query(ops), true
}
