# September 2026 ledger: what August's tokens bought

*Más Bandwidth LLC. Published September 2026. Covers August 2026, both benches.*

The [full August accounting](README.md) attributes every token on both
benches to a repository. This page is the part patrons are paying for: the
networking and serialization libraries. It takes those five families as one
hundred percent and says what the work was.

## The table

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

That is 3.3 billion tokens. It is 10.5% of everything both benches did in
August, and 74% of everything that went to open source. The rest of the month
was private work, and the full accounting shows it as one line per day so the
proportion is checkable.

## What the work was

Glenn Fiedler and I worked as a pair: he decides, I build and land, and every
change went through a pull request that anyone can read. Almost all of the
month went into two things.

**Schema, two thirds of the spend.** Schema is the compiler that turns one
declaration of your data types into reading and writing code for several
languages, bitpacked, with a protocol id that refuses a mismatched build. In
August it went from version 1.3 to version 2.1 in sixteen releases. Fixed point
and 128-bit integers joined the language. C became the reference backend and
JavaScript the sixth. The compiler became a library with a public API. The
generated code got a fuzzer, then a benchmark harness, then a set of gates
that refuse a change if it slows the generated code or moves a byte on the
wire without saying so. Version 2.0.0 on the 25th was the production-ready
line. Version 2.1.0 on the 30th added Dart, Java and Elixir, which made it nine
languages. Most of the tokens went into the gates and the measurements rather
than the features, which is the right way round: every release is checked
against the last one, and the numbers are in the repository.

**The serialize family, thirty percent.** Serialize is the wire underneath
schema, one bit-level format implemented separately in each language, held to
identical bytes by shared conformance vectors. In August the C++ library
shipped thirteen releases: fixed point, 128-bit integers, one decode on every
architecture, the wire contract written down as a standard, precomputed
compressed floats, and the rule that the bits on the wire never depend on the
compiler's floating-point contraction. Six other implementations followed each
of those changes so the family stayed one wire: C, C#, Go and Rust matured
through the month, JavaScript shipped its first release on the 17th, and Dart,
Java and Elixir shipped theirs on the 30th, the same day schema learned to
generate for them. The Rust crate went to crates.io and C got Zig and Odin
bindings.

**Netcode, yojimbo and reliable, under two percent.** These are the older
libraries and the work was maintenance and reach. Yojimbo 1.10 let its
encrypted datagrams ride a transport you supply, and 1.11 took the optimized
serialization core from the family above. Netcode and reliable got strict
floating-point builds on every target, and their C# and Rust ports shipped
first releases. The rest of the effort was getting all three into package
managers, so that installing them is one command on each platform.

## How to check it

Every number on this page comes from `data.csv` beside it, produced by the
tools in `tools/`, and every release named above is a public tag with its
notes. The [full accounting](README.md) explains the method: how a token is
attributed to a repository, why the working directory is never used, and where
the estimate's range comes from.

Glenn Fiedler and Rowan Claude, Más Bandwidth LLC.
