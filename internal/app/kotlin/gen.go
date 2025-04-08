package kotlin

import (
	"c2k/internal/app/curl"
	"errors"
	"fmt"
	"log"
	"path"
	"slices"
	"strings"
	"unicode"
)

func GenAst(command *curl.Command) (file *KtFile, err error) {
	file = &KtFile{}
	builderFound := false
	methodFunc := requestRequest
	request := command.Request

	for _, fqn := range builders {
		if name := simpleName(fqn); name == strings.ToLower(request.Method) {
			methodFunc = fqn
			builderFound = true
			break
		}
	}

	imports := make(map[*Fqn]struct{})
	addImport(imports, methodFunc)

	var clientCall MethodCall
	var requestBuilder *LambdaLiteral = nil

	runBlockingScope := newScope()

	var ctorArgs []any

	if !command.FollowRedirects { // Ktor follows redirects by default
		ctorArgs = append(ctorArgs, LambdaLiteral{Statements: []any{
			PropAssignment{Prop: "followRedirects", Expr: command.FollowRedirects},
		}})
	}

	client, clientDecl := declareVal(runBlockingScope, "client", CtorInvoke{Type: UserType{simpleName(httpClient)}, ValueArgs: ctorArgs})
	if builderFound {
		clientCall = callMethod(runBlockingScope, client, simpleName(methodFunc), request.Url)
	} else {
		requestBuilder = &LambdaLiteral{Statements: []any{
			PropAssignment{Prop: "method", Expr: CtorInvoke{Type: UserType{simpleName(httpMethod)}, ValueArgs: []any{request.Method}}},
		}}

		clientCall = callMethod(runBlockingScope, client, simpleName(methodFunc), request.Url)

		addImport(imports, httpMethod)
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
		if requestBuilder == nil {
			requestBuilder = &LambdaLiteral{}
		}
		var appends []any

		for _, p := range b.Params {
			appends = append(appends, FuncCall{Name: "append", ValueArgs: []any{p.Name, p.Value}})
		}

		requestBuilder.Statements = append(requestBuilder.Statements, FuncCall{Name: simpleName(setBody), ValueArgs: []any{
			CtorInvoke{Type: UserType{simpleName(formDataContent)}, ValueArgs: []any{
				FuncCall{Name: simpleName(parameters), ValueArgs: []any{
					LambdaLiteral{Statements: appends},
				}},
			}},
		}})

		addImport(imports, formDataContent)
		addImport(imports, parameters)
		addImport(imports, setBody)
	case *curl.FormDataBody:
		if requestBuilder == nil {
			requestBuilder = &LambdaLiteral{}
		}
		var fdStatements []any

		for _, p := range b.Parts {
			switch p.Kind {
			case curl.FormPartItem:
				fdStatements = append(fdStatements, FuncCall{Name: "append", ValueArgs: []any{p.Name, p.Value}})
			case curl.FormPartFile:
				fileName := path.Base(p.FilePath)
				varName := genFormVar(fileName)
				fdStatements = append(fdStatements, VarDecl{Name: varName, Assignment: CtorInvoke{Type: UserType{simpleName(fileCtor)}, ValueArgs: []any{p.FilePath}}})

				addImport(imports, fileCtor)
				addImport(imports, channelProvider)
				addImport(imports, readChannel)

				contentType := "application/octet-stream"
				if p.ContentType != "" {
					contentType = p.ContentType
				}

				chProvider := CtorInvoke{Type: UserType{simpleName(channelProvider)}, ValueArgs: []any{
					NamedArg{Name: "size", Value: MethodCall{Receiver: varName, Method: "length"}},
					InlineLambdaLiteral{Statements: []any{MethodCall{Receiver: varName, Method: "readChannel"}}},
				}}
				headers := MethodCall{Receiver: "Headers", Method: "build", ValueArgs: []any{
					LambdaLiteral{Statements: []any{
						FuncCall{Name: "append", ValueArgs: []any{
							PropAccess{Object: "HttpHeaders", Prop: "ContentType"},
							contentType,
						}},
						FuncCall{Name: "append", ValueArgs: []any{
							PropAccess{Object: "HttpHeaders", Prop: "ContentDisposition"},
							fmt.Sprintf("filename=\"${%s.name}\"", varName),
						}},
					}},
				}}

				addImport(imports, headersObject)
				addImport(imports, httpHeadersObject)
				fdStatements = append(fdStatements, FuncCall{Name: "append", ValueArgs: []any{p.Name, chProvider, headers}})
			case curl.FormPartUnknown:
				err = errors.New("form-data-body: unrecognized form part type")
				return
			}
		}

		requestBuilder.Statements = append(requestBuilder.Statements, FuncCall{Name: simpleName(setBody), ValueArgs: []any{
			CtorInvoke{Type: UserType{simpleName(multipartContent)}, ValueArgs: []any{
				FuncCall{Name: simpleName(formData), ValueArgs: []any{
					LambdaLiteral{Statements: fdStatements},
				}},
			}},
		}})

		addImport(imports, multipartContent)
		addImport(imports, formData)
		addImport(imports, setBody)
	}

	if requestBuilder != nil {
		clientCall.ValueArgs = append(clientCall.ValueArgs, *requestBuilder)
	}

	file.TopLevels = append(file.TopLevels, FuncDecl{
		Name: "main",
		Expr: FuncCall{Name: simpleName(runBlocking), ValueArgs: []any{
			LambdaLiteral{Statements: []any{
				clientDecl,
				VarDecl{
					Name:       "response",
					Assignment: clientCall,
				},
				FuncCall{
					Name: "print",
					ValueArgs: []any{
						MethodCall{Receiver: "response", Method: simpleName(bodyAsText)},
					},
				},
			}},
		}},
	})

	addImport(imports, runBlocking)
	addImport(imports, httpClient)
	addImport(imports, bodyAsText)

	for fqn := range imports {
		file.ImportList = append(file.ImportList, Import{fqn: *fqn})
	}

	slices.SortFunc(file.ImportList, func(a, b Import) int {
		for i, p := range a.fqn {
			if ord := strings.Compare(p, b.fqn[i]); ord != 0 {
				return ord
			}
		}

		return len(a.fqn) - len(b.fqn)
	})

	return
}

func addImport(imports map[*Fqn]struct{}, fqn *Fqn) {
	imports[fqn] = struct{}{}
}

func simpleName(fqn *Fqn) (name string) {
	if len(*fqn) != 0 {
		name = (*fqn)[len(*fqn)-1]
	} else {
		log.Panicf("cannot get name from empty fqn")
	}

	return
}

func genFormVar(filename string) string {
	name := strings.TrimSuffix(filename, path.Ext(filename))

	if name == "" {
		return "file"
	} else {
		var numPrefix strings.Builder
		for _, c := range name {
			if unicode.IsNumber(c) {
				numPrefix.WriteRune(c)
			} else {
				break
			}
		}

		if numPrefix.Len() > 0 {
			return "file" + numPrefix.String()
		}

		return name
	}
}
