package main

import (
	"log"
	"net/http"
	"router"
)

func main() {
	log.Fatal(http.ListenAndServe(":80", router.Router{}))
}
