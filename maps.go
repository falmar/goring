package main

var maps map[string]*Map
var mapCmdChan map[string]chan string

func startMapServer() {

	maps = make(map[string]*Map)
	mapCmdChan = make(map[string]chan string)

	maps["prontera"] = NewMap("prontera", "Prontera", 18)

	for _, m := range maps {
		mapCmdChan[m.id] = make(chan string)
		go m.Run()
	}
}
