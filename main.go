package main

import (
	"fmt"
	"go_file_server/utils"
	"io"
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
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20+512) //limit request size to 32MB
	defer r.Body.Close()

	// Limit file size in RAM to 32MB
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	//get file from form data
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	//Read file bytes from memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}

	//validate file type
	if !utils.IsValidFileType(fileBytes) {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}

	// save locally
	dst, err := utils.CreateFile(handler.Filename)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	//Write file to disk
	_, err = dst.Write(fileBytes)
	if err != nil {
		http.Error(w, "Error saving the file to disk", http.StatusInternalServerError)
	}

	// copy uploaded file to destination file
	_, err = dst.ReadFrom(file)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "Uploaded File: %s\n", handler.Filename)
	fmt.Print(w, "File Size: %d\n", handler.Size)
	fmt.Print(w, "MIME Header: %v\n", handler.Header)
}
