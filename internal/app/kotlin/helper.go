package kotlin

//type Scope struct {
//	FuncDecls  []*funcDecl
//	Statements []any
//	Imports    map[*Fqn]struct{}
//	TopScope   *Scope
//}
//
//func newScope(topScope *Scope) *Scope {
//	sc := &Scope{}
//
//	if topScope != nil {
//		sc.TopScope = topScope
//	} else {
//		sc.TopScope = sc
//	}
//
//	sc.Imports = make(map[*Fqn]struct{})
//	return sc
//}
//
//func addImport(sc *Scope, fqn *Fqn) {
//	sc.TopScope.Imports[fqn] = struct{}{}
//}
//
//func declareFunc(sc *Scope, signature *funcSignature, builder func(sc *Scope)) {
//	childScope := newScope(sc.TopScope)
//	builder(childScope)
//	decl := &funcDecl{signature: signature, bodyScope: childScope}
//	sc.FuncDecls = append(sc.FuncDecls, decl)
//}
//
//func declareVal(sc *Scope, name string, valType *ktType, expr any) {
//	sc.Statements = append(sc.Statements, varDecl{name: name, varType: valType, expr: expr})
//}
//
//func ret(sc *Scope, expr any) {
//	sc.Statements = append(sc.Statements, &returnStat{expr})
//}
//
//func callCtor(sc *Scope, class *ktClass) ctorCall {
//	addImport(sc, class.fqn)
//	return ctorCall{class: class}
//}
//
//func symbolName(fqn *Fqn) string {
//	if len(*fqn) == 0 {
//		panic("empty FQN")
//	}
//
//	return (*fqn)[len(*fqn)-1]
//}
