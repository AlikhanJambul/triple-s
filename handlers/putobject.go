package handlers

import (
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"strconv"
	create "triple-s/helperfunc"
	"triple-s/metadata"
)

func CreateObject(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")

	bucket := r.PathValue("BucketName")
	object := r.PathValue("ObjectKey")
	if object == "object.csv" {
		resp := response{
			Code:    400,
			Status:  "error",
			Message: object + " cannot be created",
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusBadRequest)
		if err := xml.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	if valid := create.IsValidName(object); !valid {
		resp := response{
			Code:    400,
			Status:  "error",
			Message: "Object name is invalid",
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusBadRequest)
		if err := xml.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	if exists, _ := metadata.CheckBucket(bucket, Directory); !exists {
		resp := response{
			Code:    409,
			Status:  "error",
			Message: bucket + " not found",
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusConflict)
		if err := xml.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	if status := metadata.GetStatus(bucket, Directory); status == "inactive" {
		file, _ := os.Create(Directory + "/" + bucket + "/object.csv")
		defer file.Close()
	}
	var length int

	cl := r.Header.Get("Content-Length")
	if cl != "" {
		length, _ = strconv.Atoi(cl)
	}

	content := make([]byte, length)
	_, err := io.ReadAtLeast(r.Body, content, length)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err1 := create.WriteBytesToFile(object, bucket, Directory, content)
	if err1 != nil {
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
	if exists, _ := metadata.CheckObject(bucket, Directory, object); !exists {
		metadata.WriteInObjectCsv(object, ct, Directory, bucket)
	} else if exists {
		metadata.ChangeObject(bucket, Directory, object, ct)
	}
	metadata.ChangeMetadataStatus(bucket, Directory, "active")

	resp := response{
		Code:    200,
		Status:  "success",
		Message: "Object created successfully",
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	if err := xml.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
