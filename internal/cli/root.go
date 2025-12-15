// Package cli provides command-line interface handling for toepub.
package cli

import (
	"github.com/spf13/cobra"
)

var (
	// Version information set at build time
	version   = "dev"
	buildDate = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "toepub",
	Short: "Convert Markdown, HTML, and PDF to EPUB 3+",
	Long: `toepub - Convert Markdown, HTML, and PDF to EPUB 3+

A CLI tool that converts document formats to valid EPUB 3+ e-book files.
Supports GitHub Flavored Markdown, HTML with CSS, and text-based PDFs.

Examples:
  # Convert Markdown to EPUB
  toepub convert document.md

  # Convert with metadata
  toepub convert document.md --title "My Book" --author "John Doe"

  # Convert multiple files
  toepub convert chapter1.md chapter2.md chapter3.md -o book.epub`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("toepub version %s\n", version)
		cmd.Printf("Built: %s\n", buildDate)
		cmd.Println("EPUB 3.3 compliant output")
	},
}
