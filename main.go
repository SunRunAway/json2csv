package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/SunRunAway/json2csv/json2map"
	"github.com/qiniu/log"
)

var (
	inputFile   = flag.String("i", "", "/path/to/input.json (optional; default is stdin)")
	outputFile  = flag.String("o", "", "/path/to/output.json (optional; default is stdout)")
	outputDelim = flag.String("d", ",", "delimiter used for output values")
	printHeader = flag.Bool("p", true, "prints header to output")
)

func main() {
	flag.Parse()

	r := os.Stdin
	if *inputFile != "" {
		file, err := os.OpenFile(*inputFile, os.O_RDONLY, 0600)
		if err != nil {
			log.Fatalf("Error %s opening input file %v\n", err, *inputFile)
		}
		defer file.Close()
		r = file
	}

	w := csv.NewWriter(os.Stdout)
	if *outputFile != "" {
		file, err := os.OpenFile(*outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Fatalf("Error %s opening output file %v\n", err, *outputFile)
		}
		defer file.Close()
		w = csv.NewWriter(file)
	}

	delim, _ := utf8.DecodeRuneInString(*outputDelim)
	w.Comma = delim

	var collectedMaps []map[string]interface{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		b := scanner.Bytes()
		if len(b) == 0 {
			continue
		}
		m, err := json2map.New().Convert(bytes.NewReader(b))
		if err != nil {
			log.Warn("convert error:", err)
		}

		collectedMaps = append(collectedMaps, m)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln("reading standard input:", err)
	}

	err := maps2csv(w, collectedMaps)
	if err != nil {
		log.Fatalln("maps2csv error:", err)
	}
}

func maps2csv(w *csv.Writer, maps []map[string]interface{}) (err error) {
	if len(maps) == 0 {
		return
	}

	keyVisit := make(map[string]bool, len(maps[0]))
	header := make([]string, 0, len(maps[0]))
	for _, m := range maps {
		for k := range m {
			if keyVisit[k] {
				continue
			}
			keyVisit[k] = true
			header = append(header, k)
		}
	}
	sort.Strings(header)

	if *printHeader {
		// remove leading dots introducing by pakcage json2map
		header2 := make([]string, len(header))
		for i, str := range header {
			if strings.HasPrefix(str, ".") {
				str = str[1:]
			}
			header2[i] = str
		}
		w.Write(header2)
		w.Flush()
		if err = w.Error(); err != nil {
			return
		}
	}

	for _, m := range maps {
		record := make([]string, len(header))
		for i, k := range header {
			if v, ok := m[k]; ok {
				record[i] = fmt.Sprintf("%+v", v)
			}
		}
		w.Write(record)
		w.Flush()
		if err = w.Error(); err != nil {
			return
		}
	}
	return
}
