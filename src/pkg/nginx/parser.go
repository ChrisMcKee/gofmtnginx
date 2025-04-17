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
		originalLine := scanner.Text() // Keep original line to check for original emptiness
		line := strings.TrimSpace(originalLine)

		// Handle lines that are purely comments first, if not removing comments
		if !f.RemoveComments && strings.HasPrefix(line, "#") {
			indentStr := strings.Repeat(" ", indentLevel*f.IndentSize)
			formattedLines = append(formattedLines, indentStr+line)
			continue
		}

		// Remove comments if requested (handles comments starting mid-line)
		if f.RemoveComments {
			line = removeComments(line)
		}

		if line == "" {
			// Only preserve newline if requested AND it wasn't caused by comment removal
			// OR if the original line was only whitespace
			if f.PreserveNewlines && (!f.RemoveComments || strings.TrimSpace(originalLine) == "") {
				formattedLines = append(formattedLines, "")
			}
			continue // Skip processing for effectively empty lines
		}

		// Process non-empty, non-comment lines
		processed := f.processLine(line, &indentLevel)
		formattedLines = append(formattedLines, processed...)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input file: %w", err)
	}

	return formattedLines, nil
}

// removeComments removes comments (#) from a line
func removeComments(line string) string {
	var result strings.Builder
	inSingleQuotes := false
	inDoubleQuotes := false
	escaped := false

	for _, r := range line {
		if r == '\\' && !escaped {
			escaped = true
			result.WriteRune(r)
			continue
		}

		if r == '\'' && !escaped {
			inSingleQuotes = !inSingleQuotes
		} else if r == '"' && !escaped {
			inDoubleQuotes = !inDoubleQuotes
		} else if r == '#' && !inSingleQuotes && !inDoubleQuotes && !escaped {
			// Found a comment marker outside quotes, stop processing
			break
		}

		result.WriteRune(r)
		escaped = false
	}

	return strings.TrimSpace(result.String())
}

func (f *Formatter) processLine(line string, indentLevel *int) []string {
	var result []string
	line = strings.TrimSpace(line)

	// Handle empty lines after potential comment removal
	if line == "" {
		if f.PreserveNewlines {
			return []string{""}
		}
		return nil
	}

	var currentLine strings.Builder
	indentStr := func() string {
		return strings.Repeat(" ", *indentLevel*f.IndentSize)
	}

	inSingleQuotes := false
	inDoubleQuotes := false
	escaped := false

	for _, r := range line {
		if r == '\\' && !escaped {
			escaped = true
			currentLine.WriteRune(r)
			continue
		}

		// Track quote status
		if r == '\'' && !escaped {
			inSingleQuotes = !inSingleQuotes
		} else if r == '"' && !escaped {
			inDoubleQuotes = !inDoubleQuotes
		}

		// Process special characters only if not inside quotes
		if !inSingleQuotes && !inDoubleQuotes {
			switch r {
			case '{':
				// Append content before '{', add the line, then start a new line for '{'
				if currentLine.Len() > 0 {
					result = append(result, indentStr()+strings.TrimSpace(currentLine.String())+" {")
					currentLine.Reset()
				} else {
					// If '{' is the first non-space char, or follows a '}' on the same line
					result = append(result, indentStr()+"{")
				}
				*indentLevel++
				continue
			case '}':
				// Append content before '}', add the line, then start a new line for '}'
				if currentLine.Len() > 0 {
					result = append(result, indentStr()+strings.TrimSpace(currentLine.String()))
					currentLine.Reset()
				}
				if *indentLevel > 0 {
					*indentLevel--
				}
				result = append(result, indentStr()+"}")
				continue
			case ';':
				// Semicolon marks the end of a directive, but there might be inline comments.
				currentLine.WriteRune(r)
				continue
			}
		}

		currentLine.WriteRune(r)
		escaped = false
	}

	// Add any remaining content on the line
	if currentLine.Len() > 0 {
		finalLine := strings.TrimRight(currentLine.String(), " ")
		if finalLine != "" {
			result = append(result, indentStr()+finalLine)
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
