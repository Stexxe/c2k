package curl

import (
	"fmt"
	"strings"
)

type Command struct {
	FollowRedirects      bool
	ResolvedAddr         string
	PrintResponseHeaders bool
	Request              *Request
}

type Request struct {
	Url     string
	Method  string
	Headers []Header
	Body    any
}

type Header struct {
	Name, Value string
}

type curlOption int

const (
	UnknownOption curlOption = iota
	HeaderOption
	MethodOption
	DataOption
	DataUrlEncodeOption
	RawDataOption
	LocationOption
	FormOption
	ResolveOption
	IncludeOption
)

var oneArgOptions = map[string]curlOption{
	"-H": HeaderOption, "--header": HeaderOption,
	"-X": MethodOption, "--request": MethodOption,
	"-d": DataOption, "--data": DataOption,
	"--data-raw":       RawDataOption,
	"--data-urlencode": DataUrlEncodeOption,
	"-F":               FormOption, "--form": FormOption,
	"--resolve": ResolveOption,
}

var flagOptions = map[string]curlOption{
	"-L": LocationOption, "--location": LocationOption,
	"-i": IncludeOption, "--include": IncludeOption,
}

type curlOptionInstance struct {
	option curlOption
	value  []string
}

type UrlEncodedBody struct {
	Params []FormParam
}

type FormDataBody struct {
	Parts []FormPart
}

type FormParam struct {
	Name, Value string
	FilePath    string
}

type FormPartKind int

const (
	FormPartUnknown FormPartKind = iota
	FormPartItem
	FormPartFile
)

type FormPart struct {
	Kind        FormPartKind
	Name        string
	Value       string
	FilePath    string
	ContentType string
}

func ParseCommand(cmd []string) (command *Command, err error) {
	var options []curlOptionInstance
	request := &Request{}
	command = &Command{Request: request}
	var freeStandingArgs []string
	var unexpectedOptions []string

	if len(cmd) > 0 && cmd[0] == "curl" {
		if len(cmd) > 1 {
			args := cmd[1:]
			for i := 0; i < len(args); {
				arg := args[i]

				if strings.HasPrefix(arg, "-") {
					if strings.HasPrefix(arg, "--") {
						if opt, ok := oneArgOptions[arg]; ok {
							options = append(options, curlOptionInstance{option: opt, value: args[i+1 : i+2]})
							i += 2
						} else if opt, ok := flagOptions[arg]; ok {
							options = append(options, curlOptionInstance{option: opt})
							i += 1
						} else {
							unexpectedOptions = append(unexpectedOptions, arg)
							i++
						}
					} else {
						var opt curlOption
						var ok bool
						argBytes := []byte(arg)
						end := len(argBytes)
						isFlag := false

						for ; end > 0; end-- {
							if opt, ok = oneArgOptions[string(argBytes[0:end])]; ok {
								break
							}

							if opt, ok = flagOptions[string(argBytes[0:end])]; ok {
								isFlag = true
								break
							}
						}

						if isFlag {
							options = append(options, curlOptionInstance{option: opt})
							i += 1
						} else if opt != UnknownOption {
							if end == len(argBytes) {
								options = append(options, curlOptionInstance{option: opt, value: args[i+1 : i+2]})
								i += 2
							} else {
								options = append(options, curlOptionInstance{option: opt, value: []string{arg[end:]}})
								i += 1
							}
						} else {
							unexpectedOptions = append(unexpectedOptions, arg)
							i++
						}
					}
				} else {
					freeStandingArgs = append(freeStandingArgs, arg)
					i += 1
				}
			}

			if len(unexpectedOptions) > 0 {
				err = fmt.Errorf("parse curl: unexpected options: %s", strings.Join(unexpectedOptions, " "))
				return
			}

			if len(freeStandingArgs) == 0 {
				err = fmt.Errorf("parse curl: expected URL")
			} else if len(freeStandingArgs) == 1 {
				request.Url = freeStandingArgs[0]
			} else {
				err = fmt.Errorf("parse curl: unexpected arguments: %v", freeStandingArgs)
			}

			var userContentType string
			var contentTypeDefined bool
			for _, inst := range options {
				if inst.option == HeaderOption {
					header := parseHeader(inst.value[0])

					if strings.ToLower(header.Name) == "content-type" {
						userContentType = header.Value
						contentTypeDefined = true
					}

					request.Headers = append(request.Headers, header)
				}
			}

			var formBody *FormDataBody
			var urlEncodedBody *UrlEncodedBody
			for _, inst := range options {
				switch inst.option {
				case HeaderOption:
					// Already parsed
				case MethodOption:
					request.Method = inst.value[0]
				case DataOption, DataUrlEncodeOption:
					if inst.option != DataUrlEncodeOption && contentTypeDefined && !strings.HasPrefix(userContentType, "application/x-www-form-urlencoded") {
						request.Body = inst.value[0]
					} else {
						if urlEncodedBody == nil {
							urlEncodedBody = &UrlEncodedBody{}
							request.Body = urlEncodedBody
						}

						urlEncodedBody.Params = append(urlEncodedBody.Params, parseData(inst.value[0])...)
					}

					if request.Method == "" {
						request.Method = "POST"
					}
				case RawDataOption:
					if contentTypeDefined && !strings.HasPrefix(userContentType, "application/x-www-form-urlencoded") {
						request.Body = inst.value[0]
					} else {
						if urlEncodedBody == nil {
							urlEncodedBody = &UrlEncodedBody{}
							request.Body = urlEncodedBody
						}
						urlEncodedBody.Params = append(urlEncodedBody.Params, parseRawData(inst.value[0])...)
					}

					if request.Method == "" {
						request.Method = "POST"
					}
				case LocationOption:
					command.FollowRedirects = true
				case FormOption:
					if formBody == nil {
						formBody = &FormDataBody{}
						request.Body = formBody
					}

					formBody.Parts = append(formBody.Parts, parseFormPart(inst.value[0]))

					if request.Method == "" {
						request.Method = "POST"
					}
				case ResolveOption:
					parts := strings.Split(inst.value[0], ":")
					host := parts[0]
					port := parts[1]
					ip := parts[2]

					plainHostMatches := strings.HasPrefix(request.Url, "http://") && port == "80" &&
						(strings.TrimPrefix(request.Url, "http://") == host || host == "*")

					secureHostMatches := strings.HasPrefix(request.Url, "https://") && port == "443" &&
						(strings.TrimPrefix(request.Url, "https://") == host || host == "*")

					if plainHostMatches || secureHostMatches {
						command.ResolvedAddr = ip
					}
				case IncludeOption:
					command.PrintResponseHeaders = true
				case UnknownOption:
					err = fmt.Errorf("curl: unknown option")
				}
			}

			if request.Method == "" {
				request.Method = "GET"
			}
		} else {
			err = fmt.Errorf("curl: expected URL, got none")
		}
	} else {
		err = fmt.Errorf("curl: invalid curl command '%s'", strings.Join(cmd, " "))
	}

	return
}

