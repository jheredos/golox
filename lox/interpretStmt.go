package lox

import "fmt"

func (env *Environment) interpretVarDecl(stmt *Node) *Node {
	name := stmt.Left.ToString()
	if _, already := env.Values[name]; already {
		fmt.Printf("\nRuntime error: variable \"%s\" redeclared", name)
		return nil
	}
	val := env.interpretExpr(stmt.Right)
	env.Values[name] = val

	return stmt.Next
}

func (env *Environment) interpretFunDecl(stmt *Node) *Node {
	name := stmt.Left.ToString()
	if _, already := env.Values[name]; already {
		fmt.Printf("\nRuntime error: function \"%s\" redeclared", name)
		return nil
	}

	env.Values[name] = &Node{
		Type:  FunctionNT,
		Data:  stmt.Data,  // arity (number)
		Left:  stmt.Right, // params, connected by Next
		Right: stmt.Third, // function body
	}

	return stmt.Next
}

func (env *Environment) interpretBlock(stmt *Node) *Node {
	scope := &Environment{Enclosing: env, Values: make(map[string]*Node)}
	next := stmt.Right
	for next != nil {
		if next.Type == ReturnStmtNT {
			// break block for return stmts
			return &Node{
				Type:  ReturnStmtNT,
				Right: scope.interpretExpr(next.Right),
				Next:  stmt.Next,
			}
		}
		next = scope.interpretStmt(next)
	}
	return stmt.Next
}

func (env *Environment) interpretIfStmt(stmt *Node) *Node {
	cond := env.interpretExpr(stmt.Left)
	if cond.truthy() {
		stmt.Right.Next = stmt.Next
		return stmt.Right
	}
	if stmt.Third != nil {
		stmt.Third.Next = stmt.Next
		return stmt.Third
	}
	return stmt.Next
}

func (env *Environment) interpretWhileStmt(stmt *Node) *Node {
	scope := &Environment{Enclosing: env, Values: make(map[string]*Node)}
	for cond := scope.interpretExpr(stmt.Left); cond.truthy(); cond = scope.interpretExpr(stmt.Left) {
		res := scope.interpretStmt(stmt.Right)
		if res != nil && res.Type == ReturnStmtNT {
			// break loop for return stmts
			return &Node{
				Type:  ReturnStmtNT,
				Right: scope.interpretExpr(res.Right),
				Next:  stmt.Next,
			}
		}
	}
	return stmt.Next
}

func (env *Environment) interpretAssignment(stmt *Node) *Node {
	name := stmt.Left.ToString()
	val := env.interpretExpr(stmt.Right)

	for scope := env; scope != nil; scope = scope.Enclosing {
		_, ok := scope.Values[name]
		if ok {
			scope.Values[name] = val
			return stmt.Next
		}
	}

	fmt.Printf("\nRuntime error: undeclared variable \"%s\"", name)
	return nil
}

func (env *Environment) interpretCall(stmt *Node) *Node {
	// TODO: nested calls, eg foo(bar)(baz)()
	if stmt.Left.Type != IdentifierNT {
		fmt.Printf("\nRuntime error: \"%s\" is not callable", stmt.ToString())
		return nil
	}

	name := stmt.Left.ToString()
	var fun *Node
	var ok bool
	for scope := env; !ok && scope != nil; scope = scope.Enclosing {
		fun, ok = scope.Values[name]
	}
	if !ok || fun == nil {
		fmt.Printf("\nRuntime error: Function %s is undefined", name)
		return nil
	}

	// set up function's environment with param values
	funcEnv := &Environment{
		Enclosing: env,
		Values:    make(map[string]*Node),
	}
	for arg, param := stmt.Right, fun.Left; arg != nil || param != nil; arg, param = arg.Next, param.Next {
		if arg == nil && param != nil {
			fmt.Printf("\nRuntime error: Too few parameters for function %s, (expected %f)", stmt.Left.ToString(), decodeLoxNumber(fun.Data))
			return nil
		}
		if param == nil && arg != nil {
			fmt.Printf("\nRuntime error: Too many parameters for function %s, (expected %f)", stmt.Left.ToString(), decodeLoxNumber(fun.Data))
			return nil
		}
		val := funcEnv.interpretExpr(arg)
		funcEnv.Values[param.ToString()] = val
	}

	// execute function
	result := funcEnv.interpretStmt(fun.Right)
	if result.Type == ReturnStmtNT {
		// when call is expr, return return stmt, but set next stmt to stmt following call
		return &Node{
			Type:  ReturnStmtNT,
			Right: funcEnv.interpretExpr(result.Right),
			Next:  stmt.Next,
		}
	}

	return stmt.Next // when call is stmt, return next stmt
}

func (env *Environment) interpretReturnStmt(stmt *Node) *Node {
	return stmt
}
