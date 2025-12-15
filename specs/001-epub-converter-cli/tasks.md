# Tasks: EPUB Converter CLI

**Input**: Design documents from `specs/001-epub-converter-cli/`
**Prerequisites**: design.md (required), spec.md (required), data-model.md, contracts/, research.md, quickstart.md

**Tests**: Tests are included as they support the ground-rules principle V (Test-Driven Quality).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, US4, US5)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: `cmd/`, `internal/`, `tests/` at repository root
- Go module: `github.com/dauquangthanh/epub-converter`

---

## Phase 1: Setup (Project Initialization)

**Purpose**: Project structure, dependencies, and basic configuration

- [X] T001 Initialize Go module with `go mod init github.com/dauquangthanh/epub-converter`
- [X] T002 [P] Create directory structure: cmd/toepub/, internal/cli/, internal/parser/, internal/epub/, internal/converter/, internal/model/
- [X] T003 [P] Create directory structure: tests/unit/, tests/integration/, tests/contract/, tests/fixtures/
- [X] T004 [P] Add .gitignore for Go projects (binaries, test coverage, IDE files)
- [X] T005 Add dependencies to go.mod: cobra, goldmark, goldmark-frontmatter, testify
- [X] T006 [P] Create Makefile with targets: build, test, lint, clean

**Checkpoint**: Project compiles with `go build ./...`

---

## Phase 2: Foundational (Core Infrastructure)

**Purpose**: Shared infrastructure needed by ALL user stories - CLI framework and EPUB generation

**⚠️ CRITICAL**: User Story implementation cannot begin until this phase is complete

### Data Models (shared by all stories)

- [X] T007 [P] Create Document struct in internal/model/document.go
- [X] T008 [P] Create Metadata struct in internal/model/metadata.go
- [X] T009 [P] Create Chapter struct in internal/model/document.go
- [X] T010 [P] Create Resource struct in internal/model/document.go
- [X] T011 [P] Create TableOfContents and TOCEntry structs in internal/model/toc.go
- [X] T012 [P] Create ConversionResult and ConversionStats structs in internal/model/document.go

### Parser Interface

- [X] T013 Create Parser interface in internal/parser/parser.go with Parse(input) → Document method

### EPUB Generator (required for any conversion)

- [X] T014 Create EPUB Builder skeleton in internal/epub/builder.go with Build(Document) → []byte method
- [X] T015 Implement mimetype file generation in internal/epub/builder.go
- [X] T016 Implement META-INF/container.xml generation in internal/epub/builder.go
- [X] T017 Implement content.opf (package document) in internal/epub/package.go
- [X] T018 Implement nav.xhtml (navigation document) with TOC in internal/epub/navigation.go
- [X] T019 Implement nav.xhtml landmarks section in internal/epub/navigation.go
- [X] T020 Implement XHTML content document generation in internal/epub/content.go
- [X] T021 Implement Dublin Core metadata handling in internal/epub/metadata.go
- [X] T022 Implement ZIP archive creation with correct EPUB structure in internal/epub/builder.go

### CLI Framework (US4 - foundational for all)

- [X] T023 Create root command with help in internal/cli/root.go using cobra
- [X] T024 Create version command in internal/cli/root.go
- [X] T025 Create convert command skeleton in internal/cli/convert.go
- [X] T026 Add --output flag to convert command in internal/cli/convert.go
- [X] T027 Add --format flag (human/json) to convert command in internal/cli/convert.go
- [X] T028 Add --input-format flag to convert command in internal/cli/convert.go
- [X] T029 Implement human-readable output formatter in internal/cli/output.go
- [X] T030 Implement JSON output formatter in internal/cli/output.go
- [X] T031 Implement exit codes (0, 1, 2, 64, 65, 66, 70) in internal/cli/convert.go
- [X] T032 Create main.go entry point in cmd/toepub/main.go

### Converter Pipeline

- [X] T033 Create Converter interface in internal/converter/converter.go
- [X] T034 Implement format detection from file extension in internal/converter/converter.go
- [X] T035 Implement input file validation in internal/converter/converter.go
- [X] T036 Implement output path resolution (default to CWD + .epub) in internal/converter/converter.go

### Tests for Foundational

- [ ] T037 [P] Unit test for Document model validation in tests/unit/model/document_test.go
- [ ] T038 [P] Unit test for EPUB builder (generates valid ZIP) in tests/unit/epub/builder_test.go
- [ ] T039 [P] Unit test for CLI flag parsing in tests/unit/cli/convert_test.go
- [ ] T040 Integration test for CLI help output in tests/integration/cli_test.go

