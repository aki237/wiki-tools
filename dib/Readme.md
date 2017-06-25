# dib - Dibba packager tool

This tools is used to package files into a 
single dibba package.

## Building

This is a simple golang source. So just, a
`go build` will do.
Before the `build` step get the dependencies.

```shell
$ go get -u github.com/aki237/dibba
```

Then run `go build` to get the tool's binary.

## Usage

A sample usage is given below.

```shell
$ ls
a b c
$ dib
Usage: dib -o [output] <files>
  -o string
    	name of the dibba output file
  <files> strings
  	list of all files to be packaged
$ dib -o files.dib a b c
$ ls
a b c files.dib
```
