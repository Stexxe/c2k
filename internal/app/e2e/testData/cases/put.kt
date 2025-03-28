// curl -X PUT https://httpbin.org/put
import io.ktor.client.HttpClient
import io.ktor.client.request.put
import io.ktor.client.statement.bodyAsText
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient {
        followRedirects = false
    }
    val response = client.put("https://httpbin.org/put")
    print(response.bodyAsText())
}