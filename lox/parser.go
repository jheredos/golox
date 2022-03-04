package lox

import (
	"fmt"
)

// recursive descent descends through the grammar with each token

// program			-> declaration* EOF ;
// declaration	-> funDecl | varDecl | statement ;
// varDecl			-> "var" IDENTIFIER ( "=" expression )? ";" ;
// funDecl			-> "fun" function ;
// function			-> IDENTIFIER "(" parameters? ")" block ;
// parameters		-> IDENTIFIER ( "," IDENTIFIER )* ;
// statement		-> exprStmt | ifStmt | printStmt | forStmt | whileStmt | returnStmt | block ;
// block				-> "{" declaration* "}" ;
// returnStmt 	-> "return" expression? ";" ;
// forStmt			-> "for" "(" varDecl | exprStmt | ";" ) expression? ";" expression? ")" statement ;
// whileStmt		-> "while" "(" expression ")" statement ;
// ifStmt				-> "if" "(" expression ")" statement ( "else" statement )? ;
// exprStmt			-> expression ";" ;
// printStmt		-> "print" expression ";" ;

// expression 	-> equality ;
// assignment		-> IDENTIFIER "=" ( assignment | logicOr ) ;
// logicOr			-> logicAnd ( "or" logicAnd )* ;
// logicAnd		-> equality ( "and" equality)* ;
// equality 		-> comparison ( ( "!=" | "==" ) comparison )* ;
// comparison 	-> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term					-> factor ( ( "-" | "+" ) factor )* ;
// factor				-> unary ( ( "/" | "*" ) unary )* ;
// unary				-> ( "!" | "-" ) unary | call ;
// call					-> primary ( "(" arguments? ")" )* ; TODO
// primary			-> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER ;

