// dibba packger tool
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aki237/dibba"
)

// Usage prints the usage of the tool
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s -o [output] <files>\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "  <files> strings\n  \tlist of all files to be packaged\n")
}

// packageFiles packages all the files passed into a dibba package
// file (output)
func packageFiles(output string, files []string) error {
	dbFile, err := os.Create(output)
	db := dibba.NewWriter(dbFile)
	if err != nil {
		os.Remove(output)
		return err
	}
	for _, val := range files {
		f, err := os.Open(val)
		if err != nil {
			os.Remove(output)
			return err
		}
		file := dibba.NewFile(filepath.Base(f.Name()), f)
		err = db.Add(file)
		if err != nil {
			os.Remove(output)
			return err
		}
	}
	return db.Commit()
}

func main() {
	var outFileName = flag.String("o", "", "name of the dibba output file")
	flag.Parse()
	if *outFileName == "" {
		Usage()
		return
	}
	if len(flag.Args()) < 1 {
		Usage()
		return
	}

	if err := packageFiles(*outFileName, flag.Args()); err != nil {
		fmt.Println(err)
	}
}
