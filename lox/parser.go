package lox

import "fmt"

// recursive descent descends through the grammar with each token

// expression 	-> equality ;
// equality 		-> comparison ( ( "!=" | "==" ) comparison )* ;
// comparison 	-> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term					-> factor ( ( "-" | "+" ) factor )* ;
// factor				-> unary ( ( "/" | "*" ) unary )* ;
// unary				-> ( "!" | "-" ) unary | primary ;
// primary			-> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;

// Parse takes a slice of Token and creates an Abstract Syntax Tree of Expr using the Recursive Descent method
func Parse(tokens []Token) (Expr, error) {
	var expression, equality, comparison, term, factor, unary, primary func() (Expr, error)
	current := 0

	match := func(types ...TokenType) bool {
		if current >= len(tokens) {
			return false
		}
		for _, t := range types {
			if tokens[current].Type == t {
				current++
				return true
			}
		}
		return false
	}

	previous := func() Token {
		return tokens[current-1]
	}

	// expression -> equality ;
	expression = func() (Expr, error) {
		return equality()
	}

	// equality -> comparison ( ( "!=" | "==" ) comparison )* ;
	equality = func() (Expr, error) {
		expr, err := comparison()
		for match(BangEqual, EqualEqual) {
			operator := previous()
			right, _ := comparison()
			expr = BinaryExpr{
				Left:     expr,
				Operator: operator,
				Right:    right,
			}
		}
		return expr, err
	}

	// comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
	comparison = func() (Expr, error) {
		expr, err := term()
		for match(Greater, GreaterEqual, Less, LessEqual) {
			operator := previous()
			right, _ := term()
			expr = BinaryExpr{
				Left:     expr,
				Operator: operator,
				Right:    right,
			}
		}
		return expr, err
	}

	// term	-> factor ( ( "-" | "+" ) factor )* ;
	term = func() (Expr, error) {
		expr, err := factor()
		for match(Minus, Plus) {
			operator := previous()
			right, _ := factor()
			expr = BinaryExpr{
				Left:     expr,
				Operator: operator,
				Right:    right,
			}
		}
		return expr, err
	}

	// factor	-> unary ( ( "/" | "*" ) unary )* ;
	factor = func() (Expr, error) {
		expr, err := unary()
		for match(Slash, Star) {
			operator := previous()
			right, _ := unary()
			expr = BinaryExpr{
				Left:     expr,
				Operator: operator,
				Right:    right,
			}
		}
		return expr, err
	}

	// unary -> ( "!" | "-" ) unary | primary ;
	unary = func() (Expr, error) {
		if match(Bang, Minus) {
			operator := previous()
			right, err := unary()
			return UnaryExpr{
				Operator: operator,
				Right:    right,
			}, err
		}
		return primary()
	}

	// primary -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
	primary = func() (Expr, error) {
		if match(Number, String, True, False, Nil) {
			return LiteralExpr{Value: previous()}, nil
		}
		if match(LeftParen) {
			expr, _ := expression()
			if match(RightParen) {
				current++
			}
			return GroupingExpr{Expr: expr}, nil
		}
		return nil, fmt.Errorf("Parsing error on line %d: Unexpected token \"%s\"", tokens[current].Line, tokens[current].Lexeme)
	}

	var root Expr
	var err error
	for current < len(tokens) && tokens[current].Type != EOF {
		root, err = expression()
	}

	return root, err
}
