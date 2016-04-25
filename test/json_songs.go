package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type playlist struct {
	Title string   `json:"title"`
	URL   string   `json:"url"`
	Songs []string `json:"songs"`
}

func main() {

	file, err := os.Open("songs.txt")

	if err != nil {
		fmt.Println(err)
		return
	}

	i := 0
	scanner := bufio.NewScanner(file)
	playlist := &playlist{
		Songs: make([]string, 0),
	}

	for scanner.Scan() {
		line := scanner.Text()

		if i == 0 {
			playlist.URL = line
		} else if i == 1 {
			playlist.Title = line
		} else {
			line = line[strings.LastIndex(line, ":")+4:]
			line = strings.Replace(line, string('\u0026'), "&", -1)
			playlist.Songs = append(playlist.Songs, line)
		}

		i++
	}

	jsn, err := json.Marshal(playlist)

	if err != nil {
		fmt.Println(err)
		return
	}

	file, err = os.Create("json_songs.txt")

	if err != nil {
		return
	}

	//writer := bufio.NewWriter(file)

	file.Write(jsn)
}
