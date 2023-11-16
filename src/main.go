package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	removeComment := flag.Bool("removecomments", false, "Remove comments from the configuration file")
	indent := flag.Int("indent", 2, "Number of spaces for indentation")
	flag.Parse()

	// Check if a directory is provided
	if len(flag.Args()) < 1 {
		fmt.Println("Usage: gofmtnginx [flags] <directory>")
		flag.PrintDefaults()
		return
	}

	directory := flag.Arg(0)

	fmt.Printf("Directory: %s\n", directory)
	fmt.Printf("Remove comments: %t\n", *removeComment)
	fmt.Printf("Indentation: %d\n", *indent)

	indentSize := strings.Repeat(" ", *indent)
	removeComments := *removeComment

	// recursive
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() && info.Name() != ".git" {

			ext := filepath.Ext(path)

			if ext == ".conf" || ext == ".proxy" {
				formatConfig(path, indentSize, removeComments)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}
}

func formatConfig(fileName string, indentSize string, removeComments bool) {
	inputFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	indentLevel := 0
	indent := indentSize // 2 spaces
	var formattedLines []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Remove line comments
		if removeComments {
			if strings.HasPrefix(line, "#") {
				continue
			}

			if strings.Contains(line, "#") {
				before, _, found := strings.Cut(line, "#")
				if found && len(before) > 0 {
					line = before
				} else {
					continue
				}
			}
		}

		if strings.HasSuffix(line, "{") {
			formattedLine := strings.TrimRight(fmt.Sprintf("%s%s", strings.Repeat(indent, indentLevel), line), " ")
			formattedLines = append(formattedLines, formattedLine)
			indentLevel++
		} else if strings.HasPrefix(line, "}") {
			if indentLevel > 0 {
				indentLevel--
			}
			formattedLine := strings.TrimRight(fmt.Sprintf("%s%s", strings.Repeat(indent, indentLevel), line), " ")
			formattedLines = append(formattedLines, formattedLine)
		} else {
			formattedLine := strings.TrimRight(fmt.Sprintf("%s%s", strings.Repeat(indent, indentLevel), line), " ")
			formattedLines = append(formattedLines, formattedLine)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input file:", err)
	}

	// Reopen the input file for writing
	inputFile, err = os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error reopening input file for writing:", err)
		return
	}
	defer inputFile.Close()

	writer := bufio.NewWriter(inputFile)

	for _, formattedLine := range formattedLines {
		_, _ = writer.WriteString(formattedLine + "\n")
	}

	_ = writer.Flush()
}