// Parse takes a slice of Token and creates an Abstract Syntax Tree of Expr using the Recursive Descent method
func Parse(tokens []Token) (*Node, error) {
	var program, declaration, funDecl, varDecl, statement, function, parameters, block, returnStmt, forStmt, whileStmt, ifStmt, exprStmt, printStmt, expression, assignment, logicOr, logicAnd, equality, comparison, term, factor, unary, call, primary func() (*Node, error)
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

	// program -> declaration* EOF ;
	program = func() (*Node, error) {
		stmt, err := declaration()
		prgm := &Node{Type: ProgramNT, Right: stmt}
		for current < len(tokens) && !match(EOF) {
			decl, err := declaration()
			if err != nil {
				return prgm, err
			}
			if stmt != nil {
				stmt.Next = decl
			}
			stmt = decl

		}
		return prgm, err
	}

	// declaration -> varDecl | funDecl | statement ;
	declaration = func() (*Node, error) {
		if match(Var) {
			return varDecl()
		}
		if match(Fun) {
			return funDecl()
		}
		return statement()
	}

	// funDecl -> "fun" function ;
	funDecl = func() (*Node, error) {
		return function()
	}

	// function -> IDENTIFIER "(" parameters? ")" block ;
	function = func() (*Node, error) {
		var name Token
		var param *Node
		var err error
		var arity float32

		if match(Identifier) {
			name = previous()
		} else {
			prev := previous()
			return nil, fmt.Errorf("Parsing error on line %d: Expected function name after token \"%s\"", prev.Line, prev.Lexeme)
		}

		// params
		if match(LeftParen) {
			param, err = parameters()
			if err != nil {
				return nil, err
			}

			// check arity
			p := param
			for p != nil {
				p = p.Next
				arity++
			}
			if arity >= 255 {
				return nil, fmt.Errorf("Parsing error on line %d: Maximum argument count (254) exceeded with %d arguments", name.Line, int(arity))
			}
		} else {
			return nil, fmt.Errorf("Parsing error on line %d: Expected argument list after token \"%s\"", name.Line, name.Lexeme)
		}
		if !match(RightParen) {
			return nil, fmt.Errorf("Parsing error on line %d: Expected closing parenthesis after argument list", name.Line)
		}

		// body
		var body *Node
		if match(LeftBrace) {
			body, err = block()
		} else {
			return nil, fmt.Errorf("Parsing error on line %d: Expected function body", name.Line)
		}
		if err != nil {
			return nil, err
		}

		return &Node{
			Type: FunDeclNT,
			Data: encodeLoxNumber(arity),
			Left: &Node{
				Type: IdentifierNT,
				Data: encodeString(name.Lexeme),
			}, // name
			Right: param, // param list
			Third: body,  // function body
		}, err
	}

	// parameters -> IDENTIFIER ( "," IDENTIFIER )* ;
	parameters = func() (*Node, error) {
		var first *Node
		if match(Identifier) {
			first = &Node{Type: ParamNT, Data: encodeString(previous().Lexeme)}
		} else {
			return nil, nil // function takes zero parameters
		}
		param := first
		for {
			if match(Comma) && match(Identifier) {
				param.Next = &Node{Type: ParamNT, Data: encodeString(previous().Lexeme)}
				param = param.Next
			} else {
				break
			}
		}
		return first, nil
	}

	// varDecl -> "var" IDENTIFIER ( "=" expression )? ";" ;
	varDecl = func() (*Node, error) {
		ident, err := primary()
		var expr *Node
		if match(Equal) {
			expr, err = expression()
		}
		if match(Semicolon) {
			return &Node{
				Type:  VarDeclNT,
				Left:  ident,
				Right: expr,
			}, err
		}
		return nil, fmt.Errorf("Parsing error on line %d: Expected semicolon after token \"%s\"", tokens[current].Line, tokens[current].Lexeme)
	}

	// statement -> exprStmt | ifStmt | printStmt | block | returnStmt ;
	statement = func() (*Node, error) {
		if match(Print) {
			return printStmt()
		}
		if match(If) {
			return ifStmt()
		}
		if match(While) {
			return whileStmt()
		}
		if match(For) {
			return forStmt()
		}
		if match(LeftBrace) {
			return block()
		}
		if match(Return) {
			return returnStmt()
		}
		return exprStmt()
	}

	// block -> "{" declaration* "}" ;
	block = func() (*Node, error) {
		var prev *Node
		blk := &Node{Type: BlockNT}
		for !match(RightBrace) {
			decl, err := declaration()
			if err != nil {
				return nil, err
			}
			if prev == nil {
				blk.Right = decl
			} else {
				prev.Next = decl
			}
			prev = decl
		}

		if previous().Type == RightBrace {
			return blk, nil
		}
		fmt.Println(blk.ToSExpression())
		fmt.Println(previous().ToString())
		return nil, fmt.Errorf("Parsing error on line %d: Expected closing brace", tokens[current].Line)
	}

	// returnStmt -> "return" expression? ";" ;
	returnStmt = func() (*Node, error) {
		expr, err := expression()
		if err != nil {
			return nil, err
		}
		if match(Semicolon) {
			return &Node{
				Type:  ReturnStmtNT,
				Right: expr,
			}, err
		}
		return nil, fmt.Errorf("Parsing error on line %d: Expected semicolon after return statement", tokens[current].Line)
	}

	// forStmt -> "for" "(" varDecl | exprStmt | ";" ) expression? ";" expression? ")" statement ;
	forStmt = func() (*Node, error) {
		var init, cond, incr, body *Node
		var err error
		if !match(LeftParen) {
			return nil, fmt.Errorf("Parsing error on line %d: Expected left parenthesis", tokens[current].Line)
		}

		// initializer
		if match(Semicolon) {
			// leave initializer empty
		} else if match(Var) {
			init, err = varDecl()
		} else {
			init, err = exprStmt()
		}
		if err != nil {
			return nil, err
		}

		// condition
		cond, err = expression()
		if err != nil {
			return nil, err
		}
		if !match(Semicolon) {
			return nil, fmt.Errorf("Parsing error on line %d: Expected semicolon in for statement", tokens[current].Line)
		}

		// increment
		incr, err = expression()
		if err != nil {
			return nil, err
		}
		if !match(RightParen) {
			return nil, fmt.Errorf("Parsing error on line %d: Expected closing parenthesis in for statement", tokens[current].Line)
		}

		// body
		body, err = statement()
		if err != nil {
			return nil, err
		}
		if init == nil && cond == nil && incr == nil && body == nil {
			return nil, fmt.Errorf("Parsing error on line %d: For loop can not be entirely empty", tokens[current].Line)
		}

		// desugar into a while loop
		bodyWithIncr := &Node{
			Type:  BlockNT,
			Right: body,
		}
		body.Next = incr

		while := &Node{
			Type:  WhileStmtNT,
			Left:  cond,
			Right: bodyWithIncr,
		}
		if cond == nil {
			while.Left = &Node{Type: BoolNT, Data: encodeBool(true)} // nil condition means always true
		}

		forStmt := &Node{
			Type:  BlockNT,
			Right: init,
		}
		if init == nil {
			forStmt.Right = while
		} else {
			init.Next = while
		}

		return forStmt, nil
	}

	// whileStmt -> "while" "(" expression ")" statement ;
	whileStmt = func() (*Node, error) {
		var cond, body *Node
		var err error
		if match(LeftParen) {
			cond, err = expression()
			if err != nil {
				return nil, err
			}
			if match(RightParen) {
				body, err = statement()
				if err != nil {
					return nil, err
				}
			}
			if cond != nil && body != nil {
				return &Node{
					Type:  WhileStmtNT,
					Left:  cond,
					Right: body,
				}, err
			}
		}
		return nil, fmt.Errorf("Parsing error on line %d: Malformed \"while\" statement", tokens[current].Line)
	}

	// ifStmt	-> "if" "(" expression ")" statement ( "else" statement )? ;
	ifStmt = func() (*Node, error) {
		var cond, thenBranch, elseBranch *Node
		var err error
		if match(LeftParen) {
			cond, err = expression()
			if err != nil {
				return nil, err
			}
			if match(RightParen) {
				thenBranch, err = statement()
				if err != nil {
					return nil, err
				}
			}
			if match(Else) {
				elseBranch, err = statement()
				if err != nil {
					return nil, err
				}
			}
			if cond != nil && thenBranch != nil {
				n := &Node{
					Type:  IfStmtNT,
					Left:  cond,
					Right: thenBranch,
				}
				if elseBranch != nil {
					n.Third = elseBranch
				}
				return n, err
			}
			return nil, fmt.Errorf("Parsing error on line %d: Malformed \"if\" statement", tokens[current].Line)
		}
		return nil, fmt.Errorf("Parsing error on line %d: Expected parentheses after \"if\" token", tokens[current].Line)
	}

	// exprStmt -> expression ";" ;
	exprStmt = func() (*Node, error) {
		expr, err := expression()
		if match(Semicolon) {
			return &Node{Type: ExprStmtNT, Right: expr}, err
		}
		return nil, fmt.Errorf("Parsing error on line %d: Expected semicolon after token \"%s\"", tokens[current].Line, tokens[current].Lexeme)
	}

	// printStmt -> "print" expression ";" ;
	printStmt = func() (*Node, error) {
		expr, err := expression()
		if match(Semicolon) {
			return &Node{Type: PrintStmtNT, Right: expr}, err
		}
		return nil, fmt.Errorf("Parsing error on line %d: Expected semicolon after token \"%s\"", tokens[current].Line, tokens[current].Lexeme)
	}

	// expression -> assignment ;
	expression = func() (*Node, error) {
		return assignment()
	}

	// assignment -> IDENTIFIER "=" ( assignment | logicOr ) ;
	assignment = func() (*Node, error) {
		expr, err := logicOr()
		if match(Equal) {
			operator := previous()
			right, err := assignment()
			if err != nil {
				return nil, fmt.Errorf("Parsing error on line %d: Invalid r-value for assignment", tokens[current].Line)
			}
			if expr.Type == IdentifierNT {
				return &Node{
					Type:  AssignmentNT,
					Left:  expr,
					Data:  operator.toValue(),
					Right: right,
				}, err
			}
		}
		return expr, err
	}

	// logicOr	-> logicAnd ( "or" logicAnd )* ;
	logicOr = func() (*Node, error) {
		expr, err := logicAnd()
		for match(Or) {
			operator := previous()
			right, err := logicAnd()
			if err != nil {
				break
			}
			expr = &Node{
				Type:  LogicOrNT,
				Left:  expr,
				Data:  operator.toValue(),
				Right: right,
			}
		}
		return expr, err
	}

	// logicAnd -> equality ( "and" equality)* ;
	logicAnd = func() (*Node, error) {
		expr, err := equality()
		for match(And) {
			operator := previous()
			right, err := equality()
			if err != nil {
				break
			}
			expr = &Node{
				Type:  LogicAndNT,
				Left:  expr,
				Data:  operator.toValue(),
				Right: right,
			}
		}
		return expr, err
	}

	// equality -> comparison ( ( "!=" | "==" ) comparison )* ;
	equality = func() (*Node, error) {
		expr, err := comparison()
		for match(BangEqual, EqualEqual) {
			operator := previous()
			right, err := comparison()
			if err != nil {
				break
			}
			expr = &Node{
				Type:  EqualityNT,
				Left:  expr,
				Data:  operator.toValue(),
				Right: right,
			}
		}
		return expr, err
	}

	// comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
	comparison = func() (*Node, error) {
		expr, err := term()
		for match(Greater, GreaterEqual, Less, LessEqual) {
			operator := previous()
			right, err := term()
			if err != nil {
				break
			}
			expr = &Node{
				Type:  ComparisonNT,
				Left:  expr,
				Data:  operator.toValue(),
				Right: right,
			}
		}
		return expr, err
	}

	// term	-> factor ( ( "-" | "+" ) factor )* ;
	term = func() (*Node, error) {
		expr, err := factor()
		for match(Minus, Plus) {
			operator := previous()
			right, err := factor()
			if err != nil {
				break
			}
			expr = &Node{
				Type:  TermNT,
				Left:  expr,
				Data:  operator.toValue(),
				Right: right,
			}
		}
		return expr, err
	}

	// factor	-> unary ( ( "/" | "*" ) unary )* ;
	factor = func() (*Node, error) {
		expr, err := unary()
		for match(Slash, Star) {
			operator := previous()
			right, err := unary()
			if err != nil {
				break
			}
			expr = &Node{
				Type:  FactorNT,
				Left:  expr,
				Data:  operator.toValue(),
				Right: right,
			}
		}
		return expr, err
	}

	// unary -> ( "!" | "-" ) unary | call ;
	unary = func() (*Node, error) {
		if match(Bang, Minus) {
			operator := previous()
			right, err := unary()
			return &Node{
				Type:  UnaryNT,
				Data:  operator.toValue(),
				Right: right,
			}, err
		}
		return call()
	}

	var finishCall func() (*Node, float32, error)
	// call -> primary ( "(" arguments? ")" )* ;
	call = func() (*Node, error) {
		expr, err := primary()
		for {
			if match(LeftParen) {
				arg, arity, err := finishCall()
				if err != nil {
					return nil, err
				}
				expr = &Node{
					Type:  CallNT,
					Data:  encodeLoxNumber(arity),
					Left:  expr, // just IdentifierNT now. CallableNT as wrapper later, for object methods
					Right: arg,  // arg list (ArgNT?), tied together through Next
				}
				if !match(RightParen) {
					return nil, fmt.Errorf("Parsing error on line %d: Expected closing parenthesis after argument list", previous().Line)
				}
			} else {
				break
			}
		}
		return expr, err
	}

	// arguments
	finishCall = func() (*Node, float32, error) {
		var first *Node
		var err error
		var count float32

		first, err = expression()
		if err != nil {
			return nil, count, err
		}
		if first != nil {
			count = 1
		}

		var arg, next *Node
		for match(Comma) {
			count++
			arg, err = expression()
			if err != nil {
				return nil, count, err
			}

			if first.Next == nil {
				first.Next = arg
				next = first.Next
			} else {
				next.Next = arg
				next = next.Next
			}
		}

		if count >= 255 {
			return nil, count, fmt.Errorf("Parsing error: Maximum argument count (254) exceeded with %f arguments", count)
		}
		return first, count, err
	}

	// primary -> IDENTIFIER | NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
	primary = func() (*Node, error) {
		if match(Identifier) {
			return &Node{Type: IdentifierNT, Data: previous().toValue()}, nil
		}
		if match(Number) {
			return &Node{Type: NumberNT, Data: previous().toValue()}, nil
		}
		if match(String) {
			return &Node{Type: StringNT, Data: previous().toValue()}, nil
		}
		if match(True, False) {
			return &Node{Type: BoolNT, Data: previous().toValue()}, nil
		}
		if match(Nil) {
			return &Node{Type: NilNT, Data: previous().toValue()}, nil
		}
		if match(LeftParen) {
			expr, err := expression()
			if match(RightParen) {
				return &Node{
					Type:  GroupNT,
					Right: expr}, err
			}
			return nil, fmt.Errorf("Parsing error on line %d: Expected closing parenthesis following token \"%s\"", tokens[current].Line, tokens[current].Lexeme)
		}
		return nil, fmt.Errorf("Parsing error on line %d: Unexpected token \"%s\"", tokens[current].Line, tokens[current].Lexeme)
	}

	return program()
}
