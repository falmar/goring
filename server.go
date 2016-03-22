package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/net/websocket"

	"github.com/julienschmidt/httprouter"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	startMapServer()
	go startWebSocketServer()
	startHTTPServer()
}

func startHTTPServer() {
	router := httprouter.New()
	router.GET("/", HomeHandler)
	router.ServeFiles("/src/*filepath", http.Dir("./public"))
	log.Fatal(http.ListenAndServe(":8080", router))
}

func startWebSocketServer() {
	WSMux := http.NewServeMux()
	WSMux.Handle("/getMap", websocket.Handler(getMap))
	log.Fatal(http.ListenAndServe(":9020", WSMux))
}
