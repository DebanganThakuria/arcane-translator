# Changelog

Notable changes to Arcane Translator. Format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

- SvelteKit frontend replacing the previous React app. Ships under 25 KB of
  JavaScript on first load.
- Two-page book view in the reader, alongside the scrolling view.
- Fullscreen reading, with the title bar hidden until the top edge is hovered.
- Background translation of the next chapter, started from the reader while the
  current chapter is still being read.
- Server-side filtering, sorting and paging for the novel list, so counts
  reflect the filters and views are linkable.
- Paginated chapter listing. A novel page now loads about 50 rows instead of
  every chapter.
- Configuration through the environment, covering AWS Bedrock, Gemini and
  OpenAI-compatible providers.
- The binary serves the built frontend, so a packaged install is one process.

### Fixed

- The scraper registered a new response callback on a shared collector for every
  request. Handlers accumulated for the lifetime of the process, which leaked
  memory, made each request slower than the last, and let a later response
  overwrite an earlier caller's buffer.
- SQLite ran without WAL, without a busy timeout and without foreign key
  enforcement, so concurrent writes could fail and deleting a novel could orphan
  its chapters.
- The keyed mutex grew without bound and busy-waited.
- `date_added` was populated from `last_updated`, so sorting by it was wrong.
- Missing novels returned 500 rather than 404.
- Chapter list, novel metadata and control heights failed WCAG AA contrast or
  alignment in one or both themes.

### Security

- Model identifiers and API keys are no longer compiled into the binary. The
  Bedrock inference profile ARN, which embeds an AWS account id, now comes from
  the environment.
