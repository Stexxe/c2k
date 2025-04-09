// curl --request POST '<URL goes here>' --header 'Content-Type: text/plain'
import io.ktor.client.HttpClient
import io.ktor.client.request.post
import io.ktor.client.statement.bodyAsText
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient {
        followRedirects = false
    }
    val response = client.post("<URL goes here>") {
        headers.append("Content-Type", "text/plain")
    }
    print(response.bodyAsText())
}