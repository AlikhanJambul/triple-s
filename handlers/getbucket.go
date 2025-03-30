package handlers

import (
	"encoding/csv"
	"encoding/xml"
	"net/http"
	"os"
	"triple-s/metadata"
)

type Bucket struct {
	Name         string `xml:"Name"`
	Status       string `xml:"Status"`
	Size         string `xml:"Size"`
	CreationTime string `xml:"CreationTime"`
	LastModified string `xml:"LastModified"`
}

type ListBucketsResponse struct {
	XMLName xml.Name `xml:"Buckets"`
	Buckets []Bucket `xml:"Bucket"`
}

func GetBuckets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	empty, err := metadata.CountDirs(Directory)
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
	if len(empty) == 0 {
		resp := response{
			Code:    409,
			Status:  "error",
			Message: Directory + " is empty",
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusConflict)
		if err := xml.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	csvFilePath := Directory + "/buckets.csv"

	file, err1 := os.Open(csvFilePath)
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
	defer file.Close()

	reader := csv.NewReader(file)
	records, err2 := reader.ReadAll()
	if err2 != nil {
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

	var buckets []Bucket
	for _, record := range records {
		if len(record) < 5 {
			http.Error(w, "Invalid bucket metadata format", http.StatusInternalServerError)
			return
		}
		buckets = append(buckets, Bucket{
			Name:         record[0],
			Status:       record[1],
			Size:         record[2],
			CreationTime: record[3],
			LastModified: record[4],
		})
	}

	response := ListBucketsResponse{Buckets: buckets}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	if err := xml.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode XML response", http.StatusInternalServerError)
		return
	}
}

// LastModified: record[3],
// CreationTime: record[4],
