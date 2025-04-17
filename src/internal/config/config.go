package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	RemoveComments   bool
	IndentSize       int
	DryRun           bool
	Verbose          bool
	Backup           bool
	Concurrent       bool
	MaxWorkers       int
	Extensions       []string
	PreserveNewlines bool
}

func ParseFlags() *Config {
	config := &Config{}
	flag.IntVar(&config.IndentSize, "indent", 2, "Number of spaces for indentation")
	flag.BoolVar(&config.DryRun, "dry-run", false, "Show what would be done without making changes")
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose logging")
	flag.BoolVar(&config.RemoveComments, "removecomments", false, "Remove comments from the configuration file")
	flag.BoolVar(&config.Backup, "backup", false, "Create backup files before modifying")
	flag.BoolVar(&config.Concurrent, "concurrent", true, "Process files concurrently")
	flag.IntVar(&config.MaxWorkers, "workers", 4, "Number of concurrent workers")
	flag.BoolVar(&config.PreserveNewlines, "preserve-newlines", false, "Preserve existing newlines between blocks")

	extensions := flag.String("extensions", ".conf,.proxy", "Comma-separated list of file extensions to process")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: gofmtnginx [flags] <directory>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	config.Extensions = strings.Split(*extensions, ",")
	for i, ext := range config.Extensions {
		if !strings.HasPrefix(ext, ".") {
			config.Extensions[i] = "." + ext
		}
	}

	return config
}
