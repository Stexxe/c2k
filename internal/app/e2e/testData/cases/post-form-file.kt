// curl -d @filename https://httpbin.org/post
import io.ktor.client.HttpClient
import io.ktor.client.request.forms.FormDataContent
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.client.statement.bodyAsText
import io.ktor.http.parameters
import io.ktor.http.parseUrlEncodedParameters
import java.io.File
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient {
        followRedirects = false
    }
    val response = client.post("https://httpbin.org/post") {
        setBody(FormDataContent(parameters {
            val filenameParams = File("filename").readText().parseUrlEncodedParameters()
            for ((name, values) in filenameParams.entries()) {
                for (v in values) {
                    append(name, v)
                }
            }
        }))
    }
    print(response.bodyAsText())
}