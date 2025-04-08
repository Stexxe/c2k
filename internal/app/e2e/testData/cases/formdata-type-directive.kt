// curl -F "web=@index.html;type=text/html" https://httpbin.org/post
import io.ktor.client.HttpClient
import io.ktor.client.request.forms.ChannelProvider
import io.ktor.client.request.forms.MultiPartFormDataContent
import io.ktor.client.request.forms.formData
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.client.statement.bodyAsText
import io.ktor.http.Headers
import io.ktor.http.HttpHeaders
import io.ktor.util.cio.readChannel
import java.io.File
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient {
        followRedirects = false
    }
    val response = client.post("https://httpbin.org/post") {
        setBody(MultiPartFormDataContent(formData {
            val index = File("index.html")
            append("web", ChannelProvider(size = index.length()) { index.readChannel() }, Headers.build {
                append(HttpHeaders.ContentType, "text/html")
                append(HttpHeaders.ContentDisposition, "filename=\"${index.name}\"")
            })
        }))
    }
    print(response.bodyAsText())
}