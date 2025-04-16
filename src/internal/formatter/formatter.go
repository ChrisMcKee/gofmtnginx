package formatter

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ChrisMcKee/gofmtnginx/internal/config"
	"github.com/ChrisMcKee/gofmtnginx/internal/stats"
	"github.com/ChrisMcKee/gofmtnginx/pkg/nginx"
)

// Formatter handles the formatting process
type Formatter struct {
	config *config.Config
	stats  *stats.Stats
	nginx  *nginx.Formatter
}

// New creates a new Formatter instance
func New(cfg *config.Config) *Formatter {
	return &Formatter{
		config: cfg,
		stats:  stats.New(),
		nginx:  nginx.New(cfg.IndentSize, cfg.RemoveComments, cfg.PreserveNewlines),
	}
}

// ProcessDirectory processes all nginx configuration files in a directory
func (f *Formatter) ProcessDirectory(directory string) error {
	if f.config.Concurrent {
		return f.processDirectoryConcurrent(directory)
	}
	return f.processDirectorySequential(directory)
}

// processDirectorySequential processes files sequentially
func (f *Formatter) processDirectorySequential(directory string) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() && info.Name() != ".git" {
			if f.shouldProcessFile(path) {
				if err := f.processFile(path); err != nil {
					log.Printf("Error processing file %s: %v\n", path, err)
				}
			} else {
				f.stats.IncrementSkipped()
				if f.config.Verbose {
					log.Printf("Skipping non-nginx file: %s\n", path)
				}
			}
		}
		return nil
	})
}

// processDirectoryConcurrent processes files concurrently
func (f *Formatter) processDirectoryConcurrent(directory string) error {
	files := make(chan string, 100)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < f.config.MaxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range files {
				if f.shouldProcessFile(file) {
					if err := f.processFile(file); err != nil {
						log.Printf("Error processing file %s: %v\n", file, err)
					}
				} else {
					f.stats.IncrementSkipped()
					if f.config.Verbose {
						log.Printf("Skipping non-nginx file: %s\n", file)
					}
				}
			}
		}()
	}

	// Walk directory and send files to channel
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() && info.Name() != ".git" {
			files <- path
		}
		return nil
	})

	close(files)
	wg.Wait()

	return err
}

// shouldProcessFile checks if a file should be processed
func (f *Formatter) shouldProcessFile(path string) bool {
	ext := filepath.Ext(path)
	for _, validExt := range f.config.Extensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

// processFile processes a single file
func (f *Formatter) processFile(fileName string) error {
	if f.config.Verbose {
		log.Printf("Processing file: %s\n", fileName)
	}

	// Read and format the file first
	formatted, err := f.nginx.FormatFile(fileName)
	if err != nil {
		f.stats.IncrementFailed()
		return err
	}

	if !f.config.DryRun {
		// Create backup if enabled
		if f.config.Backup {
			backupFile := fileName + ".bak"
			backupContent := []byte(strings.Join(formatted, "\n") + "\n")
			if err := os.WriteFile(backupFile, backupContent, 0o644); err != nil {
				f.stats.IncrementFailed()
				return err
			}
		}

		// Write the formatted content
		if err := nginx.WriteFormatted(fileName, formatted); err != nil {
			f.stats.IncrementFailed()
			return err
		}
	}

	f.stats.IncrementProcessed()
	return nil
}

// Stats returns the current statistics
func (f *Formatter) Stats() *stats.Stats {
	return f.stats
}
