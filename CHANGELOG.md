# Changelog

Notable changes to Arcane Translator. Format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

- A `litellm` provider for any model behind a LiteLLM gateway. Set
  `ARCANE_LLM_PROVIDER=litellm`, `LITELLM_API_KEY` and `ARCANE_LITELLM_BASE_URL`;
  `ARCANE_LITELLM_MODEL` defaults to `claude-sonnet-5-5`.

### Changed

- Every provider now uses the same chapter and novel-details prompts. The
  novel-details prompt follows the chapter prompt's naming rules: only
  character names stay in pinyin, and titles, places, sects and other terms are
  translated into English.
- `ARCANE_TRANSLATION_TEMPERATURE` also applies to the `litellm` provider. A
  model that accepts only its default temperature, such as gpt-5, is retried
  without one. The `claude` provider still sends none, since current Claude
  models reject it.

### Removed

- The `gemini` and `openai` providers. Use `litellm` for models outside AWS
  Bedrock.

## [1.1.0] - 2026-10-16

### Added

- A `bedrock` provider that reaches any chat model on AWS Bedrock through the
  Converse API, using the same translation prompts as the `claude` provider.
  Set `ARCANE_LLM_PROVIDER=bedrock` and a model id such as `zai.glm-4.7`.

### Changed

- `.env.example` and the README default to the `bedrock` provider. The
  `claude` provider is unchanged.

## [1.0.0]

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
- novel543 (novel543.com) as a Chinese source. The site splits each chapter
  across pages, and each page is stored as its own chapter.

### Fixed

- ixdzs pages fetched by the server came back with the wrong details while
  still looking like a normal page, so nothing flagged them. The server no
  longer scrapes ixdzs; the reader is always asked to paste the page source.
- A failed novel extraction appended a stray `null` to its error message.
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
