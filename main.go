package main

import (
	"bufio"
	bytes2 "bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var todos map[string][]string

func walkDir(path string, dirEntry fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if dirEntry.IsDir() {
		return nil
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	reader := bytes2.NewReader(bytes)
	scanner := bufio.NewScanner(reader)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for i := range lines {
		line := lines[i]
		if strings.Contains(line, "TODO:") {
			fileTodos := todos[path]
			if fileTodos == nil {
				fileTodos = []string{}
			}

			fileTodos = append(fileTodos, line)
			todos[path] = fileTodos
		}
	}

	return nil
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Not enough arguments!")
		fmt.Printf("Usage: %s <directory>", args[0])
		return
	}

	dir := args[1]
	dirInfo, err := os.Stat(dir)
	if err != nil {
		fmt.Printf("Error finding directory %s!\n", dir)
		fmt.Println(err)
		return
	}

	if !dirInfo.IsDir() {
		fmt.Printf("Path \"%s\" is not a directory!\n", dir)
		return
	}

	todos = make(map[string][]string)

	err = filepath.WalkDir(dir, walkDir)
	if err != nil {
		fmt.Printf("Error walking directory \"%s\"!\n", dir)
		fmt.Println(err)
		return
	}

	for s := range todos {
		fileTodos := todos[s]
		fmt.Println(fileTodos)
	}
}
