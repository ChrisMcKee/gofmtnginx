package nginx

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFormatFile(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expected       []string
		indentSize     int
		removeComments bool
	}{
		{
			name: "basic indentation",
			input: `http {
server {
    listen 80;
}
}`,
			expected: []string{
				"http {",
				"  server {",
				"    listen 80;",
				"  }",
				"}",
			},
			indentSize: 2,
		},
		{
			name: "with comments",
			input: `# Main server block
server {
    # Listen on port 80
    listen 80;
}`,
			expected: []string{
				"# Main server block",
				"server {",
				"  # Listen on port 80",
				"  listen 80;",
				"}",
			},
			indentSize: 2,
		},
		{
			name: "remove comments",
			input: `# Main server block
server {
    # Listen on port 80
    listen 80;
}`,
			expected: []string{
				"server {",
				"  listen 80;",
				"}",
			},
			indentSize:     2,
			removeComments: true,
		},
		{
			name: "inline comments",
			input: `server {
    listen 80; # HTTP port
    server_name example.com; # Domain name
}`,
			expected: []string{
				"server {",
				"  listen 80; # HTTP port",
				"  server_name example.com; # Domain name",
				"}",
			},
			indentSize: 2,
		},
		{
			name: "remove inline comments",
			input: `server {
    listen 80; # HTTP port
    server_name example.com; # Domain name
}`,
			expected: []string{
				"server {",
				"  listen 80;",
				"  server_name example.com;",
				"}",
			},
			indentSize:     2,
			removeComments: true,
		},
		{
			name: "inline comments",
			input: `server {
    listen 80; # HTTP port
    server_name example.com; # Domain name
}`,
			expected: []string{
				"server {",
				"  listen 80; # HTTP port",
				"  server_name example.com; # Domain name",
				"}",
			},
			indentSize: 2,
		},
		{
			name: "remove inline comments with closing brace",
			input: `server {
    listen 80; # HTTP port
    server_name example.com; # Domain name
	#location / {
	#return 200;
	#}
}`,
			expected: []string{
				"server {",
				"  listen 80; # HTTP port",
				"  server_name example.com; # Domain name",
				"  #location / {",
				"  #return 200;",
				"  #}",
				"}",
			},
			indentSize:     2,
			removeComments: false,
		},
		{
			name: "nested blocks",
			input: `http {
    server {
        location / {
            root /var/www/html;
        }
    }
}`,
			expected: []string{
				"http {",
				"  server {",
				"    location / {",
				"      root /var/www/html;",
				"    }",
				"  }",
				"}",
			},
			indentSize: 2,
		},
		{
			name: "empty lines",
			input: `http {

server {

    listen 80;

}
	}`,
			expected: []string{
				"http {",
				"",
				"  server {",
				"",
				"    listen 80;",
				"",
				"  }",
				"}",
			},
			indentSize: 2,
		},
		{
			name:  "multiple closing braces",
			input: `http{ server{ location / { return 200; }}}`,
			expected: []string{
				"http {",
				"  server {",
				"    location / {",
				"      return 200;",
				"    }",
				"  }",
				"}",
			},
			indentSize: 2,
		},
		{
			name:  "multiple closing braces with content",
			input: `http{ server{ location / { return 200; }} server_name example.com;}`,
			expected: []string{
				"http {",
				"  server {",
				"    location / {",
				"      return 200;",
				"    }",
				"  }",
				"  server_name example.com;",
				"}",
			},
			indentSize: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary input file
			tmpDir := t.TempDir()
			inputFile := filepath.Join(tmpDir, "test.conf")
			if err := os.WriteFile(inputFile, []byte(tt.input), 0o644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Create formatter and format file
			f := New(tt.indentSize, tt.removeComments, true)
			formatted, err := f.FormatFile(inputFile)
			if err != nil {
				t.Fatalf("FormatFile() error = %v", err)
			}

			// Compare results
			if len(formatted) != len(tt.expected) {
				t.Errorf("FormatFile() got %d lines, want %d lines", len(formatted), len(tt.expected))
				t.Logf("Got:\n%s", formatLines(formatted))
				t.Logf("Want:\n%s", formatLines(tt.expected))
				return
			}

			for i, line := range formatted {
				if line != tt.expected[i] {
					t.Errorf("FormatFile() line %d = %q, want %q", i+1, line, tt.expected[i])
					t.Logf("Got:\n%s", formatLines(formatted))
					t.Logf("Want:\n%s", formatLines(tt.expected))
					return
				}
			}
		})
	}
}

func TestWriteFormatted(t *testing.T) {
	tests := []struct {
		name     string
		content  []string
		expected string
	}{
		{
			name: "basic content",
			content: []string{
				"http {",
				"  server {",
				"    listen 80;",
				"  }",
				"}",
			},
			expected: "http {\n  server {\n    listen 80;\n  }\n}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary output file
			tmpDir := t.TempDir()
			outputFile := filepath.Join(tmpDir, "output.conf")

			// Write formatted content
			if err := WriteFormatted(outputFile, tt.content); err != nil {
				t.Fatalf("WriteFormatted() error = %v", err)
			}

			// Read back and compare
			content, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}

			if string(content) != tt.expected {
				t.Errorf("WriteFormatted() got:\n%s\nwant:\n%s", string(content), tt.expected)
			}
		})
	}
}

// Helper function to format lines for test output
func formatLines(lines []string) string {
	var result string
	for _, line := range lines {
		if line == "" {
			result += "<empty>\n"
		} else {
			result += line + "\n"
		}
	}
	return result
}
