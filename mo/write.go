package mo

import (
	"fmt"
	"os"
)

// todo manage interfaces with type
type fileManager struct {
	f *os.File
}

var paretoPATH = os.Getenv("HOME") + "/.go-de/mode/bolt.db"

// todo: maybe remove this and do a separate subcommand to write the result in a .csv file!
func writeHeader(pop []Elem, f *os.File) {
	for i := range pop {
		fmt.Fprintf(f, "elem[%d]\t", i)
	}
	fmt.Fprintf(f, "\n")
}

// todo: maybe remove this and do a separate subcommand to write the result in a .csv file!
func writeGeneration(elems Elements, f *os.File) {
	if len(elems) == 0 {
		return
	}
	objs := len(elems[0].objs)
	for i := 0; i < objs; i++ {
		for _, p := range elems {
			fmt.Fprintf(f, "%10.3f\t", p.objs[i])
		}
		fmt.Fprintf(f, "\n")
	}
}
