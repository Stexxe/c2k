package kotlin

import "log"

type Scope struct {
	Locals map[*Object]string
}

func newScope() *Scope {
	scope := &Scope{}
	scope.Locals = make(map[*Object]string)
	return scope
}

func newLocal(sc *Scope, name string) *Object {
	obj := &Object{}
	sc.Locals[obj] = name
	return obj
}

type Object struct {
}

func declareVal(sc *Scope, varName string, expr any) (obj *Object, decl VarDecl) {
	obj = newLocal(sc, varName)
	decl = VarDecl{
		Name:       varName,
		Assignment: expr,
	}
	return
}

func callMethod(sc *Scope, obj *Object, method string, args ...any) (call MethodCall) {
	if name, ok := sc.Locals[obj]; ok {
		call = MethodCall{Receiver: name, Method: method, ValueArgs: args}
	} else {
		log.Panicf("Local %v not found in scope for the method call %s", obj, method)
	}

	return
}

func callPropMethod(sc *Scope, prop, method string, args ...any) MethodCall {
	return MethodCall{Receiver: prop, Method: method, ValueArgs: args}
}
