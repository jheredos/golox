package lox

import (
	"fmt"
	"math"
)

// Node represents a node in the AST. Left and Right refer to the next branches of the AST, and Type tells you what to expect in each place. Leaf nodes store the Token's Literal value in Node.Val
type Node struct {
	Type  NodeType
	Left  *Node
	Right *Node
	Third *Node
	Next  *Node
	Data  Value
}

// Value wraps disparate values
type Value []byte

// NodeType represents the types of AST Nodes, from top-level program nodes to literals like Bool and Number
type NodeType uint8

// NodeType values
const (
	ProgramNT NodeType = iota
	DeclarationNT
	VarDeclNT
	FunDeclNT
	FunctionNT
	StmtNT
	BlockNT
	ReturnStmtNT
	ExprStmtNT
	PrintStmtNT
	WhileStmtNT // For loops are desugared into while loops
	IfStmtNT
	AssignmentNT
	LogicOrNT
	LogicAndNT
	EqualityNT
	ComparisonNT
	TermNT   // addition/subtraction
	FactorNT // multiplication/division
	UnaryNT
	ArgNT
	ParamNT
	CallNT
	CallableNT
	IdentifierNT
	NumberNT
	StringNT
	BoolNT
	GroupNT
	NilNT
	EOFNT
)

func (t Token) toValue() Value {
	var val Value
	switch t.Type {
	case Number:
		val = encodeLoxNumberFromString(t.Lexeme)
	case String:
		val = []byte(t.Lexeme)
	case True:
		val = []byte{1}
	case False:
		val = []byte{0}
	default:
		val = []byte(t.Lexeme)
	}
	return val
}

func encodeBool(b bool) Value {
	if b {
		return []byte{1}
	}
	return []byte{0}
}

func encodeString(s string) Value {
	return []byte(s)
}

func compareValues(a Value, b Value) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func encodeLoxNumberFromString(s string) Value {
	var n float32 = 0
	var dec float32 = 0

	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			dec = 10
			continue
		}
		if dec == 0 {
			n *= 10
			n += float32(s[i] - '0')
		} else {
			n += float32(s[i]-'0') / dec
			dec *= 10
		}
	}

	return encodeLoxNumber(n)
}

func encodeLoxNumber(n float32) Value {
	f := math.Float32bits(n)
	return []byte{
		byte(f >> 24),
		byte(f >> 16),
		byte(f >> 8),
		byte(f),
	}
}

func decodeLoxNumber(v Value) float32 {
	var u uint32
	for i := 0; i < 4; i++ {
		u |= uint32(v[i]) << ((3 - i) * 8)
	}
	return math.Float32frombits(u)
}

func trimNumber(s string) string {
	for i := len(s) - 1; i > 0; i-- {
		if s[i] != '0' {
			if s[i] == '.' {
				return s[:i]
			}
			return s[:i+1]
		}
	}
	return s
}

// ToSExpression converts an AST into parenthesized S-expressions
func (n *Node) ToSExpression() string {
	if n == nil {
		return ""
	}
	switch n.Type {
	case NumberNT, StringNT, BoolNT, NilNT, ParamNT:
		return n.ToString()
	default:
		s := "(" + n.ToString()
		if n.Left != nil {
			s += " " + n.Left.ToSExpression()
		}
		if n.Right != nil {
			s += " " + n.Right.ToSExpression()
		}
		if n.Third != nil {
			s += " " + n.Third.ToSExpression()
		}
		s += ")"
		if n.Next == nil {
			return s
		}
		return s + "\n -> " + n.Next.ToSExpression()
	}
}

// ToString represents a AST Node as a string
func (n *Node) ToString() string {
	if n == nil {
		return ""
	}
	switch n.Type {
	case ProgramNT:
		return "<program>"
	case DeclarationNT:
		return "<declaration>"
	case VarDeclNT:
		return "<variable declaration>"
	case FunDeclNT:
		return "<function declaration \"" + n.Left.ToString() + "\">"
	case FunctionNT:
		return "<function object>"
	case BlockNT:
		return "<block>"
	case ReturnStmtNT:
		return "<return>"
	case WhileStmtNT:
		return "<while>"
	case IfStmtNT:
		return "<if>"
	case AssignmentNT:
		return "<assignment>"
	case LogicOrNT:
		return "<or>"
	case LogicAndNT:
		return "<and>"
	case ArgNT:
		return "<argument" + string(n.Data) + ">"
	case CallNT:
		return "<\"" + n.Left.ToString() + "\" call>"
	case CallableNT:
		return "<callable>"
	case StmtNT:
		return "<statement>"
	case ExprStmtNT:
		return "<expression statement>"
	case PrintStmtNT:
		return "print"
	case EqualityNT:
		return string(n.Data)
	case ComparisonNT:
		return string(n.Data)
	case TermNT:
		return string(n.Data)
	case FactorNT:
		return string(n.Data)
	case UnaryNT:
		return string(n.Data)
	case IdentifierNT, ParamNT:
		return string(n.Data)
	case GroupNT:
		return "<group>"
	case EOFNT:
		return "<end-of-file>"
	case NumberNT:
		return trimNumber(fmt.Sprintf("%f", decodeLoxNumber(n.Data)))
	case BoolNT:
		if n.Data[0] == 1 {
			return "true"
		}
		return "false"
	case NilNT:
		return "nil"
	default:
		if n.Data != nil {
			return string(n.Data)
		}
		return "<unknown>"
	}
}

func (n *Node) truthy() bool {
	if n.Type == BoolNT && n.Data[0] == 0 {
		return false
	}
	if n.Type == NilNT {
		return false
	}
	return true
}
