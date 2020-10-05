package main

import (
	"log"
	"net/http"
	"os"

	dbin "github.com/TCP_Proxy/interfacedb"
	tcp "github.com/TCP_Proxy/tcpproxy"
)

func main() {

	arguments := os.Args
	if len(arguments) == 1 {
		log.Println("Please provide port number")
		return
	}

	PORT := arguments[1]

	log.Println("Server started on https://", PORT)

	tcp.Startup()
	go tcp.Calcdata(20)

	router := dbin.NewRouter()
	http.ListenAndServeTLS(PORT, "localhost.crt", "localhost.key", router)
}
