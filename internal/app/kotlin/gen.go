package kotlin

import "c2k/internal/app/curl"

const (
	httpClient  SimpleId = "HttpClient"
	requestGet  SimpleId = "get"
	bodyAsText  SimpleId = "bodyAsText"
	runBlocking SimpleId = "runBlocking"
)

var clientSymbolsMap = map[SimpleId]Fqn{
	httpClient: {},
	requestGet: {"request"},
	bodyAsText: {"statement"},
}

var coroutinesSymbolsMap = map[SimpleId]Fqn{
	runBlocking: {},
}

func GenAst(request *curl.Request) (file KtFile, err error) {
	file.ImportList = append(file.ImportList, Import{Id: clientPackageFor(httpClient)})
	file.ImportList = append(file.ImportList, Import{Id: clientPackageFor(requestGet)})
	file.ImportList = append(file.ImportList, Import{Id: clientPackageFor(bodyAsText)})
	file.ImportList = append(file.ImportList, Import{Id: coroutinesPackageFor(runBlocking)})

	file.TopLevels = append(file.TopLevels, FuncDecl{
		Name: "main",
		Expr: CallExpr{Method: "runBlocking", ValueArgs: []any{
			LambdaLiteral{Statements: []any{
				VarDecl{
					Name:       "client",
					Assignment: CtorInvoke{Type: UserType{"HttpClient"}},
				},
				VarDecl{
					Name: "response",
					Assignment: CallExpr{Receiver: "client", Method: "get", ValueArgs: []any{
						StringLiteral(request.Url),
					}},
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

	return
}

func coroutinesPackageFor(symbol SimpleId) (result Fqn) {
	// TODO: Introduce asserts which only work in the debug mode (see go:generate)
	if fqn, ok := coroutinesSymbolsMap[symbol]; ok {
		clientPackage := []SimpleId{"kotlinx", "coroutines"}
		result = append(clientPackage, fqn...)
		result = append(result, symbol)
	} else {
		panic("coroutines: package for " + symbol + " not found")
	}

	return
}

func clientPackageFor(symbol SimpleId) (result Fqn) {
	// TODO: Introduce asserts which only work in the debug mode (see go:generate)
	if fqn, ok := clientSymbolsMap[symbol]; ok {
		clientPackage := []SimpleId{"io", "ktor", "client"}
		result = append(clientPackage, fqn...)
		result = append(result, symbol)
	} else {
		panic("client: package for " + symbol + " not found")
	}

	return
}
