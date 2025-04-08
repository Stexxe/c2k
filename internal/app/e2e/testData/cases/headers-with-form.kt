// curl -H 'accept: application/json' --form name=value https://httpbin.org/post
import io.ktor.client.HttpClient
import io.ktor.client.request.forms.MultiPartFormDataContent
import io.ktor.client.request.forms.formData
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.client.statement.bodyAsText
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient {
        followRedirects = false
    }
    val response = client.post("https://httpbin.org/post") {
        headers.append("accept", "application/json")

        setBody(MultiPartFormDataContent(formData {
            append("name", "value")
        }))
    }
    print(response.bodyAsText())
}