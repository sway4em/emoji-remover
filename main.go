package main

import (
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

var (
	// Comprehensive emoji regex covering all major Unicode emoji ranges.
	emojiRegex = regexp.MustCompile(`[\x{1F600}-\x{1F64F}]|` + // Emoticons
		`[\x{1F300}-\x{1F5FF}]|` + // Misc Symbols and Pictographs
		`[\x{1F680}-\x{1F6FF}]|` + // Transport and Map
		`[\x{1F1E0}-\x{1F1FF}]|` + // Flags
		`[\x{1F900}-\x{1F9FF}]|` + // Supplemental Symbols and Pictographs
		`[\x{1FA00}-\x{1FA6F}]|` + // Chess Symbols
		`[\x{1FA70}-\x{1FAFF}]|` + // Symbols and Pictographs Extended-A
		`[\x{2600}-\x{26FF}]|` + // Misc Symbols
		`[\x{2700}-\x{27BF}]|` + // Dingbats
		`[\x{231A}-\x{23FF}]|` + // Misc Technical
		`[\x{200D}]|` + // Zero Width Joiner
		`[\x{FE0F}]|` + // Variation Selector-16
		`[\x{20E3}]|` + // Combining Enclosing Keycap
		`[\x{2934}-\x{2935}]|` + // Arrows
		`[\x{25AA}-\x{25AB}]|` + // Small squares
		`[\x{25B6}]|` + // Play button
		`[\x{25C0}]|` + // Reverse button
		`[\x{25FB}-\x{25FE}]|` + // Medium squares
		`[\x{2B05}-\x{2B07}]|` + // Arrows
		`[\x{2B1B}-\x{2B1C}]|` + // Large squares
		`[\x{2B50}]|` + // Star
		`[\x{2B55}]|` + // Circle
		`[\x{3030}]|` + // Wavy dash
		`[\x{303D}]|` + // Part alternation mark
		`[\x{3297}]|` + // Circled Ideograph Congratulation
		`[\x{3299}]|` + // Circled Ideograph Secret
		`[\x{1F004}]|` + // Mahjong tile
		`[\x{1F0CF}]|` + // Playing card joker
		`[\x{E0020}-\x{E007F}]`) // Tags

	// Directories to always skip
	skipDirs = map[string]bool{
		".git":         true,
		"node_modules": true,
		"vendor":       true,
		".venv":        true,
		"venv":         true,
		"__pycache__":  true,
		"build":        true,
		"dist":         true,
		".next":        true,
		".nuxt":        true,
		".output":      true,
		"target":       true,
	}

	// Command-line flags
	dryRun = flag.Bool("dry-run", false, "Show which files would be modified, without changing them")
	check  = flag.Bool("check", false, "Exit with non-zero status if emojis are found (for CI)")

	// Global state for check mode
	emojisFound = false
)

func removeEmojis(content []byte) []byte {
	return emojiRegex.ReplaceAll(content, []byte(""))
}

func isBinary(content []byte) bool {
	// Use http.DetectContentType which checks the first 512 bytes
	contentType := http.DetectContentType(content)
	// Text types start with "text/" or are "application/json", "application/xml", etc.
	switch {
	case len(content) == 0:
		return false
	case contentType == "application/octet-stream":
		return true
	case len(contentType) >= 5 && contentType[:5] == "text/":
		return false
	case contentType == "application/json":
		return false
	case contentType == "application/xml":
		return false
	default:
		// Images, audio, video, application/zip, etc. are binary
		return true
	}
}

func processFile(path string, info fs.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("Error accessing path %s: %v\n", path, err)
		return nil
	}

	if info.IsDir() {
		if skipDirs[info.Name()] {
			return filepath.SkipDir
		}
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", path, err)
		return nil
	}

	if isBinary(content) {
		return nil
	}

	cleaned := emojiRegex.ReplaceAll(content, []byte(""))
	if len(cleaned) == len(content) {
		return nil // No emojis found
	}

	if *check {
		emojisFound = true
		fmt.Printf("[check] Emojis found in: %s\n", path)
		return nil
	}

	if *dryRun {
		fmt.Printf("[dry-run] Emojis found in: %s\n", path)
		return nil
	}

	if err := os.WriteFile(path, cleaned, info.Mode()); err != nil {
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
