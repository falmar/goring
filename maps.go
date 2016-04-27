package main

var maps map[string]*Map

func startMapServer(mapChan chan bool) {
	maps = make(map[string]*Map)
	var loadedMaps int
	loadedMap := make(chan bool)

	//TODO: Loop to load mobs from db or json file
	maps["prontera"] = NewMap("prontera", "Prontera", [2]int64{18, 15})

	totalMaps = len(maps)

	for _, m := range maps {
		go m.Run(loadedMap)
	}

checkMapLoop:
	for {
		select {
		case <-loadedMap:
			loadedMaps++
			if loadedMaps == totalMaps {
				break checkMapLoop
			}
		}
	}

	mapChan <- true
}