func parseHeader(header string) (h Header) {
	parts := strings.Split(header, ":")

	if len(parts) == 2 {
		h.Name = strings.TrimSpace(parts[0])
		h.Value = strings.TrimSpace(parts[1])
	}

	return
}

func isQuote(r rune) bool {
	return r == '\'' || r == '"' || r == '”' || r == '“' || r == '‘' || r == '’'
}

func parseFormPart(str string) (param FormPart) {
	mainParts := strings.Split(str, ";")

	var nameValue string
	var typeInfo string

	if len(mainParts) == 2 {
		nameValue = mainParts[0]
		typeInfo = mainParts[1]
	} else if len(mainParts) == 1 {
		nameValue = mainParts[0]
	}

	if nameValue != "" {
		parts := strings.Split(nameValue, "=")
		if len(parts) == 2 {
			param.Name = parts[0]

			if strings.HasPrefix(parts[1], "@") {
				param.Kind = FormPartFile
				param.FilePath = strings.TrimPrefix(parts[1], "@")
				param.FilePath = strings.TrimLeftFunc(param.FilePath, isQuote)
				param.FilePath = strings.TrimRightFunc(param.FilePath, isQuote)
			} else {
				param.Kind = FormPartItem
				param.Value = parts[1]
			}
		} else {
			param.Name = parts[0]
		}
	}

	if typeInfo != "" {
		parts := strings.Split(typeInfo, "=")

		if len(parts) == 2 && strings.TrimSpace(parts[0]) == "type" {
			param.ContentType = strings.TrimSpace(parts[1])
		}
	}

	return
}

func parseData(str string) (params []FormParam) {
	if strings.Contains(str, "@") {
		parts := strings.Split(str, "=")

		param := FormParam{}
		if len(parts) == 2 {
			param.Name = parts[0]
			param.FilePath = strings.TrimPrefix(parts[1], "@")
		} else {
			param.FilePath = strings.TrimPrefix(str, "@")
		}

		params = append(params, param)
	} else {
		params = parseRawData(str)
	}

	return
}

func parseRawData(str string) (params []FormParam) {
	for _, kv := range strings.Split(str, "&") {
		param := FormParam{}
		parts := strings.Split(kv, "=")

		if len(parts) == 2 {
			if parts[0] == "" {
				param.Name = parts[1]
			} else {
				param.Name = parts[0]
				param.Value = parts[1]
			}
		} else {
			param.Name = parts[0]
		}

		params = append(params, param)
	}

	return
}
