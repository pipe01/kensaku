package internal

import "testing"

func TestLexer(t *testing.T) {
	tests := map[string][]Token{
		"single":             {{TokenText, "single"}},
		"two words":          {{TokenText, "two words"}},
		"word = word":        {{TokenText, "word"}, {TokenEquals, "="}, {TokenText, "word"}},
		"word > 123":         {{TokenText, "word"}, {TokenGreater, ">"}, {TokenText, "123"}},
		"word >= 123":        {{TokenText, "word"}, {TokenGreaterEquals, ">="}, {TokenText, "123"}},
		"word < 123":         {{TokenText, "word"}, {TokenLess, "<"}, {TokenText, "123"}},
		"word <= 123":        {{TokenText, "word"}, {TokenLessEquals, "<="}, {TokenText, "123"}},
		"word = 123":         {{TokenText, "word"}, {TokenEquals, "="}, {TokenText, "123"}},
		"word (field:value)": {{TokenText, "word"}, {TokenOpenParen, "("}, {TokenText, "field"}, {TokenColon, ":"}, {TokenText, "value"}, {TokenCloseParen, ")"}},
	}

	for str, exptks := range tests {
		t.Run(str, func(t *testing.T) {
			ch := make(chan Token, 30)

			lexer := NewLexer(str, ch)
			lexer.Lex()

			i := 0
			for tk := range ch {
				if i > len(exptks)-1 {
					t.Logf("expected %d tokens, got %d", len(exptks), i+1)
					t.FailNow()
				}

				if tk != exptks[i] {
					t.Logf("at %d: expected token %v, got %v", i, exptks[i], tk)
					t.FailNow()
				}

				i++
			}
		})
	}
}
