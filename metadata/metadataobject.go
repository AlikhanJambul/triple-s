package metadata

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"
)

func WriteInObjectCsv(nameObject, contentType, dir, bucket string) {
	fileCsv, _ := os.OpenFile(dir+"/"+bucket+"/object.csv", os.O_APPEND|os.O_WRONLY, 0o644)
	defer fileCsv.Close()
	info, _ := os.Stat(dir + "/" + bucket + "/" + nameObject)
	size := strconv.Itoa(int(info.Size()))
	writer := csv.NewWriter(fileCsv)

	err := writer.Write([]string{
		nameObject,
		size,
		contentType,
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

func ChangeObject(name, dir, object, contentType string) bool {
	path := dir + "/" + name + "/object.csv"
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return false
	}

	info, _ := os.Stat(dir + "/" + name + "/" + object)
	newSize := strconv.Itoa(int(info.Size()))

	var filteredRecords [][]string
	for _, record := range records {
		if record[0] == object {
			record[1] = newSize
			record[2] = contentType
			record[3] = time.Now().Format("2006-01-02 15:04:05.000")
		}
		filteredRecords = append(filteredRecords, record)
	}

	file, err = os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0o644)
	if err != nil {
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

	return true
}

func CheckObject(name, dir, object string) (bool, error) {
	file, err := os.Open(dir + "/" + name + "/object.csv")
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
		if record[0] == object {
			return true, nil
		}
	}

	return false, nil
}

func DeleteObjectFromCsv(name, dir, object string) (bool, bool) {
	isEmpty := 0
	path := dir + "/" + name + "/object.csv"
	file, err := os.Open(path)
	if err != nil {
		return false, false
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return false, false
	}

	var filteredRecords [][]string
	for _, record := range records {
		if record[0] == object {
			continue
		}
		isEmpty += len(records)
		filteredRecords = append(filteredRecords, record)
	}
	file.Close()

	fileCsv, err1 := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0o644)
	if err1 != nil {
		return false, false
	}
	defer fileCsv.Close()

	writer := csv.NewWriter(fileCsv)
	defer writer.Flush()

	for _, record := range filteredRecords {
		err := writer.Write(record)
		if err != nil {
			return false, false
		}
	}
	return true, isEmpty == 0
}

func CheckObjectCsvFormat(name, dir, object string) string {
	file, err := os.Open(dir + "/" + name + "/object.csv")
	if err != nil {
		return ""
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err1 := reader.ReadAll()
	if err1 != nil {
		return ""
	}

	for _, record := range records {
		if record[0] == object {
			return record[2]
		}
	}

	return ""
}
