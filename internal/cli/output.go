package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dauquangthanh/epub-converter/internal/model"
)

// Human-readable output symbols
const (
	symbolSuccess = "✓"
	symbolWarning = "⚠"
	symbolError   = "✗"
)

// outputHuman prints human-readable output to stdout
func outputHuman(cmd *cobra.Command, result *model.ConversionResult) {
	if !result.Success {
		outputHumanError(cmd, result.Error)
		return
	}

	// Print warnings first
	for _, warning := range result.Warnings {
		cmd.PrintErrf("%s Warning: %s\n", symbolWarning, warning)
	}

	// Print success message
	sizeKB := result.Stats.OutputSize / 1024
	cmd.Printf("%s Created %s (%d KB)\n", symbolSuccess, result.OutputPath, sizeKB)
	cmd.Printf("  - %d chapters\n", result.Stats.ChapterCount)
	cmd.Printf("  - %d images\n", result.Stats.ImageCount)
	cmd.Printf("  - Duration: %.1fs\n", result.Stats.Duration.Seconds())
}

// outputProgress prints progress message during conversion
func outputProgress(cmd *cobra.Command, message string) {
	cmd.PrintErrf("%s\n", message)
}

// outputHumanError prints human-readable error to stderr
func outputHumanError(cmd *cobra.Command, err error) {
	cmd.PrintErrln()
	cmd.PrintErrf("%s Error: %s\n", symbolError, err.Error())
	cmd.PrintErrln()
}

// outputJSON prints JSON output to stdout
func outputJSON(cmd *cobra.Command, result *model.ConversionResult) {
	output := jsonOutput{
		Success: result.Success,
	}

	if result.Success {
		output.Output = result.OutputPath
		output.Stats = &jsonStats{
			InputFormat: result.Stats.InputFormat,
			InputFiles:  result.Stats.InputFiles,
			Chapters:    result.Stats.ChapterCount,
			Images:      result.Stats.ImageCount,
			OutputSize:  result.Stats.OutputSize,
			DurationMS:  result.Stats.Duration.Milliseconds(),
		}
		output.Warnings = result.Warnings
	} else {
		output.Error = &jsonError{
			Code:    determineExitCode(result.Error),
			Message: result.Error.Error(),
		}
	}

	data, _ := json.MarshalIndent(output, "", "  ")
	cmd.Println(string(data))
}

// JSON output structures

type jsonOutput struct {
	Success  bool        `json:"success"`
	Output   string      `json:"output,omitempty"`
	Stats    *jsonStats  `json:"stats,omitempty"`
	Warnings []string    `json:"warnings,omitempty"`
	Error    *jsonError  `json:"error,omitempty"`
}

type jsonStats struct {
	InputFormat string `json:"input_format"`
	InputFiles  int    `json:"input_files"`
	Chapters    int    `json:"chapters"`
	Images      int    `json:"images"`
	OutputSize  int64  `json:"output_size"`
	DurationMS  int64  `json:"duration_ms"`
}

type jsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// FormatFileSize formats bytes into human-readable size
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
