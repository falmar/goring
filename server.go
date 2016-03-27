package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/net/websocket"

	"github.com/julienschmidt/httprouter"
)

var totalMobs int
var totalMaps int

func main() {
	rand.Seed(time.Now().UnixNano())
	mapChan := make(chan bool)
	go startMapServer(mapChan)
	<-mapChan
	fmt.Println("Total Maps loaded:", totalMaps)
	fmt.Println("Total Monsters loaded:", totalMobs)
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
	WSMux.Handle("/getMob", websocket.Handler(getMob))
	log.Fatal(http.ListenAndServe(":9020", WSMux))
}
