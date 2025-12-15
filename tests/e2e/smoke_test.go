//go:build e2e
// +build e2e

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSmoke_SimpleMarkdownConversion tests Scenario 5.2.1
func TestSmoke_SimpleMarkdownConversion(t *testing.T) {
	// ARRANGE
	binary := buildBinary(t)
	inputPath := filepath.Join("..", "fixtures", "markdown", "simple.md")
	outputDir := t.TempDir()
	outputPath := filepath.Join(outputDir, "output.epub")

	// ACT
	result := runToepub(t, binary, "convert", inputPath, "--output", outputPath)

	// ASSERT
	assert.Equal(t, 0, result.ExitCode, "Exit code should be 0")
	// Progress and success messages go to stderr per Unix conventions
	assert.Contains(t, result.Stderr, "Converting", "Should show progress on stderr")
	assert.Contains(t, result.Stderr, "Created", "Should show success on stderr")
	assert.FileExists(t, outputPath, "EPUB file should be created")

	// Verify EPUB is a valid ZIP
	isValidZip := verifyEPUBIsZip(t, outputPath)
	assert.True(t, isValidZip, "EPUB should be a valid ZIP file")
}

// TestSmoke_CLIHelpDisplay tests Scenario 5.2.5
func TestSmoke_CLIHelpDisplay(t *testing.T) {
	// ARRANGE
	binary := buildBinary(t)

	// ACT - Test root help
	result := runToepub(t, binary, "--help")

	// ASSERT
	assert.Equal(t, 0, result.ExitCode, "Exit code should be 0")
	assert.Contains(t, result.Stdout, "Usage:", "Should contain Usage section")
	assert.Contains(t, result.Stdout, "convert", "Should mention convert command")
	assert.Empty(t, result.Stderr, "Stderr should be empty for help")

	// ACT - Test convert help to verify --output flag is documented
	resultConvert := runToepub(t, binary, "convert", "--help")

	// ASSERT
	assert.Equal(t, 0, resultConvert.ExitCode, "Exit code should be 0")
	assert.Contains(t, resultConvert.Stdout, "--output", "Convert help should mention --output flag")
	assert.Contains(t, resultConvert.Stdout, "--title", "Convert help should mention --title flag")
	assert.Contains(t, resultConvert.Stdout, "--author", "Convert help should mention --author flag")
}

// TestSmoke_CLINoArguments tests Scenario 5.2.6
func TestSmoke_CLINoArguments(t *testing.T) {
	// ARRANGE
	binary := buildBinary(t)

	// ACT
	result := runToepub(t, binary)

	// ASSERT
	assert.Equal(t, 0, result.ExitCode, "Exit code should be 0")
	assert.Contains(t, result.Stdout, "Usage:", "Should show usage")
	assert.Contains(t, result.Stdout, "Available Commands", "Should list commands")
}

// TestSmoke_InvalidInputFile tests Scenario 5.2.7
func TestSmoke_InvalidInputFile(t *testing.T) {
	// ARRANGE
	binary := buildBinary(t)
	nonExistentFile := "nonexistent.md"

	// ACT
	result := runToepub(t, binary, "convert", nonExistentFile)

	// ASSERT
	assert.NotEqual(t, 0, result.ExitCode, "Exit code should be non-zero")
	assert.Contains(t, result.Stderr, nonExistentFile, "Error should mention the file")
	assert.True(t,
		strings.Contains(result.Stderr, "not found") ||
			strings.Contains(result.Stderr, "no such file") ||
			strings.Contains(result.Stderr, "does not exist"),
		"Error should indicate file not found")
}

// TestSmoke_MarkdownWithGFMTables tests Scenario 5.2.2
func TestSmoke_MarkdownWithGFMTables(t *testing.T) {
	// ARRANGE
	binary := buildBinary(t)
	inputPath := filepath.Join("..", "fixtures", "markdown", "gfm-features.md")
	outputDir := t.TempDir()
	outputPath := filepath.Join(outputDir, "output.epub")

	// ACT
	result := runToepub(t, binary, "convert", inputPath, "--output", outputPath)

	// ASSERT
	assert.Equal(t, 0, result.ExitCode, "Exit code should be 0")
	assert.FileExists(t, outputPath, "EPUB should be created")

	// Verify EPUB contains table markup
	content := extractEPUBContent(t, outputPath)
	assert.Contains(t, content, "<table", "EPUB should contain table element")
	assert.Contains(t, content, "<th", "EPUB should contain table header")
	assert.Contains(t, content, "<td", "EPUB should contain table data")
}

// TestSmoke_HTMLToEPUBBasic tests Scenario 5.2.4
func TestSmoke_HTMLToEPUBBasic(t *testing.T) {
	// ARRANGE
	binary := buildBinary(t)
	inputPath := filepath.Join("..", "fixtures", "html", "simple.html")
	outputDir := t.TempDir()
	outputPath := filepath.Join(outputDir, "output.epub")

	// ACT
	result := runToepub(t, binary, "convert", inputPath, "--output", outputPath)

	// ASSERT
	assert.Equal(t, 0, result.ExitCode, "Exit code should be 0")
	assert.FileExists(t, outputPath, "EPUB should be created")

	// Verify EPUB is valid ZIP
	isValidZip := verifyEPUBIsZip(t, outputPath)
	assert.True(t, isValidZip, "EPUB should be a valid ZIP file")
}

// Helper types and functions

type Result struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

func buildBinary(t *testing.T) string {
	t.Helper()

	binaryPath := filepath.Join(t.TempDir(), "toepub")
	if strings.Contains(os.Getenv("GOOS"), "windows") {
		binaryPath += ".exe"
	}

	// Build the binary
	cmd := exec.Command("go", "build", "-o", binaryPath, "../../cmd/toepub")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Failed to build binary: %s", string(output))

	return binaryPath
}

func runToepub(t *testing.T, binary string, args ...string) *Result {
	t.Helper()

	cmd := exec.Command(binary, args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	}

	return &Result{
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}
}

func verifyEPUBIsZip(t *testing.T, epubPath string) bool {
	t.Helper()

	// Try to open as ZIP and check for mimetype
	cmd := exec.Command("unzip", "-l", epubPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Failed to unzip: %s", string(output))
		return false
	}

	// Check for mimetype file
	return strings.Contains(string(output), "mimetype")
}

func extractEPUBContent(t *testing.T, epubPath string) string {
	t.Helper()

	// Extract EPUB to temp dir
	tmpDir := t.TempDir()
	cmd := exec.Command("unzip", "-q", epubPath, "-d", tmpDir)
	err := cmd.Run()
	require.NoError(t, err, "Failed to extract EPUB")

	// Read all XHTML files
	var content strings.Builder
	err = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".xhtml") || strings.HasSuffix(path, ".html")) {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			content.Write(data)
		}
		return nil
	})
	require.NoError(t, err, "Failed to read EPUB content")

	return content.String()
}
