// curl --resolve example.com:443:127.0.0.1 https://example.com
import io.ktor.client.HttpClient
import io.ktor.client.request.get
import io.ktor.client.statement.bodyAsText
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient {
        followRedirects = false
    }
    val response = client.get("https://127.0.0.1") {
        headers.append("Host", "example.com")
    }
    print(response.bodyAsText())
}