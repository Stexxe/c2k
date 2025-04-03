package ir

import "c2k/internal/app/kotlin"

var clientPackage = &kotlin.Fqn{"io", "ktor", "client"}
var requestPackage = &kotlin.Fqn{"io", "ktor", "client", "request"}
var formsPackage = &kotlin.Fqn{"io", "ktor", "client", "request", "forms"}
var clientStatementPackage = &kotlin.Fqn{"io", "ktor", "client", "statement"}
var httpPackage = &kotlin.Fqn{"io", "ktor", "http"}
var coroutinesPackage = &kotlin.Fqn{"kotlinx", "coroutines"}
var javaIoPackage = &kotlin.Fqn{"java", "io"}
var cioUtilsPackage = &kotlin.Fqn{"io", "ktor", "util", "cio"}

var mainSign = &funcSignature{
	Fqn:        &kotlin.Fqn{"main"},
	returnType: unitType,
}

var runBlockingSign = &funcSignature{
	Fqn:        buildFqn("runBlocking", coroutinesPackage),
	Params:     []parameter{{voidBlock}},
	returnType: anyType,
}

var printSign = &funcSignature{
	Fqn:        &kotlin.Fqn{"print"},
	Params:     []parameter{{anyType}},
	returnType: unitType,
}

func runBlockingCall(parentScope *Scope, builder func(sc *Scope)) *FuncCall {
	blockScope := newScope(parentScope)
	builder(blockScope)
	addImport(parentScope, runBlockingSign.Fqn)
	return &FuncCall{Signature: runBlockingSign, Args: []any{&Block{Scope: blockScope}}}
}

func printCall(sc *Scope, arg any) *FuncCall {
	call := &FuncCall{Signature: printSign, Args: []any{arg}}
	sc.Statements = append(sc.Statements, call)
	return call
}

func httpClientCtorBlock(sc *Scope, builder func(sc *Scope)) *CtorCall {
	blockScope := newScope(sc)
	builder(blockScope)
	return callCtor(sc, httpClientType.class, &Block{Scope: blockScope})
}

var getRequest = buildFqn("get", requestPackage)
var bodyAsText = buildFqn("bodyAsText", clientStatementPackage)

var httpClientType = &ktType{
	kind: ktUserType,
	class: &Class{
		Fqn: buildFqn("HttpClient", clientPackage),
		constructors: []*funcSignature{
			{},
		},
	},
}

var httpResponseType = &ktType{
	kind: ktUserType,
}

func buildFqn(name string, pack *kotlin.Fqn) *kotlin.Fqn {
	fqn := append(kotlin.Fqn{}, *pack...)
	fqn = append(fqn, name)
	return &fqn
}
