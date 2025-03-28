// curl --location -X GET https://httpbin.org/status/302
import io.ktor.client.HttpClient
import io.ktor.client.request.get
import io.ktor.client.statement.bodyAsText
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient()
    val response = client.get("https://httpbin.org/status/302")
    print(response.bodyAsText())
}