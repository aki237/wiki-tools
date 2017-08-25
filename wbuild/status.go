package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/shivam07a/tparse"
)

type Status struct {
	C string
}

func (s *Status) Exec() error {
	d, err := os.Getwd()
	if err != nil {
		return err
	}
	dir := filepath.Base(d)
	untracked, changed, err := s.status()
	if err != nil {
		return err
	}
	if len(untracked) != 0 || len(changed) != 0 {
		fmt.Println("In directory '" + dir + "'\n")
	}
	if len(untracked) != 0 {
		fmt.Println("Untracked files :")
		for _, val := range untracked {
			fmt.Println("\t" + val)
		}
	}
	if len(changed) != 0 {
		fmt.Println("Modified files :")
		for _, val := range changed {
			fmt.Println("\t" + val)
		}
	}
	return nil
}

func (s *Status) status() ([]string, []string, error) {
	dir := "./"
	if s.C != "" {
		dir = s.C
	}
	if !Exists(dir + "/.wbuild/state") {
		return nil, nil, errors.New("'" + dir + "' : Not a wbuild repo.")
	}
	contents, err := ioutil.ReadFile(dir + "/.wbuild/state")
	if err != nil {
		return nil, nil, err
	}
	config := tparse.NewDict()
	if err = config.Parse(string(contents)); err != nil {
		return nil, nil, err
	}

	files, err := config.Find("files")
	if err != nil {
		return nil, nil, err
	}
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}

	untracked := make([]string, 0)
	changed := make([]string, 0)
	for _, val := range infos {
		if val.IsDir() {
			continue
		}
		sum, ok := files[val.Name()]
		if !ok {
			untracked = append(untracked, val.Name())
			continue
		}
		actualSum, err := getMD5Sum(dir + "/" + val.Name())
		if err != nil {
			return nil, nil, err
		}
		if actualSum != sum {
			changed = append(changed, val.Name())
		}
	}
	return untracked, changed, nil
}

func getMD5Sum(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
