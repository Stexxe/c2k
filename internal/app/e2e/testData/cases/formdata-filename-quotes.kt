// curl --form 'upload=@"/data/1.XML”' https://httpbin.org/post
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
            val file = File("/data/1.XML")
            append("upload", ChannelProvider(size = file.length()) { file.readChannel() }, Headers.build {
                append(HttpHeaders.ContentType, "application/octet-stream")
                append(HttpHeaders.ContentDisposition, "filename=\"1.XML\"")
            })
        }))
    }
    print(response.bodyAsText())
}