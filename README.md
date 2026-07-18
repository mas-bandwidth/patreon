# Más Bandwidth on Patreon

> ## ⚠️ Security notice: upgrade your libraries
>
> The 2026 hardening work found and fixed real security bugs in code that had
> shipped for years, including a **remotely reachable heap overflow in yojimbo
> present in every release since 2019**. If you use any of these libraries in a
> product, upgrade to the latest release now:
>
> | Library | Upgrade | What it fixes |
> |---|---|---|
> | **yojimbo** | [→ latest release](https://github.com/mas-bandwidth/yojimbo/releases/latest) | remote heap overflow, wire-reachable asserts, union misread (v1.5.0); AEAD nonce-reuse on restart (v1.7.0) |
> | **netcode** | [→ latest release](https://github.com/mas-bandwidth/netcode/releases/latest) | AEAD nonce-reuse on server restart (v1.4.0), replay-protection overflow, memory-safety hardening |
> | **reliable** | [→ latest release](https://github.com/mas-bandwidth/reliable/releases/latest) | read-buffer over-read, integer overflow |
> | **serialize** | [→ latest release](https://github.com/mas-bandwidth/serialize/releases/latest) | fuzz-hardened, wire format pinned on every platform |
>
> Every fix has a regression test, and the full honest accounting is here:
> **[the bugs found and fixed with the help of AI →](BUGS.md)**

*[patreon.com/MasBandwidth](https://patreon.com/MasBandwidth)*

---

## Introduction

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

## Just landed: fixed3d

The newest library: **[fixed3d](https://github.com/mas-bandwidth/fixed3d)**
takes Erin Catto's [Box3D](https://github.com/erincatto/box3d) and tears every
`float` out of the simulation, replacing them with Q48.16 fixed point. Uniform
1/65536 resolution across a ±1.4×10¹⁴ meter world. Bit-exact on every platform.
All 22 test suites pass, and the benchmarks are published honestly — about 2×
slower than float, and the README tells you straight when you should keep
using Box3D instead.

Here's the part worth noticing: fixed3d went from fork to v1.3.0 in four days
(created 2026-07-12, [v1.3.0](https://github.com/mas-bandwidth/fixed3d/releases/tag/v1.3.0)
released 2026-07-16). A full physics engine converted to fixed point, tested,
benchmarked, and released — that pace is what the AI collaboration makes
possible. Every commit is public. Check the history yourself.

There will be more work like this. Support the Patreon and you'll see more of it.

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

## The craft, plainly

These days the open source happens as a disclosed collaboration between Glenn
and an AI assistant. Two promises, so nobody has to wonder:

**Glenn writes his own articles. Always.** Every word at gafferongames.com and
mas-bandwidth.com is his — no AI, ever. If you read an article, he wrote it.

**AI-assisted code gets more care than the hand-written code did, not less.**
That's a claim you can check, not take on faith: this year's collaboration put
five fuzz targets over every untrusted parser in yojimbo, ran
multi-million-iteration sanitizer soaks before every release, and found real
bugs in code that had shipped in real games for close to a decade — including
a remotely reachable heap overflow present in every yojimbo release since
2019. Every fix has a regression test. The token cost of the collaboration is
a line item in the [ledger](https://github.com/mas-bandwidth/open-ledger)
like everything else.

**[The full list of bugs found and fixed with the help of AI →](BUGS.md)**
If you are using older versions of the libraries, upgrade now.

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

## Crediting

If you ship a product using any of these libraries, please credit each one you
use in your product credits:

> **Más Bandwidth LLC**
> yojimbo — Glenn Fiedler

For fixed3d, please also credit **Box3D — Erin Catto** — it's a conversion of
her engine. The licenses don't require any of this (the BSD copyright notice
in your distribution is the only legal requirement, and that stands). The
credits listing is an official request, made here and in each library's
README. Fair credit keeps open source honest.

**Credit what you ship.** Some libraries include others — we track this for
you:

| If you use | Credit |
|---|---|
| **yojimbo** | yojimbo, netcode, reliable, serialize — it bundles all four |
| **netcode**, **reliable**, or **serialize** standalone | just that library |
| **fixed3d** | fixed3d — Glenn Fiedler, plus **Box3D — Erin Catto** |

**New libraries ship under the [Más Bandwidth Source License (MBSL)](MBSL.md)**
— BSD 3-Clause plus one clause: the credit above becomes part of the license.
The existing libraries keep their BSD-3/MIT licenses unchanged, forever; for
them the credit is an official request and the expected standard.

## The lines that don't move

- The articles stay free. The existing libraries stay MIT/BSD — forever.
  (New libraries may ship under the [MBSL](MBSL.md): source open, free to
  use, credit required.)
- Glenn always writes his own articles — never AI-ghostwritten.
- The work funded here is public. The transparency is the product.