**Checkpoint**: `toepub --help` works, `toepub version` works, empty EPUB structure can be generated

---

## Phase 3: User Story 1 - Convert Markdown to EPUB (Priority: P1) 🎯 MVP

**Goal**: Convert single Markdown file to valid EPUB 3+ with TOC

**Independent Test**: `toepub convert sample.md` produces valid EPUB that passes epubcheck

### Implementation for User Story 1

- [X] T041 [US1] Implement Markdown parser skeleton in internal/parser/markdown.go
- [X] T042 [US1] Configure goldmark with GFM extension in internal/parser/markdown.go
- [X] T043 [US1] Configure goldmark with XHTML output in internal/parser/markdown.go
- [X] T044 [US1] Implement YAML front matter parsing with goldmark-frontmatter in internal/parser/markdown.go
- [X] T045 [US1] Extract metadata (title, author, language) from front matter in internal/parser/markdown.go
- [X] T046 [US1] Extract heading structure for TOC generation in internal/parser/markdown.go
- [X] T047 [US1] Handle relative image paths in Markdown in internal/parser/markdown.go
- [X] T048 [US1] Implement GFM tables rendering in internal/parser/markdown.go
- [X] T049 [US1] Implement GFM task lists rendering in internal/parser/markdown.go
- [X] T050 [US1] Implement code blocks with syntax styling in internal/parser/markdown.go
- [X] T051 [US1] Wire Markdown parser to convert command in internal/cli/convert.go
- [X] T052 [US1] Implement multi-file Markdown input (alphabetical order) in internal/converter/converter.go
- [X] T053 [US1] Implement directory input for Markdown files in internal/converter/converter.go

### Image Handling (needed for US1)

- [X] T054 [US1] Create image handler in internal/converter/image.go
- [X] T055 [US1] Implement local image file reading in internal/converter/image.go
- [X] T056 [US1] Implement JPEG/PNG/GIF/SVG format detection in internal/converter/image.go
- [X] T057 [US1] Implement WebP to PNG conversion in internal/converter/image.go
- [ ] T058 [US1] Generate warnings for missing/unsupported images in internal/converter/image.go

### Tests for User Story 1

- [X] T059 [P] [US1] Create test fixture: simple.md in tests/fixtures/markdown/
- [ ] T060 [P] [US1] Create test fixture: with-images.md in tests/fixtures/markdown/
- [X] T061 [P] [US1] Create test fixture: with-frontmatter.md in tests/fixtures/markdown/
- [X] T062 [P] [US1] Create test fixture: gfm-features.md (tables, task lists) in tests/fixtures/markdown/
- [ ] T063 [US1] Unit test for Markdown parser in tests/unit/parser/markdown_test.go
- [ ] T064 [US1] Integration test for Markdown → EPUB in tests/integration/markdown_test.go
- [ ] T065 [US1] Contract test: output passes epubcheck in tests/contract/epubcheck_test.go

**Checkpoint**: User Story 1 complete - Markdown conversion produces valid EPUB

---

## Phase 4: User Story 4 - CLI Completeness (Priority: P1)

**Goal**: All CLI options work correctly with Markdown conversion

**Independent Test**: All CLI flags work as documented, exit codes correct

### Implementation for User Story 4

- [X] T066 [US4] Implement stdin input reading ("-" argument) in internal/cli/convert.go
- [X] T067 [US4] Implement error messages on stderr in internal/cli/output.go
- [X] T068 [US4] Implement file not found error (exit 64) in internal/cli/convert.go
- [X] T069 [US4] Implement format detection error (exit 65) in internal/cli/convert.go
- [X] T070 [US4] Implement output not writable error (exit 66) in internal/cli/convert.go
- [X] T071 [US4] Implement progress output for human format in internal/cli/output.go
- [X] T072 [US4] Implement warning output in human format in internal/cli/output.go

### Tests for User Story 4

- [ ] T073 [P] [US4] Test --help output matches contract in tests/integration/cli_test.go
- [ ] T074 [P] [US4] Test --version output in tests/integration/cli_test.go
- [ ] T075 [P] [US4] Test exit code 64 for missing file in tests/integration/cli_test.go
- [ ] T076 [P] [US4] Test exit code 0 for success in tests/integration/cli_test.go
- [ ] T077 [US4] Test JSON output format in tests/integration/cli_test.go
- [ ] T078 [US4] Test stdin input in tests/integration/cli_test.go

**Checkpoint**: CLI fully functional with Markdown - MVP complete

---

## Phase 5: User Story 5 - Metadata Configuration (Priority: P2)

**Goal**: CLI flags override metadata from source files

**Independent Test**: `toepub convert doc.md --title "Custom" --author "Name"` shows metadata in EPUB

