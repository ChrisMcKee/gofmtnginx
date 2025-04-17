package formatter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ChrisMcKee/gofmtnginx/internal/config"
)

func TestProcessFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gofmtnginx-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	files := map[string]string{
		"test1.conf": `server {
listen 80;
}`,
		"test2.conf": `# Comment
server {
  listen 443;
}`,
		"test3.txt": `server {
  listen 80;
}`,
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		err := os.WriteFile(path, []byte(content), 0o644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", name, err)
		}
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("Failed to verify test file %s exists: %v", name, err)
		}
	}

	cfg := &config.Config{
		RemoveComments: true,
		IndentSize:     2,
		Backup:         true,
		Extensions:     []string{".conf"},
	}

	f := New(cfg)

	for name := range files {
		path := filepath.Join(tmpDir, name)
		if f.shouldProcessFile(path) {
			if err := f.processFile(path); err != nil {
				t.Errorf("Error processing file %s: %v", path, err)
			}
		} else {
			f.stats.IncrementSkipped()
		}
	}

	stats := f.Stats()
	if stats.FilesProcessed != 2 { // Only .conf files should be processed
		t.Errorf("Expected 2 files processed, got %d", stats.FilesProcessed)
	}
	if stats.FilesSkipped != 1 { // .txt file should be skipped
		t.Errorf("Expected 1 file skipped, got %d", stats.FilesSkipped)
	}
	if stats.FilesFailed != 0 {
		t.Errorf("Expected 0 files failed, got %d", stats.FilesFailed)
	}

	// Verify backup files were created
	for name := range files {
		if strings.HasSuffix(name, ".conf") {
			backupPath := filepath.Join(tmpDir, name+".bak")
			if _, err := os.Stat(backupPath); err != nil {
				t.Errorf("Backup file %s was not created: %v", backupPath, err)
			}
		}
	}
}

func TestProcessDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gofmtnginx-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	for i := 0; i < 10; i++ {
		content := fmt.Sprintf(`server {
  listen 80;
  server_name test%d.com;
}`, i)

		ext := ".conf"
		if i%2 == 0 {
			ext = ".proxy"
		}

		path := filepath.Join(tmpDir, fmt.Sprintf("test%d%s", i, ext))
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("Failed to create test file %d: %v", i, err)
		}

		if _, err := os.Stat(path); err != nil {
			t.Fatalf("Failed to verify test file %d exists: %v", i, err)
		}
	}

	cfg := &config.Config{
		Concurrent: true,
		MaxWorkers: 4,
		Extensions: []string{".conf", ".proxy"},
		IndentSize: 2,
	}

	f := New(cfg)

	if err := f.ProcessDirectory(tmpDir); err != nil {
		t.Errorf("Error processing directory: %v", err)
	}

	stats := f.Stats()
	if stats.FilesProcessed != 10 {
		t.Errorf("Expected 10 files processed, got %d", stats.FilesProcessed)
	}
	if stats.FilesFailed != 0 {
		t.Errorf("Expected 0 files failed, got %d", stats.FilesFailed)
	}
}
