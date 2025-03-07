package curl

import (
	"fmt"
	"strings"
)

type Request struct {
	Url string
}

func ParseCommand(cmd []string) (request Request, err error) {
	if len(cmd) > 0 && cmd[0] == "curl" {
		if len(cmd) > 1 {
			request.Url = cmd[1]
		} else {
			err = fmt.Errorf("curl: expected URL, got none")
		}
	} else {
		err = fmt.Errorf("curl: invalid curl command '%s'", strings.Join(cmd, " "))
	}

	return
}
