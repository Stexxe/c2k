// curl --data-urlencode 'name=value' --header 'Content-Type: application/json' https://httpbin.org/post
import io.ktor.client.HttpClient
import io.ktor.client.request.forms.FormDataContent
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.client.statement.bodyAsText
import io.ktor.http.parameters
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient {
        followRedirects = false
    }
    val response = client.post("https://httpbin.org/post") {
        headers.append("Content-Type", "application/json")

        setBody(FormDataContent(parameters {
            append("name", "value")
        }))
    }
    print(response.bodyAsText())
}