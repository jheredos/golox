package lox

import "fmt"

// Interpret is the main function called on a Lox program
func (prgm *Node) Interpret() {
	if prgm.Type != ProgramNT {
		fmt.Printf("\nRuntime error: ...")
		return
	}
	global := &Environment{Values: make(map[string]*Node)}
	global.setNativeFunctions()
	stmt := prgm.Right

	// fmt.Println("Program S-expression:")
	// fmt.Println(stmt.ToSExpression(), "\n\n")

	for stmt != nil {
		stmt = global.interpretStmt(stmt)
	}

}

// interpretStmt dispatches statement nodes to functions that handle particular types of statements
func (env *Environment) interpretStmt(stmt *Node) *Node {
	var next *Node
	switch stmt.Type {
	case DeclarationNT, StmtNT, ExprStmtNT:
		_ = env.interpretStmt(stmt.Right)
		next = stmt.Next
	case VarDeclNT:
		next = env.interpretVarDecl(stmt)
	case FunDeclNT:
		next = env.interpretFunDecl(stmt)
	case BlockNT:
		next = env.interpretBlock(stmt)
	case IfStmtNT:
		next = env.interpretIfStmt(stmt)
	case WhileStmtNT:
		next = env.interpretWhileStmt(stmt)
	case PrintStmtNT:
		val := env.interpretExpr(stmt.Right)
		fmt.Println(val.ToString())
		next = stmt.Next
	case AssignmentNT:
		next = env.interpretAssignment(stmt)
	case CallNT:
		next = env.interpretCall(stmt)
	case ReturnStmtNT:
		next = env.interpretReturnStmt(stmt)
	default:
		fmt.Printf("\nRuntime error: \"%s\" is not a statement", stmt.ToString())
		return nil
	}
	return next
}

// interpretExpr dispatches expression nodes to functions that evaluate particular types of expressions
func (env *Environment) interpretExpr(expr *Node) *Node {
	result := &Node{Type: NilNT}
	switch expr.Type {
	case CallNT:
		// call can be stmt or expr
		result = env.interpretCall(expr).Right
	// case CallableNT:
	// 	result = env.interpretCallable(n, env)
	case LogicOrNT:
		result = env.interpretOr(expr)
	case LogicAndNT:
		result = env.interpretAnd(expr)
	case EqualityNT:
		result = env.interpretEquality(expr)
	case ComparisonNT:
		result = env.interpretComparison(expr)
	case TermNT:
		result = env.interpretTerm(expr)
	case FactorNT:
		result = env.interpretFactor(expr)
	case UnaryNT:
		result = env.interpretUnary(expr)
	case IdentifierNT, ParamNT:
		result = env.interpretIdentifier(expr)
	case NumberNT, StringNT, BoolNT, NilNT, FunctionNT:
		result = expr
	}

	return result
}

func (env *Environment) setNativeFunctions() {
	env.Values["clock"] = &Node{
		Type:  CallableNT,
		Data:  []byte{0}, // arity
		Left:  &Node{Type: IdentifierNT, Data: encodeString("clock")},
		Right: nil,
		// TODO
	}

}
