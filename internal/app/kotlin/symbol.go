package kotlin

type symbol struct {
	Name    SimpleId
	Package *Fqn
}

var clientPackage = Fqn{"io", "ktor", "client"}
var requestPackage = Fqn{"io", "ktor", "client", "request"}
var formsPackage = Fqn{"io", "ktor", "client", "request", "forms"}
var clientStatementPackage = Fqn{"io", "ktor", "client", "statement"}
var httpPackage = Fqn{"io", "ktor", "http"}
var coroutinesPackage = Fqn{"kotlinx", "coroutines"}

var getRequest = &symbol{"get", &requestPackage}
var postRequest = &symbol{"post", &requestPackage}
var patchRequest = &symbol{"patch", &requestPackage}
var headRequest = &symbol{"head", &requestPackage}
var optionsRequest = &symbol{"options", &requestPackage}
var deleteRequest = &symbol{"delete", &requestPackage}
var putRequest = &symbol{"put", &requestPackage}
var requestRequest = &symbol{"request", &requestPackage}

var builders = []*symbol{getRequest, postRequest, patchRequest, headRequest, optionsRequest, deleteRequest, putRequest}

var setBody = &symbol{"setBody", &requestPackage}

var httpClient = &symbol{"HttpClient", &clientPackage}

var httpMethod = &symbol{"HttpMethod", &httpPackage}
var parameters = &symbol{"parameters", &httpPackage}

var runBlocking = &symbol{"runBlocking", &coroutinesPackage}

var formDataContent = &symbol{"FormDataContent", &formsPackage}

var bodyAsText = &symbol{"bodyAsText", &clientStatementPackage}
