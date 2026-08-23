# Contributing

Thanks for taking a look. This is a small project, so the process is light.

## Getting set up

```bash
make install
cp .env.example .env   # fill in one provider
make dev
```

`make dev` runs the Go API on `:8088` and the SvelteKit frontend on `:8080`.

## Before you open a pull request

```bash
cd backend && go build ./... && go vet ./... && gofumpt -l .
cd web && npm run check
```

CI runs the same checks. `gofumpt -l .` should print nothing.

## House rules

- **No secrets in the tree.** Model identifiers, keys and account identifiers
  are configuration, read through `backend/config`. Nothing goes in a
  constant, a default, or a test fixture.
- **No generated or binary artefacts.** `data/`, `bin/`, `node_modules/` and
  compiled binaries stay out of git.
- **Match the surrounding code.** Comments explain why, not what.
- **Keep the frontend dependency-light.** The whole app ships under 25 KB of
  JavaScript; a new runtime dependency needs a reason.

## Design constraints

Two are deliberate and shouldn't be changed without discussion:

- The reader's sepia background is `#FEF7CD`, and Verdana is the default
  reading face. They are what makes long sessions comfortable.
- The interface stays in one script. Source badges are `CN` / `KR` / `JP`; a
  novel's original title is shown because it identifies the work, but the
  chrome itself is not decorated with CJK characters.

All colours, spacing and control sizes come from tokens in
`web/src/styles/tokens.css`. In particular, every control shares `--control-h`;
never set a control height locally.

## Adding a source site

Sources live in `backend/provider/sources`. Each implements the interface in
`interfaces.go`: given a novel URL it extracts an id, and given a chapter page
it finds the content and the link to the next chapter.

To add one:

1. Create `backend/provider/sources/<name>.go` implementing the interface.
2. Register it in the source list in `backend/service/novel_service.go`,
   including the language, so it appears on the right shelf.
3. Add a novel from that site locally and read a couple of chapters.

Say which site you tested against in the pull request. Sites change their
markup often, so a fix for a broken source is always welcome.

## Reporting bugs

Include the source site, what you expected, what happened, and the relevant
server log lines. Scrub anything identifying from logs before pasting them.

## Security

Please don't open a public issue for a security problem. See
[SECURITY.md](SECURITY.md).

## License

Contributions are accepted under the [Apache 2.0](LICENSE) license.
