package main

import (
	"log"
	"net/http"
)

func main() {
	r := router.NewRouter()

	log.Fatal(http.ListenAndServe(":80", r))
}
