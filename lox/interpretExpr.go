package lox

import "fmt"

func (env *Environment) interpretOr(expr *Node) *Node {
	left := env.interpretExpr(expr.Left)
	if left.truthy() {
		return left
	}
	right := env.interpretExpr(expr.Right)
	if right.truthy() {
		return right
	}
	return &Node{
		Type: BoolNT,
		Data: encodeBool(false),
	}
}

func (env *Environment) interpretAnd(expr *Node) *Node {
	left := env.interpretExpr(expr.Left)
	if left.truthy() {
		right := env.interpretExpr(expr.Right)
		if right.truthy() {
			return &Node{
				Type: BoolNT,
				Data: encodeBool(true),
			}
		}
	}
	return &Node{
		Type: BoolNT,
		Data: encodeBool(false),
	}
}

func (env *Environment) interpretEquality(expr *Node) *Node {
	left := env.interpretExpr(expr.Left)
	right := env.interpretExpr(expr.Right)
	switch expr.ToString() {
	case "==":
		if left.Type != right.Type {
			return &Node{
				Type: BoolNT,
				Data: encodeBool(false),
			}
		}
		return &Node{
			Type: BoolNT,
			Data: encodeBool(compareValues(left.Data, right.Data)),
		}
	case "!=":
		if left.Type != right.Type {
			return &Node{
				Type: BoolNT,
				Data: encodeBool(true),
			}
		}
		return &Node{
			Type: BoolNT,
			Data: encodeBool(!compareValues(left.Data, right.Data)),
		}
	}
	fmt.Printf("Runtime error: expected equality expression, instead found \"%s\"", expr.ToString())
	return nil
}

func (env *Environment) interpretComparison(expr *Node) *Node {
	left := env.interpretExpr(expr.Left)
	right := env.interpretExpr(expr.Right)
	if left.Type != NumberNT || right.Type != NumberNT {
		fmt.Printf("\nRuntime error: cannot compare type \"%s\" with type \"%s\"", left.ToString(), right.ToString())
		return nil
	}
	numL, numR := decodeLoxNumber(left.Data), decodeLoxNumber(right.Data)
	switch expr.ToString() {
	case "<":
		return &Node{
			Type: BoolNT,
			Data: encodeBool(numL < numR),
		}
	case "<=":
		return &Node{
			Type: BoolNT,
			Data: encodeBool(numL <= numR),
		}
	case ">":
		return &Node{
			Type: BoolNT,
			Data: encodeBool(numL > numR),
		}
	case ">=":
		return &Node{
			Type: BoolNT,
			Data: encodeBool(numL >= numR),
		}
	}
	fmt.Printf("\nRuntime error: expected comparison expression, instead found \"%s\"", expr.ToString())
	return nil
}

func (env *Environment) interpretTerm(expr *Node) *Node {
	switch expr.ToString() {
	case "+":
		left := env.interpretExpr(expr.Left)
		right := env.interpretExpr(expr.Right)
		if left.Type == NumberNT && right.Type == NumberNT {
			numL, numR := decodeLoxNumber(left.Data), decodeLoxNumber(right.Data)
			return &Node{
				Type: NumberNT,
				Data: encodeLoxNumber(numL + numR),
			}
		}
		if left.Type == StringNT && right.Type == StringNT {
			// string concatenation
			return &Node{
				Type: StringNT,
				Data: append(left.Data, right.Data...),
			}
		}
		fmt.Printf("\nRuntime error: cannot add \"%s\" and \"%s\"", left.ToString(), right.ToString())
		return nil
	case "-":
		left := env.interpretExpr(expr.Left)
		right := env.interpretExpr(expr.Right)
		if left.Type != NumberNT || right.Type != NumberNT {
			fmt.Printf("\nRuntime error: cannot subtract type \"%s\" and type \"%s\"", left.ToString(), right.ToString())
			return nil
		}
		numL, numR := decodeLoxNumber(left.Data), decodeLoxNumber(right.Data)
		return &Node{
			Type: NumberNT,
			Data: encodeLoxNumber(numL - numR),
		}
	}
	fmt.Printf("\nRuntime error: expected addition/subtraction expression, instead found \"%s\"", expr.ToString())
	return nil
}

func (env *Environment) interpretFactor(expr *Node) *Node {
	switch expr.ToString() {
	case "*":
		left := env.interpretExpr(expr.Left)
		right := env.interpretExpr(expr.Right)
		if left.Type != NumberNT || right.Type != NumberNT {
			fmt.Printf("\nRuntime error: cannot multiply type \"%s\" and type \"%s\"", left.ToString(), right.ToString())
			return nil
		}
		numL, numR := decodeLoxNumber(left.Data), decodeLoxNumber(right.Data)
		return &Node{
			Type: NumberNT,
			Data: encodeLoxNumber(numL * numR),
		}
	case "/":
		left := env.interpretExpr(expr.Left)
		right := env.interpretExpr(expr.Right)
		if left.Type != NumberNT || right.Type != NumberNT {
			fmt.Printf("\nRuntime error: cannot multiply type \"%s\" and type \"%s\"", left.ToString(), right.ToString())
			return nil
		}
		numL, numR := decodeLoxNumber(left.Data), decodeLoxNumber(right.Data)
		return &Node{
			Type: NumberNT,
			Data: encodeLoxNumber(numL / numR),
		}
	}
	fmt.Printf("\nRuntime error: expected multiplication/division expression, instead found \"%s\"", expr.ToString())
	return nil
}

func (env *Environment) interpretUnary(expr *Node) *Node {
	switch expr.ToString() {
	case "!":
		right := env.interpretExpr(expr.Right)
		return &Node{
			Type: BoolNT,
			Data: encodeBool(!right.truthy()),
		}
	case "-":
		right := env.interpretExpr(expr.Right)
		if expr.Type != NumberNT {
			fmt.Printf("\nRuntime error: operator \"-\" undefined for \"%s\"", expr.ToString())
			return nil
		}
		num := right.Data
		num[0] ^= 1 << 7
		return &Node{
			Type: NumberNT,
			Data: num,
		}
	}
	fmt.Printf("\nRuntime error: expected unary expression, instead found \"%s\"", expr.ToString())
	return nil
}

func (env *Environment) interpretIdentifier(expr *Node) *Node {
	var val *Node
	var ok bool
	name := expr.ToString()
	for scope := env; !ok && scope != nil; scope = scope.Enclosing {
		val, ok = scope.Values[name]
	}
	if !ok || val == nil {
		fmt.Printf("\nRuntime error: undefined variable \"%s\"", name)
		return nil
	}
	return val
}
