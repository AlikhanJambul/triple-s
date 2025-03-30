package main

import (
	"fmt"
	"net/http"
	"strconv"
	"triple-s/handlers"
	helps "triple-s/helperfunc"
)

func main() {
	portFlag, dirFlag := helps.AllFlags()
	newPort := strconv.Itoa(portFlag)

	helps.CreateDir(dirFlag)

	handlers.Directory = dirFlag

	http.HandleFunc("PUT /{BucketName}", handlers.CreateBucket) // PUT /test, PUT /alem
	http.HandleFunc("GET /", handlers.GetBuckets)
	http.HandleFunc("DELETE /{BucketName}", handlers.DeleteBucket)
	http.HandleFunc("PUT /{BucketName}/{ObjectKey}", handlers.CreateObject)
	http.HandleFunc("DELETE /{BucketName}/{ObjectKey}", handlers.DeleteObject)
	http.HandleFunc("GET /{BucketName}/{ObjectKey}", handlers.GetObject)

	fmt.Printf("Server is running on port %d...\n", portFlag)
	if err := http.ListenAndServe(":"+newPort, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
