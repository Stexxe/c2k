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

type MethodCall struct {
	Receiver  SimpleId
	Method    SimpleId
	ValueArgs []any // Value argument -> Expression
}

type FuncCall struct {
	Name      SimpleId
	ValueArgs []any // Value argument -> Expression
}

type NamedArg struct {
	Name  SimpleId
	Value any // Expression
}

type PropAccess struct {
	Object string
	Prop   string
}

type UserType Fqn

type Import struct {
	fqn Fqn
}

type Fqn []SimpleId
type SimpleId string

type StringLiteral string
type BoolLiteral bool

type LambdaLiteral struct {
	Statements []any // Statement
}

type InlineLambdaLiteral struct {
	Statements []any // Statement
}

type PropAssignment struct {
	Prop SimpleId
	Expr any
}