### Implementation for User Story 5

- [X] T079 [US5] Add --title flag to convert command in internal/cli/convert.go
- [X] T080 [US5] Add --author flag to convert command in internal/cli/convert.go
- [X] T081 [US5] Add --language flag to convert command in internal/cli/convert.go
- [X] T082 [US5] Add --cover flag to convert command in internal/cli/convert.go
- [X] T083 [US5] Implement metadata precedence (CLI > front matter > default) in internal/converter/converter.go
- [X] T084 [US5] Implement cover image embedding in internal/epub/builder.go
- [X] T085 [US5] Generate UUID identifier when not provided in internal/model/metadata.go

### Tests for User Story 5

- [ ] T086 [P] [US5] Test --title overrides front matter in tests/integration/metadata_test.go
- [ ] T087 [P] [US5] Test --author overrides front matter in tests/integration/metadata_test.go
- [ ] T088 [P] [US5] Test --cover embeds image in tests/integration/metadata_test.go
- [ ] T089 [US5] Test metadata appears in content.opf in tests/unit/epub/metadata_test.go

**Checkpoint**: User Story 5 complete - Metadata configuration works

---

## Phase 6: User Story 2 - Convert HTML to EPUB (Priority: P2)

**Goal**: Convert HTML files to valid EPUB 3+

**Independent Test**: `toepub convert article.html` produces valid EPUB with styles preserved

### Implementation for User Story 2

- [ ] T090 [US2] Implement HTML parser skeleton in internal/parser/html.go
- [ ] T091 [US2] Parse HTML using golang.org/x/net/html in internal/parser/html.go
- [ ] T092 [US2] Convert HTML5 to valid XHTML in internal/parser/html.go
- [ ] T093 [US2] Extract title from <title> tag in internal/parser/html.go
- [ ] T094 [US2] Extract metadata from <meta> tags in internal/parser/html.go
- [ ] T095 [US2] Extract heading structure for TOC in internal/parser/html.go
- [ ] T096 [US2] Handle inline CSS in internal/parser/html.go
- [ ] T097 [US2] Handle linked CSS (download and embed) in internal/parser/html.go
- [ ] T098 [US2] Handle local image references in internal/parser/html.go
- [ ] T099 [US2] Handle remote image URLs (download and embed) in internal/converter/image.go
- [ ] T100 [US2] Strip JavaScript from HTML in internal/parser/html.go
- [ ] T101 [US2] Wire HTML parser to convert command in internal/cli/convert.go
- [ ] T102 [US2] Implement multi-file HTML input in internal/converter/converter.go

### Tests for User Story 2

- [ ] T103 [P] [US2] Create test fixture: simple.html in tests/fixtures/html/
- [ ] T104 [P] [US2] Create test fixture: with-css.html in tests/fixtures/html/
- [ ] T105 [P] [US2] Create test fixture: with-images.html in tests/fixtures/html/
- [ ] T106 [US2] Unit test for HTML parser in tests/unit/parser/html_test.go
- [ ] T107 [US2] Integration test for HTML → EPUB in tests/integration/html_test.go
- [ ] T108 [US2] Contract test: HTML EPUB passes epubcheck in tests/contract/epubcheck_test.go

**Checkpoint**: User Story 2 complete - HTML conversion works

---

## Phase 7: User Story 3 - Convert PDF to EPUB (Priority: P3)

**Goal**: Extract text from PDF and generate EPUB

**Independent Test**: `toepub convert document.pdf` produces readable EPUB from text-based PDF

### Implementation for User Story 3

- [ ] T109 [US3] Add ledongthuc/pdf and pdfcpu dependencies to go.mod
- [ ] T110 [US3] Implement PDF parser skeleton in internal/parser/pdf.go
- [ ] T111 [US3] Implement text extraction using ledongthuc/pdf in internal/parser/pdf.go
- [ ] T112 [US3] Implement heading detection using font size heuristics in internal/parser/pdf.go
- [ ] T113 [US3] Implement paragraph structure detection in internal/parser/pdf.go
- [ ] T114 [US3] Implement image extraction using pdfcpu in internal/parser/pdf.go
- [ ] T115 [US3] Detect scanned/image-based PDFs and return error in internal/parser/pdf.go
- [ ] T116 [US3] Wire PDF parser to convert command in internal/cli/convert.go

### Tests for User Story 3

