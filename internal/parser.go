package internal

import "github.com/pipe01/kensaku/query"

func Parse(tokench <-chan Token) {
	ops := make([]query.Operator, 0)

	for tk := range tokench {
		switch tk.Type {
		case TokenQuoted:
			ops = append(ops, &query.TextOperator{})
		}
	}
}
