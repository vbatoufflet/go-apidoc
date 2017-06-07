package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var (
	outputExt  string
	outputPath string
	pkgPath    string
	tmplPath   string
)

func main() {
	flag.StringVar(&outputExt, "ext", "", "output files extension")
	flag.StringVar(&outputPath, "o", ".", "output directory path")
	flag.StringVar(&tmplPath, "template", "", "path to template file")
	flag.Usage = printUsage
	flag.Parse()

	args := flag.Args()
	if n := len(args); n > 1 {
		printUsage()
		os.Exit(1)
	} else if n == 1 {
		pkgPath = args[0]
	} else {
		pkgPath = "."
	}

	if err := render(); err != nil {
		printError(err)
		os.Exit(1)
	}
}

func render() error {
	// Parse package for documentation data, or stop if none found
	sections, err := parse(pkgPath)
	if err != nil {
		return err
	} else if len(sections.Sections) == 0 {
		fmt.Println("No data found")
		return nil
	}

	// Ensure requested template exists
	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		return fmt.Errorf("unable to use template file: %s", err)
	}

	// Loop through sections to generate output files
	if outputPath != "." {
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			if err := os.MkdirAll(outputPath, 0755); err != nil {
				return errors.Wrap(err, "failed to create output directory")
			}
		}
	}

	for _, section := range sections.Sections {
		fileName := strings.ToLower(section.Name)
		if fileName == "" {
			fileName = "index"
		}
		if outputExt != "" {
			fileName += outputExt
		}

		file, err := os.OpenFile(filepath.Join(outputPath, fileName), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return errors.Wrap(err, "failed to open output file")
		}
		defer file.Close()

		if err := section.Render(tmplPath, file); err != nil {
			printError(err)
			continue
		}
		fmt.Printf("generated %q\n", file.Name())
	}

	return nil
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [PACKAGE]\n\nOptions:\n", filepath.Base(os.Args[0]))
	flag.PrintDefaults()
}

func printError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
}
