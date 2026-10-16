# Arcane Translator

Read Chinese, Korean and Japanese webnovels in English. Point it at a novel's
page, and it translates chapters with an LLM as you read them, keeping your
library and reading position locally.

Everything runs on your own machine: a Go API with SQLite, and a SvelteKit
frontend. No account, no server, no telemetry.

## Screens

A sepia reader with a two-page book mode, a paginated library with server-side
filtering, and per-source shelves.

## Requirements

- Go 1.26+
- Node 20+
- An LLM provider: AWS Bedrock, or a LiteLLM gateway

## Install

From source:

```bash
git clone https://github.com/DebanganThakuria/arcane-translator.git
cd arcane-translator
make install
```

With Homebrew, from a local tap:

```bash
brew tap-new <you>/tap --no-git
cp packaging/homebrew/arcane-translator.rb "$(brew --repository)/Library/Taps/<you>/homebrew-tap/Formula/"
brew install <you>/tap/arcane-translator
```

This builds the frontend and the binary, and installs a single process that
serves both. `brew services start arcane-translator` runs it in the background
and starts it again at login.

## Configure

Copy the example file and fill in the provider you use:

```bash
cp .env.example .env
```

The minimum for AWS Bedrock, which reads credentials from the standard AWS
chain (environment, a named profile, or an instance role):

```bash
ARCANE_LLM_PROVIDER=bedrock
AWS_REGION=us-east-1
ARCANE_BEDROCK_MODEL_ID=zai.glm-4.7
```

The `bedrock` provider calls any Bedrock chat model through the Converse API,
so switching models is a change to `ARCANE_BEDROCK_MODEL_ID` alone. The older
`claude` provider sends Anthropic's native request and works only with
Anthropic models.

For a LiteLLM gateway, set `ARCANE_LLM_PROVIDER=litellm`, `LITELLM_API_KEY` and
`ARCANE_LITELLM_BASE_URL`, and optionally `ARCANE_LITELLM_MODEL` (default
`claude-sonnet-5-5`). Every option is documented in [.env.example](.env.example).

No keys or model identifiers are compiled into the binary.

## Run

```bash
make dev
```

The API listens on `:8088` and the frontend on `:8080`. Open
<http://localhost:8080>. Both bind to all interfaces, so the app also opens from
a phone on the same network.

## Adding a novel

Paste the URL of a novel's index page and pick its source site. Only the details
are fetched up front; each chapter is translated the first time you open it.

Some sources block server-side requests. When that happens the app opens the
page in a tab and asks you to paste its source, then translates from that.

Supported sources are listed in
[`backend/provider/sources`](backend/provider/sources). Adding one means
implementing a small interface; see [CONTRIBUTING.md](CONTRIBUTING.md).

## Layout

```
backend/     Go API: handlers, service layer, SQLite repo, LLM and scraper providers
web/         SvelteKit frontend (TypeScript)
packaging/   Homebrew formula
data/        SQLite database, created on first run and never committed
```

A production build serves the frontend from the Go binary, so a packaged
install is a single process on a single port.

## Commands

```bash
make dev      # run the API and the frontend dev server together
make check    # vet the backend and typecheck the frontend
make build    # build the binary and the frontend bundle
make run      # build, then serve everything from the binary on :8088
make stop     # stop running dev processes
```

`make clean` removes dependencies and build output. It leaves `data/` alone.

## Legal

Arcane Translator is a personal reading tool. It fetches pages you point it at
and sends them to a translation model you pay for. You are responsible for
respecting the terms of the sites you read and the copyright of the works you
translate. Do not use it to redistribute translations.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Bug reports and fixes for broken source
sites are the most useful contributions, since sites change their markup often.

## License

[Apache 2.0](LICENSE).
