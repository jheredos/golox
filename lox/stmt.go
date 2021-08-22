package lox

type Stmt interface {
	Accept(v StmtVisitor) interface{}
}

type StmtVisitor interface {
}
