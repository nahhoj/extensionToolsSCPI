/* Refer the link below to learn more about the use cases of script.
https://help.sap.com/viewer/368c481cd6954bdfa5d0435479fd4eaf/Cloud/en-US/148851bf8192412cba1f9d2c17f4bd25.html

If you want to know more about the SCRIPT APIs, refer the link below
https://help.sap.com/doc/a56f52e1a58e4e2bac7f7adbf45b2e26/Cloud/en-US/index.html */
import com.sap.gateway.ip.core.customdev.util.Message
import java.util.HashMap
import org.apache.camel.impl.*

def Message processData(Message message) {
    def camelContext = new DefaultCamelContext()
    StringBuilder versions = new StringBuilder()
    versions << "Groovy: ${GroovySystem.getVersion()}\n"
    versions << "Java: ${System.getProperty('java.version')}\n"
    versions << "Camel: ${camelContext.getVersion()}"
    message.setBody(versions.toString())
    if (versions.length < 0) {
        println('hello')
    }

    return message
}
