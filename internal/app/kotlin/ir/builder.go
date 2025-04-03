package ir

import "c2k/internal/app/kotlin"

func newScope(parent *Scope) *Scope {
	sc := &Scope{}

	if parent != nil {
		sc.TopLevel = parent.TopLevel
	} else {
		sc.TopLevel = sc
	}

	sc.Imports = make(map[*kotlin.Fqn]struct{})
	sc.Variables = make(map[*Object]string)
	return sc
}

func declareFunc(sc *Scope, signature *funcSignature, builder func(sc *Scope)) {
	childScope := newScope(sc)
	builder(childScope)
	decl := &FuncDecl{Signature: signature, BodyScope: childScope}
	sc.FuncDecls = append(sc.FuncDecls, decl)
}

func ret(sc *Scope, expr any) {
	sc.Statements = append(sc.Statements, &ReturnStat{expr})
}

func addImport(sc *Scope, fqn *kotlin.Fqn) {
	sc.TopLevel.Imports[fqn] = struct{}{}
}

func declareVal(sc *Scope, name string, valType *ktType, expr any) *Object {
	decl := &VarDecl{Name: name, varType: valType, Expr: expr}
	sc.Statements = append(sc.Statements, decl)

	obj := &Object{KtType: valType}
	sc.Variables[obj] = name

	return obj
}

func callCtor(sc *Scope, class *Class, args ...any) *CtorCall {
	addImport(sc, class.Fqn)
	return &CtorCall{Class: class, Args: args}
}

func callMethod(sc *Scope, obj *Object, method *kotlin.Fqn, args ...any) *MethodCall {
	addImport(sc, method)
	return &MethodCall{Object: obj, Method: method, Args: args}
}

func assignProp(sc *Scope, name string, expr any) *PropAssign {
	assign := &PropAssign{Prop: name, Expr: expr}
	sc.Statements = append(sc.Statements, assign)
	return assign
}
