// curl -X SOME https://httpbin.org/anything
import io.ktor.client.HttpClient
import io.ktor.client.request.request
import io.ktor.client.statement.bodyAsText
import io.ktor.http.HttpMethod
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient {
        followRedirects = false
    }
    val response = client.request("https://httpbin.org/anything") {
        method = HttpMethod("SOME")
    }
    print(response.bodyAsText())
}