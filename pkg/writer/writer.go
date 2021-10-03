// package writer is responsible for writing data related
// to the Vector struct that can be found on the models pkg

package writer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/nicholaspcr/gde3/pkg/models"
)

// Writer is the custom writer provided by the gode package, it contains the
// methods used to write the information regarding Population and Vector
type Writer struct {
	*csv.Writer
}

const (
	floatFormat = "%.8f"
)

// NewWriter returns a Writer pointer that contains the methods to write
// Population and Vector into a file with a specific path
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Writer: csv.NewWriter(w),
	}
}

// CheckFilePath checks the existance of filePath
func CheckFilePath(basePath, filePath string) {
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
func (w *Writer) WriteHeader(sz int) error {
	tmpData := []string{"header"}
	headName := "A"
	for i := 0; i < sz; i++ {
		tmpData = append(tmpData, headName)
		headName = incrementHeader(headName)
	}
	err := w.Write(tmpData)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func incrementHeader(h string) string {
	str := []rune(h)
	pos := len(str) - 1

	if h[pos] < 'Z' {
		str[pos]++
		return string(str)
	}
	for ; len(str) > 0 && str[pos] >= 'Z'; pos-- {
		str[pos] = 'A'
		if pos == 0 {
			str = append([]rune{'A'}, str...)
			pos++
		} else {
			str[pos-1] = rune(str[pos-1] + 1)
		}
	}
	return string(str)
}

// WriteGeneration writes the objectives in the csv writer file
func (w *Writer) ElementsObjs(elems models.Population) error {
	if len(elems) == 0 {
		return errors.New("empty slice of elements")
	}

	// matrix of slices
	data := make([][]string, len(elems))
	objs := len(elems[0].Objs)

	for ind, p := range elems {
		// allocates space for the current slice
		data[ind] = make([]string, len(p.Objs)+1)
		data[ind][0] = fmt.Sprintf("elem[%d]", ind)

		for i := 0; i < objs; i++ {
			data[ind][i+1] = fmt.Sprintf(floatFormat, p.Objs[i])
		}
	}

	err := w.WriteAll(data)
	if err != nil {
		log.Fatal("Couldn't write file")
	}

	w.Flush()
	return nil
}

func (w *Writer) ElementsVectors(elems models.Population) error {
	if len(elems) == 0 {
		return errors.New("empty slice of elements")
	}

	// matrix of slices
	data := make([][]string, len(elems))
	dim := len(elems[0].X)

	for ind, p := range elems {
		// allocates space for the current slice
		data[ind] = make([]string, len(p.X)+1)
		data[ind][0] = fmt.Sprintf("elem[%d]", ind)

		for i := 0; i < dim; i++ {
			data[ind][i+1] = fmt.Sprintf(floatFormat, p.X[i])
		}
	}

	err := w.WriteAll(data)
	if err != nil {
		log.Fatal("Couldn't write file")
	}

	w.Flush()
	return nil
}
