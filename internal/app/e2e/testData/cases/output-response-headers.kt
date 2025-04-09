// curl -i https://httpbin.org/get
import io.ktor.client.HttpClient
import io.ktor.client.request.get
import io.ktor.client.statement.bodyAsText
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient {
        followRedirects = false
    }
    val response = client.get("https://httpbin.org/get")

    println("${response.version} ${response.status.value}")

    for ((name, values) in response.headers.entries()) {
        for (v in values) {
            println("$name: $v")
        }
    }

    println()

    print(response.bodyAsText())
}