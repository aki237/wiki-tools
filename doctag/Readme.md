# doctag - Tagging tool for any generalised document

*based on semantic comparison weightage or words*

This tools is used to generate tags for a given document

## Building

This is a simple golang source. So just, a
`go build` will do.
Before the `build` step get the dependencies.

```shell
$ go get -u golang.org/x/net/html
```

Then run `go build` to get the tool's binary.

## Usage

A sample usage is given below.

```shell
$ ls
doc.md common.txt
$ doctag
Usage: wctag -o [output] <file>
  -common string
    	file which contains the less significant words
  -o string
    	file in which the output is to be dumped (default "TAGS")
  -omit string
    	extra omit words separated by a space
  <file> strings
  	file to be tagged
$ doctag -o TAGS -common common.txt doc.md
$ ls
doc.md common.txt TAGS
```
