package internal

import "strings"

type Token struct {
	Type    TokenType
	Content string
}

type TokenType byte

const (
	TokenNone TokenType = iota
	TokenText
	TokenQuoted
	TokenOpenParen
	TokenCloseParen
	TokenColon
)

type Lexer struct {
	str string
	pos int
	ch  chan<- Token

	textb strings.Builder
}

func NewLexer(str string, ch chan<- Token) *Lexer {
	return &Lexer{
		str: str,
		ch:  ch,
	}
}

func (l *Lexer) Lex() {
	for l.pos < len(l.str) {
		c := l.str[l.pos]

		switch c {
		case '(':
			l.putToken(Token{TokenOpenParen, "("})
		case ')':
			l.putToken(Token{TokenCloseParen, ")"})
		case ':':
			l.putToken(Token{TokenColon, ":"})
		case '"':
			l.pos++
			l.putToken(l.takeQuoted())
		default:
			l.textb.WriteByte(c)
		}

		l.pos++
	}

	l.putString()
	close(l.ch)
}

func (l *Lexer) putToken(tk Token) {
	l.putString()

	tk.Content = strings.TrimSpace(tk.Content)
	if len(tk.Content) > 0 {
		l.ch <- tk
	}
}

func (l *Lexer) putString() {
	str := strings.TrimSpace(l.textb.String())
	if len(str) > 0 {
		l.ch <- Token{TokenText, str}
		l.textb.Reset()
	}
}

func (l *Lexer) takeQuoted() Token {
	start := l.pos
	end := 0

	for {
		c := l.str[l.pos]

		if c == '"' || l.pos == len(l.str)-1 {
			end = l.pos
			break
		}

		l.pos++
	}

	return Token{TokenQuoted, l.str[start:end]}
}
