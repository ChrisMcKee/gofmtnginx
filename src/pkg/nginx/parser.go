package nginx

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Formatter handles nginx configuration formatting
type Formatter struct {
	IndentSize       int
	RemoveComments   bool
	PreserveNewlines bool
}

// New creates a new Formatter instance
func New(indentSize int, removeComments bool, preserveNewlines bool) *Formatter {
	return &Formatter{
		IndentSize:       indentSize,
		RemoveComments:   removeComments,
		PreserveNewlines: preserveNewlines,
	}
}

// FormatFile formats a nginx configuration file
func (f *Formatter) FormatFile(fileName string) ([]string, error) {
	inputFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening input file: %w", err)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	indentLevel := 0
	var formattedLines []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Remove line comments
		if f.RemoveComments {
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
			formattedLine := strings.TrimRight(fmt.Sprintf("%s%s", strings.Repeat(" ", indentLevel*f.IndentSize), line), " ")
			formattedLines = append(formattedLines, formattedLine)
			indentLevel++
		} else if strings.HasPrefix(line, "}") {
			if indentLevel > 0 {
				indentLevel--
			}
			formattedLine := strings.TrimRight(fmt.Sprintf("%s%s", strings.Repeat(" ", indentLevel*f.IndentSize), line), " ")
			formattedLines = append(formattedLines, formattedLine)
		} else {
			formattedLine := strings.TrimRight(fmt.Sprintf("%s%s", strings.Repeat(" ", indentLevel*f.IndentSize), line), " ")
			formattedLines = append(formattedLines, formattedLine)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input file: %w", err)
	}

	return formattedLines, nil
}

// WriteFormatted writes formatted lines to a file
func WriteFormatted(fileName string, formatted []string) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("error opening file for writing: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range formatted {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("error writing line: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	return nil
}
