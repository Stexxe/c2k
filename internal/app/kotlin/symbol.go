package kotlin

var clientPackage = &Fqn{"io", "ktor", "client"}
var requestPackage = &Fqn{"io", "ktor", "client", "request"}
var formsPackage = &Fqn{"io", "ktor", "client", "request", "forms"}
var clientStatementPackage = &Fqn{"io", "ktor", "client", "statement"}
var httpPackage = &Fqn{"io", "ktor", "http"}
var coroutinesPackage = &Fqn{"kotlinx", "coroutines"}
var javaIoPackage = &Fqn{"java", "io"}
var cioUtilsPackage = &Fqn{"io", "ktor", "util", "cio"}

var getRequest = buildFqn("get", requestPackage)
var postRequest = buildFqn("post", requestPackage)
var patchRequest = buildFqn("patch", requestPackage)
var headRequest = buildFqn("head", requestPackage)
var optionsRequest = buildFqn("options", requestPackage)
var deleteRequest = buildFqn("delete", requestPackage)
var putRequest = buildFqn("put", requestPackage)
var requestRequest = buildFqn("request", requestPackage)

var builders = []*Fqn{getRequest, postRequest, patchRequest, headRequest, optionsRequest, deleteRequest, putRequest}

var setBody = buildFqn("setBody", requestPackage)

var httpClient = buildFqn("HttpClient", clientPackage)

var httpMethod = buildFqn("HttpMethod", httpPackage)
var parameters = buildFqn("parameters", httpPackage)
var headersObject = buildFqn("Headers", httpPackage)
var httpHeadersObject = buildFqn("HttpHeaders", httpPackage)

var runBlocking = buildFqn("runBlocking", coroutinesPackage)

var formDataContent = buildFqn("FormDataContent", formsPackage)
var multipartContent = buildFqn("MultiPartFormDataContent", formsPackage)
var formData = buildFqn("formData", formsPackage)
var channelProvider = buildFqn("ChannelProvider", formsPackage)

var bodyAsText = buildFqn("bodyAsText", clientStatementPackage)

var fileCtor = buildFqn("File", javaIoPackage)
var readChannel = buildFqn("readChannel", cioUtilsPackage)

func buildFqn(name string, pack *Fqn) *Fqn {
	fqn := append(Fqn{}, *pack...)
	fqn = append(fqn, name)
	return &fqn
}
