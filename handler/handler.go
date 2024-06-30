package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
