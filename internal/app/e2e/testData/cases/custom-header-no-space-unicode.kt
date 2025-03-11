// curl -H'Custom:привет' https://httpbin.org/get
import io.ktor.client.HttpClient
import io.ktor.client.request.get
import io.ktor.client.statement.bodyAsText
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient()
    val response = client.get("https://httpbin.org/get") {
        headers.append("Custom", "привет")
    }
    print(response.bodyAsText())
}