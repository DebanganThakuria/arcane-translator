# Security

## Reporting a vulnerability

Please report security issues privately through GitHub's
[private vulnerability reporting](https://github.com/DebanganThakuria/arcane-translator/security/advisories/new),
rather than opening a public issue.

Include what you found, how to reproduce it, and what an attacker could do with
it. You can expect an acknowledgement within a week.

## Threat model

Arcane Translator is a local-first tool. The API binds to all interfaces so the
app can be opened from another device on the same network, and it has **no
authentication**. Anyone who can reach port 8088 can read and delete your
library.

Do not expose it to the public internet. If you need remote access, put it
behind a VPN or an authenticating reverse proxy.

## What the app handles

- **LLM credentials.** Read from the environment at startup. AWS credentials
  come from the standard AWS chain and are never read from application config.
- **Scraped pages.** Fetched from source sites you nominate, then sent to your
  translation provider.
- **Translated content.** Stored in a local SQLite database.

Translated chapter HTML is rendered without sanitisation, because it is produced
by your own translation backend from pages you chose to fetch. A source site
that injects markup, combined with a model that passes it through, could get
script into the reader. If you point the app at sites you do not trust, be aware
of that path.

## Secrets

No keys, model identifiers or account identifiers belong in the repository. If
you find one committed, report it privately using the link above.
