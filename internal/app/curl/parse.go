package curl

import (
	"fmt"
	"strings"
)

type Request struct {
	Url     string
	Method  string
	Headers []Header
}

type Header struct {
	Name, Value string
}

type curlOption int

const (
	UnknownOption curlOption = iota
	HeaderOption
	MethodOption
)

var oneArgOptions = map[string]curlOption{
	"-H": HeaderOption, "--header": HeaderOption,
	"-X": MethodOption, "--request": MethodOption,
}

type curlOptionInstance struct {
	option curlOption
	value  []string
}

func ParseCommand(cmd []string) (request Request, err error) {
	var options []curlOptionInstance

	if len(cmd) > 0 && cmd[0] == "curl" {
		if len(cmd) > 1 {
			args := cmd[1:]
			for i := 0; i < len(args); {
				arg := args[i]

				if strings.HasPrefix(arg, "-") {
					if strings.HasPrefix(arg, "--") {
						if opt, ok := oneArgOptions[arg]; ok {
							options = append(options, curlOptionInstance{option: opt, value: args[i+1 : i+2]})
						}

						i += 2
					} else {
						var opt curlOption
						var ok bool
						argBytes := []byte(arg)
						end := len(argBytes)

						for ; end > 0; end-- {
							if opt, ok = oneArgOptions[string(argBytes[0:end])]; ok {
								break
							}
						}

						if opt != UnknownOption {
							if end == len(argBytes) {
								options = append(options, curlOptionInstance{option: opt, value: args[i+1 : i+2]})
								i += 2
							} else {
								options = append(options, curlOptionInstance{option: opt, value: []string{arg[end:]}})
								i += 1
							}
						} else {
							err = fmt.Errorf("curl: unexpected option %s", arg)
						}
					}
				} else {
					if i == len(args)-1 {
						request.Url = arg
					} else {
						err = fmt.Errorf("curl: unexpected argument %s", arg)
					}

					i++
				}
			}

			for _, inst := range options {
				switch inst.option {
				case HeaderOption:
					request.Headers = append(request.Headers, parseHeader(inst.value[0]))
				case MethodOption:
					request.Method = inst.value[0]
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
