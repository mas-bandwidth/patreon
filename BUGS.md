# The bugs found and fixed with the help of AI

In 2026 the open source under [github.com/mas-bandwidth](https://github.com/mas-bandwidth)
became a disclosed collaboration between Glenn Fiedler and an AI collaborator
(Claude Code, Fable 5). One of the first things the collaboration did was turn
serious tooling on code that had shipped for years: libFuzzer targets over
every untrusted parser, AddressSanitizer / UBSan / MemorySanitizer CI,
multi-million-iteration soaks, and a line-by-line security audit.

This page is the honest record of what that found, **built from the actual
commit diffs**, not release notes. Every row links to the commit that fixed it
and to the release it shipped in, so you can verify each one yourself.

**43 fixes across the four libraries, 27 of them reachable from the network** (a
crafted packet or untrusted input can trigger them). The most serious is a
remotely reachable heap overflow in yojimbo present in every release since 2019.
Every fix has a regression test.

## If you use these libraries in a product, upgrade now

| Library | Upgrade to | The reason |
|---|---|---|
| **yojimbo** | [v1.7.0](https://github.com/mas-bandwidth/yojimbo/releases/latest) | remote heap overflow, uninitialized-union reads, wire-reachable DoS asserts; also ships the netcode AEAD fix |
| **netcode** | [v1.4.0](https://github.com/mas-bandwidth/netcode/releases/latest) | AEAD nonce-reuse on server restart; release-build bounds guards |
| **reliable** | [v1.3.5](https://github.com/mas-bandwidth/reliable/releases/latest) | read-buffer over-read on the fragment path; reassembly integer overflow |
| **serialize** | [v1.4.4](https://github.com/mas-bandwidth/serialize/releases/latest) | read-past-end guard, wire-value validation, unaligned-read UB |

*No double counting: each bug is listed once, under the library whose code
contained it. yojimbo ships netcode, reliable and serialize inside it, so
upgrading yojimbo picks up their fixes too.*

**Column key.** *Type:* Write OOB / Read OOB = out-of-bounds memory access · UB =
undefined behavior · DoS / crash = a packet can crash or hang the process ·
Validation = missing check on untrusted input · Crypto · Correctness ·
Robustness = leak / allocator-exhaustion hardening. *Wire:* ✓ means an attacker
can trigger it with a crafted packet or untrusted input over the network; — means
it needs local API misuse or is on the send/write path.

---

## netcode

*7 fixes, 2 reachable from the network. netcode's read path was already solid;
most of these are defense-in-depth for release builds and integrators.*

| Commit | Bug | Type | Severity | Wire | Fixed in |
|---|---|---|---|:--:|---|
| [`dc21b70`](https://github.com/mas-bandwidth/netcode/commit/dc21b70) | AEAD nonce reuse: a stopped-and-restarted server re-encrypts global packets at sequence 0 under the same key as per-client packets | Crypto | High | ✓ | [v1.4.0](https://github.com/mas-bandwidth/netcode/releases/tag/v1.4.0) |
| [`7c638bd`](https://github.com/mas-bandwidth/netcode/commit/7c638bd) | Replay-protection add overflowed near UINT64_MAX, falsely rejecting legitimate high-sequence packets | Correctness | Low | ✓ | [v1.3.0](https://github.com/mas-bandwidth/netcode/releases/tag/v1.3.0) |
| [`7ef3604`](https://github.com/mas-bandwidth/netcode/commit/7ef3604) | Release builds skipped assert-only bounds checks, so an out-of-range client_index / max_clients indexed per-client arrays out of bounds | Write OOB | Medium | — | [v1.3.0](https://github.com/mas-bandwidth/netcode/releases/tag/v1.3.0) |
| [`8283fae`](https://github.com/mas-bandwidth/netcode/commit/8283fae) | An oversized packet to send_packet overflowed a fixed on-stack buffer in release builds (guard was assert-only) | Write OOB | Medium | — | [v1.3.0](https://github.com/mas-bandwidth/netcode/releases/tag/v1.3.0) |
| [`478f5e7`](https://github.com/mas-bandwidth/netcode/commit/478f5e7) | Port parsing used atoi(), silently truncating out-of-range or non-numeric ports instead of rejecting them | Validation | Low | — | [v1.3.0](https://github.com/mas-bandwidth/netcode/releases/tag/v1.3.0) |
| [`8283fae`](https://github.com/mas-bandwidth/netcode/commit/8283fae) | Bracketed-IPv6 port scan matched the colons inside the address, corrupting the parsed address string | Correctness | Low | — | [v1.3.0](https://github.com/mas-bandwidth/netcode/releases/tag/v1.3.0) |
| [`c14229c`](https://github.com/mas-bandwidth/netcode/commit/c14229c) | Public process_packet lacked a size/running guard, reading packet_data[0] on a 0-byte packet (affects direct API callers only) | DoS / crash | Low | — | [v1.3.0](https://github.com/mas-bandwidth/netcode/releases/tag/v1.3.0) |

## reliable

*7 fixes, 6 reachable from the network.*

| Commit | Bug | Type | Severity | Wire | Fixed in |
|---|---|---|---|:--:|---|
| [`f0e3be1`](https://github.com/mas-bandwidth/reliable/commit/f0e3be1) | Fragment-header read passed the full packet length after advancing past the 5-byte header, over-reading the receive buffer | Read OOB | High | ✓ | [v1.3.0](https://github.com/mas-bandwidth/reliable/releases/tag/v1.3.0) |
| [`d747162`](https://github.com/mas-bandwidth/reliable/commit/d747162) | Reassembly buffer size computed in signed int could overflow (reported by [@orbisai0security](https://github.com/orbisai0security)) | UB | Medium | ✓ | [v1.3.2](https://github.com/mas-bandwidth/reliable/releases/tag/v1.3.2) |
| [`42e6725`](https://github.com/mas-bandwidth/reliable/commit/42e6725) | Fragment-0 payload length validated against the wire byte-count instead of the re-encoded header size | Correctness | Low | ✓ | [v1.3.0](https://github.com/mas-bandwidth/reliable/releases/tag/v1.3.0) |
| [`2ec8b0f`](https://github.com/mas-bandwidth/reliable/commit/2ec8b0f) | Packets duplicated or replayed within the receive window were delivered twice | Correctness | Low | ✓ | [v1.3.0](https://github.com/mas-bandwidth/reliable/releases/tag/v1.3.0) |
| [`1f40e42`](https://github.com/mas-bandwidth/reliable/commit/1f40e42) | A duplicate final fragment arriving after reassembly bypassed the duplicate check | Correctness | Low | ✓ | [v1.3.0](https://github.com/mas-bandwidth/reliable/releases/tag/v1.3.0) |
| [`f65b33b`](https://github.com/mas-bandwidth/reliable/commit/f65b33b) | Bandwidth stats divided by a near-zero interval, producing NaN/Inf (the guard only checked exactly zero) | Correctness | Low | ✓ | [v1.3.0](https://github.com/mas-bandwidth/reliable/releases/tag/v1.3.0) |
| [`f65b33b`](https://github.com/mas-bandwidth/reliable/commit/f65b33b) | A reassembly entry's packet_header_bytes was read before fragment 0 initialized it | Read OOB | Low | — | [v1.3.0](https://github.com/mas-bandwidth/reliable/releases/tag/v1.3.0) |

## serialize

*11 fixes, 7 reachable from the network. No known exploit in shipped code; these
are the read-path validation and undefined-behavior fixes the audit and fuzzing
turned up, on a 2016-era bit-packing library.*

| Commit | Bug | Type | Severity | Wire | Fixed in |
|---|---|---|---|:--:|---|
| [`e7c7c33`](https://github.com/mas-bandwidth/serialize/commit/e7c7c33) | SerializeBytes computed its read-past-end guard with a signed byte count and no lower-bound check | Read OOB | High | ✓ | [v1.3.0](https://github.com/mas-bandwidth/serialize/releases/tag/v1.3.0) |
| [`049d42b`](https://github.com/mas-bandwidth/serialize/commit/049d42b) | BitReader did aligned 32-bit loads on a possibly-unaligned packet pointer | UB | Medium | ✓ | [v1.3.0](https://github.com/mas-bandwidth/serialize/releases/tag/v1.3.0) |
| [`e7c7c33`](https://github.com/mas-bandwidth/serialize/commit/e7c7c33) | Compressed-float read never validated the quantized integer taken from the wire | Validation | Medium | ✓ | [v1.3.0](https://github.com/mas-bandwidth/serialize/releases/tag/v1.3.0) |
| [`2daaba3`](https://github.com/mas-bandwidth/serialize/commit/2daaba3) | SerializeInteger reconstructed values with signed arithmetic that could overflow | UB | Medium | ✓ | [v1.3.0](https://github.com/mas-bandwidth/serialize/releases/tag/v1.3.0) |
| [`e7c7c33`](https://github.com/mas-bandwidth/serialize/commit/e7c7c33) | SerializeInteger returned out-of-range values without rejecting them | Validation | Low | ✓ | [v1.3.0](https://github.com/mas-bandwidth/serialize/releases/tag/v1.3.0) |
| [`e7c7c33`](https://github.com/mas-bandwidth/serialize/commit/e7c7c33) | int_relative's 32-bit fallback read a raw value without enforcing the encoding's range | Validation | Low | ✓ | [v1.3.0](https://github.com/mas-bandwidth/serialize/releases/tag/v1.3.0) |
| [`2c9501b`](https://github.com/mas-bandwidth/serialize/commit/2c9501b) | int_relative reconstructed values in signed arithmetic that could overflow | UB | Low | ✓ | [v1.3.0](https://github.com/mas-bandwidth/serialize/releases/tag/v1.3.0) |
| [`2daaba3`](https://github.com/mas-bandwidth/serialize/commit/2daaba3) | Compressed-float max integer value computed with no clamp (write path) | UB | Low | — | [v1.3.0](https://github.com/mas-bandwidth/serialize/releases/tag/v1.3.0) |
| [`3124bbb`](https://github.com/mas-bandwidth/serialize/commit/3124bbb) | Compressed-float normalized value not fully bounded before scaling (write path) | UB | Low | — | [v1.3.0](https://github.com/mas-bandwidth/serialize/releases/tag/v1.3.0) |
| [`e7c7c33`](https://github.com/mas-bandwidth/serialize/commit/e7c7c33) | Zig-zag helpers left-shifted a signed int (signed-shift UB) | UB | Low | — | [v1.3.0](https://github.com/mas-bandwidth/serialize/releases/tag/v1.3.0) |
| [`b732fd8`](https://github.com/mas-bandwidth/serialize/commit/b732fd8) | BitWriter stored through an aligned 32-bit pointer that could be unaligned | UB | Low | — | [v1.3.0](https://github.com/mas-bandwidth/serialize/releases/tag/v1.3.0) |

## yojimbo

*18 fixes, 12 reachable from the network.*

| Commit | Bug | Type | Severity | Wire | Fixed in |
|---|---|---|---|:--:|---|
| [`d3d4f33`](https://github.com/mas-bandwidth/yojimbo/commit/d3d4f33) | Remote heap overflow: a block fragment was memcpy'd into the reassembly buffer before its size was validated (present since 2019) | Write OOB | Critical | ✓ | [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0) |
| [`50d2ae0`](https://github.com/mas-bandwidth/yojimbo/commit/50d2ae0) | Unreliable channel read the messages array without checking the block-message flag, writing past it | Write OOB | High | ✓ | [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0) |
| [`b2fccb2`](https://github.com/mas-bandwidth/yojimbo/commit/b2fccb2) | ChannelPacketData::Initialize zeroed only part of the message/block union, leaving a stale pointer to be read | Read OOB | High | ✓ | [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0) |
| [`fc90d41`](https://github.com/mas-bandwidth/yojimbo/commit/fc90d41) | BitReader did 32-bit word loads on an unaligned reliable-payload pointer | UB | Medium | ✓ | [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0) |
| [`1b6bec0`](https://github.com/mas-bandwidth/yojimbo/commit/1b6bec0) | An over-budget packet on the read path tripped a debug assert (remote DoS) | DoS / crash | Medium | ✓ | [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0) |
| [`50d2ae0`](https://github.com/mas-bandwidth/yojimbo/commit/50d2ae0) | A zero-payload packet tripped a bufferSize>0 assert in the bit reader (remote DoS) | DoS / crash | Medium | ✓ | [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0) |
| [`760ebc2`](https://github.com/mas-bandwidth/yojimbo/commit/760ebc2) | Fragments-per-block used floor division, under-allocating the fragment array by one | Write OOB | Medium | — | [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0) |
| [`663acf5`](https://github.com/mas-bandwidth/yojimbo/commit/663acf5) | Packet write size not rounded to a multiple of 8, so the qword flush wrote up to 7 bytes past the buffer | Write OOB | Medium | — | [v1.6.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.6.0) |
| [`b722e9a`](https://github.com/mas-bandwidth/yojimbo/commit/b722e9a) | Non-zero block fragments read block.messageType uninitialized (only fragment 0 serializes it) | UB | Low | ✓ | [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0) |
| [`c27a0dd`](https://github.com/mas-bandwidth/yojimbo/commit/c27a0dd) | Read-path message-array allocation not NULL-checked before use (crash under allocator exhaustion) | DoS / crash | Low | ✓ | [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0) |
| [`f73569a`](https://github.com/mas-bandwidth/yojimbo/commit/f73569a) | Placement-new ran on a NULL pointer when allocation failed, before the caller's NULL check | DoS / crash | Low | ✓ | [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0) |
| [`da31341`](https://github.com/mas-bandwidth/yojimbo/commit/da31341) | Reassembly-buffer allocation result unchecked; a failed alloc led to a bad store | DoS / crash | Low | ✓ | [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0) |
| [`50d2ae0`](https://github.com/mas-bandwidth/yojimbo/commit/50d2ae0) | A read-path error branch leaked a just-allocated Message | Robustness | Low | ✓ | [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0) |
| [`696bf36`](https://github.com/mas-bandwidth/yojimbo/commit/696bf36) | ProcessAck never set acked=true, leaving the duplicate-ack guard assert reachable | Correctness | Low | ✓ | [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0) |
| [`26f51aa`](https://github.com/mas-bandwidth/yojimbo/commit/26f51aa) | Send-path allocator-exhaustion leaks in GeneratePacket and channel-data allocation | Robustness | Low | — | [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0) |
| [`760ebc2`](https://github.com/mas-bandwidth/yojimbo/commit/760ebc2) | The network simulator's duplicate-packet branch leaked an existing packet buffer | Robustness | Low | — | [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0) |
| [`06de25b`](https://github.com/mas-bandwidth/yojimbo/commit/06de25b) | maxPacketFragments used integer division inside ceil(), so ceil() was a no-op | Correctness | Low | — | [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0) |
| [`0dfef69`](https://github.com/mas-bandwidth/yojimbo/commit/0dfef69) | Address::Parse used atoi(), silently truncating out-of-range ports | Validation | Low | — | [v1.5.0](https://github.com/mas-bandwidth/yojimbo/releases/tag/v1.5.0) |

---

## How this was done

Fuzzing (five libFuzzer targets on yojimbo alone, including stateful and
structure-aware targets over the connection deserializer), sanitizers in CI
(ASan, UBSan, LSan, MemorySanitizer), multi-million-iteration soaks at high
packet loss before every tag, golden wire-format tests, and a line-by-line
security audit. All of it is in the public CI configuration and commit history
of each repo. This page is built from the commits themselves, so nothing on it
is taken on trust: click any hash and read the diff.

Fair credit runs both directions. Where a bug was reported from outside the
collaboration, the reporter is named in the table (the reliable reassembly
overflow was reported by @orbisai0security). And the hand-written code this page
scrutinizes held up remarkably well for its era: it shipped in real games for a
decade. The point is not that it was bad. The point is that tested is better, and
now it is tested.

*This work is funded by the [Más Bandwidth Patreon](https://patreon.com/MasBandwidth).
Supporting it is how more of it happens.*
