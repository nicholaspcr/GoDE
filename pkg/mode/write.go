package mo

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

// todo manage interfaces with type
type fileManager struct {
	f *os.File
}

// todo: maybe remove this and do a separate subcommand to write the result in a .csv file!
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

// todo: maybe remove this and do a separate subcommand to write the result in a .csv file!
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
