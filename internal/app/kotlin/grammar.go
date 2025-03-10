package kotlin

type KtFile struct {
	ImportList []Import
	TopLevels  []any // Top-level
}

type FuncDecl struct {
	Name SimpleId
	Body Block
	Expr any // = expression
}

type Block struct {
	Statements []any // Statement
}

type VarDecl struct {
	Name           SimpleId
	Mutable        bool
	TypeAnnotation UserType
	Assignment     any // Expression
}

type CtorInvoke struct {
	Type      UserType
	ValueArgs []any // Value argument -> Expression
}

type CallExpr struct {
	Receiver  SimpleId
	Method    SimpleId
	ValueArgs []any // Value argument -> Expression
}

type UserType Fqn

type Import struct {
	fqn Fqn
}

type Fqn []SimpleId
type SimpleId string

type StringLiteral string

type LambdaLiteral struct {
	Statements []any // Statement
}

type PropAssignment struct {
	Prop SimpleId
	Expr any
}
