package kotlin

import (
	"c2k/internal/app/curl"
	"slices"
	"strings"
)

func GenAst(command *curl.Command) (file *KtFile, err error) {
	file = &KtFile{}
	builderFound := false
	methodFunc := requestRequest
	request := command.Request

	for _, sym := range builders {
		if sym.Name == SimpleId(strings.ToLower(request.Method)) {
			methodFunc = sym
			builderFound = true
			break
		}
	}

	useSymbol(file, methodFunc)

	var clientCall MethodCall
	var requestBuilder *LambdaLiteral = nil

	if builderFound {
		clientCall = MethodCall{Receiver: "client", Method: methodFunc.Name, ValueArgs: []any{
			StringLiteral(request.Url),
		}}
	} else {
		requestBuilder = &LambdaLiteral{Statements: []any{
			PropAssignment{Prop: "method", Expr: CtorInvoke{Type: UserType{httpMethod.Name}, ValueArgs: []any{StringLiteral(request.Method)}}},
		}}

		clientCall = MethodCall{Receiver: "client", Method: methodFunc.Name, ValueArgs: []any{
			StringLiteral(request.Url),
		}}

		useSymbol(file, httpMethod)
	}

	if len(request.Headers) > 0 {
		if requestBuilder == nil {
			requestBuilder = &LambdaLiteral{}
		}

		for _, h := range request.Headers {
			requestBuilder.Statements = append(requestBuilder.Statements, MethodCall{Receiver: "headers", Method: "append", ValueArgs: []any{
				StringLiteral(h.Name), StringLiteral(h.Value),
			}})
		}
	}

	switch b := request.Body.(type) {
	case curl.UrlEncodedBody:
		requestBuilder = &LambdaLiteral{}
		var appends []any

		for _, p := range b.Params {
			appends = append(appends, FuncCall{Name: "append", ValueArgs: []any{StringLiteral(p.Name), StringLiteral(p.Value)}})
		}

		requestBuilder.Statements = append(requestBuilder.Statements, FuncCall{Name: setBody.Name, ValueArgs: []any{
			CtorInvoke{Type: UserType{formDataContent.Name}, ValueArgs: []any{
				FuncCall{Name: parameters.Name, ValueArgs: []any{
					LambdaLiteral{Statements: appends},
				}},
			}},
		}})

		useSymbol(file, formDataContent)
		useSymbol(file, parameters)
		useSymbol(file, setBody)
	case *curl.FormDataBody:
		requestBuilder = &LambdaLiteral{}
		var appends []any

		for _, p := range b.Parts {
			appends = append(appends, FuncCall{Name: "append", ValueArgs: []any{StringLiteral(p.Name), StringLiteral(p.Value)}})
		}

		requestBuilder.Statements = append(requestBuilder.Statements, FuncCall{Name: setBody.Name, ValueArgs: []any{
			CtorInvoke{Type: UserType{multipartContent.Name}, ValueArgs: []any{
				FuncCall{Name: formData.Name, ValueArgs: []any{
					LambdaLiteral{Statements: appends},
				}},
			}},
		}})

		useSymbol(file, multipartContent)
		useSymbol(file, formData)
		useSymbol(file, setBody)
	}

	if requestBuilder != nil {
		clientCall.ValueArgs = append(clientCall.ValueArgs, *requestBuilder)
	}

	var ctorArgs []any

	if !command.FollowRedirects { // Ktor follows redirects by default
		ctorArgs = append(ctorArgs, LambdaLiteral{Statements: []any{
			PropAssignment{Prop: "followRedirects", Expr: BoolLiteral(command.FollowRedirects)},
		}})
	}

	file.TopLevels = append(file.TopLevels, FuncDecl{
		Name: "main",
		Expr: FuncCall{Name: runBlocking.Name, ValueArgs: []any{
			LambdaLiteral{Statements: []any{
				VarDecl{
					Name:       "client",
					Assignment: CtorInvoke{Type: UserType{httpClient.Name}, ValueArgs: ctorArgs},
				},
				VarDecl{
					Name:       "response",
					Assignment: clientCall,
				},
				FuncCall{
					Name: "print",
					ValueArgs: []any{
						MethodCall{Receiver: "response", Method: bodyAsText.Name},
					},
				},
			}},
		}},
	})

	useSymbol(file, runBlocking)
	useSymbol(file, httpClient)
	useSymbol(file, bodyAsText)

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

func useSymbol(file *KtFile, symbol *symbol) {
	fqn := append(Fqn{}, *symbol.Package...)
	fqn = append(fqn, symbol.Name)
	file.ImportList = append(file.ImportList, Import{fqn: fqn})
}
