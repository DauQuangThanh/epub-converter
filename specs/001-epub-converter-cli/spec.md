# Feature Specification: EPUB Converter CLI

**Feature Branch**: `001-epub-converter-cli`
**Created**: 2025-12-15
**Status**: Draft
**Input**: User description: "Create a CLI app with golang to convert markdown, pdf, html files to epub files (support epub 3+ for modern features)"

## Clarifications

### Session 2025-12-15

- Q: Which Markdown flavor should be the primary target? → A: GitHub Flavored Markdown (GFM) - supports tables, task lists, strikethrough, autolinks
- Q: How should multi-file input be handled? → A: Multiple files OR directory input supported, combined into single EPUB in alphabetical order
- Q: Should the CLI support verbose or debug output modes? → A: No verbose mode; errors only on stderr (keep CLI simple)
- Q: Which image formats should be supported? → A: EPUB core (JPEG, PNG, GIF, SVG) + auto-convert WebP to PNG; warn on other unsupported formats
- Q: What should the default output filename be? → A: Current working directory, same base name + .epub (e.g., doc.md → doc.epub)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Convert Markdown to EPUB (Priority: P1)

As a technical writer, I want to convert my Markdown documentation into a professionally formatted EPUB file so that readers can consume it on any e-reader device.

**Why this priority**: Markdown is the most common source format for technical documentation and books. This represents the core value proposition and simplest conversion path to validate the EPUB generation pipeline.

**Independent Test**: Can be fully tested by providing a sample Markdown file and verifying the output EPUB opens correctly in an e-reader application and passes validation.

**Acceptance Scenarios**:

1. **Given** a valid Markdown file with headings, paragraphs, and lists, **When** I run the converter with the file path, **Then** a valid EPUB 3+ file is generated with proper table of contents reflecting heading structure.
2. **Given** a Markdown file with embedded images (relative paths), **When** I run the converter, **Then** images are included in the EPUB and display correctly.
3. **Given** a Markdown file with code blocks and tables, **When** I run the converter, **Then** these elements are properly formatted in the EPUB output.
4. **Given** a Markdown file with front matter (title, author), **When** I run the converter, **Then** this metadata appears in the EPUB metadata.
5. **Given** multiple Markdown files or a directory containing Markdown files, **When** I run the converter, **Then** files are combined into a single EPUB in alphabetical order with proper chapter navigation.

---

### User Story 2 - Convert HTML to EPUB (Priority: P2)

As a web content creator, I want to convert my HTML articles or web pages into EPUB format so that I can distribute them as e-books.

**Why this priority**: HTML is a natural fit for EPUB (which uses XHTML internally) and represents the second most requested conversion path. Shares significant infrastructure with Markdown conversion.

**Independent Test**: Can be fully tested by providing a sample HTML file and verifying the output EPUB renders correctly with styles preserved.

**Acceptance Scenarios**:

1. **Given** a valid HTML file with semantic markup (h1-h6, p, ul, ol), **When** I run the converter, **Then** a valid EPUB 3+ file is generated preserving document structure.
2. **Given** an HTML file with inline and linked CSS, **When** I run the converter, **Then** styling is appropriately included or converted in the EPUB.
3. **Given** an HTML file with images (local and remote URLs), **When** I run the converter, **Then** images are embedded in the EPUB (remote images downloaded and included).
4. **Given** multiple HTML files specified as input, **When** I run the converter, **Then** they are combined into a single EPUB with proper chapter navigation.

---

### User Story 3 - Convert PDF to EPUB (Priority: P3)

As a publisher, I want to convert existing PDF documents to EPUB format so that I can offer reflowable e-book versions to readers.

**Why this priority**: PDF conversion is the most complex due to PDF's fixed-layout nature. It provides high value but requires sophisticated text extraction. Delivered after core EPUB generation is proven stable.

**Independent Test**: Can be fully tested by providing a sample PDF and verifying extracted text and structure appear correctly in the generated EPUB.

**Acceptance Scenarios**:

1. **Given** a text-based PDF with clear paragraph structure, **When** I run the converter, **Then** text is extracted and formatted as readable EPUB content.
2. **Given** a PDF with chapter headings, **When** I run the converter, **Then** the EPUB table of contents reflects detected chapter structure.
3. **Given** a PDF with embedded images, **When** I run the converter, **Then** images are extracted and included in the EPUB.
4. **Given** a scanned/image-based PDF, **When** I run the converter, **Then** an appropriate error message indicates OCR is required (out of scope for initial version).

---

### User Story 4 - CLI Usage and Options (Priority: P1)

As a developer or automation engineer, I want comprehensive CLI options so that I can integrate the converter into scripts and publishing pipelines.

**Why this priority**: CLI usability is fundamental to all conversion scenarios. Users need clear commands, helpful documentation, and scriptable interfaces.

**Independent Test**: Can be fully tested by running help commands and verifying all options are documented and functional.

**Acceptance Scenarios**:

1. **Given** no arguments, **When** I run the converter, **Then** usage help is displayed with available commands and options.
2. **Given** the --help flag, **When** I run the converter, **Then** detailed help text explains all options with examples.
3. **Given** an invalid input file, **When** I run the converter, **Then** a clear error message is displayed on stderr and exit code is non-zero.
4. **Given** a valid conversion with --output flag, **When** I run the converter, **Then** the EPUB is written to the specified path.
5. **Given** --format json flag, **When** I run the converter, **Then** output (success/error) is formatted as JSON for programmatic parsing.

