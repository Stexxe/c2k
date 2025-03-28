package curl

import (
	"fmt"
	"strings"
)

type Command struct {
	FollowRedirects bool
	Request         *Request
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
	LocationOption
)

var oneArgOptions = map[string]curlOption{
	"-H": HeaderOption, "--header": HeaderOption,
	"-X": MethodOption, "--request": MethodOption,
	"-d": DataOption, "--data": DataOption,
}

var flagOptions = map[string]curlOption{
	"-L": LocationOption, "--location": LocationOption,
}

type curlOptionInstance struct {
	option curlOption
	value  []string
}

type FormParam struct {
	Name, Value string
}

func ParseCommand(cmd []string) (command *Command, err error) {
	var options []curlOptionInstance
	request := &Request{}
	command = &Command{Request: request}

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
							err = fmt.Errorf("curl: unexpected option %s", arg)
							i += len(args) - i
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
							err = fmt.Errorf("curl: unexpected option %s", arg)
							i += len(args) - i
						}
					}
				} else if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
					request.Url = arg
					i += 1
				} else {
					err = fmt.Errorf("curl: unexpected argument %s", arg)
					i += len(args) - i
				}
			}

			for _, inst := range options {
				switch inst.option {
				case HeaderOption:
					request.Headers = append(request.Headers, parseHeader(inst.value[0]))
				case MethodOption:
					request.Method = inst.value[0]
				case DataOption:
					request.Body = parseFormData(inst.value[0])

					if request.Method == "" {
						request.Method = "POST"
					}
				case LocationOption:
					command.FollowRedirects = true
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

func parseFormData(str string) (params []FormParam) {
	for _, kv := range strings.Split(str, "&") {
		param := FormParam{}
		parts := strings.Split(kv, "=")

		if len(parts) == 2 {
			param.Name = parts[0]
			param.Value = parts[1]
		} else {
			param.Name = parts[0]
		}

		params = append(params, param)
	}

	return
}
