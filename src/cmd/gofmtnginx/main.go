package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ChrisMcKee/gofmtnginx/internal/config"
	"github.com/ChrisMcKee/gofmtnginx/internal/formatter"
)

func main() {
	cfg := config.ParseFlags()
	setupLogging(cfg.Verbose)

	f := formatter.New(cfg)
	if err := f.ProcessDirectory(flag.Arg(0)); err != nil {
		log.Printf("Error processing directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(f.Stats())
}

func setupLogging(verbose bool) {
	if verbose {
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	} else {
		log.SetFlags(log.LstdFlags)
	}
}
