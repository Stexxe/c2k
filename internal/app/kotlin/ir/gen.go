package ir

import (
	"c2k/internal/app/curl"
)

func GenScope(command *curl.Command) (fileScope *Scope, err error) {
	fileScope = newScope(nil)

	declareFunc(fileScope, mainSign, func(sc *Scope) {
		ret(sc, runBlockingCall(sc, func(sc *Scope) {
			client := declareVal(sc, "client", httpClientType, httpClientCtorBlock(sc, func(sc *Scope) {
				assignProp(sc, "followRedirects", false)
			}))
			response := declareVal(sc, "response", httpResponseType, callMethod(sc, client, getRequest, command.Request.Url))
			printCall(sc, callMethod(sc, response, bodyAsText))
		}))
	})

	return
}
