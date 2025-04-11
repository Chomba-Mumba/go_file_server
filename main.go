package main

import (
	"fmt"
	"go_file_server/utils"
	"net/http"
	"os"
)

func main() {
	directoryPath := "./public"

	// check if directory exists
	_, err := os.Stat(directoryPath)

	// check if error was returned from os.Stat and if error was IsNotExist
	if os.IsNotExist(err) {
		fmt.Printf("Directory '%s' not found.\n", directoryPath)
		return
	}

	//create file server handler to serve directory contents
	fileServer := http.FileServer(http.Dir(directoryPath))

	//create http server to handle requests to '/' using fileServer
	http.Handle("/", fileServer)
	http.HandleFunc("/upload", fileUploadHandler)

	port := 8080
	fmt.Printf("Server started at localhost:%d \n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {

	// Limit file size to 10MB
	r.ParseMultipartForm(10 << 20)

	//get file from form data
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Fprintf(w, "Uploaded File: %s\n", handler.Filename)
	fmt.Fprint(w, "File Size: %d\n", handler.Size)
	fmt.Fprint(w, "MIME Header: %v\n", handler.Header)

	// save locally
	dst, err := utils.CreateFile(handler.Filename)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// copy uploaded file to destination file
	if _, err := dst.ReadFrom(file); err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
	}
}
