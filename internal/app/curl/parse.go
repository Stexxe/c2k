package curl

import (
	"fmt"
	"strings"
)

type HttpMethod string

const (
	GetMethod     HttpMethod = "GET"
	PostMethod    HttpMethod = "POST"
	DeleteMethod  HttpMethod = "DELETE"
	HeadMethod    HttpMethod = "HEAD"
	PutMethod     HttpMethod = "PUT"
	OptionsMethod HttpMethod = "OPTIONS"
	PathMethod    HttpMethod = "PATCH"
)

type Request struct {
	Url    string
	Method HttpMethod
}

func ParseCommand(cmd []string) (request Request, err error) {
	if len(cmd) > 0 && cmd[0] == "curl" {
		if len(cmd) > 1 {
			i := 0
			args := cmd[1:]
			for i < len(args) {
				arg := args[i]
				if arg == "-X" && i+1 < len(args) {
					request.Method = HttpMethod(args[i+1])
					i += 2
				} else {
					break
				}
			}

			if request.Method == "" {
				request.Method = GetMethod
			}

			request.Url = args[i]
		} else {
			err = fmt.Errorf("curl: expected URL, got none")
		}
	} else {
		err = fmt.Errorf("curl: invalid curl command '%s'", strings.Join(cmd, " "))
	}

	return
}
