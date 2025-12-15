<!--
SYNC IMPACT REPORT
==================
Version change: N/A → 1.0.0
Modified principles: N/A (initial creation)
Added sections:
  - Core Principles (5 principles)
  - Quality Standards
  - Development Workflow
  - Governance
Removed sections: N/A
Templates requiring updates:
  - .rainbow/templates/templates-for-commands/design-template.md ✅ (no changes needed - Ground-rules Check section already generic)
  - .rainbow/templates/templates-for-commands/spec-template.md ✅ (no changes needed - template is generic)
  - .rainbow/templates/templates-for-commands/tasks-template.md ✅ (no changes needed - template is generic)
Follow-up TODOs: None
-->

# to-epub Ground-rules

## Core Principles

### I. EPUB 3+ Standards Compliance

All generated EPUB files MUST conform to EPUB 3.0+ specifications as defined by the W3C.

- Output MUST pass epubcheck validation with zero errors
- MUST support EPUB 3 semantic markup (epub:type attributes)
- MUST generate valid XHTML5 content documents
- MUST include proper navigation document (nav.xhtml) with toc, landmarks
- MUST support accessibility metadata (WCAG 2.0 AA compliance where applicable)

**Rationale**: EPUB 3+ is the current standard; compliance ensures compatibility with all major
e-readers and future-proofs the output.

### II. CLI-First Design

The application MUST expose all functionality through a clean command-line interface.

- All operations available via CLI commands with clear --help documentation
- Support stdin/stdout pipelines for Unix composability
- Exit codes MUST follow conventions (0=success, non-zero=error with specific codes)
- Error messages MUST go to stderr, output data to stdout
- Support both human-readable and JSON output formats (--format flag)

**Rationale**: CLI-first enables automation, scripting, and integration into publishing workflows
while maintaining ease of use for direct invocation.

### III. Input Format Fidelity

Conversion MUST preserve semantic meaning and structure from source formats.

- Markdown: Preserve heading hierarchy, lists, code blocks, tables, images, links
- PDF: Extract text with layout awareness; preserve chapter structure when detectable
- HTML: Convert semantic HTML elements to EPUB equivalents; handle CSS appropriately
- MUST handle embedded images (copy, convert if needed, maintain references)
- MUST preserve metadata (title, author, language) when available in source

**Rationale**: The value of conversion lies in maintaining the author's intent and document
structure, not just raw text extraction.

### IV. Maintainability-First Architecture

Code MUST be structured for long-term maintenance and extension.

- Single responsibility: Each module handles one format or one concern
- Clear interfaces between format parsers, EPUB generator, and CLI layer
- No circular dependencies between modules
- Configuration via well-documented options (no magic defaults)
- Comprehensive error handling with actionable error messages

**Rationale**: A converter supporting multiple formats will grow over time; clean separation
enables adding formats or features without destabilizing existing functionality.

### V. Test-Driven Quality

All features MUST have automated tests before implementation is considered complete.

- Unit tests for parsing logic (each input format)
- Integration tests for end-to-end conversion (input → valid EPUB)
- Contract tests: Output MUST pass epubcheck for all test cases
- Test fixtures MUST include edge cases (malformed input, large files, unicode)
- Tests run in CI on every commit

**Rationale**: EPUB generation involves many subtle requirements; tests catch regressions and
validate compliance that manual testing cannot reliably cover.

## Quality Standards

### Output Quality Requirements

- Generated EPUB files MUST be readable in major e-readers (Kindle, Apple Books, Calibre, Kobo)
- Table of contents MUST accurately reflect document structure
- Images MUST be properly embedded and display correctly
- Text encoding MUST be UTF-8 throughout
- File sizes SHOULD be optimized (compress images, minify where appropriate)

### Performance Requirements

- Conversion of typical documents (<100 pages) MUST complete in under 10 seconds
- Memory usage MUST remain bounded (streaming processing for large files)
- Progress indication MUST be provided for long-running operations

### Error Handling Requirements

- Invalid input MUST produce clear error messages identifying the problem
- Partial conversion MUST NOT produce corrupt output files
- Recovery suggestions SHOULD be provided where applicable

## Development Workflow

### Code Quality Gates

- All PRs MUST pass linting and formatting checks
- All PRs MUST pass the full test suite
- Generated EPUB files MUST pass epubcheck validation
- Code MUST include type hints/annotations for static analysis

### Documentation Requirements

- CLI commands MUST be self-documenting via --help
- Public APIs MUST have docstrings with usage examples
- README MUST include quickstart guide with common use cases

### Versioning

- Follow semantic versioning (MAJOR.MINOR.PATCH)
- MAJOR: Breaking CLI interface changes or output format changes
- MINOR: New input formats, new features, new CLI options
- PATCH: Bug fixes, performance improvements, internal refactoring

## Governance

These ground-rules supersede all other development practices for this project.

**Amendment Process**:
1. Propose change with rationale in a PR modifying this document
2. Change MUST include migration plan if affecting existing behavior
3. Version increment according to semantic versioning rules in this document

**Compliance**:
- All code reviews MUST verify adherence to these principles
- Deviations MUST be justified in PR description with clear rationale
- Unjustified deviations block merge

**Version**: 1.0.0 | **Ratified**: 2025-12-15 | **Last Amended**: 2025-12-15
