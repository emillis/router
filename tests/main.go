package main

import (
	"log"
	"net/http"
	"router"
)

func main() {
	r := router.NewRouter()

	log.Fatal(http.ListenAndServe(":80", r))
}
