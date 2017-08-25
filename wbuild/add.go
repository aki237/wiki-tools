package main

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/shivam07a/tparse"
)

type Add struct {
	C        string
	RestArgs []string
}

func (a *Add) Exec() error {
	if len(a.RestArgs) < 1 {
		return errors.New("add requires one or more files as arguments.")
	}
	dir := "./"
	if a.C != "" {
		dir = a.C
	}
	if !Exists(dir + "/.wbuild/state") {
		return errors.New("'" + dir + "' : Not a wbuild repo.")
	}
	contents, err := ioutil.ReadFile(dir + "/.wbuild/state")
	if err != nil {
		return err
	}
	config := tparse.NewDict()
	if err = config.Parse(string(contents)); err != nil {
		return err
	}

	files, err := config.Find("files")
	if err != nil {
		return err
	}

	for _, val := range a.RestArgs {
		if !Exists(dir + "/" + val) {
			return errors.New("The specified file doesn't exist : '" + val + "'")
		}
		sum, err := getMD5Sum(dir + "/" + val)
		if err != nil {
			return err
		}
		files[val] = sum
	}

	configFile, err := os.OpenFile(dir+"/.wbuild/state", os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	if err := config.UnMarshal(configFile); err != nil {
		return err
	}
	return configFile.Close()
}
