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
	Receiver  any
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
type Id string
type Fqn []string

type EmptyStatement struct{}

type LambdaLiteral struct {
	Statements []any // Statement
}

type InlineLambdaLiteral struct {
	Statements []any // Statement
}

type ForInLoop struct {
	Bind       any
	Expr       any
	Statements []any
}

type PairDestruct []string

type PropAssignment struct {
	Prop string
	Expr any
}
