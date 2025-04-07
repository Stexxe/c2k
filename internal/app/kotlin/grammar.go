package kotlin

type KtFile struct {
	ImportList []Import
	TopLevels  []any // Top-level
}

type FuncDecl struct {
	Name string
	Body Block
	Expr any // = expression
}

type Block struct {
	Statements []any // Statement
}

type VarDecl struct {
	Name           string
	Mutable        bool
	TypeAnnotation UserType
	Assignment     any // Expression
}

type CtorInvoke struct {
	Type      UserType
	ValueArgs []any // Value argument -> Expression
}

type MethodCall struct {
	Receiver  string
	Method    string
	ValueArgs []any // Value argument -> Expression
}

type FuncCall struct {
	Name      string
	ValueArgs []any // Value argument -> Expression
}

type NamedArg struct {
	Name  string
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

type Fqn []string

//type StringLiteral string
//type BoolLiteral bool

type LambdaLiteral struct {
	Statements []any // Statement
}

type InlineLambdaLiteral struct {
	Statements []any // Statement
}

type PropAssignment struct {
	Prop string
	Expr any
}
