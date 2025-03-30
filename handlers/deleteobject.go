package handlers

import (
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"triple-s/metadata"
)

func DeleteObject(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("BucketName")
	object := r.PathValue("ObjectKey")
	// fmt.Println(bucket, object)
	mu.Lock()
	defer mu.Unlock()

	if exists, _ := metadata.CheckBucket(bucket, Directory); !exists {
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

	if existsObject, _ := metadata.CheckObject(bucket, Directory, object); !existsObject {
		resp := response{
			Code:    404,
			Status:  "error",
			Message: "This object does not exist",
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusNotFound)
		if err := xml.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	err := os.RemoveAll(Directory + "/" + bucket + "/" + object)
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

	delete, isEmpty := metadata.DeleteObjectFromCsv(bucket, Directory, object)
	if !delete {
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
	if isEmpty {
		metadata.ChangeMetadataStatus(bucket, Directory, "inactive")
		err = os.Remove(Directory + "/" + bucket + "/object.csv")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		metadata.ChangeMetadataStatus(bucket, Directory, "active")
	}

	resp := response{
		Code:    204,
		Status:  "success",
		Message: "Object delete",
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	if err := xml.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}
