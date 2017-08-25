# wbuild - git like wiki repo manager tool

This tool is used to maitain, build wiki bundles from a wiki repo

## Building

This is a simple golang source. So just, a
`go build` will do.
Before the `build` step get the dependencies.

```shell
$ go get -u github.com/aki237/dibba
```

Then run `go build` to get the tool's binary.

## Usage

### init

Initialize a wiki wbuild repository

```shell
$ cd SOMEDIR
$ wbuild init
```

### status

Check the status whether some file has been modified or whether a build should be updated

```shell
$ cd $SOME_WBUILD_REPO_DIR
$ wbuild status
In directory $'SOME_WBUILD_REPO_DIR'

Modified files :
	doc.md

```

### add

Add a file to the tracking stream.

```shell
$ cd $SOME_WBUILD_REPO_DIR
$ wbuild add doc.md
```

### build

Build the repo into a wiki dibba package.

```shell
$ cd $SOME_WBUILD_REPO_DIR
$ wbuild build -o $OUTPUT/$SOME_NAME.wdib
```
