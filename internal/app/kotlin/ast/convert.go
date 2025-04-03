package ast

import (
	"c2k/internal/app/kotlin"
	"c2k/internal/app/kotlin/ir"
	"log"
	"slices"
	"strings"
)

func GenAst(fileScope *ir.Scope) (file *KtFile, err error) {
	file = &KtFile{}

	for fqn := range fileScope.TopLevel.Imports {
		file.ImportList = append(file.ImportList, KtImport{Fqn: fqn})
	}

	slices.SortFunc(file.ImportList, func(a, b KtImport) int {
		for i, p := range *a.Fqn {
			if ord := strings.Compare(p, (*b.Fqn)[i]); ord != 0 {
				return ord
			}
		}

		return len(*a.Fqn) - len(*b.Fqn)
	})

	for _, decl := range fileScope.FuncDecls {
		bodyScope := decl.BodyScope

		astDecl := KtFuncDecl{
			Name: symbolName(decl.Signature.Fqn),
		}

		if len(bodyScope.Statements) == 1 {
			if r, ok := bodyScope.Statements[0].(*ir.ReturnStat); ok {
				astDecl.Expr = convertExpr(bodyScope, r.Expr)
			}
		}

		file.TopLevels = append(file.TopLevels, astDecl)
	}

	return
}

func convertExpr(sc *ir.Scope, in any) (out any) {
	switch expr := in.(type) {
	case *ir.FuncCall:
		call := KtInvocation{Method: symbolName(expr.Signature.Fqn)}
		call.ValueArgs = convertArgs(sc, expr.Args)
		out = call
	case *ir.Block:
		lit := KtLambdaLiteral{}
		var statements []any

		for _, st := range expr.Scope.Statements {
			statements = append(statements, convertStatement(expr.Scope, st))
		}

		lit.Statements = statements
		out = lit
	case *ir.CtorCall:
		call := KtInvocation{Method: symbolName(expr.Class.Fqn)}
		call.ValueArgs = convertArgs(sc, expr.Args)

		out = call
	case *ir.MethodCall:
		call := KtInvocation{Receiver: sc.Variables[expr.Object], Method: symbolName(expr.Method)}
		call.ValueArgs = convertArgs(sc, expr.Args)
		out = call
	case string:
		out = KtStringLiteral(expr)
	case bool:
		out = KtBoolLiteral(expr)
	default:
		log.Panicf("unexpected type %q", expr)
	}

	return
}

func convertArgs(sc *ir.Scope, args []any) (out []any) {
	for _, a := range args {
		out = append(out, convertExpr(sc, a))
	}

	return
}

func convertStatement(sc *ir.Scope, in any) (out any) {
	switch stmt := in.(type) {
	case *ir.VarDecl:
		out = KtVarDecl{Name: stmt.Name, Assignment: convertExpr(sc, stmt.Expr)}
	case *ir.FuncCall:
		out = convertExpr(sc, in)
	case *ir.PropAssign:
		out = KtPropAssignment{Prop: stmt.Prop, Expr: convertExpr(sc, stmt.Expr)}
	default:
		log.Panicf("invalid statement %q", in)
	}

	return
}

func symbolName(fqn *kotlin.Fqn) string {
	if len(*fqn) == 0 {
		panic("empty FQN")
	}

	return (*fqn)[len(*fqn)-1]
}
