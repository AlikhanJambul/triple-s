package handlers

import (
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"triple-s/metadata"
)

func GetObject(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("BucketName")
	object := r.PathValue("ObjectKey")
	exists, _ := metadata.CheckBucket(bucket, Directory)
	if !exists {
		resp := response{
			Code:    409,
			Status:  "error",
			Message: "Bucket does not exist",
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusBadRequest)
		if err := xml.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	existsObject, _ := metadata.CheckObject(bucket, Directory, object)
	if !existsObject {
		resp := response{
			Code:    409,
			Status:  "error",
			Message: "Object does not exist",
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusBadRequest)
		if err := xml.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	ct := metadata.CheckObjectCsvFormat(bucket, Directory, object)

	w.Header().Set("Content-Type", ct)
	filePath := Directory + "/" + bucket + "/" + object

	file, err := os.Open(filePath)
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
	defer file.Close()

	_, err = io.Copy(w, file)
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
}
