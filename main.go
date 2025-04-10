package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	directoryPath := "./public"

	_, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
		fmt.Printf("Directory '%s' not found.\n", directoryPath)
		return
	}

	//create file server handler to serve directory contents
	fileServer := http.FileServer(http.Dir(directoryPath))

	//create http server to handle requests to '/' using fileServer
	http.Handle("/", fileServer)

	port := 8080
	fmt.Printf("Server started at localhost:%d \n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