- [ ] T117 [P] [US3] Create test fixture: text-based.pdf in tests/fixtures/pdf/
- [ ] T118 [P] [US3] Create test fixture: with-headings.pdf in tests/fixtures/pdf/
- [ ] T119 [P] [US3] Create test fixture: with-images.pdf in tests/fixtures/pdf/
- [ ] T120 [P] [US3] Create test fixture: scanned.pdf in tests/fixtures/pdf/
- [ ] T121 [US3] Unit test for PDF parser in tests/unit/parser/pdf_test.go
- [ ] T122 [US3] Integration test for PDF → EPUB in tests/integration/pdf_test.go
- [ ] T123 [US3] Test error message for scanned PDF in tests/integration/pdf_test.go
- [ ] T124 [US3] Contract test: PDF EPUB passes epubcheck in tests/contract/epubcheck_test.go

**Checkpoint**: User Story 3 complete - PDF conversion works for text-based PDFs

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Quality improvements, documentation, final validation

- [ ] T125 [P] Add README.md with installation and usage instructions
- [ ] T126 [P] Add LICENSE file
- [ ] T127 [P] Configure golangci-lint for code quality
- [ ] T128 Run all tests and fix any failures
- [ ] T129 Run epubcheck on all test output EPUBs
- [ ] T130 Performance test: 50-page Markdown in <5s
- [ ] T131 Performance test: 50-page HTML in <5s
- [ ] T132 Performance test: 50-page PDF in <30s
- [ ] T133 Memory test: 100MB file doesn't exceed limits
- [ ] T134 Cross-platform build test (Linux, macOS, Windows)
- [ ] T135 Create GitHub Actions CI workflow in .github/workflows/ci.yml

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Phase 2 - Markdown conversion
- **User Story 4 (Phase 4)**: Depends on Phase 3 - CLI completeness
- **User Story 5 (Phase 5)**: Depends on Phase 4 - Metadata flags
- **User Story 2 (Phase 6)**: Depends on Phase 2 - Can start parallel to US1 if needed
- **User Story 3 (Phase 7)**: Depends on Phase 2 - Can start parallel to US1/US2 if needed
- **Polish (Phase 8)**: Depends on all user stories

### MVP Scope

**Minimum Viable Product = Phase 1 + Phase 2 + Phase 3 + Phase 4**

- Project setup ✓
- Core EPUB generation ✓
- Markdown → EPUB conversion ✓
- CLI with all flags working ✓

### User Story Independence

| Story | Can Start After | Dependencies on Other Stories |
|-------|-----------------|-------------------------------|
| US1 (Markdown) | Phase 2 | None |
| US4 (CLI) | Phase 3 (uses US1 for testing) | US1 |
| US5 (Metadata) | Phase 4 | US4 (needs flags) |
| US2 (HTML) | Phase 2 | None (shares image handler with US1) |
| US3 (PDF) | Phase 2 | None |

### Parallel Opportunities

Within each phase, tasks marked [P] can run in parallel:

```bash
# Phase 1 - All setup tasks parallel
T002, T003, T004, T006

# Phase 2 - Models parallel, then EPUB components
T007, T008, T009, T010, T011, T012 (models)
T014-T022 sequential (EPUB builder)
T023-T032 sequential (CLI)

# Phase 3 - Test fixtures parallel, then implementation
T059, T060, T061, T062 (fixtures)
T041-T058 mostly sequential (parser depends on previous)

# Multiple user stories can be parallelized by different developers
```

---

## Implementation Strategy

### MVP First (Recommended)

1. Complete Phase 1: Setup (~1 hour)
2. Complete Phase 2: Foundational (~1-2 days)
3. Complete Phase 3: User Story 1 - Markdown (~1 day)
4. Complete Phase 4: User Story 4 - CLI (~0.5 day)
5. **STOP and VALIDATE**: Test with real Markdown files, run epubcheck
6. Demo/Release MVP

### Incremental Delivery

After MVP:
1. Add User Story 5 (Metadata) → minor release
2. Add User Story 2 (HTML) → minor release
3. Add User Story 3 (PDF) → minor release
4. Polish → patch releases

### Task Count Summary

| Phase | Tasks | Parallel Opportunities |
|-------|-------|------------------------|
| Phase 1: Setup | 6 | 4 |
| Phase 2: Foundational | 34 | 10 |
| Phase 3: US1 Markdown | 27 | 4 |
| Phase 4: US4 CLI | 13 | 4 |
| Phase 5: US5 Metadata | 11 | 3 |
| Phase 6: US2 HTML | 19 | 3 |
| Phase 7: US3 PDF | 16 | 4 |
| Phase 8: Polish | 11 | 3 |
| **Total** | **137** | **35** |

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story
- Verify tests fail before implementing (TDD)
- Commit after each task or logical group
- Run epubcheck after each user story completion
- Stop at any checkpoint to validate independently
