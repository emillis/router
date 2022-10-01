package main

import (
	"fmt"
	"github.com/emillis/veryFastRouter"
	"log"
	"net/http"
)

func main() {
	router := veryFastRouter.NewRouter()

	router.HttpStatusCodeHandler(404, func(w http.ResponseWriter, r *http.Request, info *veryFastRouter.AdditionalInfo) {
		w.Write([]byte("404 Not Found"))
	})

	router.HttpStatusCodeHandler(405, func(w http.ResponseWriter, r *http.Request, info *veryFastRouter.AdditionalInfo) {
		w.Write([]byte("I don't know what are you trying to do here, but this method is not allowed!!!"))
	})

	router.HandleFunc("/test", []string{http.MethodGet}, func(w http.ResponseWriter, r *http.Request, info *veryFastRouter.AdditionalInfo) {
		fmt.Println("TEST!")
	})

	router.HandleFunc("/one", []string{http.MethodGet}, func(w http.ResponseWriter, r *http.Request, info *veryFastRouter.AdditionalInfo) {
		w.Write([]byte("Hello, this is coming from /one"))
	})

	router.HandleFunc("/one/two/", []string{http.MethodGet}, func(w http.ResponseWriter, r *http.Request, info *veryFastRouter.AdditionalInfo) {
		w.Write([]byte("Hello, this is coming from /one/two/"))
	})

	router.HandleFunc("/one/two/three", []string{http.MethodPost, http.MethodGet}, func(w http.ResponseWriter, r *http.Request, info *veryFastRouter.AdditionalInfo) {
		w.Write([]byte("Hello, this is coming from /one/two/three (the static one)"))
	})

	router.HandleFunc("/one/two/:three", []string{http.MethodGet}, func(w http.ResponseWriter, r *http.Request, info *veryFastRouter.AdditionalInfo) {
		w.Write([]byte("Hello, this is coming from /one/two/:three"))
	})

	log.Fatal(http.ListenAndServe(":80", router))
}
