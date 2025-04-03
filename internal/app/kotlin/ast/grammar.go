package ast

import "c2k/internal/app/kotlin"

type KtFile struct {
	ImportList []KtImport
	TopLevels  []any // Top-level
}

type KtFuncDecl struct {
	Name string
	Body KtBlock
	Expr any // = expression
}

type KtBlock struct {
	Statements []any // Statement
}

type KtVarDecl struct {
	Name       string
	Mutable    bool
	Assignment any // Expression
}

type KtInvocation struct {
	Receiver  string
	Method    string
	ValueArgs []any // Value argument -> Expression
}

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

type KtNamedArg struct {
	Name  string
	Value any // Expression
}

type KtPropAccess struct {
	Object string
	Prop   string
}

//type UserType Fqn

type KtImport struct {
	Fqn *kotlin.Fqn
}

type KtStringLiteral string
type KtBoolLiteral bool

type KtLambdaLiteral struct {
	Statements []any // Statement
}

type KtInlineLambdaLiteral struct {
	Statements []any // Statement
}

type KtPropAssignment struct {
	Prop string
	Expr any
}

//type returnStat struct {
//	expr any
//}

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
//	Fqn        *Fqn
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
//	Fqn          *Fqn
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
