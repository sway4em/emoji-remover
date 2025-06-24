package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

// In a real-world scenario, you'd want a more comprehensive regex.
// This one covers a wide range of common emojis.
var emojiRegex = regexp.MustCompile(`[\x{1F600}-\x{1F64F}\x{1F300}-\x{1F5FF}\x{1F680}-\x{1F6FF}\x{1F1E0}-\x{1F1FF}]`)

func removeEmojis(content []byte) []byte {
	return emojiRegex.ReplaceAll(content, []byte(""))
}

func main() {
	// Define and parse the command-line flags
	filePath := flag.String("file", "", "Path to the file to process")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Error: --file flag is required")
		os.Exit(1)
	}

	// Read the file
	content, err := ioutil.ReadFile(*filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Remove emojis
	newContent := removeEmojis(content)

	// Write the modified content back to the file
	if err := ioutil.WriteFile(*filePath, newContent, 0644); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully removed emojis from %s\n", *filePath)
} 