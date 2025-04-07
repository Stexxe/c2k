package kotlin

import (
	"c2k/internal/app/curl"
	"errors"
	"fmt"
	"path"
	"slices"
	"strings"
)

func GenAst(command *curl.Command) (file *KtFile, err error) {
	file = &KtFile{}
	builderFound := false
	methodFunc := requestRequest
	request := command.Request

	for _, sym := range builders {
		if sym.Name == strings.ToLower(request.Method) {
			methodFunc = sym
			builderFound = true
			break
		}
	}

	useSymbol(file, methodFunc)

	var clientCall MethodCall
	var requestBuilder *LambdaLiteral = nil

	runBlockingScope := newScope()

	var ctorArgs []any

	if !command.FollowRedirects { // Ktor follows redirects by default
		ctorArgs = append(ctorArgs, LambdaLiteral{Statements: []any{
			PropAssignment{Prop: "followRedirects", Expr: command.FollowRedirects},
		}})
	}

	client, clientDecl := declareVal(runBlockingScope, "client", CtorInvoke{Type: UserType{httpClient.Name}, ValueArgs: ctorArgs})
	if builderFound {
		clientCall = callMethod(runBlockingScope, client, methodFunc.Name, request.Url)
	} else {
		requestBuilder = &LambdaLiteral{Statements: []any{
			PropAssignment{Prop: "method", Expr: CtorInvoke{Type: UserType{httpMethod.Name}, ValueArgs: []any{request.Method}}},
		}}

		clientCall = callMethod(runBlockingScope, client, methodFunc.Name, request.Url)

		useSymbol(file, httpMethod)
	}

	requestScope := newScope()
	if len(request.Headers) > 0 {
		if requestBuilder == nil {
			requestBuilder = &LambdaLiteral{}
		}

		for _, h := range request.Headers {
			requestBuilder.Statements = append(
				requestBuilder.Statements,
				callPropMethod(requestScope, "headers", "append", h.Name, h.Value),
			)
		}
	}

	switch b := request.Body.(type) {
	case curl.UrlEncodedBody:
		requestBuilder = &LambdaLiteral{}
		var appends []any

		for _, p := range b.Params {
			appends = append(appends, FuncCall{Name: "append", ValueArgs: []any{p.Name, p.Value}})
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
		var fdStatements []any

		for _, p := range b.Parts {
			switch p.Kind {
			case curl.FormPartItem:
				fdStatements = append(fdStatements, FuncCall{Name: "append", ValueArgs: []any{p.Name, p.Value}})
			case curl.FormPartFile:
				fdStatements = append(fdStatements, VarDecl{Name: "file", Assignment: CtorInvoke{Type: UserType{fileCtor.Name}, ValueArgs: []any{p.FilePath}}})
				useSymbol(file, fileCtor)
				useSymbol(file, channelProvider)
				useSymbol(file, readChannel)
				chProvider := CtorInvoke{Type: UserType{channelProvider.Name}, ValueArgs: []any{
					NamedArg{Name: "size", Value: MethodCall{Receiver: "file", Method: "length"}},
					InlineLambdaLiteral{Statements: []any{MethodCall{Receiver: "file", Method: "readChannel"}}},
				}}
				headers := MethodCall{Receiver: "Headers", Method: "build", ValueArgs: []any{
					LambdaLiteral{Statements: []any{
						FuncCall{Name: "append", ValueArgs: []any{
							PropAccess{Object: "HttpHeaders", Prop: "ContentType"},
							"application/octet-stream",
						}},
						FuncCall{Name: "append", ValueArgs: []any{
							PropAccess{Object: "HttpHeaders", Prop: "ContentDisposition"},
							fmt.Sprintf("filename=\"%s\"", path.Base(p.FilePath)),
						}},
					}},
				}}
				useSymbol(file, headersObject)
				useSymbol(file, httpHeadersObject)
				fdStatements = append(fdStatements, FuncCall{Name: "append", ValueArgs: []any{p.Name, chProvider, headers}})
			case curl.FormPartUnknown:
				err = errors.New("unrecognized form part type")
				return
			}

		}

		requestBuilder.Statements = append(requestBuilder.Statements, FuncCall{Name: setBody.Name, ValueArgs: []any{
			CtorInvoke{Type: UserType{multipartContent.Name}, ValueArgs: []any{
				FuncCall{Name: formData.Name, ValueArgs: []any{
					LambdaLiteral{Statements: fdStatements},
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

	file.TopLevels = append(file.TopLevels, FuncDecl{
		Name: "main",
		Expr: FuncCall{Name: runBlocking.Name, ValueArgs: []any{
			LambdaLiteral{Statements: []any{
				clientDecl,
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
