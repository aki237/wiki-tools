package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/shivam07a/tparse"
)

type Build struct {
	O string
	C string
}

var j = filepath.Join

func (b *Build) Exec() error {
	dir := "./"
	if b.C != "" {
		dir = b.C
	} else {
		b.C = dir
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

	sfiles, err := config.Find("files")
	if err != nil {
		return err
	}

	s := Status{C: dir}
	_, changed, err := s.status()
	if err != nil {
		return err
	}

	if len(changed) != 0 {
		return errors.New("Cannot build when there are pending changes to be added.")
	}

	if b.O == "" {
		d, err := os.Getwd()
		if err != nil {
			return err
		}
		b.O = filepath.Base(d)
	}

	if !executableFound("dib") {
		return errors.New("It seems, 'dib' is not installed in your system. Either that or it is not in $PATH. Make sure it is and try again.")
	}

	files := make([]string, 0)
	files = append(files, "-o", b.O)

	if !Exists(j(dir, "doc.md")) {
		fmt.Println("Warning : 'doc.md' file not found. Skipping Tag generation")
	} else {
		ent, err := config.Find("config")
		if err != nil {
			fmt.Println("Warning : common tags file not found. Skipping Tag generation")
		} else {
			if common, ok := ent["common"]; ok {
				tagFile, err := b.generateTags(common, j(dir, "doc.md"))
				if err != nil {
					fmt.Println("Warning: ", err, ", Skipping Tag Generation")
				} else {
					files = append(files, tagFile)
				}
			} else {
				fmt.Println("Warning : common tags file not found. Skipping Tag generation")
			}
		}
	}

	for key, _ := range sfiles {
		if key != ".omit" {
			files = append(files, j(dir, key))
		}
	}
	return exec.Command("dib", files...).Run()
}

func (b *Build) generateTags(commonFile string, file string) (string, error) {
	if !executableFound("doctag") {
		return "", errors.New("It seems, 'doctag' is not installed in your system. Either that or it is not in $PATH. Make sure it is and try again.")
	}
	omitString := ""
	rand.Seed(time.Now().UnixNano())
	tmpdir := fmt.Sprintf("%s-%d", b.C, rand.Int63())
	if err := os.MkdirAll("/tmp/"+tmpdir, 0755); err != nil {
		return "", err
	}
	if Exists(j(b.C, ".omit")) {
		c, err := ioutil.ReadFile(j(b.C, ".omit"))
		if err != nil {
			return "", err
		}
		omitString = strings.Replace(strings.TrimSpace(string(c)), "\n", " ", -1)
	}
	args := []string{
		"-common", commonFile,
		"-o", "/tmp/" + tmpdir + "/TAGS",
	}
	if omitString != "" {
		args = append(args, "-omit", omitString)
	}
	args = append(args, file)
	err := exec.Command("doctag", args...).Run()
	if err != nil {
		return "", err
	}
	return "/tmp/" + tmpdir + "/TAGS", nil
}

func executableFound(execName string) bool {
	paths := strings.Split(os.Getenv("PATH"), ":")
	exists := false
	for _, val := range paths {
		if !Exists(j(val, execName)) {
			continue
		}
		exists = true
		break
	}
	return exists
}
