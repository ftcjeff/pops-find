package main

import (
  "bufio"
	"fmt"
  "log"
	"os"
	"strings"

	"io/ioutil"
)

func isIgnored(name string) bool {
	dirsToIgnore := []string{"Godeps", ".svn", ".git"}

  for _, dir := range dirsToIgnore {
    if name == dir {
  	  return true
    }
  }

  return false
}

func processFile(path string, patterns []string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	headerPrinted := false

	lineNum := 1

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, pattern := range patterns {
			if strings.Contains(strings.ToLower(scanner.Text()), strings.ToLower(pattern)) {
				if !headerPrinted {
				  fmt.Printf("##### %s\n", path)
					headerPrinted = true
				}

				fmt.Printf("%5d: %s\n", lineNum, scanner.Text())
				continue
			}
		}
		lineNum++
	}

	if headerPrinted {
		fmt.Printf("\n")
	}
}

func processDirectory(path string, fileType string, patterns []string) {
  files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if f.IsDir() {
			if isIgnored(f.Name()) {
				continue
			}

			dirName := path + "/" + f.Name()
			processDirectory(dirName, fileType, patterns)
		}

		if !strings.HasSuffix(f.Name(), fileType) {
			continue
		}

		fileName := path + "/" + f.Name()
		processFile(fileName, patterns)
	}
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("usage: pops-find <filetype> string1 [string2...stringn]")
	}

	fileType := "." + os.Args[1]
	patterns := os.Args[2:]

	log.Printf("### Searching '%s' files for: %s\n", fileType, patterns)

	processDirectory(".", fileType, patterns)
}
