package lox

// Parse converts a slice of Token to a slice of statements (?), returning the Abstract Syntax Tree
func Parse(tokens []Token) []Stmt {
	current := 0
	stmts := make([]Stmt, 0)

	check := func(expected TokenType) bool {
		if current >= len(tokens) {
			return false
		}
		return tokens[current].Type == expected
	}

	// peek := func() Token {
	// 	return tokens[current]
	// }

	previous := func() Token {
		return tokens[current-1]
	}

	advance := func() Token {
		if current < len(tokens) {
			current++
		}
		return previous()
	}

	match := func(types ...TokenType) bool {
		for _, t := range types {
			if check(t) {
				advance()
				return true
			}
		}
		return false
	}

	consume := func(t TokenType) Token {
		// if !check(t) {
		// 	fmt.Println("Unexpected Token?")
		// 	return nil
		// }
		return advance()
	}

	// recursive descent descends through the grammar with each token
	var expression func() Expr // expression 	-> equality ;
	var equality func() Expr   // equality 		-> comparison ( ( "!=" | "==" ) comparison )* ;
	var comparison func() Expr // comparison 	-> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
	var term func() Expr       // term					-> factor ( ( "-" | "+" ) factor )* ;
	var factor func() Expr     // factor				-> unary ( ( "/" | "*" ) unary )* ;
	var unary func() Expr      // unary				-> ( "!" | "-" ) unary | primary ;
	var primary func() Expr    // primary			-> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;

	expression = func() Expr {
		return equality()
	}

	equality = func() Expr {
		expr := comparison()
		for match(BangEqual, EqualEqual) {
			expr = &BinaryExpr{
				Left:     expr,
				Operator: previous(),
				Right:    comparison(),
			}
		}
		return expr
	}

	comparison = func() Expr {
		expr := term()
		for match(Greater, GreaterEqual, Less, LessEqual) {
			expr = &BinaryExpr{
				Left:     expr,
				Operator: previous(),
				Right:    term(),
			}
		}
		return expr
	}

	term = func() Expr {
		expr := factor()
		for match(Minus, Plus) {
			expr = &BinaryExpr{
				Left:     expr,
				Operator: previous(),
				Right:    factor(),
			}
		}
		return expr
	}

	factor = func() Expr {
		expr := unary()
		for match(Slash, Star) {
			expr = &BinaryExpr{
				Left:     expr,
				Operator: previous(),
				Right:    unary(),
			}
		}
		return expr
	}

	unary = func() Expr {
		if match(Minus, Bang) {
			return &UnaryExpr{
				Operator: previous(),
				Right:    unary(),
			}
		}
		return primary()
	}

	primary = func() Expr {
		if match(False) {
			return &LiteralExpr{Value: False}
		} else if match(True) {
			return &LiteralExpr{Value: True}
		} else if match(Nil) {
			return &LiteralExpr{Value: Nil}
		} else if match(Number, String) {
			return &LiteralExpr{Value: previous().Literal}
		} else if match(LeftParen) {
			expr := expression() // loop back to the top of the grammar if a "(" is found
			consume(RightParen)
			return &GroupingExpr{Expr: expr}
		}
		return nil
	}

	// main parsing loop
	for current < len(tokens) {
		// stmts = append(stmts, expression())
		break
	}

	return stmts
}
