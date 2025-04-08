// curl -F 'meta={
//  "images": [
//    {
//      "name": "file"
//    }
//  ]
//}' https://httpbin.org/post
import io.ktor.client.HttpClient
import io.ktor.client.request.forms.MultiPartFormDataContent
import io.ktor.client.request.forms.formData
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.client.statement.bodyAsText
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = HttpClient {
        followRedirects = false
    }
    val response = client.post("https://httpbin.org/post") {
        setBody(MultiPartFormDataContent(formData {
            append("meta", """{
              "images": [
                {
                  "name": "file"
                }
              ]
            }""".trimIndent())
        }))
    }
    print(response.bodyAsText())
}