package handlers

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"triple-s/metadata"
)

func DeleteBucket(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("BucketName")
	fmt.Println(name)

	mu.Lock()
	defer mu.Unlock()

	exists, err := metadata.CheckBucket(name, Directory)
	if err != nil {
		resp := response{
			Code:    500,
			Status:  "error",
			Message: "Error in server",
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusInternalServerError)
		if err := xml.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	if !exists {
		resp := response{
			Code:    404,
			Status:  "error",
			Message: "This bucket does not exist",
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusNotFound)
		if err := xml.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	result := metadata.CheckDir(name, Directory)
	if result == false {
		resp := response{
			Code:    409,
			Status:  "error",
			Message: "Bucket is not empty",
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusConflict)
		if err := xml.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	err = os.Remove(Directory + "/" + name)
	if err != nil {
		log.Fatal(err)
	}

	resp := response{
		Code:    204,
		Status:  "success",
		Message: "Bucket delete",
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	if err := xml.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}
