package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

var (
	// This regex covers a wide range of common emojis.
	emojiRegex = regexp.MustCompile(`[\x{1F600}-\x{1F64F}\x{1F300}-\x{1F5FF}\x{1F680}-\x{1F6FF}\x{1F1E0}-\x{1F1FF}]`)

	// Command-line flags
	dryRun = flag.Bool("dry-run", false, "Show which files would be modified, without changing them")
	check  = flag.Bool("check", false, "Exit with non-zero status if emojis are found (for CI)")

	// Global state for check mode
	emojisFound = false
)

func removeEmojis(content []byte) []byte {
	return emojiRegex.ReplaceAll(content, []byte(""))
}

func processFile(path string, info fs.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("Error accessing path %s: %v\n", path, err)
		return nil // Continue walking
	}

	// Skip .git directory
	if info.IsDir() && info.Name() == ".git" {
		return filepath.SkipDir
	}

	// Skip other directories
	if info.IsDir() {
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", path, err)
		return nil // Continue walking
	}

	if !emojiRegex.Match(content) {
		return nil // No emojis found, nothing to do
	}

	if *check {
		emojisFound = true
		fmt.Printf("[check] Emojis found in: %s\n", path)
		return nil // Continue checking other files
	}

	if *dryRun {
		fmt.Printf("[dry-run] Emojis found in: %s\n", path)
		return nil // Don't modify the file
	}

	// If not check or dry-run, remove emojis and write back to the file
	newContent := removeEmojis(content)
	if err := os.WriteFile(path, newContent, info.Mode()); err != nil {
		fmt.Printf("Error writing to file %s: %v\n", path, err)
	} else {
		fmt.Printf("Removed emojis from: %s\n", path)
	}

	return nil
}

func main() {
	flag.Parse()
	paths := flag.Args()

	if len(paths) == 0 {
		fmt.Println("Usage: emoji-remover [flags] <file-or-directory1> [<file-or-directory2> ...]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	for _, path := range paths {
		if err := filepath.Walk(path, processFile); err != nil {
			fmt.Printf("Error walking path %s: %v\n", path, err)
		}
	}

	if *check && emojisFound {
		fmt.Println("\nCheck failed: Emojis were found in the files listed above.")
		os.Exit(1)
	}

	if *check {
		fmt.Println("\nCheck passed: No emojis found.")
	}
} 