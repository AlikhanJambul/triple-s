package handlers

import (
	"encoding/xml"
	"net/http"
	"os"
	"sync"
	validate "triple-s/helperfunc"
	"triple-s/metadata"
)

var (
	mu        sync.Mutex
	Directory string
)

type response struct {
	XMLName xml.Name `xml:"response"` // Kорневой элемент XML
	Code    uint16   `xml:"code"`
	Status  string   `xml:"status"`
	Message string   `xml:"message"`
	Bucket  string   `xml:"bucket,omitempty"`
}

func CreateBucket(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("BucketName")

	if !validate.IsValidName(name) {
		resp := response{
			Code:    400,
			Status:  "error",
			Message: "Bucket name is invalid",
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusBadRequest)
		if err := xml.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	mu.Lock()
	defer mu.Unlock()

	exists, err1 := metadata.CheckBucket(name, Directory)
	if err1 != nil {
		resp := response{
			Code:    409,
			Status:  "error",
			Message: "Error in server",
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusConflict)
		if err := xml.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if exists == true {
		resp := response{
			Code:    409,
			Status:  "error",
			Message: "Bucket name already exists",
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusConflict)
		if err := xml.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err := os.Mkdir(Directory+"/"+name, 0o766)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	dirSize, _ := metadata.GetFolderSize(Directory + "/" + name)

	metadata.WriteInBucketCsv(name, dirSize, "inactive", Directory)

	resp := response{
		Code:    200,
		Status:  "success",
		Message: "Bucket created successfully",
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	if err := xml.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
