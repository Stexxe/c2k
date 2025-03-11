package kotlin

import (
	"c2k/internal/app/curl"
	"slices"
	"strings"
)

var clientPackage = Fqn{"io", "ktor", "client"}
var clientRequestPackage = Fqn{"io", "ktor", "client", "request"}
var clientStatementPackage = Fqn{"io", "ktor", "client", "statement"}
var httpPackage = Fqn{"io", "ktor", "http"}
var coroutinesPackage = Fqn{"kotlinx", "coroutines"}

var funcMethods = map[SimpleId]struct{}{"get": {}, "post": {}, "patch": {}, "head": {}, "options": {}, "delete": {}, "put": {}}

func GenAst(request *curl.Request) (file KtFile, err error) {
	methodFunc := SimpleId(strings.ToLower(request.Method))
	customMethod := false
	if _, ok := funcMethods[SimpleId(strings.ToLower(request.Method))]; !ok {
		customMethod = true
		methodFunc = "request"
	}

	addImportFor(&file, methodFunc, clientRequestPackage)

	var clientCall CallExpr
	var requestBuilder *LambdaLiteral = nil

	if !customMethod {
		clientCall = CallExpr{Receiver: "client", Method: methodFunc, ValueArgs: []any{
			StringLiteral(request.Url),
		}}
	} else {
		requestBuilder = &LambdaLiteral{Statements: []any{
			PropAssignment{Prop: "method", Expr: CtorInvoke{Type: UserType{"HttpMethod"}, ValueArgs: []any{StringLiteral(request.Method)}}},
		}}

		clientCall = CallExpr{Receiver: "client", Method: methodFunc, ValueArgs: []any{
			StringLiteral(request.Url),
			//LambdaLiteral{Statements: []any{
			//	PropAssignment{Prop: "method", Expr: CtorInvoke{Type: UserType{"HttpMethod"}, ValueArgs: []any{StringLiteral(request.Method)}}},
			//}},
		}}

		addImportFor(&file, "HttpMethod", httpPackage)
	}

	if len(request.Headers) > 0 {
		if requestBuilder == nil {
			requestBuilder = &LambdaLiteral{}
		}

		for _, h := range request.Headers {
			requestBuilder.Statements = append(requestBuilder.Statements, CallExpr{Receiver: "headers", Method: "append", ValueArgs: []any{
				StringLiteral(h.Name), StringLiteral(h.Value),
			}})
		}
	}

	if requestBuilder != nil {
		clientCall.ValueArgs = append(clientCall.ValueArgs, *requestBuilder)
	}

	file.TopLevels = append(file.TopLevels, FuncDecl{
		Name: "main",
		Expr: CallExpr{Method: "runBlocking", ValueArgs: []any{
			LambdaLiteral{Statements: []any{
				VarDecl{
					Name:       "client",
					Assignment: CtorInvoke{Type: UserType{"HttpClient"}},
				},
				VarDecl{
					Name:       "response",
					Assignment: clientCall,
				},
				CallExpr{
					Method: "print",
					ValueArgs: []any{
						CallExpr{Receiver: "response", Method: "bodyAsText"},
					},
				},
			}},
		}},
	})

	addImportFor(&file, "runBlocking", coroutinesPackage)
	addImportFor(&file, "HttpClient", clientPackage)
	addImportFor(&file, "bodyAsText", clientStatementPackage)

	slices.SortFunc(file.ImportList, func(a, b Import) int {
		for i, p := range a.fqn {
			if ord := strings.Compare(string(p), string(b.fqn[i])); ord != 0 {
				return ord
			}
		}

		return len(a.fqn) - len(b.fqn)
	})

	return
}

func addImportFor(file *KtFile, symbol SimpleId, pack Fqn) {
	fqn := append([]SimpleId{}, pack...)
	fqn = append(fqn, symbol)
	file.ImportList = append(file.ImportList, Import{fqn: fqn})
}
