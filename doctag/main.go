package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

func wordWeights(words []string, r io.Reader) (map[int][]string, []int, int) {
	m := make(map[string]int, 0)
	w := bytes.NewBuffer(nil)
	io.Copy(w, r)
	str := string(w.Bytes())
	current := ""
	for _, val := range str {
		switch {
		case (val >= 65 && val <= 90) || (val >= 97 && val <= 122):
			current += string(val)
		default:
			if current != "" && !wordIn(words, strings.ToLower(current)) {
				m[strings.ToLower(current)] += 1
			}
			current = ""
		}
	}
	n := map[int][]string{}
	var a []int
	for k, v := range m {
		n[v] = append(n[v], k)
	}
	for k := range n {
		a = append(a, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	return n, a, len(m)
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s -o [output] <file>\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "  <file> strings\n  \tfile to be tagged\n")
}

func main() {
	out := flag.String("o", "TAGS", "file in which the output is to be dumped")
	common := flag.String("common", "", "file which contains the less significant words")
	omit := flag.String("omit", "", "extra omit words separated by a space")

	flag.Parse()

	if *common == "" || len(flag.Args()) != 1 {
		Usage()
		return
	}

	outFile, err := os.OpenFile(*out, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	b, err := ioutil.ReadFile(*common)
	if err != nil {
		fmt.Println(err)
		return
	}
	words := strings.Split(string(b), "\n")
	words = append(words, strings.Split(strings.TrimSpace(*omit), " ")...)
	// Parse and weigh the local file.
	f, err := os.Open(flag.Args()[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	lPos, la, lwc := wordWeights(words, f)
	lMap := getMap(lPos, la)
	//	fmt.Println(lMap)
	if len(la) < 2 {
		fmt.Println("Cannot tag with less results")
		return
	}
	strList := make([]string, 0)
	for _, val := range la {
		strList = append(strList, lPos[val]...)
		if len(strList) > 5 {
			strList = strList[:5]
			break
		}
	}
	rrd, err := getWikipediaReader(strList)
	if err != nil {
		fmt.Println(err)
		return
	}
	// get the remote tags.
	rPos, ra, rwc := wordWeights(words, rrd)
	rMap := getMap(rPos, ra)
	vals := make(map[float64]string, 0)
	nums := make([]float64, 0)
	for term, reps := range rMap {
		rReps, ok := lMap[term]
		if !ok {
			vals[float64(reps)/float64(rwc)] = term
			continue
		}
		weightage := (float64(reps) / float64(rwc)) + (float64(rReps) / float64(lwc))
		vals[weightage] = term
	}
	for k := range vals {
		nums = append(nums, k)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(nums)))
	for i, reps := range nums {
		fmt.Fprint(outFile, vals[reps])
		if i == 10 {
			break
		}
		fmt.Fprint(outFile, " ")
	}
	fmt.Fprint(outFile, "\n")
	outFile.Close()
}

func wordIn(slist []string, word string) bool {
	for _, val := range slist {
		if strings.Contains(val, word) {
			return true
		}
	}
	return false
}

func getMap(Pos map[int][]string, a []int) map[string]int {
	Map := make(map[string]int, 0)
	i := 0
	for _, k := range a {
		for _, s := range Pos[k] {
			Map[s] = k
			i += 1
			if i >= 50 {
				break
			}
		}
		if i >= 50 {
			break
		}
	}
	return Map
}

func removeTags(input, startTag, endTag string) string {
	start := strings.Index(input, startTag)
	for start != -1 {
		end := strings.Index(input, endTag)
		if end == -1 {
			input = input[start:]
			break
		}
		if start > end {
			input = input[:end] + input[start+len(startTag):]
		} else {
			input = input[:start] + input[end+len(endTag):]
		}
		start = strings.Index(input, startTag)
	}
	return input
}
