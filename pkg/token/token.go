package token

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func NewToken(tokenType TokenType, lexeme string, literal any, line int) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t Token) String() string {
	if t.Type == IDENTIFIER {
		return fmt.Sprintf("%s<%s>", t.Type, t.Lexeme)
	}

	if t.Literal == nil {
		return t.Type.String()
	}

	return fmt.Sprintf("%s<%v>", t.Type, t.Literal)
}
