package main

import (
	"fmt"
	"labs/redis/router"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	router := router.NewRouter()
	fmt.Println("Listening....")
	log.Fatal(http.ListenAndServe(":80", router))
}
