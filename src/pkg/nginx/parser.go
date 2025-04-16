package nginx

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Formatter struct {
	IndentSize       int
	RemoveComments   bool
	PreserveNewlines bool
}

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

		formattedLines = append(formattedLines, f.processLine(line, &indentLevel)...)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input file: %w", err)
	}

	return formattedLines, nil
}

func (f *Formatter) processLine(line string, indentLevel *int) []string {
	var result []string

	// Handle empty lines
	if line == "" {
		if f.PreserveNewlines {
			return []string{""}
		}
		return nil
	}

	parts := strings.Split(line, "{")
	for i, part := range parts {
		if i > 0 {
			// For parts after the first one, we need to add the opening brace to the previous line
			if len(result) > 0 {
				result[len(result)-1] = result[len(result)-1] + " {"
			} else {
				result = append(result, strings.TrimRight(fmt.Sprintf("%s{", strings.Repeat(" ", *indentLevel*f.IndentSize)), " "))
			}
			*indentLevel++
		}

		// Process each part for closing braces in case of multiple closing braces
		if strings.TrimSpace(part) != "" {
			closeParts := strings.Split(part, "}")

			// Process the content before any closing braces
			if strings.TrimSpace(closeParts[0]) != "" {
				formattedLine := strings.TrimRight(fmt.Sprintf("%s%s", strings.Repeat(" ", *indentLevel*f.IndentSize), strings.TrimSpace(closeParts[0])), " ")
				result = append(result, formattedLine)
			}

			// Process each closing brace
			for j := 1; j < len(closeParts); j++ {
				if *indentLevel > 0 {
					*indentLevel--
				}
				formattedLine := strings.TrimRight(fmt.Sprintf("%s}", strings.Repeat(" ", *indentLevel*f.IndentSize)), " ")
				result = append(result, formattedLine)

				// Process any content after the closing brace
				if strings.TrimSpace(closeParts[j]) != "" {
					formattedLine := strings.TrimRight(fmt.Sprintf("%s%s", strings.Repeat(" ", *indentLevel*f.IndentSize), strings.TrimSpace(closeParts[j])), " ")
					result = append(result, formattedLine)
				}
			}
		}
	}

	return result
}

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
