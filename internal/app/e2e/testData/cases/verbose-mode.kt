// curl -v https://httpbin.org/get
// Dependencies:
// implementation("io.ktor:ktor-client-logging")
// implementation("org.slf4j:slf4j-simple")
import io.ktor.client.HttpClient
import io.ktor.client.plugins.logging.LogLevel
import io.ktor.client.plugins.logging.Logging
import io.ktor.client.plugins.logging.LoggingFormat
import io.ktor.client.request.get
import io.ktor.client.statement.bodyAsText
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient {
        install(Logging) {
            level = LogLevel.HEADERS
            format = LoggingFormat.OkHttp
        }
        followRedirects = false
    }
    val response = client.get("https://httpbin.org/get")
    print(response.bodyAsText())
}