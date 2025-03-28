// curl --request POST 'https://httpbin.org/post' --header 'Content-Type: text/plain'
import io.ktor.client.HttpClient
import io.ktor.client.request.post
import io.ktor.client.statement.bodyAsText
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient {
        followRedirects = false
    }
    val response = client.post("https://httpbin.org/post") {
        headers.append("Content-Type", "text/plain")
    }
    print(response.bodyAsText())
}