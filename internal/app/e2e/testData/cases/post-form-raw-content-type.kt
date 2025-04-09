// curl --data-raw '@he@llo@' -H 'Content-Type: application/json' https://httpbin.org/post
import io.ktor.client.HttpClient
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.client.statement.bodyAsText
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient {
        followRedirects = false
    }
    val response = client.post("https://httpbin.org/post") {
        headers.append("Content-Type", "application/json")

        setBody("@he@llo@")
    }
    print(response.bodyAsText())
}