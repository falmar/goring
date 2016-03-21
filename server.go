package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {

	rand.Seed(time.Now().UnixNano())
	startMapServer()
	startHTTPServer()
}

func startHTTPServer() {
	router := httprouter.New()
	router.GET("/", HomeHandler)
	router.GET("/getMap", getMap)
	router.ServeFiles("/src/*filepath", http.Dir("./public"))
	log.Fatal(http.ListenAndServe(":8080", router))
}
