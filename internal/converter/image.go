package converter

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/webp"

	"github.com/dauquangthanh/epub-converter/internal/model"
)

// Image handling errors
var (
	ErrImageNotFound    = errors.New("image file not found")
	ErrUnsupportedImage = errors.New("unsupported image format")
)

// ImageHandler processes images for EPUB embedding.
type ImageHandler struct{}

// NewImageHandler creates a new image handler.
func NewImageHandler() *ImageHandler {
	return &ImageHandler{}
}

// ProcessImage reads and validates an image file.
func (h *ImageHandler) ProcessImage(path string, basePath string) (*model.Resource, error) {
	// Resolve relative path
	fullPath := path
	if !filepath.IsAbs(path) {
		fullPath = filepath.Join(basePath, path)
	}

	// Read file
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrImageNotFound, path)
	}

	// Detect and validate format
	mediaType, needsConversion := h.detectImageFormat(data, path)
	if mediaType == "" {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedImage, path)
	}

	// Convert WebP to PNG if needed
	if needsConversion {
		var convertErr error
		data, convertErr = h.convertWebPToPNG(data)
		if convertErr != nil {
			return nil, fmt.Errorf("converting WebP to PNG: %w", convertErr)
		}
		mediaType = "image/png"
	}

	// Generate resource ID and filename
	baseName := filepath.Base(path)
	ext := filepath.Ext(baseName)
	name := strings.TrimSuffix(baseName, ext)

	// Update extension if converted
	if needsConversion {
		baseName = name + ".png"
	}

	resource := &model.Resource{
		ID:        "img-" + sanitizeID(name),
		FileName:  "images/" + baseName,
		MediaType: mediaType,
		Data:      data,
	}

	return resource, nil
}

// detectImageFormat determines the image MIME type from content and filename.
// Returns media type and whether conversion is needed.
func (h *ImageHandler) detectImageFormat(data []byte, filename string) (string, bool) {
	// Check magic bytes first
	if len(data) >= 8 {
		// PNG: 89 50 4E 47 0D 0A 1A 0A
		if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
			return "image/png", false
		}
		// JPEG: FF D8 FF
		if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
			return "image/jpeg", false
		}
		// GIF: GIF87a or GIF89a
		if string(data[:6]) == "GIF87a" || string(data[:6]) == "GIF89a" {
			return "image/gif", false
		}
		// WebP: RIFF....WEBP
		if len(data) >= 12 && string(data[:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
			return "image/webp", true // Needs conversion
		}
	}

	// SVG detection by content (starts with <?xml or <svg)
	content := strings.TrimSpace(string(data[:min(len(data), 1024)]))
	if strings.HasPrefix(content, "<?xml") || strings.HasPrefix(content, "<svg") ||
		strings.Contains(content[:min(len(content), 256)], "<svg") {
		return "image/svg+xml", false
	}

	// Fallback to extension
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".png":
		return "image/png", false
	case ".jpg", ".jpeg":
		return "image/jpeg", false
	case ".gif":
		return "image/gif", false
	case ".svg":
		return "image/svg+xml", false
	case ".webp":
		return "image/webp", true
	default:
		return "", false
	}
}

// convertWebPToPNG converts WebP image data to PNG format.
func (h *ImageHandler) convertWebPToPNG(data []byte) ([]byte, error) {
	// Decode WebP
	img, err := webp.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decoding WebP: %w", err)
	}

	// Encode as PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("encoding PNG: %w", err)
	}

	return buf.Bytes(), nil
}

// ValidateImage checks if image data is valid.
func (h *ImageHandler) ValidateImage(data []byte) error {
	_, _, err := image.Decode(bytes.NewReader(data))
	return err
}

// EncodeImage re-encodes an image in the specified format.
func (h *ImageHandler) EncodeImage(img image.Image, format string) ([]byte, error) {
	var buf bytes.Buffer

	switch format {
	case "png", "image/png":
		if err := png.Encode(&buf, img); err != nil {
			return nil, err
		}
	case "jpeg", "jpg", "image/jpeg":
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85}); err != nil {
			return nil, err
		}
	case "gif", "image/gif":
		if err := gif.Encode(&buf, img, nil); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}

	return buf.Bytes(), nil
}

// sanitizeID converts a filename to a valid XML ID.
func sanitizeID(s string) string {
	// Replace non-alphanumeric characters with hyphens
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		} else {
			result.WriteRune('-')
		}
	}
	return strings.Trim(result.String(), "-")
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
