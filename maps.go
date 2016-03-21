package main

func startMapServer() {

	maps := []*Map{
		NewMap("prontera", 18),
	}

	for _, m := range maps {
		go m.Run()
	}
}
