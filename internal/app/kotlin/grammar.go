package kotlin

type Fqn []string

//type KtFile struct {
//	ImportList []Import
//	TopLevels  []any // Top-level
//}
//
//type FuncDecl struct {
//	Name string
//	Body Block
//	Expr any // = expression
//}
//
//type Block struct {
//	Statements []any // Statement
//}
//
//type VarDecl struct {
//	Name           string
//	Mutable        bool
//	TypeAnnotation UserType
//	Assignment     any // Expression
//}
//
//type CtorInvoke struct {
//	Type      UserType
//	ValueArgs []any // Value argument -> Expression
//}
//
//type MethodCall struct {
//	Receiver  string
//	Method    string
//	ValueArgs []any // Value argument -> Expression
//}
//
//type FuncCall struct {
//	Name      string
//	ValueArgs []any // Value argument -> Expression
//}
//
//type NamedArg struct {
//	Name  string
//	Value any // Expression
//}
//
//type PropAccess struct {
//	Object string
//	Prop   string
//}
//
//type UserType Fqn
//
//type Import struct {
//	fqn *Fqn
//}
//
//type Fqn []string
//
//type StringLiteral string
//type BoolLiteral bool
//
//type LambdaLiteral struct {
//	Statements []any // Statement
//}
//
//type InlineLambdaLiteral struct {
//	Statements []any // Statement
//}
//
//type PropAssignment struct {
//	Prop string
//	Expr any
//}
//
//type returnStat struct {
//	expr any
//}
//
//type statementKind int
//
//const (
//	statUnknown statementKind = iota
//	statReturn
//)
//
//type funcDecl struct {
//	signature *funcSignature
//	bodyScope *Scope
//}
//
//type funcSignature struct {
//	fqn        *Fqn
//	params     []parameter
//	returnType *ktType
//}
//
//type parameter struct {
//	paramType *ktType
//}
//
//type ktType struct {
//	kind  ktTypeKind
//	class *ktClass
//}
//
//type ktClass struct {
//	fqn          *Fqn
//	constructors []*funcSignature
//}
//
//type ktTypeKind int
//
//const (
//	ktUnknownType ktTypeKind = iota
//	ktUnit
//	ktAny
//	ktBlock
//	ktUserType
//)
//
//type funcCall struct {
//	signature *funcSignature
//	args      []any
//}
//
//type block struct {
//	scope *Scope
//}
//
//type ctorCall struct {
//	class *ktClass
//}
//
//type varDecl struct {
//	name    string
//	varType *ktType
//	expr    any
//}
