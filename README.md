# Más Bandwidth on Patreon

*[patreon.com/MasBandwidth](https://patreon.com/MasBandwidth)*

**Creating open source game networking code, and the articles that explain it.**

For twenty years I've written about how multiplayer games work and given the
code away. The articles are free. The libraries are MIT/BSD and running in real
games right now. None of that is changing. This Patreon supports the continued
work — and every month an open ledger shows exactly where the money goes.

This repo is the public reference for what the Patreon funds. The monthly books
live in **[open-ledger](https://github.com/mas-bandwidth/open-ledger)**.

## The articles

- **[gafferongames.com](https://gafferongames.com)** — twenty years of articles
  on game networking, physics, and how multiplayer actually works. Free, and
  staying free.
- **[mas-bandwidth.com](https://mas-bandwidth.com)** — where the new writing
  lives.

## The work so far — the open source

The reference libraries under [github.com/mas-bandwidth](https://github.com/mas-bandwidth),
MIT/BSD, used in shipped games. *Star counts as of July 2026 — they only go up.*

| Library | What it does | Latest | ★ |
|---|---|---|---|
| **[yojimbo](https://github.com/mas-bandwidth/yojimbo)** | Client/server network protocol for games — encrypted, dedicated-server ready | v1.6.3 | ~2.7k |
| **[netcode](https://github.com/mas-bandwidth/netcode)** | Secure client/server connection over UDP (connect tokens, encryption) | v1.3.5 | ~2.6k |
| **[reliable](https://github.com/mas-bandwidth/reliable)** | Reliable-ordered messages and acks over UDP | v1.3.4 | ~650 |
| **[serialize](https://github.com/mas-bandwidth/serialize)** | Bitpacking and serialization — one unified read/write path | v1.4.3 | ~130 |
| **[fixed3d](https://github.com/mas-bandwidth/fixed3d)** | Cross-platform **deterministic** fixed-point physics | v1.3.0 | new |

**Ports to other languages** (so the protocols aren't C-only): Go and Rust
ports of netcode, reliable, and serialize —
[netcode.go](https://github.com/mas-bandwidth/netcode.go),
[netcode.rs](https://github.com/mas-bandwidth/netcode.rs),
[reliable.go](https://github.com/mas-bandwidth/reliable.go),
[reliable.rs](https://github.com/mas-bandwidth/reliable.rs),
[serialize.go](https://github.com/mas-bandwidth/serialize.go),
[serialize.rs](https://github.com/mas-bandwidth/serialize.rs).

## Getting the libraries into every package manager

Part of the work nobody sees: getting these libraries into the package managers
so people can just `install` them instead of vendoring source. Each ecosystem
has its own submission process, maintainers, and review — this is the ongoing
push. *Status as of July 2026; the links are the live source of truth.*

| Package manager | Status | Where it stands |
|---|---|---|
| **Homebrew** | ✅ **Shipped** | `serialize` and `libyojimbo` merged into homebrew-core (PRs [#292317](https://github.com/Homebrew/homebrew-core/pull/292317), [#292681](https://github.com/Homebrew/homebrew-core/pull/292681)). Also a tap: [mas-bandwidth/homebrew-tap](https://github.com/mas-bandwidth/homebrew-tap). |
| **vcpkg** | 🔄 **In review** | All four libraries in one PR ([microsoft/vcpkg#52858](https://github.com/microsoft/vcpkg/pull/52858)) — CI green, working through maintainer review. |
| **Debian** | 🔄 **In progress** | ITP/RFS filed and uploaded via mentors.debian.net; awaiting a sponsoring Debian Developer to push to unstable. |
| **FreeBSD** | 🔄 **In review** | Four port submissions filed and assigned to a committer (bugs [296779–296782](https://bugs.freebsd.org/bugzilla/show_bug.cgi?id=296779)). |
| **OpenBSD** | 🔄 **In progress** | `[NEW]` port submission on the ports@ mailing list, awaiting a committer. |

Homebrew is done; the rest are moving through their review pipelines. Progress
shows up in the monthly [ledger](https://github.com/mas-bandwidth/open-ledger).

Roughly **6,000 GitHub stars** across the core libraries, and a steady stream of
releases — the [open ledger](https://github.com/mas-bandwidth/open-ledger)
tracks each month's work with commit ranges you can verify.

These days the open source happens as a disclosed collaboration between Glenn
and an AI assistant. The code review, the tests, the wire-compatibility checks —
real engineering, in the open. The token cost of that collaboration is a line
item in the ledger like everything else.

## The tiers

| | | |
|---|---|---|
| **Free** | — | The posts: new articles + the monthly ledger, as they land. |
| **Thanks** | $1 | For twenty years of free articles. |
| **Supporter** | $5 | Keeps the articles and libraries funded. |
| **Ledger** | $10 | The full itemized monthly books. |
| **Workshop** | $25 | Early draft access + office hours. No treadmill. |
| **Patron** | $50 | For those who want the work to exist and can say so. |
| **Commercial** | $100 | Making money on the libraries? This is how you give back. |
| **Sponsor** | $15,000 | Fund the work directly and help steer priorities — email first. |

*Scope: everything funded here is open source under github.com/mas-bandwidth,
plus the articles. Nothing private, nothing from other organizations.*

## Sponsors

Direct sponsors are credited, with their consent, in
**[SPONSORS.md](https://github.com/mas-bandwidth/open-ledger/blob/main/SPONSORS.md)**.
Larger sponsorships are arranged directly with Más Bandwidth LLC
(glenn@mas-bandwidth.com), off Patreon.

## The lines that don't move

- The articles stay free. The libraries stay MIT/BSD.
- Glenn always writes his own articles — never AI-ghostwritten.
- The work funded here is public. The transparency is the product.
