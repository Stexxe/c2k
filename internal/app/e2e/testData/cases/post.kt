// curl -X POST https://httpbin.org
import io.ktor.client.HttpClient
import io.ktor.client.request.post
import io.ktor.client.statement.bodyAsText
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient()
    val response = client.post("https://httpbin.org")
    print(response.bodyAsText())
}