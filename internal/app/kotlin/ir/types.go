package ir

import "c2k/internal/app/kotlin"

type Scope struct {
	FuncDecls  []*FuncDecl
	Variables  map[*Object]string
	Statements []any
	Imports    map[*kotlin.Fqn]struct{}
	TopLevel   *Scope
}

type Object struct {
	KtType *ktType
}

type FuncDecl struct {
	Signature *funcSignature
	BodyScope *Scope
}

type funcSignature struct {
	Fqn        *kotlin.Fqn
	Params     []parameter
	returnType *ktType
}

type parameter struct {
	paramType *ktType
}

type Class struct {
	Fqn          *kotlin.Fqn
	constructors []*funcSignature
}

type ktType struct {
	kind  ktTypeKind
	class *Class
}

type ktTypeKind int

const (
	ktUnknownType ktTypeKind = iota
	ktUnit
	ktAny
	ktBlock
	ktUserType
)

var unitType = &ktType{
	kind: ktUnit,
}

var anyType = &ktType{
	kind: ktAny,
}

var voidBlock = &ktType{
	kind: ktBlock,
}

type FuncCall struct {
	Signature *funcSignature
	Args      []any
}

type Block struct {
	Scope *Scope
}

type CtorCall struct {
	Class *Class
	Args  []any
}

type MethodCall struct {
	Object *Object
	Method *kotlin.Fqn
	Args   []any
}

type VarDecl struct {
	Name    string
	varType *ktType
	Expr    any
}

type ReturnStat struct {
	Expr any
}

type PropAssign struct {
	Prop string
	Expr any
}
