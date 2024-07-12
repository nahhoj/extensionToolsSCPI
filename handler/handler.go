package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"github.com/nahhoj/extensionToolsSCPI/datatypes"
	"github.com/nahhoj/extensionToolsSCPI/tools"
)

func SecurutyMaterial(w http.ResponseWriter, r *http.Request) {
	cookie := r.Header.Get("cookie")
	tenant := r.Header.Get("tenant")
	credential := r.URL.Query().Get("credential")
	response := tools.GetSecurutyMaterial(tenant, cookie, credential)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	jsonBytes, _ := json.Marshal(response)
	w.Write(jsonBytes)
}

func KeyPair(w http.ResponseWriter, r *http.Request) {
	cookie := r.Header.Get("cookie")
	tenant := r.Header.Get("tenant")
	key := r.URL.Query().Get("key")
	response := tools.GetPairKey(tenant, cookie, key)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	jsonBytes, _ := json.Marshal(response)
	w.Write(jsonBytes)
}

func WebServiceTest(w http.ResponseWriter, r *http.Request) {
	var responseBody string
	responseBody = fmt.Sprintln(r.Method)
	responseBody += fmt.Sprintln(r.RequestURI)
	responseBody += fmt.Sprintln("")
	for i, m := range r.Header {
		responseBody += fmt.Sprintln(i, m)
	}
	responseBody += fmt.Sprintln("")
	body, _ := io.ReadAll(r.Body)
	responseBody += fmt.Sprintln(string(body))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseBody))
}

func FormatCode(w http.ResponseWriter, r *http.Request) {
	var formatCode datatypes.FormatCode
	fileName := uuid.New()
	bodyBytes, _ := io.ReadAll(r.Body)
	json.Unmarshal(bodyBytes, &formatCode)
	bodyDecoder, _ := base64.StdEncoding.DecodeString(string(formatCode.Code))
	os.WriteFile(fileName.String()+".groovy", bodyDecoder, 0644)
	cmd := exec.Command("npm-groovy-lint", "--format", fileName.String()+".groovy")
	err := cmd.Run()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		bodyDecoder, _ = os.ReadFile(fileName.String() + ".groovy")
		formatCode = datatypes.FormatCode{
			Code: base64.StdEncoding.EncodeToString(bodyDecoder),
		}
		jsonOutput, _ := json.Marshal(&formatCode)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonOutput)
	}
	os.Remove(fileName.String() + ".groovy")
}

func GroovyLog(w http.ResponseWriter, r *http.Request) {
	var groovyLog datatypes.GroovyLog
	bodyBytes, _ := io.ReadAll(r.Body)
	json.Unmarshal(bodyBytes, &groovyLog)
	cmd := exec.Command("groovy", "utils/main.groovy", groovyLog.Script, groovyLog.Body, groovyLog.Headers, groovyLog.Properties, groovyLog.Method)
	output, err := cmd.Output()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		outputArray := strings.Split(string(output), "\r\n")
		groovyLog = datatypes.GroovyLog{
			Log:        string(string(output)[strings.Index(string(output), "-start-")+5 : strings.Index(string(output), "-end-")]),
			Body:       outputArray[4],
			Headers:    outputArray[5],
			Properties: outputArray[6],
		}
		jsonOutput, _ := json.Marshal(&groovyLog)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonOutput)
	}
}
