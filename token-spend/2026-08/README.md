# August 2026

In August, schema went from its first release to production ready to nine
languages. The serialize wire went from six implementations to nine, all held
to identical bits. The networking libraries, yojimbo, netcode and reliable,
got strict builds on every compiler, first C# releases, and landed in Homebrew,
apt and vcpkg. If you write a game in more than one language, or ship a
client and server that have to agree on every bit, that is the work.

Where the AI tokens went, by repository, and what got done in each one:

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

## schema

Schema went from its first release to production ready to nine languages in
one month. Version 1.0.0 shipped on August 11: a language for declaring a
game's constants, enums and data types, compiled to bit-packed reading and
writing code in C++, C#, Go and Rust. Fixed point and 128-bit integers joined
the language, with an unsigned fixed-point type beside them, and a ranged
field costs exactly the bits its range needs. C became the reference backend
and JavaScript the sixth. The compiler became a library with a public API.
The generated code got a fuzzer, a benchmark harness, and gates that refuse
any change that makes it slower or moves a byte on the wire without saying so.
Version 2.0.0 on August 25 was the production-ready release. Version 2.1.0 on
August 30 added Dart, Java and Elixir: nine languages, one schema, identical
bits.

## serialize

The bit-level wire format underneath schema, in C++. Thirteen releases in
August. Fixed point and 128-bit integers on the wire. One decode on every
architecture, with the wire contract written down as a standard that the
other implementations are held to. Readers now refuse malformed strings. The
read and write paths force inlining, which took a bulk write from 314
instructions to 40. Compressed floats are precomputed, and the bits on the
wire no longer depend on the compiler's floating-point contraction settings.

## serialize.c

The whole wire in a header-only C library, first released August 13, nine
releases by the end of the month. A validating reader, the same forced
inlining as C++ for a 34 to 49 percent gain on x86 and about three times on
arm64, and Zig and Odin bindings.

## serialize.rs

Twelve releases and a 2.0 line. Reads inline end to end in safe Rust, writes
that cannot fail in safe Rust, and the crate is on crates.io as
serialize-official.

## serialize.js

First released August 17, three releases. The sixth implementation of the
wire, with a production mode that made it 1.76 times faster.

## serialize.go

Ten releases. Fixed point and 128-bit integers, validating string readers,
and the family's cross-language range clamp proven with witness bytes on the
wire.

## serialize.cs

The wire in C#, first released August 13, nine releases. Both of C#'s
rounding modes handled correctly, UTF-16 wide strings, and API-misuse checks
that compile out of release builds to match the C++ reference.

## serialize.dart, serialize.java, serialize.elixir

First releases on August 30. The same wire, native on Dart, on the JVM, and on
the BEAM, each one byte-identical to the rest of the family.

## netcode

netcode 1.4.4: strict floating-point builds on every target, so the same source
behaves the same on every compiler. netcode.cs 1.0.0, the first C# release, and
the Rust port brought current at 1.1.1. Homebrew, apt and vcpkg now ship it,
and the Conan, FreeBSD, Debian and OpenBSD submissions moved through review.

## yojimbo

Five releases. Version 1.10 lets yojimbo's encrypted datagrams ride a
transport you supply, so a relay or an ICE route works without touching the
connection, channel or encryption code. Version 1.11 takes the optimized
serialization core from the serialize work above. The wire now builds strict
on every target.

## reliable

reliable 1.4.1: strict floating-point builds on every target. reliable.cs 1.0.0,
the first C# release, and reliable.rs took the same work in seven pull requests
without a release yet. Shipping in Homebrew, apt and vcpkg beside the others.
