package kotlin

import (
	"c2k/internal/app/curl"
	"c2k/internal/app/utils"
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

	if !builderFound {
		requestBuilder = &LambdaLiteral{Statements: []any{
			PropAssignment{Prop: "method", Expr: CtorInvoke{Type: UserType{simpleName(httpMethod)}, ValueArgs: []any{request.Method}}},
		}}

		addImport(imports, httpMethod)
	}

	requestAddr := request.Url
	if command.ResolvedAddr != "" {
		requestAddr = fmt.Sprintf("%s://%s", utils.UrlProto(requestAddr), command.ResolvedAddr)
	}

	clientCall = callMethod(runBlockingScope, client, simpleName(methodFunc), requestAddr)

	requestScope := newScope()

	if command.ResolvedAddr != "" {
		if requestBuilder == nil {
			requestBuilder = &LambdaLiteral{}
		}

		requestBuilder.Statements = append(
			requestBuilder.Statements,
			callPropMethod("headers", "append", "Host", utils.UrlHost(request.Url)),
		)
	}

	if len(request.Headers) > 0 {
		if requestBuilder == nil {
			requestBuilder = &LambdaLiteral{}
		}

		for _, h := range request.Headers {
			requestBuilder.Statements = append(
				requestBuilder.Statements,
				callPropMethod("headers", "append", h.Name, h.Value),
			)
		}

		if request.Body != nil {
			requestBuilder.Statements = append(requestBuilder.Statements, EmptyStatement{})
		}
	}

	if request.Body != nil {
		addImport(imports, setBody)
	}

	switch b := request.Body.(type) {
	case string:
		if requestBuilder == nil {
			requestBuilder = &LambdaLiteral{}
		}

		requestBuilder.Statements = append(requestBuilder.Statements, FuncCall{Name: "setBody", ValueArgs: []any{b}})
	case *curl.UrlEncodedBody:
		if requestBuilder == nil {
			requestBuilder = &LambdaLiteral{}
		}

		var statements []any
		for _, p := range b.Params {
			if p.FilePath == "" {
				statements = append(statements, FuncCall{Name: "append", ValueArgs: []any{p.Name, p.Value}})
			} else if p.Name != "" {
				addImport(imports, readText)
				addImport(imports, fileCtor)

				statements = append(
					statements,
					FuncCall{Name: "append", ValueArgs: []any{p.Name, MethodCall{Receiver: callCtor(fileCtor, p.FilePath), Method: "readText"}}},
				)
			} else {
				varName := genFormVar(path.Base(p.FilePath)) + "Params"
				addImport(imports, fileCtor)
				addImport(imports, readText)
				addImport(imports, parseUrlEncodedParameters)

				params, decl := declareVal(
					requestScope,
					varName,
					MethodCall{Receiver: MethodCall{Receiver: callCtor(fileCtor, p.FilePath), Method: simpleName(readText)}, Method: simpleName(parseUrlEncodedParameters)},
				)
				statements = append(statements, decl)

				loop := ForInLoop{Bind: PairDestruct{"name", "values"}, Expr: callMethod(requestScope, params, "entries"), Statements: []any{
					ForInLoop{Bind: "v", Expr: Id("values"), Statements: []any{
						FuncCall{Name: "append", ValueArgs: []any{Id("name"), Id("v")}},
					}},
				}}

				statements = append(statements, loop)
			}
		}

		requestBuilder.Statements = append(requestBuilder.Statements, FuncCall{Name: simpleName(setBody), ValueArgs: []any{
			CtorInvoke{Type: UserType{simpleName(formDataContent)}, ValueArgs: []any{
				FuncCall{Name: simpleName(parameters), ValueArgs: []any{
					LambdaLiteral{Statements: statements},
				}},
			}},
		}})

		addImport(imports, formDataContent)
		addImport(imports, parameters)
	case *curl.FormDataBody:
		if requestBuilder == nil {
			requestBuilder = &LambdaLiteral{}
		}
		var fdStatements []any

		for i, p := range b.Parts {
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

				if i != len(b.Parts)-1 {
					fdStatements = append(fdStatements, EmptyStatement{})
				}
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
	}

	if requestBuilder != nil {
		clientCall.ValueArgs = append(clientCall.ValueArgs, *requestBuilder)
	}

	var mainStatements []any

	mainStatements = append(mainStatements, clientDecl)
	mainStatements = append(mainStatements, VarDecl{
		Name:       "response",
		Assignment: clientCall,
	})

	var outputStatements []any

	if command.PrintResponseHeaders {
		outputStatements = append(outputStatements, EmptyStatement{})
		outputStatements = append(
			outputStatements,
			FuncCall{Name: "println", ValueArgs: []any{fmt.Sprintf("%s %s", interp("response.version"), interp("response.status.value"))}},
		)

		outputStatements = append(outputStatements, EmptyStatement{})
		headersLoop := ForInLoop{Bind: PairDestruct{"name", "values"}, Expr: MethodCall{Receiver: Id("response.headers"), Method: "entries"}, Statements: []any{
			ForInLoop{Bind: "v", Expr: Id("values"), Statements: []any{
				FuncCall{Name: "println", ValueArgs: []any{fmt.Sprintf("%s: %s", interp("name"), interp("v"))}},
			}},
		}}

		outputStatements = append(outputStatements, headersLoop)

		outputStatements = append(outputStatements, EmptyStatement{})
		outputStatements = append(outputStatements, FuncCall{Name: "println"})

		outputStatements = append(outputStatements, EmptyStatement{})
	}

	outputStatements = append(outputStatements, FuncCall{
		Name: "print",
		ValueArgs: []any{
			MethodCall{Receiver: "response", Method: simpleName(bodyAsText)},
		},
	})

	mainStatements = append(mainStatements, outputStatements...)

	file.TopLevels = append(file.TopLevels, FuncDecl{
		Name: "main",
		Expr: FuncCall{Name: simpleName(runBlocking), ValueArgs: []any{
			LambdaLiteral{Statements: mainStatements},
		}},
	})

	addImport(imports, runBlocking)
	addImport(imports, httpClient)
	addImport(imports, bodyAsText)

	for fqn := range imports {
		autoImported := false
		for _, pack := range autoImportedPackages {
			if hasPackage(pack, fqn) {
				autoImported = true
				break
			}
		}

		if !autoImported {
			file.ImportList = append(file.ImportList, Import{fqn: *fqn})
		}
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

func interp(s string) string {
	if strings.Contains(s, "(") || strings.Contains(s, ".") {
		return fmt.Sprintf("${%s}", s)
	}

	return "$" + s
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

func hasPackage(pack *Fqn, fqn *Fqn) bool {
	if len(*pack) > len(*fqn) {
		return false
	}

	for i, part := range *pack {
		if part != (*fqn)[i] {
			return false
		}
	}

	return true
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
