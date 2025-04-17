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

		// Remove line comments if requested
		if f.RemoveComments {
			line = removeCommentsConsideringQuotes(line)
			if line == "" {
				continue // Skip lines that become empty after comment removal
			}
		}

		formattedLines = append(formattedLines, f.processLine(line, &indentLevel)...)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input file: %w", err)
	}

	return formattedLines, nil
}

// removeCommentsConsideringQuotes removes comments (#) from a line,
// correctly handling comments inside single or double quotes.
func removeCommentsConsideringQuotes(line string) string {
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
		escaped = false // Reset escaped flag after processing the character
	}

	return strings.TrimSpace(result.String())
}

func (f *Formatter) processLine(line string, indentLevel *int) []string {
	var result []string
	line = strings.TrimSpace(line) // Ensure leading/trailing spaces are removed

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

	// Scan through the line character by character
	for i := 0; i < len(line); i++ {
		r := rune(line[i])

		// Handle escape character
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
				continue // Continue to next char after handling brace
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
				continue // Continue to next char after handling brace
			case ';':
				// Append content before ';', add the line including ';'
				currentLine.WriteRune(r)
				result = append(result, indentStr()+strings.TrimSpace(currentLine.String()))
				currentLine.Reset()
				continue // Continue to next char after handling semicolon
			}
		}

		// Append the current character if it wasn't a special char handled above
		currentLine.WriteRune(r)
		escaped = false // Reset escaped flag after processing the character
	}

	// Add any remaining content on the line (if no ';' or '}' was found)
	if currentLine.Len() > 0 {
		result = append(result, indentStr()+strings.TrimSpace(currentLine.String()))
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
