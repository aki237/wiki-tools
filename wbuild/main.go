package main

import (
	"fmt"
	"os"

	"github.com/aki237/clc"
)

func Exists(dir string) bool {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return true
	}
	return false
}

func main() {
	var initArgs = Init{C: "./", OmitFile: false}
	var addArgs = Add{C: "./"}
	var buildArgs = Build{C: "./"}
	var statusArgs Status
	app := clc.NewApp("wbuild", "wiki page build tool", "v0.0.1")
	app.AddCommand("init", "initialize a wbuild directory with all skeleton files.", &initArgs)
	app.AddCommand("status", "track all the changes in a wiki project", &statusArgs)
	app.AddCommand("add", "add a new file or a modified file to the wbuild project", &addArgs)
	app.AddCommand("build", "build the project into a dibba file", &buildArgs)
	if err := app.Run(); err != nil {
		fmt.Fprint(os.Stderr, err, "\n")
	}
}
