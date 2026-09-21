package scanner

import (
	"fmt"
	"strconv"

	"glox/pkg/token"
)

type Scanner struct {
	source  string
	tokens  []token.Token
	start   int
	current int
	line    int
	_error  error // Scanning Phase Detected Error
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  make([]token.Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

// Scanner Public Methods
func (s *Scanner) Scan() ([]token.Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()

		// Stop scanning if an error is detected
		if s._error != nil {
			return nil, s._error
		}
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
		if s.match('=') {
			s.addToken(token.BANG_EQUAL)
		} else {
			s.addToken(token.BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL)
		} else {
			s.addToken(token.EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LESS_EQUAL)
		} else {
			s.addToken(token.LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL)
		} else {
			s.addToken(token.GREATER)
		}

	// Strings
	case '"', '\'':
		s.stringLiteral(c)

	default:
		if isDigit(c) {
			s.numberLiteral()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			s.recordError(s.line, fmt.Sprintf("Unexpected character: `%c`", c))
		}
	}
}

func (s *Scanner) stringLiteral(quote byte) {
	for !s.isAtEnd() && s.peek() != quote {
		if s.peek() == '\n' {
			if quote == '\'' {
				s.recordError(s.line, fmt.Sprintf("Unterminated string: `%s`", s.source[s.start:s.current]))
				return
			}
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.recordError(s.line, fmt.Sprintf("Unterminated string: `%s`", s.source[s.start:s.current]))
		return
	}

	// Consume the closing quote
	s.advance()

	// The value of the string without the quotes
	value := s.source[s.start+1 : s.current-1]
	s.addTokenLiteral(token.STRING, value)
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
		s.recordError(s.line, fmt.Sprintf("Invalid number: `%s`", text))
		return
	}

	s.addTokenLiteral(token.NUMBER, val)
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

func (s *Scanner) addToken(tokenType token.TokenType) {
	s.addTokenLiteral(tokenType, nil)
}

func (s *Scanner) addTokenLiteral(tokenType token.TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.NewToken(tokenType, text, literal, s.line))
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

func (s *Scanner) recordError(line int, message string) {
	s._error = fmt.Errorf("[line %d] Error: %s\n", line, message)
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
