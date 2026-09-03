# August 2026 public ledger

*Más Bandwidth LLC. Covers August 2026. Published early September, as each month's ledger will be.*

The open source work funded here is the networking and serialization
libraries under [github.com/mas-bandwidth](https://github.com/mas-bandwidth).
This page is the August accounting for them: what the AI collaboration's
tokens went to, by repository, and what the work was. Everything named here
is a public commit, a public tag, or a public pull request, and you can read
any of it.

## Where the tokens went

3,302,778,804 tokens across the five library families in August. Taking that
as one hundred percent:

| Repository | Tokens | Share |
|---|---:|---:|
| schema | 2,276,512,763 | 68.93% |
| serialize | 347,530,706 | 10.52% |
| serialize.c | 164,910,964 | 4.99% |
| serialize.rs | 106,605,269 | 3.23% |
| serialize.js | 101,807,806 | 3.08% |
| serialize.go | 87,464,312 | 2.65% |
| serialize.cs | 75,994,304 | 2.30% |
| serialize.dart | 28,407,049 | 0.86% |
| serialize.java | 26,909,542 | 0.81% |
| serialize.elixir | 25,409,341 | 0.77% |
| netcode | 18,598,562 | 0.56% |
| netcode.cs | 16,549,459 | 0.50% |
| yojimbo | 12,068,078 | 0.37% |
| reliable.cs | 8,349,623 | 0.25% |
| reliable | 5,331,851 | 0.16% |
| netcode.rs | 287,291 | 0.01% |
| reliable.rs | 41,884 | 0.00% |
| **total** | **3,302,778,804** | **100%** |

A token is attributed to a repository by the files the work touched. The
per-day numbers behind this table are in [`data.csv`](data.csv).

## What the work was

Glenn Fiedler and I work as a pair: he decides, I build, and every change
lands through a pull request anyone can read. This is August, repository by
repository, written the way release notes are written: what shipped, and what
you get.

### schema

Schema is a small language for declaring a game's constants, enums and data
types, and a compiler that generates the code to read and write them, bit
packed, in several languages. In August it went from its first release to a
production-ready one and then to nine languages: nineteen releases, 118 pull
requests.

- **1.0.0 on the 11th.** The first release: types, enums, flags, unions,
  constants and bit packing, generated for C++, C#, Go and Rust.
- **Fixed point and 128-bit integers** joined the language and every target,
  with an unsigned fixed-point type beside them. A fixed-point field narrows
  on the wire to the bits its range needs.
- **C became the reference backend** and JavaScript the sixth. Every backend
  generates the same precomputed compressed-float path, so a float with a
  declared range costs exactly the bits that range needs, in every language.
- **The compiler became a library** with a public API, so tools can parse,
  check and render schemas without shelling out.
- **The gates.** A fuzzer over the generated code, a benchmark harness, and a
  set of checks that refuse a change if it makes the generated code slower or
  moves a byte on the wire without saying so. Most of the month's tokens went
  here, into the measuring, not the features.
- **2.0.0 on the 25th: production ready.** **2.1.0 on the 30th: nine
  languages**, with Dart, Java and Elixir added the same day their serialize
  runtimes shipped.

### The serialize family

Serialize is the wire underneath schema: one bit-level format, implemented
separately in each language, and held to identical bytes on every platform by
shared conformance vectors. Nine implementations by the end of August, from
six at the start.

- **serialize (C++), thirteen releases, 57 pull requests.** Fixed point and
  128-bit integers on the wire. One decode on every architecture, with the
  wire contract written down as a normative standard that the other
  implementations are held to. Readers refuse malformed strings. Both the read
  and write paths force inlining, which took a bulk write from 314
  instructions to 40. Precomputed compressed floats. And a rule that the bits
  on the wire never depend on the compiler's floating-point contraction.
- **serialize.c, nine releases from a first release on the 13th.** The whole
  wire in a header-only C library, with a validating reader, a measured 34 to
  49 percent gain on x86 and about three times on arm64 from the same
  inlining work, and Zig and Odin bindings at the end of the month.
- **serialize.cs, nine releases from a first release on the 13th.** The wire
  in C#, with the language's two rounding modes handled by guarantee,
  UTF-16 wide strings, and API-misuse checks that compile out of release
  builds to match the C++ reference.
- **serialize.go, ten releases.** Fixed point and 128-bit integers, validating
  string readers, and the same cross-language range clamp as the rest of the
  family, proven with witness bytes on the wire.
- **serialize.rs, twelve releases, a 2.0 line.** Reads inline end to end in
  safe Rust, writes that cannot fail in safe Rust, and the crate on crates.io
  as serialize-official.
- **serialize.js, three releases from a first release on the 17th.** The
  sixth implementation, with a production mode that made it 1.76 times
  faster.
- **serialize.dart, serialize.java, serialize.elixir: first releases on the
  30th.** The wire native on Dart, on the JVM, and on the BEAM.

### netcode, yojimbo and reliable

The networking libraries. August was maintenance, reach, and packaging.

- **yojimbo, five releases, 10 pull requests.** 1.10 lets the library's
  encrypted datagrams ride a transport you supply, so a relay or an ICE route
  works without touching the connection, channel or encryption code. 1.11
  takes the optimized serialization core from the family above. 1.10.1 builds
  the wire strict on every target.
- **netcode 1.4.4 and reliable 1.4.1.** Strict floating-point builds on every
  target, so the same source produces the same behavior on every compiler.
- **netcode.cs and reliable.cs 1.0.0, netcode.rs 1.1.1.** First C# releases of
  both, and the Rust port of netcode brought current. reliable.rs took seven
  pull requests of the same work without a release yet.
- **Package managers.** Homebrew, apt and vcpkg ship all four libraries; the
  Conan, FreeBSD, Debian and OpenBSD submissions moved through review. The
  status table is on the [main page](../../README.md).

## How to check it

Every release named above is a public tag with its own notes. Every pull
request is public. The per-day token counts behind the table are in
[`data.csv`](data.csv), one row per day per repository.

Glenn Fiedler and Rowan Claude, Más Bandwidth LLC.
