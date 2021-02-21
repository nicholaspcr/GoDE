package mo

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"gitlab.com/nicholaspcr/go-de/pkg/problems/models"
)

type fileManager struct {
	f *os.File
}

// checks existance of filePath
func checkFilePath(basePath, filePath string) {
	folders := strings.Split(filePath, "/")
	for _, folder := range folders {
		basePath += "/" + folder
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			err = os.Mkdir(basePath, os.ModePerm)
			if err != nil {
				fmt.Println(basePath, folder)
				log.Fatal(err)
			}
		}
	}
}

// writeHeader writes the header of the csv writer file
func writeHeader(pop []models.Elem, w *csv.Writer) {
	tmpData := []string{}
	for i := range pop {
		tmpData = append(tmpData, fmt.Sprintf("elem[%d]", i))
	}
	err := w.Write(tmpData)
	if err != nil {
		log.Fatal("Couldn't write file")
	}
	w.Flush()
}

// writeGeneration writes the objectives in the csv writer file
func writeGeneration(elems models.Elements, w *csv.Writer) {
	if len(elems) == 0 {
		return
	}
	data := [][]string{}
	objs := len(elems[0].Objs)
	for i := 0; i < objs; i++ {
		tmpData := []string{}
		for _, p := range elems {
			tmpData = append(tmpData, fmt.Sprintf("%5.3f", p.Objs[i]))
		}
		data = append(data, tmpData)
	}
	err := w.WriteAll(data)
	if err != nil {
		log.Fatal("Couldn't write file")
	}
	w.Flush()
}

// writeResult creates a file and writes all the elements in it
// it should be used to write a single time a specific result
// in the given path
func writeResult(path string, elems models.Elements) {
	f, err := os.Create(path)
	writer := csv.NewWriter(f)
	writer.Comma = '\t'
	checkError(err)

	// header
	headerData := []string{"elems"}
	collumn := 'A'
	for range elems[0].Objs {
		headerData = append(headerData, string(collumn))
		collumn++
	}
	err = writer.Write(headerData)
	if err != nil {
		log.Fatal("Couldn't write file")
	}
	writer.Flush()

	bodyData := [][]string{}
	for i := range elems {
		tmpData := []string{}
		tmpData = append(tmpData, fmt.Sprintf("elem[%d]", i))
		for _, p := range elems[i].Objs {
			tmpData = append(tmpData, fmt.Sprint(p))
		}
		bodyData = append(bodyData, tmpData)
	}
	err = writer.WriteAll(bodyData)
	if err != nil {
		log.Fatal("Couldn't write file")
	}
	writer.Flush()
	f.Close()
}
