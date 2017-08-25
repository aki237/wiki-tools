package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/shivam07a/tparse"
)

type Init struct {
	OmitFile   bool
	C          string
	CommonTags string
	RestArgs   []string
}

func (init *Init) Exec() error {
	dir := "./"
	if len(init.RestArgs) == 1 {
		dir = init.RestArgs[0]
		if Exists(dir) {
			return errors.New("A file or folder of the same name already exists")
		}
	}
	dir = init.C + dir
	if err := os.MkdirAll(dir+"/.wbuild/", 0755); err != nil {
		return err
	}

	config, err := os.OpenFile(dir+"/.wbuild/state", os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	stateData := make(tparse.Dict, 0)
	stateEntries := make(tparse.Entries, 0)

	files := []string{
		"doc.md",
		"DESCRIPTION",
		"AUTHOR",
	}

	if init.OmitFile {
		files = append(files, ".omit")
	}

	for _, file := range files {
		f, err := os.Create(dir + "/" + file)
		if err != nil {
			return err
		}
		h := md5.New()
		if _, err := io.Copy(h, f); err != nil {
			return err
		}
		stateEntries[file] = hex.EncodeToString(h.Sum(nil))
		f.Close()
	}
	stateData["files"] = stateEntries
	if init.CommonTags != "" {
		commonTags, err := filepath.Abs(init.CommonTags)
		if err != nil {
			fmt.Println("Warning : ", err, ",Path : ", init.CommonTags)
		} else {
			stateData["config"] = tparse.Entries{"common": commonTags}
		}
	}

	if err := stateData.UnMarshal(config); err != nil {
		return err
	}

	return config.Close()
}
