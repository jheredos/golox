package lox

// Expr interface for the visitor pattern
type Expr interface {
	Accept(v ExprVisitor) interface{}
}

// BinaryExpr - binary expression with an infix operator and a left and right argument
type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

// Accept ...
func (be BinaryExpr) Accept(v ExprVisitor) interface{} {
	return nil
}

// UnaryExpr - expression with a unary operator, ! or -
type UnaryExpr struct {
	Right    Expr
	Operator Token
}

// Accept ...
func (ue UnaryExpr) Accept(v ExprVisitor) interface{} {
	return nil
}

// GroupingExpr denotes an expression in parentheses
type GroupingExpr struct {
	Expr Expr
}

// Accept ...
func (ge GroupingExpr) Accept(v ExprVisitor) interface{} {
	return nil
}

// LiteralExpr ... not sure if this one is needed
type LiteralExpr struct {
	Value interface{}
}

// Accept ...
func (le LiteralExpr) Accept(v ExprVisitor) interface{} {
	return nil
}

// ExprVisitor ...
type ExprVisitor interface {
	VisitBinaryExpr(be BinaryExpr) interface{}
	VisitUnaryExpr(ue UnaryExpr) interface{}
	VisitGroupingExpr(ge GroupingExpr) interface{}
	VisitLiteralExpr(le LiteralExpr) interface{}
}
