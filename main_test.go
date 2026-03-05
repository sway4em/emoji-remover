package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveEmojis(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no emojis", "hello world", "hello world"},
		{"empty string", "", ""},
		{"smiley face", "hello world", "hello world"},
		{"fire emoji", "code is fire", "code is fire"},
		{"rocket", "deploying", "deploying"},
		{"checkmark", "done", "done"},
		{"star", "favorite", "favorite"},
		{"warning", "caution", "caution"},
		{"sparkles", "magic", "magic"},
		{"thinking face", "hmm", "hmm"},
		{"robot", "bot", "bot"},
		{"mixed text and emojis", "Hello World! How are you?", "Hello World! How are you?"},
		{"emoji only", "", ""},
		{"flag emoji", "flag", "flag"},
		{"compound emoji preserved text", "hello world test", "hello world test"},
		{"preserves newlines", "line1\nline2\n", "line1\nline2\n"},
		{"preserves tabs", "\tindented", "\tindented"},
		{"code with emoji comment", "// todo: fix this\nfunc main() {}", "// todo: fix this\nfunc main() {}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := string(removeEmojis([]byte(tt.input)))
			if result != tt.expected {
				t.Errorf("removeEmojis(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsBinary(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{"empty file", []byte{}, false},
		{"plain text", []byte("hello world"), false},
		{"go source", []byte("package main\n\nfunc main() {}\n"), false},
		{"png header", []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, true},
		{"null bytes", []byte{0x00, 0x00, 0x00, 0x00}, true},
		{"gif header", []byte("GIF89a" + "\x00\x01\x00\x01"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBinary(tt.content)
			if result != tt.expected {
				t.Errorf("isBinary(%v) = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestProcessFileSkipsBinary(t *testing.T) {
	dir := t.TempDir()

	// Create a binary file with emoji-range bytes
	binPath := filepath.Join(dir, "image.png")
	// PNG header followed by some data
	pngData := append([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, []byte("some data")...)
	if err := os.WriteFile(binPath, pngData, 0644); err != nil {
		t.Fatal(err)
	}

	info, _ := os.Stat(binPath)
	err := processFile(binPath, info, nil)
	if err != nil {
		t.Fatalf("processFile returned error: %v", err)
	}

	// File should be unchanged
	content, _ := os.ReadFile(binPath)
	if string(content) != string(pngData) {
		t.Error("binary file was modified")
	}
}

func TestProcessFileRemovesEmojis(t *testing.T) {
	// Save and restore global flag state
	oldDryRun := *dryRun
	oldCheck := *check
	defer func() {
		*dryRun = oldDryRun
		*check = oldCheck
	}()
	*dryRun = false
	*check = false

	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte("hello world"), 0644); err != nil {
		t.Fatal(err)
	}

	info, _ := os.Stat(path)
	err := processFile(path, info, nil)
	if err != nil {
		t.Fatalf("processFile returned error: %v", err)
	}

	content, _ := os.ReadFile(path)
	if string(content) != "hello world" {
		t.Errorf("expected 'hello world', got %q", string(content))
	}
}
