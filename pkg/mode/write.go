package mo

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

// todo manage interfaces with type
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
func writeHeader(pop []Elem, w *csv.Writer) {
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
func writeGeneration(elems Elements, w *csv.Writer) {
	if len(elems) == 0 {
		return
	}
	data := [][]string{}
	objs := len(elems[0].objs)
	for i := 0; i < objs; i++ {
		tmpData := []string{}
		for _, p := range elems {
			tmpData = append(tmpData, fmt.Sprintf("%5.3f", p.objs[i]))
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
func writeResult(path string, elems Elements) {
	f, err := os.Create(path)
	writer := csv.NewWriter(f)
	writer.Comma = '\t'
	checkError(err)
	writeHeader(elems, writer)
	writeGeneration(elems, writer)
	f.Close()
}
