package main

import (
	"io/ioutil"
	"os"
	"testing"
)

// Helper function to create a temporary file with initial content
func createTempFile(t *testing.T, content string) string {
	t.Helper()
	tmpFile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpFile.Name()
}

// TestFormatConfigBalancedBraces tests formatting of a file with balanced braces
func TestFormatConfigBalancedBraces(t *testing.T) {
	content := `server {
    listen 80;
    server_name example.com;
}`
	fileName := createTempFile(t, content)
	defer os.Remove(fileName)

	formatConfig(fileName, "  ", false)

	result, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("Failed to read result file: %v", err)
	}

	expected := `server {
  listen 80;
  server_name example.com;
}
`
	if string(result) != expected {
		t.Errorf("Formatted content does not match expected content.\nExpected:\n%s\nGot:\n%s", expected, result)
	}
}

// TestFormatConfigUnbalancedBraces tests formatting of a file with more closing braces than opening ones
func TestFormatConfigUnbalancedBraces(t *testing.T) {
	content := `
server {
location / {
}
}
}
`
	fileName := createTempFile(t, content)
	defer os.Remove(fileName)

	formatConfig(fileName, "  ", false)

	result, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("Failed to read result file: %v", err)
	}

	expected := `
server {
  location / {
  }
}
}
`

	if string(result) != expected {
		t.Errorf("Formatted content does not match expected content.\nExpected:\n%s\nGot:\n%s", expected, result)
	}
}

// TestFormatConfigWithComments tests the handling of comments in the file
func TestFormatConfigWithComments(t *testing.T) {
	content := `
# This is a comment
server {
    # Another comment
    listen 80; # Inline comment
}
`
	fileName := createTempFile(t, content)
	defer os.Remove(fileName)

	formatConfig(fileName, "  ", true) // Assuming true removes comments

	result, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("Failed to read result file: %v", err)
	}

	expected := `
server {
  listen 80;
}
`
	if string(result) != expected {
		t.Errorf("Formatted content does not match expected content.\nExpected:\n%s\nGot:\n%s", expected, result)
	}
}

// TestFormatConfigEmptyFile tests formatting of an empty file
func TestFormatConfigEmptyFile(t *testing.T) {
	fileName := createTempFile(t, "")
	defer os.Remove(fileName)

	formatConfig(fileName, "  ", false)

	result, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("Failed to read result file: %v", err)
	}

	expected := ``
	if string(result) != expected {
		t.Errorf("Formatted content does not match expected content.\nExpected:\n%s\nGot:\n%s", expected, result)
	}
}

// TestFormatConfigLargeIndent tests formatting with a large indent level
func TestFormatConfigLargeIndent(t *testing.T) {
	content := `
server {
  listen 80;
   location / {
   }
}
`
	fileName := createTempFile(t, content)
	defer os.Remove(fileName)

	formatConfig(fileName, "        ", false) // Large indent

	result, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("Failed to read result file: %v", err)
	}

	expected := `
server {
        listen 80;
        location / {
        }
}
`
	if string(result) != expected {
		t.Errorf("Formatted content does not match expected content.\nExpected:\n%s\nGot:\n%s", expected, result)
	}
}
