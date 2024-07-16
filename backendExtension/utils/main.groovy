package com.co.johan.cpi;

import com.sap.gateway.ip.core.customdev.processor.MessageImpl
import com.sap.gateway.ip.core.customdev.util.Message
import groovy.json.*

def scriptCode = this.args[0]
def body = this.args[1]
def headers = this.args[2]
def properties = this.args[3]
def method = this.args[4]

GroovyShell shell = new GroovyShell()
def script = shell.parse(new String(scriptCode.decodeBase64()))

Message msg = new MessageImpl()
msg.setBody(new String(body.decodeBase64()));

def jsonSlurper = new JsonSlurper()
def headersJSON=jsonSlurper.parseText(new String(headers.decodeBase64()))
def propertiesJOSN=jsonSlurper.parseText(new String(properties.decodeBase64()))

headersJSON.each{
    it.each{
        msg.setHeader(it.key, it.value);
    }
}

propertiesJOSN.each{
    it.each{
        msg.setProperty(it.key, it.value);    
    }
}

println("-start-")
script."$method"(msg)
println("-end-")

headers = msg.getHeaders();
def h=[];
headers.each{
    h.add("$it.key":"$it.value");
}
def p=[];
properties = msg.getProperties();
properties.each{
    p.add("$it.key":"$it.value");
}

println(msg.getBody().bytes.encodeBase64().toString())
println(groovy.json.JsonOutput.toJson(h).bytes.encodeBase64().toString());
println(groovy.json.JsonOutput.toJson(p).bytes.encodeBase64().toString());