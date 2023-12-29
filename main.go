package main

import (
	"fmt"
	"net/http"
	"pedrolucaspalma/go-media-center/handlers"
)

func main() {
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/example", handlers.ExampleHandler)

	http.HandleFunc("/selected", handlers.HandleVideo)
	http.HandleFunc("/player", handlers.PlayerHandler)

	fmt.Println("Levantando server")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
