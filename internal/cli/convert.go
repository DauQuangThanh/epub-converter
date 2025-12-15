package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dauquangthanh/epub-converter/internal/converter"
	"github.com/dauquangthanh/epub-converter/internal/model"
)

// Exit codes following BSD sysexits.h conventions
const (
	ExitSuccess       = 0
	ExitGeneralError  = 1
	ExitInvalidArgs   = 2
	ExitFileNotFound  = 64
	ExitFormatError   = 65
	ExitNotWritable   = 66
	ExitInternalError = 70
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert <input>... [flags]",
	Short: "Convert input file(s) to EPUB format",
	Long: `Convert input file(s) to EPUB 3+ format.

Supports Markdown (.md), HTML (.html, .htm), and PDF (.pdf) input.
Multiple files or directories are combined into a single EPUB.`,
	Example: `  # Convert single Markdown file
  toepub convert document.md

  # Convert with custom output path
  toepub convert document.md -o book.epub

  # Convert multiple files
  toepub convert chapter1.md chapter2.md chapter3.md

  # Convert directory
  toepub convert ./docs/

  # Set metadata
  toepub convert document.md --title "My Book" --author "John Doe"

  # Add cover image
  toepub convert document.md --cover cover.jpg

  # JSON output for scripting
  toepub convert document.md --format json

  # From stdin
  cat document.md | toepub convert -`,
	Args: cobra.MinimumNArgs(1),
	RunE: runConvert,
}

// Command flags
var (
	outputPath  string
	outputFmt   string
	title       string
	author      string
	language    string
	coverImage  string
	inputFormat string
)

func init() {
	rootCmd.AddCommand(convertCmd)

	// Define flags
	convertCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path")
	convertCmd.Flags().StringVarP(&outputFmt, "format", "f", "human", "Output format: human or json")
	convertCmd.Flags().StringVarP(&title, "title", "t", "", "Override book title")
	convertCmd.Flags().StringVarP(&author, "author", "a", "", "Override author name")
	convertCmd.Flags().StringVarP(&language, "language", "l", "", "Book language (BCP 47 code)")
	convertCmd.Flags().StringVarP(&coverImage, "cover", "c", "", "Cover image path")
	convertCmd.Flags().StringVar(&inputFormat, "input-format", "", "Force input format: md, html, pdf")
}

// runConvert executes the convert command
func runConvert(cmd *cobra.Command, args []string) error {
	// Build CLI metadata overrides
	cliMeta := buildCLIMetadata()

	// Build converter options
	opts := converter.Options{
		OutputPath:  outputPath,
		InputFormat: inputFormat,
		CLIMetadata: cliMeta,
	}

	// Handle stdin input
	if len(args) == 1 && args[0] == "-" {
		return handleStdinInput(cmd, opts)
	}

	// Resolve output path if not specified
	if opts.OutputPath == "" {
		opts.OutputPath = resolveDefaultOutputPath(args)
	}

	// Print progress for human output
	if outputFmt != "json" {
		printInputSummary(cmd, args)
	}

	// Create converter and run conversion
	conv := converter.New()
	result, err := conv.Convert(args, opts)
	if err != nil {
		return handleConvertError(cmd, err)
	}

	// Output result
	return outputResult(cmd, result)
}

// printInputSummary shows what files are being converted
func printInputSummary(cmd *cobra.Command, inputs []string) {
	if len(inputs) == 1 {
		info, err := os.Stat(inputs[0])
		if err == nil && info.IsDir() {
			cmd.PrintErrf("Converting directory: %s\n", inputs[0])
		} else {
			cmd.PrintErrf("Converting: %s\n", inputs[0])
		}
	} else {
		cmd.PrintErrf("Converting %d files...\n", len(inputs))
	}
}

// buildCLIMetadata creates metadata from CLI flags
func buildCLIMetadata() *model.Metadata {
	meta := model.NewMetadata()

	if title != "" {
		meta.Title = title
	}
	if author != "" {
		meta.Authors = []string{author}
	}
	if language != "" {
		meta.Language = language
	}
	if coverImage != "" {
		meta.CoverImage = coverImage
	}

	return meta
}

// handleStdinInput handles conversion from stdin
func handleStdinInput(cmd *cobra.Command, opts converter.Options) error {
	// Read all stdin
	content, err := readStdin()
	if err != nil {
		return handleConvertError(cmd, err)
	}

	// Require input format for stdin
	if opts.InputFormat == "" {
		opts.InputFormat = "md" // Default to markdown for stdin
	}

	// Set default output path for stdin
	if opts.OutputPath == "" {
		opts.OutputPath = "output.epub"
	}

	conv := converter.New()
	result, err := conv.ConvertContent(content, opts)
	if err != nil {
		return handleConvertError(cmd, err)
	}

	return outputResult(cmd, result)
}

// readStdin reads all content from stdin
func readStdin() ([]byte, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil, fmt.Errorf("no input provided on stdin")
	}

	var content []byte
	buf := make([]byte, 4096)
	for {
		n, err := os.Stdin.Read(buf)
		if n > 0 {
			content = append(content, buf[:n]...)
		}
		if err != nil {
			break
		}
	}
	return content, nil
}

// resolveDefaultOutputPath determines output path from input
func resolveDefaultOutputPath(inputs []string) string {
	if len(inputs) == 0 {
		return "output.epub"
	}

	// For single file, use its name
	if len(inputs) == 1 {
		input := inputs[0]
		if info, err := os.Stat(input); err == nil && info.IsDir() {
			// Directory: use directory name
			return filepath.Base(input) + ".epub"
		}
		// File: replace extension
		ext := filepath.Ext(input)
		return strings.TrimSuffix(input, ext) + ".epub"
	}

	// Multiple files: use "output.epub"
	return "output.epub"
}

// handleConvertError formats and returns conversion errors
func handleConvertError(cmd *cobra.Command, err error) error {
	result := &model.ConversionResult{
		Success: false,
		Error:   err,
	}

	// Map error to exit code
	exitCode := determineExitCode(err)

	if outputFmt == "json" {
		outputJSON(cmd, result)
	} else {
		outputHumanError(cmd, err)
	}

	os.Exit(exitCode)
	return nil // Won't reach here
}

// determineExitCode maps errors to appropriate exit codes
func determineExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}

	errStr := err.Error()

	if strings.Contains(errStr, "file not found") ||
		strings.Contains(errStr, "no such file") {
		return ExitFileNotFound
	}

	if strings.Contains(errStr, "unsupported format") ||
		strings.Contains(errStr, "unknown format") {
		return ExitFormatError
	}

	if strings.Contains(errStr, "permission denied") ||
		strings.Contains(errStr, "not writable") {
		return ExitNotWritable
	}

	return ExitGeneralError
}

// outputResult outputs the conversion result in the appropriate format
func outputResult(cmd *cobra.Command, result *model.ConversionResult) error {
	if outputFmt == "json" {
		outputJSON(cmd, result)
	} else {
		outputHuman(cmd, result)
	}
	return nil
}
