package authzgen

import "fmt"

// TokenType represents the type of a token
type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	IDENTIFIER
	DEFINITION
	RELATION
	PERMISSION
	EQUAL
	PLUS
	MINUS
	PIPE
	WILDCARD
	LBRACE
	RBRACE
	COLON
	SLASH
	HASH
	MINUS_ARROW
)

var keywords = map[string]TokenType{
	"definition": DEFINITION,
	"relation":   RELATION,
	"permission": PERMISSION,
}

var tokenNames = map[TokenType]string{
	ILLEGAL:     "ILLEGAL",
	EOF:         "EOF",
	IDENTIFIER:  "IDENTIFIER",
	DEFINITION:  "DEFINITION",
	RELATION:    "RELATION",
	PERMISSION:  "PERMISSION",
	EQUAL:       "EQUAL",
	PLUS:        "PLUS",
	MINUS:       "MINUS",
	PIPE:        "PIPE",
	WILDCARD:    "WILDCARD",
	LBRACE:      "LBRACE",
	RBRACE:      "RBRACE",
	COLON:       "COLON",
	SLASH:       "SLASH",
	HASH:        "HASH",
	MINUS_ARROW: "MINUS_ARROW",
}

func (t TokenType) String() string {
	if name, ok := tokenNames[t]; ok {
		return name
	}
	return fmt.Sprintf("UNKNOWN(%d)", int(t))
}

// Token represents a lexical token
type Token struct {
	Type     TokenType
	Literal  string
	Line     int
	Position int
}

// Lexer performs lexical analysis on input
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	linePosition int
}

// NewLexer creates a new lexer for the given input
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
		line:  1,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	if l.ch == '\n' {
		l.line++
		l.linePosition = 0
	} else {
		l.linePosition++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	if l.ch == '/' && l.peekChar() == '/' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	for {
		if l.ch == '/' && l.peekChar() == '/' {
			l.skipComment()
			l.skipWhitespace()
			continue
		}
		break
	}

	switch l.ch {
	case '=':
		tok = Token{Type: EQUAL, Literal: string(l.ch), Line: l.line, Position: l.linePosition}
	case '+':
		tok = Token{Type: PLUS, Literal: string(l.ch), Line: l.line, Position: l.linePosition}
	case '|':
		tok = Token{Type: PIPE, Literal: string(l.ch), Line: l.line, Position: l.linePosition}
	case '*':
		tok = Token{Type: WILDCARD, Literal: string(l.ch), Line: l.line, Position: l.linePosition}
	case '{':
		tok = Token{Type: LBRACE, Literal: string(l.ch), Line: l.line, Position: l.linePosition}
	case '}':
		tok = Token{Type: RBRACE, Literal: string(l.ch), Line: l.line, Position: l.linePosition}
	case ':':
		tok = Token{Type: COLON, Literal: string(l.ch), Line: l.line, Position: l.linePosition}
	case '/':
		if l.peekChar() == '/' {
			l.skipComment()
			l.skipWhitespace()
			return l.NextToken()
		}
		tok = Token{Type: SLASH, Literal: string(l.ch), Line: l.line, Position: l.linePosition}
	case '#':
		tok = Token{Type: HASH, Literal: string(l.ch), Line: l.line, Position: l.linePosition}
	case '-':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: MINUS_ARROW, Literal: literal, Line: l.line, Position: l.linePosition - 1}
		} else {
			tok = Token{Type: MINUS, Literal: string(l.ch), Line: l.line, Position: l.linePosition}
		}
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		tok.Line = l.line
		tok.Position = l.linePosition
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			tok.Line = l.line
			tok.Position = l.linePosition - len(tok.Literal)
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Line: l.line, Position: l.linePosition}
		}
	}

	l.readChar()
	return tok
}

func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}

// TokenizeAll returns all tokens from the input
func (l *Lexer) TokenizeAll() []Token {
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}
	return tokens
}
