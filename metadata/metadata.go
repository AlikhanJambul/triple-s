package metadata

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func WriteInBucketCsv(nameOfDir, size, status, dir string) {
	file, _ := os.OpenFile(dir+"/buckets.csv", os.O_APPEND|os.O_WRONLY, 0o644)
	defer file.Close()

	writer := csv.NewWriter(file)

	err := writer.Write([]string{
		nameOfDir,
		status,
		size,
		time.Now().Format("2006-01-02 15:04:05.000"),
		time.Now().Format("2006-01-02 15:04:05.000"),
	})
	if err != nil {
		log.Fatal(err)
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Fatal(err)
	}
}

func CheckDir(name, dir string) bool {
	dirFile, err := os.Open(dir + "/buckets.csv")
	if err != nil {
		log.Fatal("Ошибка при открытии файла:", err)
	}
	defer dirFile.Close()

	reader := csv.NewReader(dirFile)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Ошибка при чтении файла:", err)
	}

	var filteredRecords [][]string
	for _, record := range records {
		if record[0] == name && record[1] == "inactive" {
			continue
		} else if record[0] == name && record[1] == "active" {
			return false
		}
		filteredRecords = append(filteredRecords, record)
	}

	file, err := os.OpenFile(dir+"/buckets.csv", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0o644)
	if err != nil {
		log.Fatal("Ошибка при открытии файла для записи:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range filteredRecords {
		err := writer.Write(record)
		if err != nil {
			log.Fatal("Ошибка при записи в файл:", err)
		}
	}

	return true
}

func GetStatus(name, dir string) string {
	file, err := os.Open(dir + "/buckets.csv")
	if err != nil {
		log.Fatal("Ошибка при открытии файла:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Ошибка при чтении файла:", err)
	}

	for _, record := range records {
		if record[0] == name {
			return record[1]
		}
	}

	return ""
}

func ChangeMetadataStatus(name, dir, status string) bool {
	fileCsv, err := os.Open(dir + "/buckets.csv")
	if err != nil {
		return false
	}

	reader := csv.NewReader(fileCsv)
	records, err := reader.ReadAll()
	if err != nil {
		return false
	}

	newSize, _ := GetFolderSize(dir + "/" + name)

	var filteredRecords [][]string
	for _, record := range records {
		if record[0] == name {
			record[1] = status
			record[2] = newSize
			record[4] = time.Now().Format("2006-01-02 15:04:05.000")
		}
		filteredRecords = append(filteredRecords, record)
	}
	fileCsv.Close()

	file, err1 := os.OpenFile(dir+"/buckets.csv", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0o644)
	if err1 != nil {
		return false
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range filteredRecords {
		err := writer.Write(record)
		if err != nil {
			return false
		}
	}
	// file.Close()
	return true
}

func CheckBucket(name, dir string) (bool, error) {
	file, err := os.Open(dir + "/buckets.csv")
	if err != nil {
		return false, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err1 := reader.ReadAll()
	if err1 != nil {
		return false, err1
	}

	for _, record := range records {
		if record[0] == name {
			return true, nil
		}
	}

	return false, nil
}

func IsDirEmpty(path string) (bool, error) {
	dir, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer dir.Close()

	files, err := dir.Readdirnames(0)
	if err != nil {
		return false, err
	}

	return len(files) == 0, nil
}

func CountDirs(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	getBucketCount := []string{}

	for _, file := range files {
		if file.IsDir() {
			getBucketCount = append(getBucketCount, file.Name())
		}
	}

	return getBucketCount, nil
}

func GetFolderSize(path string) (string, error) {
	var size int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return strconv.Itoa(int(size)), err
}
