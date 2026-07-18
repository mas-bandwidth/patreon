# The bugs found and fixed with the help of AI

In 2026 the open source under [github.com/mas-bandwidth](https://github.com/mas-bandwidth)
became a disclosed collaboration between Glenn Fiedler and an AI collaborator
(Claude Code, Fable 5). One of the first things the collaboration did was turn
serious tooling on code that had shipped for years: libFuzzer targets over
every untrusted parser, AddressSanitizer/UBSan/MemorySanitizer CI, multi-million
iteration soaks, and line-by-line review.

This page is the honest record of what that found. Every item links to the
release or pull request that fixed it, so you can verify each one yourself.

**If you are using older versions of these libraries, upgrade now.** The
strongest items below are real security vulnerabilities that shipped in
hand-written code for years — including a remotely reachable heap overflow
that existed in every yojimbo release since at least 2019. Every fix has a
regression test.

| Library | Upgrade to | Why |
|---|---|---|
| **yojimbo** | v1.5.0 or later (latest is best) | remote heap overflow, wire-reachable asserts, union misread — all pre-1.5.0 versions affected |
| **netcode** | v1.3.0 or later | replay-protection integer overflow, Debug-vs-Release memory-safety gap |
| **reliable** | v1.3.2 or later | fragment-reassembly integer overflow |
| **serialize** | v1.4.3 or later | fuzz-verified + golden wire format pinned on every platform |

*No double counting: each bug is listed once, under the library whose code
contained it. yojimbo ships netcode, reliable, and serialize inside it — so
upgrading yojimbo picks up their fixes too.*

---

## yojimbo

Found by the fuzzing campaign and review in [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0)
(July 2026) — five libFuzzer targets over every untrusted parser, validated by
a 5.5-million-iteration ASan/UBSan soak at high packet loss:

- **Remote heap overflow in block reassembly** — a connected peer could send
  an over-long final fragment and corrupt heap memory (the copy ran before the
  size check). Confirmed with AddressSanitizer; the vulnerable code existed in
  every release since at least v1.2.1 (2019). Fixed with a regression test
  that trips ASan against the pre-fix source.
  ([#283](https://github.com/mas-bandwidth/yojimbo/pull/283)) — *security*
- **`ChannelPacketData` union misread** — message/block union confusion on the
  packet read path; made safe by construction (de-unioned).
  ([#280](https://github.com/mas-bandwidth/yojimbo/pull/280)) — *security*
- **Remote DoS via debug asserts reachable from the wire** — an over-budget
  packet on the read path could assert-kill a debug server.
  ([#281](https://github.com/mas-bandwidth/yojimbo/pull/281)) — *security (debug builds)*
- **Uninitialized read of `block.messageType`** on a disabled-blocks channel.
  ([#272](https://github.com/mas-bandwidth/yojimbo/pull/272)) — *memory safety*
- **Undefined behavior from unaligned bit-reader reads** on the packet receive
  path. ([#251](https://github.com/mas-bandwidth/yojimbo/pull/251)) — *UB*
- **Allocation-failure crashes** — two OOM crashes under adversarial block
  fragments and allocator exhaustion, found by fuzzing allocation failures.
  ([#286](https://github.com/mas-bandwidth/yojimbo/pull/286),
  [#284](https://github.com/mas-bandwidth/yojimbo/pull/284),
  [#285](https://github.com/mas-bandwidth/yojimbo/pull/285)) — *robustness*
- **Block fragment-count overflow** and a **network-simulator packet leak**.
  ([#254](https://github.com/mas-bandwidth/yojimbo/pull/254)) — *correctness*
- **Matcher expiry check was always true** — token expiry was effectively not
  enforced in the matcher example; also hardened error/key defaults.
  ([#261](https://github.com/mas-bandwidth/yojimbo/pull/261)) — *security (example code)*
- **NULL-client crashes** in `Client::Connect` / `ConnectLoopback` on socket
  failure and `Client::IsLoopback` when disconnected.
  ([#287](https://github.com/mas-bandwidth/yojimbo/pull/287),
  [#288](https://github.com/mas-bandwidth/yojimbo/pull/288)) — *robustness*
- **Packet-size budgeting overflow (caught before it ever shipped)** — during
  the v1.6.0 serialize upgrade, a `maxPacketSize` not divisible by 8 would
  have overflowed the packet buffer by up to 7 bytes in release builds; caught
  in development by serialize's new debug assert and fixed in the same
  release. Listed for completeness — no shipped version was affected.
  ([v1.6.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.6.0)) — *correctness*

## netcode

Fixed in [v1.3.0](https://github.com/mas-bandwidth/netcode/releases/tag/v1.3.0) (July 2026):

- **Integer overflow in replay protection** — sequence numbers near the top of
  the 64-bit space could be falsely rejected. Found by fuzzing; rewritten
  overflow-safe. The buggy line is verifiable in the v1.2.1 tag from 2019 —
  it shipped for roughly seven years. — *correctness*
- **Debug-vs-Release memory-safety gap** — public entry points guarded only by
  asserts (e.g. an out-of-range `max_clients` to `netcode_server_start`)
  compiled to no protection at all in release builds; now runtime bounds
  guards. — *security hardening*
- **Port parsing truncated instead of validating**; a zeroed config crashed
  instead of getting default allocators; `netcode_init`/`netcode_term` are now
  reference counted. — *correctness*

## reliable

- **Integer overflow in the fragment reassembly buffer-size computation** —
  now computed in `size_t` arithmetic. Reported by **@orbisai0security**
  (external report — credited where credit is due; fixed and shipped by the
  collaboration in
  [v1.3.2](https://github.com/mas-bandwidth/reliable/releases/tag/v1.3.2)). — *security*
- Hardening in [v1.3.0](https://github.com/mas-bandwidth/reliable/releases/tag/v1.3.0):
  **duplicate packets within the receive window are now detected and dropped**,
  fragments with non-canonical headers are rejected, fragments for
  already-received packets are dropped on arrival, and the send path no longer
  allocates. — *robustness*

## serialize

The honest note first: the campaign found no security vulnerability in
serialize's own shipped code. What it added is the machinery that *proves*
that, on every push:

- A **golden wire-format test** pins the exact bytes the serializer produces —
  Linux, macOS, Windows and big-endian s390x (GCC under QEMU) all produce and
  decode byte-identical wire data.
  ([v1.3.0](https://github.com/mas-bandwidth/serialize/releases/tag/v1.3.0))
- A **libFuzzer harness** runs hostile reads of arbitrary bytes through every
  ReadStream primitive plus a differential write→read round trip — 60 seconds
  on every push, 1 hour nightly with a cumulative corpus.
- The 2026 line also removed the 256 MB buffer limit, made reads ~4× faster,
  and tightened the buffer contracts so misuse fails loudly in debug builds
  ([v1.4.0](https://github.com/mas-bandwidth/serialize/releases/tag/v1.4.0)–[v1.4.3](https://github.com/mas-bandwidth/serialize/releases/tag/v1.4.3)).

---

## How this was done

Fuzzing (five libFuzzer targets on yojimbo alone, including stateful and
structure-aware targets over the connection deserializer), sanitizers in CI
(ASan, UBSan, LSan, MemorySanitizer), multi-million-iteration soaks at high
packet loss before every tag, golden wire-format tests, and review. All of it
is in the public CI configuration and release notes of each repo — nothing on
this page is taken on trust.

Fair credit runs both directions: where a bug was reported from outside the
collaboration, the reporter is named above. And the hand-written code this
page scrutinizes held up remarkably well for its era — it shipped in real
games for a decade. The point is not that it was bad. The point is that
tested is better, and now it is tested.

*This work is funded by the [Más Bandwidth Patreon](https://patreon.com/MasBandwidth).
Supporting it is how more of it happens.*
