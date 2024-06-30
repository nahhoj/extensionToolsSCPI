package main

import (
	"net/http"

	"github.com/nahhoj/extensionToolsSCPI/handler"
)

func main() {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/securutymaterial", handler.SecurutyMaterial)
	serverMux.HandleFunc("/keypair", handler.KeyPair)
	serverMux.HandleFunc("/webservicetest", handler.WebServiceTest)
	server := http.Server{Addr: ":8080", Handler: serverMux}
	server.ListenAndServe()
}
