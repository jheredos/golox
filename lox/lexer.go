package lox

import "fmt"

var keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"fun":    Fun,
	"for":    For,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

// Lex is the wrapper function for the tail-recursive lex()
func Lex(source string) ([]Token, error) {
	tokens := make([]Token, 0)
	return lex(tokens, source, 1, nil)
}

func newToken(ttype TokenType, value string, line int) Token {
	return Token{
		Type:   ttype,
		Lexeme: value,
		Line:   line,
	}
}

// skipComment recurses through a string until finding a newline and returns the rest of the input string
func skipComment(tail string) string {
	if len(tail) <= 0 {
		return ""
	} else if tail[0] == '\n' {
		return tail[1:]
	} else {
		return skipComment(tail[1:])
	}
}

// takes a string and recurses through it until finding a closing '"' rune
// returns the tail, current string, and number of lines
func findString(tail string, current string, lines int) (string, string, int) {
	if len(tail) <= 0 {
		return "", current, lines
	} else if tail[0] == '"' {
		return tail[1:], current, lines
	} else if tail[0] == '\n' {
		return findString(tail[1:], current+string(tail[0]), lines+1)
	} else {
		return findString(tail[1:], current+string(tail[0]), lines)
	}
}

// takes a string and recurses through it until finding a non-numeric rune or a second '.'
// returns the rest of the input string, the current string representing the number, and a bool denoting whether a decimal point has been seen
func findNumber(tail string, current string, dotSeen bool) (string, string, bool) {
	if len(tail) <= 0 {
		return "", current, dotSeen
	} else if !isDigit(tail[0]) && tail[0] != '.' {
		return tail, current, dotSeen
	} else if tail[0] == '.' && dotSeen {
		fmt.Printf("Warning: malformed number literal \"%s\"", current+".")
		return tail[1:], current, dotSeen
	} else if tail[0] == '.' && !dotSeen {
		return findNumber(tail[1:], current+string(tail[0]), true)
	} else {
		return findNumber(tail[1:], current+string(tail[0]), dotSeen)
	}
}

func findIdentifier(tail string, current string) (string, string) {
	if len(tail) <= 0 {
		return "", current
	} else if !isAlphaNumeric(tail[0]) {
		return tail, current
	} else {
		return findIdentifier(tail[1:], current+string(tail[0]))
	}
}

func isAlpha(r byte) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isDigit(r byte) bool {
	return r >= '0' && r <= '9'
}

func isAlphaNumeric(r byte) bool {
	return isAlpha(r) || isDigit(r)
}

// lex is the tail-recursive helper function for Lex()
// it is the main lexing switch, recursing through the string and matching tokens
// that it appends to the current slice of Token, along with tracking line number
func lex(current []Token, tail string, line int, err error) ([]Token, error) {
	if err != nil {
		return current, err
	}
	if len(tail) == 0 {
		return append(current, newToken(EOF, "\x00", line)), nil
	}
	r := tail[0]
	switch r {
	// whitespace
	case '\n':
		return lex(current, tail[1:], line+1, nil)
	case '\t':
		return lex(current, tail[1:], line, nil)
	case '\r':
		return lex(current, tail[1:], line, nil)
	case ' ':
		return lex(current, tail[1:], line, nil)

	// single-character tokens
	case '(':
		return lex(
			append(current, newToken(LeftParen, string(r), line)),
			tail[1:],
			line,
			nil,
		)
	case ')':
		return lex(
			append(current, newToken(RightParen, string(r), line)),
			tail[1:],
			line,
			nil,
		)
	case '{':
		return lex(
			append(current, newToken(LeftBrace, string(r), line)),
			tail[1:],
			line,
			nil,
		)
	case '}':
		return lex(
			append(current, newToken(RightBrace, string(r), line)),
			tail[1:],
			line,
			nil,
		)
	case ',':
		return lex(
			append(current, newToken(Comma, string(r), line)),
			tail[1:],
			line,
			nil,
		)
	case '.':
		return lex(
			append(current, newToken(Dot, string(r), line)),
			tail[1:],
			line,
			nil,
		)
	case '-':
		return lex(
			append(current, newToken(Minus, string(r), line)),
			tail[1:],
			line,
			nil,
		)
	case '+':
		return lex(
			append(current, newToken(Plus, string(r), line)),
			tail[1:],
			line,
			nil,
		)
	case ';':
		return lex(
			append(current, newToken(Semicolon, string(r), line)),
			tail[1:],
			line,
			nil,
		)
	case '*':
		return lex(
			append(current, newToken(Star, string(r), line)),
			tail[1:],
			line,
			nil,
		)

	// 1-2 characters
	case '!':
		{
			if tail[1] == '=' {
				return lex(
					append(current, newToken(BangEqual, "!=", line)),
					tail[2:],
					line,
					nil,
				)
			}
			return lex(
				append(current, newToken(Bang, string(r), line)),
				tail[1:],
				line,
				nil,
			)
		}
	case '=':
		{
			if tail[1] == '=' {
				return lex(
					append(current, newToken(EqualEqual, "==", line)),
					tail[2:],
					line,
					nil,
				)
			}
			return lex(
				append(current, newToken(Equal, string(r), line)),
				tail[1:],
				line,
				nil,
			)
		}
	case '<':
		{
			if tail[1] == '=' {
				return lex(
					append(current, newToken(LessEqual, "<=", line)),
					tail[2:],
					line,
					nil,
				)
			}
			return lex(
				append(current, newToken(Less, string(r), line)),
				tail[1:],
				line,
				nil,
			)
		}
	case '>':
		{
			if tail[1] == '=' {
				return lex(
					append(current, newToken(GreaterEqual, ">=", line)),
					tail[2:],
					line,
					nil,
				)
			}
			return lex(
				append(current, newToken(Greater, string(r), line)),
				tail[1:],
				line,
				nil,
			)
		}

	// slash - either Slash or Comment
	case '/':
		{
			if tail[1] == '/' {
				return lex(
					current,
					skipComment(tail[2:]),
					line+1,
					nil,
				)
			}
			return lex(
				append(current, newToken(Slash, string(r), line)),
				tail[1:],
				line,
				nil,
			)
		}

	// strings
	case '"':
		newTail, val, lines := findString(tail[1:], "", 0)
		return lex(
			append(current, newToken(String, val, line)),
			newTail,
			line+lines,
			nil,
		)

	default:
		{
			// numbers
			if isDigit(r) {
				newTail, val, _ := findNumber(tail[1:], string(tail[0]), false)
				return lex(
					append(current, newToken(Number, val, line)),
					newTail,
					line,
					nil,
				)
				// identifiers
			} else if isAlpha(r) {
				newTail, val := findIdentifier(tail[1:], string(tail[0]))
				ttype, isKeyword := keywords[val]
				if isKeyword {
					return lex(
						append(current, newToken(ttype, val, line)),
						newTail,
						line,
						nil,
					)
				}
				return lex(
					append(current, newToken(Identifier, val, line)),
					newTail,
					line,
					nil,
				)
			} else {
				err = fmt.Errorf("Lexing error at line %d: unexpected character \"%s\"", line, string(r))
				return current, err
			}
		}
	}
}
