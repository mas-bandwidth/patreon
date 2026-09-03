# Más Bandwidth on Patreon

If this work matters to you, please support it: **[Become a supporter](https://www.patreon.com/MasBandwidth/membership)**

> ## Security notice: upgrade your libraries
>
> The 2026 hardening work found and fixed real security bugs in code that had
> shipped for years, including a **remotely reachable heap overflow in yojimbo
> present in every release since 2019**. If you use any of these libraries in a
> product, upgrade to the latest release now:
>
> | Library | Upgrade | What it fixes |
> |---|---|---|
> | **yojimbo** | [latest release](https://github.com/mas-bandwidth/yojimbo/releases/latest) | remote heap overflow, wire-reachable asserts, union misread (v1.5.0), AEAD nonce-reuse on restart (v1.7.0) |
> | **netcode** | [latest release](https://github.com/mas-bandwidth/netcode/releases/latest) | AEAD nonce-reuse on server restart (v1.4.0), replay-protection overflow, memory-safety hardening |
> | **reliable** | [latest release](https://github.com/mas-bandwidth/reliable/releases/latest) | read-buffer over-read, integer overflow |
> | **serialize** | [latest release](https://github.com/mas-bandwidth/serialize/releases/latest) | fuzz-hardened, wire format pinned on every platform |
>
> Every fix has a regression test, and the full honest accounting is here:
> **[the bugs found and fixed with the help of AI](BUGS.md)**

---

## Introduction

**Creating open source game networking code, and the articles that explain it.**

For twenty years I've written about how multiplayer games work and given the
code away. The articles are free. The libraries are open source and running in
real games right now. This Patreon supports the continued
work, and every month an open ledger shows exactly where the money goes.

This repo is the public reference for what the Patreon funds. The monthly books
live in **[open-ledger](https://github.com/mas-bandwidth/open-ledger)**.

## The articles

- **[gafferongames.com](https://gafferongames.com)**, twenty years of articles
  on game networking, physics, and how multiplayer actually works. Free, and
  staying free.
- **[mas-bandwidth.com](https://mas-bandwidth.com)**, where the new writing
  lives.

## The open source work so far

The reference libraries under [github.com/mas-bandwidth](https://github.com/mas-bandwidth),
open source, used in shipped games. Each library states its license in its own
repository. *Versions and star counts as of 3 September 2026.*

| Library | What it does | Latest | Stars |
|---|---|---|---|
| **[yojimbo](https://github.com/mas-bandwidth/yojimbo)** | Client/server network protocol for games, encrypted and dedicated-server ready | v1.11.0 | about 2.7k |
| **[netcode](https://github.com/mas-bandwidth/netcode)** | Secure client/server connection over UDP (connect tokens, encryption) | v1.4.4 | about 2.6k |
| **[reliable](https://github.com/mas-bandwidth/reliable)** | Reliable-ordered messages and acks over UDP | v1.4.1 | about 650 |
| **[serialize](https://github.com/mas-bandwidth/serialize)** | Bitpacking and serialization, one unified read/write path | v1.15.0 | about 140 |
| **[schema](https://github.com/mas-bandwidth/schema)** | The schema language for games: constants, enums and data types compiled to nine languages | v2.4.0 | new |
| **[fixed](https://github.com/mas-bandwidth/fixed)** | Deterministic Q48.16 fixed-point math: scalars, vectors, quaternions, transforms | v1.4.0 | new |
| **[fixed3d](https://github.com/mas-bandwidth/fixed3d)** | Cross-platform **deterministic** fixed-point physics | v1.4.0 | new |

**Ports to other languages**, so the protocols are not C-only. netcode and
reliable each have C#, Go and Rust ports:
[netcode.cs](https://github.com/mas-bandwidth/netcode.cs),
[netcode.go](https://github.com/mas-bandwidth/netcode.go),
[netcode.rs](https://github.com/mas-bandwidth/netcode.rs),
[reliable.cs](https://github.com/mas-bandwidth/reliable.cs),
[reliable.go](https://github.com/mas-bandwidth/reliable.go),
[reliable.rs](https://github.com/mas-bandwidth/reliable.rs).
The serialize family covers nine languages, below.

## What just landed

**[schema](https://github.com/mas-bandwidth/schema)** is the schema language
for games. You declare your constants, enums and data types once, and the
compiler generates the code that reads and writes them, bitpacked, in nine
languages: C, C++, C#, Dart, Elixir, Go, Java, JavaScript and Rust. Every
language agrees on every bit, and a protocol id computed from the declaration
refuses a mismatched build before a byte is misread. The compiler is a library
with a public API, the generated code is fuzz tested and benchmarked per
language, and every release is gated against the last one so it cannot get
slower or move a byte on the wire without saying so.
[v2.0.0](https://github.com/mas-bandwidth/schema/releases/tag/v2.0.0) on
25 August was the production-ready line.
[v2.1.0](https://github.com/mas-bandwidth/schema/releases/tag/v2.1.0) made it
nine languages.
[v2.3.0](https://github.com/mas-bandwidth/schema/releases/tag/v2.3.0) added
tables, a second kind of declaration for data that has to survive schema
changes, and
[v2.4.0](https://github.com/mas-bandwidth/schema/releases/tag/v2.4.0) made
the generated Dart and JavaScript faster.

Underneath it, **serialize is nine implementations of one wire format**:
[C++](https://github.com/mas-bandwidth/serialize),
[C](https://github.com/mas-bandwidth/serialize.c),
[C#](https://github.com/mas-bandwidth/serialize.cs),
[Dart](https://github.com/mas-bandwidth/serialize.dart),
[Elixir](https://github.com/mas-bandwidth/serialize.elixir),
[Go](https://github.com/mas-bandwidth/serialize.go),
[Java](https://github.com/mas-bandwidth/serialize.java),
[JavaScript](https://github.com/mas-bandwidth/serialize.js) and
[Rust](https://github.com/mas-bandwidth/serialize.rs). The wire contract is
written down as a normative standard in the C++ repository, and byte-pinned
conformance vectors hold every implementation to identical bytes on every
platform and architecture. Fixed point and 128-bit integers ride the same
wire, and the bits never depend on the compiler's floating-point mood. Dart,
Java and Elixir shipped their first releases on 30 August, the same day
schema learned to generate for them.

**[fixed](https://github.com/mas-bandwidth/fixed)** is a standalone
deterministic Q48.16 fixed-point math library: scalar math, vectors,
quaternions, matrices, transforms, and the wide world-position and AABB
families. fixed3d vendors it at a pinned version, with a vendor-drift CI
workflow that goes red if the vendored tree ever diverges from the pin.
[v1.4.0](https://github.com/mas-bandwidth/fixed/releases) is production
ready.

And the networking libraries: **[yojimbo v1.11](https://github.com/mas-bandwidth/yojimbo/releases)**
takes the optimized serialization core from the family above, after v1.10
let its encrypted datagrams ride a transport you supply. netcode 1.4.4 and
reliable 1.4.1 build strict on every target, and their C# and Rust ports
shipped first releases in August.

## What is coming in September

Schema is becoming more than bitpacked types. The bitpacked wire is the right
tool for a packet between a client and a server that ship together. It is the
wrong tool for a save game, a message between a tool and a backend that ship
months apart, an asset file the game should map and point at, or the render
data C++ hands to C# sixty times a second. Until now those needed something
else, usually Protocol Buffers or FlatBuffers beside schema. The September
push is to make schema do all of it, so one declaration serves every data
type in a game.

The pieces are on main and being finished across all nine languages:

- **Tables** carry field ids and lengths, so any build reads any data. An
  unknown field is skipped, a missing one takes its default, and a report says
  what changed. That is the message format and the save game format.
- **A committed baseline** refuses, at compile time, the few edits that change
  what old data means without changing a byte, so a save game never breaks
  silently.
- **Enum-keyed arrays**, arrays indexed by an enum that size themselves and
  cannot be indexed by nothing.
- **The block form**, every fixed table laid out for another language to
  point at, generated on both sides and asserted at compile time. Render data
  from C++ to Unity's C# with no copy and no parse.
- **The cook**, a table's data written in one build's exact memory layout so
  the game maps the file and points at it. A gigabyte opens as fast as a
  kilobyte. This one ships as a preview.
- **JSON in, binary out**, so data authored as text becomes a compact table
  any language reads, and **reflection descriptors** so an editor can walk any
  table it has never seen.

All of it is dogfooded in a real game first, the same one the networking
libraries run in, and it ships as schema 3.0.0 when every language has it and
every language has been stress tested and profiled. The README on the
[schema repository](https://github.com/mas-bandwidth/schema) is being
rewritten around these use cases, and the release notes will say plainly
what you get.

## fixed3d

**[fixed3d](https://github.com/mas-bandwidth/fixed3d)**
takes Erin Catto's [Box3D](https://github.com/erincatto/box3d) and tears every
`float` out of the simulation, replacing them with Q48.16 fixed point. Uniform
1/65536 resolution across a vast world, about 140 trillion meters in each direction. Bit-exact on every platform.
All 22 test suites pass, and the benchmarks are published honestly, about 2x
slower than float, and the README tells you straight when you should keep
using Box3D instead.

Here's the part worth noticing: fixed3d went from fork to v1.3.0 in four days
(created 2026-07-12, [v1.3.0](https://github.com/mas-bandwidth/fixed3d/releases/tag/v1.3.0)
released 2026-07-16). A full physics engine converted to fixed point, tested,
benchmarked, and released. That pace is what the AI collaboration makes
possible. Every commit is public. Check the history yourself.

There will be more work like this. Support the Patreon and you'll see more of it.

## Getting the libraries into every package manager

Part of the work nobody sees: getting these libraries into the package managers
so people can just `install` them instead of vendoring source. Each ecosystem
has its own submission process, maintainers, and review. This is the ongoing
push. *Status as of 3 September 2026. The links are the live source of truth.*

| Package manager | Status | Where it stands |
|---|---|---|
| **Homebrew** | **Shipped** | `serialize` and `libyojimbo` merged into homebrew-core (PRs [#292317](https://github.com/Homebrew/homebrew-core/pull/292317), [#292681](https://github.com/Homebrew/homebrew-core/pull/292681)), with formulae for all four libraries. Also a tap: [mas-bandwidth/homebrew-tap](https://github.com/mas-bandwidth/homebrew-tap). |
| **apt (Debian/Ubuntu)** | **Shipped** | Our own apt repository serves `.deb` packages for all four libraries: [mas-bandwidth/apt](https://github.com/mas-bandwidth/apt). |
| **vcpkg** | **Shipped** | All four libraries merged in one PR ([microsoft/vcpkg#52858](https://github.com/microsoft/vcpkg/pull/52858), merged 2026-08-10), as the ports `mas-bandwidth-serialize`, `mas-bandwidth-reliable`, `mas-bandwidth-netcode` and `mas-bandwidth-yojimbo`. |
| **Conan** | **In review** | The netcode recipe is merged ([#30730](https://github.com/conan-io/conan-center-index/pull/30730)). Three PRs remain open on conan-center-index: yojimbo ([#30676](https://github.com/conan-io/conan-center-index/pull/30676)), serialize ([#30728](https://github.com/conan-io/conan-center-index/pull/30728)) and reliable ([#30729](https://github.com/conan-io/conan-center-index/pull/30729)). |
| **FreeBSD** | **Landed; updates in review** | All four ports accepted (bugs [296779](https://bugs.freebsd.org/bugzilla/show_bug.cgi?id=296779)-[296782](https://bugs.freebsd.org/bugzilla/show_bug.cgi?id=296782) closed FIXED). netcode is updated to 1.4.3 in the ports tree. Updates for yojimbo, reliable and serialize are open as [freebsd-ports#574](https://github.com/freebsd/freebsd-ports/pull/574)-[#576](https://github.com/freebsd/freebsd-ports/pull/576). |
| **Debian** | **In progress** | ITP/RFS filed via mentors.debian.net; reviewer feedback addressed and functional autopkgtests added; awaiting a sponsoring Debian Developer to push to unstable. |
| **OpenBSD** | **In progress** | `[NEW]` port submission on the ports@ mailing list; refreshed ports at the latest releases prepared for a re-roll. |

Homebrew, apt and vcpkg are done. The rest are moving through their review
pipelines.
Progress shows up in the monthly [ledger](https://github.com/mas-bandwidth/open-ledger).

Over **6,000 GitHub stars** across the core libraries, and a steady stream of
releases. The [open ledger](https://github.com/mas-bandwidth/open-ledger)
tracks each month's work with commit ranges you can verify.

## The craft, plainly

These days the open source happens as a disclosed collaboration between Glenn
and an AI assistant. Two promises, so nobody has to wonder:

**Glenn writes his own articles. Always.** Every word at gafferongames.com and
mas-bandwidth.com is his. No AI, ever. If you read an article, he wrote it.

**AI-assisted code gets more care than the hand-written code did, not less.**
That's a claim you can check, not take on faith: this year's collaboration put
five fuzz targets over every untrusted parser in yojimbo, ran
multi-million-iteration sanitizer soaks before every release, and found real
bugs in code that had shipped in real games for close to a decade, including
a remotely reachable heap overflow present in every yojimbo release since
2019. Every fix has a regression test. The token cost of the collaboration is
a line item in the [ledger](https://github.com/mas-bandwidth/open-ledger)
like everything else.

**[The full list of bugs found and fixed with the help of AI](BUGS.md)**
If you are using older versions of the libraries, upgrade now.

## Public ledgers

Early each month, a ledger for the previous month: where the AI collaboration's
tokens went across the library families, and what the work was, written like
release notes.

| Month | The ledger |
|---|---|
| August 2026 | [The August 2026 public ledger](token-spend/2026-08/README.md) |

## The tiers

| | | |
|---|---|---|
| **Free** | $0 | The posts: new articles and the monthly ledger, as they land. |
| **Thanks** | $1 | For twenty years of free articles. |
| **Supporter** | $5 | Keeps the articles and libraries funded. |
| **Ledger** | $10 | The full itemized monthly books. |
| **Workshop** | $25 | Early draft access and office hours. No treadmill. |
| **Patron** | $50 | For those who want the work to exist and can say so. |
| **Commercial** | $100 | Making money on the libraries? This is how you give back. |
| **Sponsor** | $15,000 | Fund the work directly and help steer priorities. Email first. |

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

> yojimbo, serialize, reliable, netcode by Glenn Fiedler and Rowan Claude

For fixed3d, please also credit **Box3D by Erin Catto**. It is a conversion of
his engine. The licenses don't require any of this. The credits listing is an
official request, made here and in each library's README. Fair credit keeps
open source honest.

**Credit what you ship.** Some libraries include others, so we track this for
you:

| If you use | Credit |
|---|---|
| **yojimbo** | yojimbo, netcode, reliable, serialize (it bundles all four) |
| **netcode**, **reliable**, or **serialize** standalone | just that library |
| **fixed3d** | fixed3d by Glenn Fiedler and Rowan Claude, plus **Box3D by Erin Catto** |

## The lines that don't move

- The articles stay free.
- Glenn always writes his own articles, never AI-ghostwritten.
- The work funded here is public. The transparency is the product.
