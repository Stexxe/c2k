// curl --data-urlencode 'name=@file.txt' https://httpbin.org/post
import io.ktor.client.HttpClient
import io.ktor.client.request.forms.FormDataContent
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.client.statement.bodyAsText
import io.ktor.http.parameters
import java.io.File
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient {
        followRedirects = false
    }
    val response = client.post("https://httpbin.org/post") {
        setBody(FormDataContent(parameters {
            append("name", File("file.txt").readText())
        }))
    }
    print(response.bodyAsText())
}