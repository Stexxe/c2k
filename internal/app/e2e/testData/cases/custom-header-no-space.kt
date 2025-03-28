// curl -HCustom:123 https://httpbin.org/get
import io.ktor.client.HttpClient
import io.ktor.client.request.get
import io.ktor.client.statement.bodyAsText
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient {
        followRedirects = false
    }
    val response = client.get("https://httpbin.org/get") {
        headers.append("Custom", "123")
    }
    print(response.bodyAsText())
}