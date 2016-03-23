package main

var maps map[string]*Map

func startMapServer() {
	maps = make(map[string]*Map)
	maps["prontera"] = NewMap("prontera", "Prontera", 18)
	for _, m := range maps {
		go m.Run()
	}
}