---

### User Story 5 - Metadata Configuration (Priority: P2)

As an author, I want to specify or override book metadata (title, author, language, cover image) so that the EPUB contains accurate publication information.

**Why this priority**: Proper metadata improves discoverability in e-reader libraries and is required for professional publishing workflows.

**Independent Test**: Can be fully tested by providing metadata via CLI flags and verifying it appears correctly in the EPUB metadata.

**Acceptance Scenarios**:

1. **Given** --title and --author flags, **When** I run the converter, **Then** these values appear in the EPUB metadata (dc:title, dc:creator).
2. **Given** a --cover flag with an image path, **When** I run the converter, **Then** the image is set as the EPUB cover.
3. **Given** a --language flag, **When** I run the converter, **Then** the EPUB language metadata is set accordingly.
4. **Given** metadata in source file (e.g., YAML front matter) AND CLI flags, **When** I run the converter, **Then** CLI flags take precedence.

---

### Edge Cases

- What happens when the input file does not exist? → Clear error message with file path and exit code 1.
- What happens when the input format cannot be detected? → Error suggesting explicit format flag or checking file extension.
- What happens when the output path is not writable? → Error message indicating permission issue.
- What happens when a Markdown file references non-existent images? → Warning message listing missing images; conversion continues with placeholders or skipped images.
- What happens when HTML contains JavaScript? → JavaScript is stripped (not executable in EPUB); optional warning.
- What happens when PDF text extraction fails? → Clear error indicating the PDF may be image-based or corrupted.
- What happens with extremely large files (>100MB)? → Progress indication shown; streaming processing to manage memory.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept Markdown (.md) files as input and produce valid EPUB 3+ output.
- **FR-002**: System MUST accept HTML (.html, .htm) files as input and produce valid EPUB 3+ output.
- **FR-003**: System MUST accept PDF (.pdf) files as input and produce valid EPUB 3+ output.
- **FR-004**: System MUST auto-detect input format from file extension, with option to override via --format flag.
- **FR-005**: System MUST generate EPUB files that pass epubcheck validation with zero errors.
- **FR-006**: System MUST generate a navigable table of contents based on document heading structure.
- **FR-007**: System MUST embed referenced images into the EPUB package.
- **FR-008**: System MUST preserve document metadata (title, author, language) when present in source.
- **FR-009**: System MUST allow CLI override of metadata via flags (--title, --author, --language, --cover).
- **FR-010**: System MUST provide clear error messages on stderr for invalid inputs or failures.
- **FR-011**: System MUST return appropriate exit codes (0 for success, non-zero for errors).
- **FR-012**: System MUST support --output flag to specify output file path.
- **FR-013**: System MUST support --format flag for output format (human-readable or JSON).
- **FR-014**: System MUST display usage help when run without arguments or with --help.
- **FR-015**: System MUST support reading input from stdin when input is "-" or piped.
- **FR-016**: System MUST handle UTF-8 encoded content correctly.
- **FR-017**: System MUST accept multiple input files and combine them into a single EPUB in alphabetical order.
- **FR-018**: System MUST accept a directory path as input and process all supported files within it (non-recursive by default).
- **FR-019**: System MUST support JPEG, PNG, GIF, and SVG image formats in EPUB output.
- **FR-020**: System MUST auto-convert WebP images to PNG format during EPUB generation.
- **FR-021**: System MUST display a warning for unsupported image formats and continue conversion without the unsupported image.
- **FR-022**: System MUST default output to current working directory with input base name + .epub extension when --output is not specified.

### Key Entities

- **Source Document**: The input file (Markdown, HTML, or PDF) containing content to convert. Key attributes: file path, detected format, encoding, embedded resources.
- **EPUB Package**: The output e-book file conforming to EPUB 3+ specification. Key attributes: content documents (XHTML), navigation document, package document (OPF), metadata, embedded media.
- **Book Metadata**: Publication information for the EPUB. Key attributes: title, author(s), language, identifier, cover image, publication date.
- **Table of Contents**: Hierarchical navigation structure. Key attributes: chapter titles, heading levels, links to content sections.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can convert a 50-page Markdown document to EPUB in under 5 seconds.
- **SC-002**: Users can convert a 50-page HTML document to EPUB in under 5 seconds.
- **SC-003**: Users can convert a 50-page text-based PDF to EPUB in under 30 seconds.
- **SC-004**: 100% of generated EPUB files pass epubcheck validation with zero errors.
- **SC-005**: Generated EPUBs open and render correctly in at least 3 major e-reader applications (e.g., Apple Books, Calibre, Kobo).
- **SC-006**: Users can understand available options and successfully convert a file within 2 minutes of first use (via --help).
- **SC-007**: All conversion errors produce actionable error messages that identify the problem and suggest resolution.
- **SC-008**: The converter handles files up to 100MB without crashing or excessive memory usage.

## Assumptions

- Users have the CLI tool installed and available in their PATH.
- Input Markdown follows GitHub Flavored Markdown (GFM) specification, supporting tables, task lists, strikethrough, and autolinks.
- Input HTML is well-formed (or at least parseable with standard HTML parsers).
- PDF text extraction targets text-based PDFs; scanned/image PDFs require external OCR (out of scope).
- Remote image URLs in HTML are accessible at conversion time; unreachable images produce warnings.
- The primary target audience is developers, technical writers, and small publishers comfortable with CLI tools.
