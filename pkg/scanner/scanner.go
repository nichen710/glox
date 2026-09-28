package scanner

import (
	"strconv"

	"glox/pkg/token"
)

type Scanner struct {
	source  string
	tokens  []token.Token
	start   int
	current int
	line    int
	errors  ScanErrors
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  make([]token.Token, 0),
		start:   0,
		current: 0,
		line:    1,
		errors:  make(ScanErrors, 0),
	}
}

// Scanner Public Methods
func (s *Scanner) Scan() ([]token.Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	if len(s.errors) > 0 {
		return nil, s.errors
	}

	s.tokens = append(s.tokens, token.NewToken(token.EOF, "", nil, s.line))
	return s.tokens, nil
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	// Skip whitespace
	case ' ', '\r', '\t':
		// Ignorar
	case '\n':
		s.line++

	// Single-character tokens
	case '(':
		s.addToken(token.LEFT_PAREN)
	case ')':
		s.addToken(token.RIGHT_PAREN)
	case '{':
		s.addToken(token.LEFT_BRACE)
	case '}':
		s.addToken(token.RIGHT_BRACE)
	case ',':
		s.addToken(token.COMMA)
	case '-':
		s.addToken(token.MINUS)
	case '+':
		s.addToken(token.PLUS)
	case ';':
		s.addToken(token.SEMICOLON)
	case '*':
		s.addToken(token.STAR)
	case '%':
		s.addToken(token.PERCENT)

	// Comments or division
	case '/':
		if s.match('/') {
			// Comment: consume until the end of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH)
		}

	// One- or two-character tokens
	case '!':
		s.matchToken('=', token.BANG_EQUAL, token.BANG)
	case '=':
		s.matchToken('=', token.EQUAL_EQUAL, token.EQUAL)
	case '<':
		s.matchToken('=', token.LESS_EQUAL, token.LESS)
	case '>':
		s.matchToken('=', token.GREATER_EQUAL, token.GREATER)

	// Strings
	case '"', '\'':
		s.stringLiteral(c)

	default:
		if isDigit(c) {
			s.numberLiteral()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			s.errors = append(s.errors, NewUnexpectedCharError(s.line, c))
		}
	}
}

func (s *Scanner) stringLiteral(quote byte) {
	for !s.isAtEnd() && s.peek() != quote {
		if s.peek() == '\n' {
			if quote == '\'' {
				s.errors = append(s.errors, NewUnterminatedStringError(s.line, s.source[s.start:s.current]))
				return
			}
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.errors = append(s.errors, NewUnterminatedStringError(s.line, s.source[s.start:s.current]))
		return
	}

	// Consume the closing quote
	s.advance()

	// The value of the string without the quotes
	value := s.source[s.start+1 : s.current-1]
	s.addToken(token.STRING, value)
}

func (s *Scanner) numberLiteral() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Handle fractional part
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the '.'
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	text := s.source[s.start:s.current]
	val, err := strconv.ParseFloat(text, 64)
	if err != nil {
		s.errors = append(s.errors, NewInvalidNumberError(s.line, text))
		return
	}

	s.addToken(token.NUMBER, val)
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, isKeyword := token.Keywords[text]
	if isKeyword {
		s.addToken(tokenType)
	} else {
		s.addToken(token.IDENTIFIER)
	}
}

func (s *Scanner) addToken(tokenType token.TokenType, literal ...any) {
	var lit any
	if len(literal) > 0 {
		lit = literal[0]
	}

	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.NewToken(tokenType, text, lit, s.line))
}

func (s *Scanner) matchToken(expected byte, matched, unmatched token.TokenType) {
	if s.match(expected) {
		s.addToken(matched)
	} else {
		s.addToken(unmatched)
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() || s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) advance() byte {
	c := s.source[s.current]
	s.current++
	return c
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

// Char Value Checkers
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}
