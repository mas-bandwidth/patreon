# August 2026 token spend, by repository

Data for the question "how much of the AI spend went into the open source?"

Both benches are measured and merged: the per-day, per-repository numbers for every open
source repo, the method that produced them, and the tools that ran it. The `macbook` bench was
measured and published first; the `studio` bench ran the same pipeline over its own transcript
store on 2026-09-02 and the union below is the month's number.

Everything below is *tokens*, not dollars, and every open source figure is a **floor** — the
attribution rule is deliberately conservative and gives ambiguous work to the bucket that
makes the open source share smaller. The caveats section says exactly how much sits on that
edge.

---

## 1. The month

Window `2026-08-01` .. `2026-08-31` inclusive, days counted in `Australia/Sydney`, two benches.

| | macbook | studio | **August 2026** |
|---|---:|---:|---:|
| Transcript files read | 1,426 (0 skipped) | 1,736 (0 skipped) | 3,162 |
| Folded turns | 79,579 | 47,567 | **127,146** |
| Total tokens | 21,408,017,663 | 9,914,497,601 | **31,322,515,264** |
| Open source, absolute-path rule (**the published floor**) | 3,812,106,080 — 17.807% | 635,479,430 — 6.410% | **4,447,585,510 — 14.199%** |
| Open source, tilde-path sensitivity arm (upper end of the range) | 4,987,293,753 — 23.296% | 1,023,483,639 — 10.323% | 6,010,777,392 — 19.190% |
| Everything else (private work), one aggregate line per day | 17,595,911,583 — 82.193% | 9,279,018,171 — 93.590% | 26,874,929,754 — 85.801% |

So the honest statement for the month is: **between 14.2% and 19.2% of August's tokens went
into open source repositories, and 14.2% is the number that survives the strictest reading.**

The two benches are not alike, and the difference is real rather than a measurement artifact.
The macbook is where the attended sessions ran — schema, the serialize family, fixed3d, all
worked directly in their checkouts. The Studio was provisioned on 2026-08-10 (its first row is
that day) and runs the unattended estate: the sleep cycle, the patrols, the consumers, the
self's own maintenance. Most of that is private by construction, and some of the open source
work it does do runs from scratch directories the fail-closed rule cannot credit (caveats,
below).

On each bench, two independent passes over the same window agree to the token: the
attribution pass and the shared accounting scanner both report the bench's total and turn
count, delta 0. The run refuses to print a table if they ever disagree.

### Open source repositories, August 2026, both benches

| Repository | Tokens | Share of month | Days | Turns | Benches |
|---|---:|---:|---:|---:|---|
| schema | 2,276,512,763 | 7.268% | 23 | 10,832 | macbook, studio |
| nova | 512,080,996 | 1.635% | 21 | 4,188 | macbook, studio |
| serialize | 347,530,706 | 1.110% | 22 | 2,204 | macbook, studio |
| fixed3d | 235,218,401 | 0.751% | 11 | 811 | macbook, studio |
| nova-tools | 175,225,124 | 0.559% | 12 | 1,305 | macbook, studio |
| serialize.c | 164,910,964 | 0.526% | 12 | 1,121 | macbook, studio |
| fixed | 146,353,057 | 0.467% | 6 | 605 | macbook |
| serialize.rs | 106,605,269 | 0.340% | 15 | 738 | macbook, studio |
| serialize.js | 101,807,806 | 0.325% | 5 | 833 | macbook, studio |
| serialize.go | 87,464,312 | 0.279% | 14 | 638 | macbook, studio |
| serialize.cs | 75,994,304 | 0.243% | 13 | 541 | macbook, studio |
| serialize.dart | 28,407,049 | 0.091% | 2 | 134 | macbook |
| serialize.java | 26,909,542 | 0.086% | 3 | 126 | macbook |
| serialize.elixir | 25,409,341 | 0.081% | 2 | 126 | macbook |
| netcode | 18,598,562 | 0.059% | 4 | 143 | macbook |
| netcode.cs | 16,549,459 | 0.053% | 3 | 57 | macbook |
| gafferongames | 14,241,826 | 0.045% | 2 | 94 | macbook |
| yojimbo | 12,068,078 | 0.039% | 6 | 49 | macbook |
| awesome-persistent-ai | 11,109,854 | 0.035% | 3 | 33 | macbook |
| reliable.cs | 8,349,623 | 0.027% | 2 | 33 | macbook |
| apt | 7,931,667 | 0.025% | 3 | 40 | macbook, studio |
| gafferongames/vcpkg | 6,865,203 | 0.022% | 6 | 47 | macbook |
| antirez/ds4 | 5,680,687 | 0.018% | 2 | 25 | studio |
| hydrogen | 5,398,443 | 0.017% | 3 | 60 | macbook |
| reliable | 5,331,851 | 0.017% | 3 | 50 | macbook |
| next | 5,131,580 | 0.016% | 1 | 10 | macbook |
| microsoft/vcpkg | 4,055,336 | 0.013% | 1 | 33 | macbook |
| proton | 3,998,250 | 0.013% | 3 | 17 | macbook |
| patreon | 3,923,677 | 0.013% | 1 | 52 | macbook |
| table | 3,132,652 | 0.010% | 1 | 15 | macbook |
| jedisct1/libhydrogen | 1,253,707 | 0.004% | 1 | 13 | macbook |
| gafferongames-new | 1,237,021 | 0.004% | 2 | 4 | macbook |
| homebrew-tap | 1,119,414 | 0.004% | 1 | 9 | macbook |
| fast3d | 849,811 | 0.003% | 2 | 2 | macbook |
| netcode.rs | 287,291 | 0.001% | 1 | 10 | macbook |
| reliable.rs | 41,884 | 0.000% | 1 | 1 | macbook |
| **all open source** | **4,447,585,510** | **14.199%** | | | |
| **everything else (aggregate)** | **26,874,929,754** | **85.801%** | | | |

`schema` is not a mistake and not double counting: it is the largest single body of work in
the month by a wide margin, and it is open source. Repos printed with an owner prefix
(`microsoft/vcpkg`, `jedisct1/libhydrogen`, `gafferongames/vcpkg`, `antirez/ds4`) are
upstreams and forks worked on directly, not repos under this org. `Days` counts distinct local
days across both benches, so a repo worked on by both benches on one day counts that day once.

### The macbook bench alone

| Repository | Tokens | Share of bench | Days | Turns |
|---|---:|---:|---:|---:|
| schema | 2,184,788,860 | 10.205% | 21 | 10,282 |
| serialize | 338,471,069 | 1.581% | 21 | 2,153 |
| fixed3d | 199,954,969 | 0.934% | 5 | 652 |
| nova | 182,897,097 | 0.854% | 8 | 1,300 |
| serialize.c | 149,195,477 | 0.697% | 8 | 1,043 |
| fixed | 146,353,057 | 0.684% | 6 | 605 |
| serialize.js | 101,745,457 | 0.475% | 5 | 831 |
| serialize.rs | 95,999,459 | 0.448% | 11 | 688 |
| serialize.cs | 75,156,171 | 0.351% | 12 | 531 |
| serialize.go | 73,272,314 | 0.342% | 12 | 559 |
| nova-tools | 52,323,723 | 0.244% | 2 | 340 |
| serialize.dart | 28,407,049 | 0.133% | 2 | 134 |
| serialize.java | 26,909,542 | 0.126% | 3 | 126 |
| serialize.elixir | 25,409,341 | 0.119% | 2 | 126 |
| netcode | 18,598,562 | 0.087% | 4 | 143 |
| netcode.cs | 16,549,459 | 0.077% | 3 | 57 |
| gafferongames | 14,241,826 | 0.067% | 2 | 94 |
| yojimbo | 12,068,078 | 0.056% | 6 | 49 |
| awesome-persistent-ai | 11,109,854 | 0.052% | 3 | 33 |
| reliable.cs | 8,349,623 | 0.039% | 2 | 33 |
| apt | 7,678,973 | 0.036% | 2 | 39 |
| gafferongames/vcpkg | 6,865,203 | 0.032% | 6 | 47 |
| hydrogen | 5,398,443 | 0.025% | 3 | 60 |
| reliable | 5,331,851 | 0.025% | 3 | 50 |
| next | 5,131,580 | 0.024% | 1 | 10 |
| microsoft/vcpkg | 4,055,336 | 0.019% | 1 | 33 |
| proton | 3,998,250 | 0.019% | 3 | 17 |
| patreon | 3,923,677 | 0.018% | 1 | 52 |
| table | 3,132,652 | 0.015% | 1 | 15 |
| jedisct1/libhydrogen | 1,253,707 | 0.006% | 1 | 13 |
| gafferongames-new | 1,237,021 | 0.006% | 2 | 4 |
| homebrew-tap | 1,119,414 | 0.005% | 1 | 9 |
| fast3d | 849,811 | 0.004% | 2 | 2 |
| netcode.rs | 287,291 | 0.001% | 1 | 10 |
| reliable.rs | 41,884 | 0.000% | 1 | 1 |
| **all open source** | **3,812,106,080** | **17.807%** | | |
| **everything else (aggregate)** | **17,595,911,583** | **82.193%** | | |

### The studio bench alone

Window `2026-08-10` .. `2026-08-31` in practice: the bench did not exist before the 10th.

| Repository | Tokens | Share of bench | Days | Turns |
|---|---:|---:|---:|---:|
| nova | 329,183,899 | 3.320% | 13 | 2,888 |
| nova-tools | 122,901,401 | 1.240% | 10 | 965 |
| schema | 91,723,903 | 0.925% | 13 | 550 |
| fixed3d | 35,263,432 | 0.356% | 9 | 159 |
| serialize.c | 15,715,487 | 0.159% | 7 | 78 |
| serialize.go | 14,191,998 | 0.143% | 4 | 79 |
| serialize.rs | 10,605,810 | 0.107% | 6 | 50 |
| serialize | 9,059,637 | 0.091% | 9 | 51 |
| antirez/ds4 | 5,680,687 | 0.057% | 2 | 25 |
| serialize.cs | 838,133 | 0.008% | 4 | 10 |
| apt | 252,694 | 0.003% | 1 | 1 |
| serialize.js | 62,349 | 0.001% | 2 | 2 |
| **all open source** | **635,479,430** | **6.410%** | | |
| **everything else (aggregate)** | **9,279,018,171** | **93.590%** | | |

### `data.csv`

One row per `(bench, day, repo)`. Columns:

```
bench,day,repo,visibility,attribution,input,output,cache_creation,cache_read,total,turns,sessions_touching,gcs_mirrored_total,gcs_coverage
```

* Every **open source** row is published verbatim, exactly as the tool emitted it.
* Every **private** row is collapsed into one aggregate line per `(bench, day)` whose `repo`
  and `visibility` are both literally `PRIVATE` and whose `attribution` is `aggregate`. Its
  token columns and `turns` are the sums of the rows it replaces; `sessions_touching` is left
  blank because summing it across repos would double count a session that touched several.
  Private repository names are not published, here or anywhere in this directory.
* `total` = `input + output + cache_creation + cache_read`. Summing `total` over every row in
  the file gives **31,322,515,264** and summing `turns` gives **127,146** — the whole month on
  both benches, matching section 1. The file is complete, not a sample. Summed per bench it
  gives each bench's own line in section 1.
* `attribution` records how the directory-to-repo mapping was obtained: `remote` (read off
  that checkout's git origin), `heuristic` (checkout gone from disk, folded onto a verified
  sibling or a known repo name), `unresolved` (none of the above; classified private
  fail-closed), `aggregate` (a collapsed private line).
* `gcs_mirrored_total` / `gcs_coverage` say how much of a **macbook** row was also present in
  the object-store mirror of that bench's sessions (`full` / `partial` / `none`, `mixed` on an
  aggregate line). Macbook rows from the second half of the month read `none` because
  mirroring had not caught up. Every `studio` row reads `0` / `none`: the mirror holds only
  the macbook's sessions (section 3), so nothing of the Studio's was ever in it. Either way
  it is a statement about the mirror, not about the measurement, and the mirror changes no
  number here.

macbook: 167 open source rows + 28 aggregate lines = 195 rows. studio: 80 open source rows +
22 aggregate lines = 102 rows. **297 rows in all.** No `(bench, day, repo)` tuple appears
twice; the merge tool refuses to write if one does.

---

## 2. The method

### The unit: `(bench, day, repo)`

The uniqueness tuple is `(bench, day, repo)`. A bench is one machine's transcript store. This
is the ruling that makes merging across benches safe: two benches never produce the same
tuple for different work, so a union is a union and not a guess.

### A turn is one folded assistant message

Token accounting is per **turn**, where a turn is one assistant message identified by its
message id, taking the **largest** usage record seen for that id.

**This matters more than anything else in the pipeline.** Streaming writes several transcript
records for the same assistant message, each carrying a cumulative usage block. Naively
summing every usage record in the transcripts overcounts August by **1.81x**. Any independent
recount that does not fold by message id will not agree with these numbers, and the
disagreement will be in the recount.

### Days are local days, from UTC timestamps

Transcript timestamps are UTC (RFC3339). Each turn is converted into `Australia/Sydney` and
booked to the local calendar day. A month boundary is a local midnight, not a UTC one. The
zone is a flag; whatever zone is chosen must be the same on both benches or the day keys will
not line up. It was.

### Attribution: paths in tool calls and tool results, never the `cwd`

A turn is attributed to a repository by the paths it actually touched:

1. A turn's **segment** is every transcript record since the previous usage-carrying record
   (exclusive) through the usage-carrying record itself (inclusive) — the tool results the
   turn read plus the tool calls it made. Every record in the file belongs to exactly one
   segment, so nothing is counted twice or dropped.
2. Inside a segment, **only tool calls and tool results are scanned.** Every occurrence of the
   literal `<working-root>/<dir>` (i.e. `/Users/<user>/rowan-working/<dir>`) is one reference
   to `<dir>`.
3. `<dir>` is folded to its **GitHub repo** before anything is compared (next section).
4. The turn is attributed **whole** to the repo with the most references. No references, or a
   tie for the most, falls back to the session's own `cwd` repo.

**Why tool calls and results only, and never the `cwd` field.** This was found by
instrumenting the run, not predicted. Scanning the whole record instead put a single private
orchestration checkout in 79,345 of 79,579 macbook turns and handed it 97.6% of every token it
was seen near. The reason: *every* transcript record carries a `cwd` field, so a raw scan
re-imports the session's working directory as if it were a thing the turn touched. Almost all
of this estate's work runs from orchestrator sessions standing in one private directory, so a
cwd-based attribution books a month of open source work as private — an artifact of where the
shell was standing, not of what was worked on. Injected attachments do the same thing more
quietly. A turn's own prose is excluded for the same reason: reasoning that merely *names* a
path is not evidence the path was worked on.

The first run of this analysis made exactly that mistake and reported the month as 99.996%
private. That number was wrong, and the correction is the whole reason this pipeline exists.

### Clone-to-upstream folding, keyed to the GitHub repo

**The GitHub repo is the key; the directory is only how it is recovered.** Every checkout
still on disk has its `origin` remote read, and the URL is reduced to `owner/name` (bare
`name` when the owner is this org). A directory that no longer exists is folded onto the
longest still-verified checkout whose name it extends on a `-` boundary (with and without a
`scratch-` prefix), and only failing that onto a table of known repo names.

Two consequences worth stating:

* Worktrees and scratch clones fold onto the repo they came from, so `schema-mm-pr6` and
  `scratch-schema-x` are both `schema`. Without this, one month's work is scattered across
  dozens of dead directory names.
* A directory whose *name* looks like a repo but whose *remote* says otherwise follows the
  remote. `apt-schema` is `apt`, not `schema`; a pure name heuristic gets that backwards.

Folding happens **before** the argmax, deliberately: three references to `schema-bench` and
two to `schema-gowrite` are five references to `schema`, not a two-way split that loses to
something else.

A name that answers to no remote, no verified sibling and no known repo is left as itself and
**reported** as unresolved — never quietly bundled into a repo it might not belong to.

### Open source vs private

Visibility comes from `gh repo view <name> --json isPrivate`, cached to a
`repo-visibility.json` beside the output (one `"<name>": true|false` entry per repo name;
`true` means public). It is a **cache of real `gh` answers, not a guess**. A repo name
missing from the cache is listed in the run output and bucketed **private, fail-closed**.
That is why every figure here is a floor: an unclassified name can only ever push the open
source share down. Each bench built its own cache from its own run's list of names; neither
cache is published, because each is a list of repository names.

### Conservation

Attribution moves tokens between buckets. It must not create or destroy any. So the run
re-derives the same window through the shared accounting scanner and **refuses to print a
table** unless the two grand totals *and* the two turn counts agree exactly. They do, on both
benches: `21,408,017,663` / `79,579` on the macbook and `9,914,497,601` / `47,567` on the
Studio, both ways.

There is a second, deliberately independent check (`tools/august-verify-attrib.go`) that
recomputes one day from scratch — its own walk, its own record struct, its own UTC-to-local
conversion, its own message-id fold, its own path scan, its own git remote read, sharing no
code with the main pass. If two programs built from the same parts agree, the agreement means
nothing; these were not. On the Studio it was pointed at 2026-08-30, the bench's busiest day
by tokens (5,878 turns, 1,237,391,611 tokens; 2026-08-11 has more turns): both programs
report that total and that turn count, and agree on every open source repo's tokens except
for **one turn** — the `(none)` reference edge described in the caveats.

The merge itself (`tools/august-merge-benches.go`) re-checks the published file's invariants
rather than trusting them (aggregate shape, and every total equal to its parts), refuses a
raw input that repeats a tuple, applies the privacy filter to the raw side, refuses to write
if the filter changed the raw grand total or if any tuple collides between the two sides —
identical or not — and then re-reads the file it wrote from disk, sums that, and refuses
unless it equals published plus raw exactly. So the figures in section 1 come from the
published file and not from the slice that produced it. Each refusal was exercised with a
deliberately broken input before the tool was trusted: a duplicated private raw row, a
duplicated open source raw row, and a duplicated published row all exit non-zero with no
file written.

---

## 3. How the second bench was completed and merged

### Step 1 — the pipeline on the Studio's own store

The two analysis tools in `tools/` were dropped unchanged into `rowan-tools/scratch/` and run
from that repo's root (they import `github.com/mas-bandwidth/rowan-tools/modules/{spend,usage}`
from the private workshop repo). The visibility cache was built first from one real
`gh repo view <name> --json isPrivate` per name the first run reported as unclassified that
looked like one of this org's repositories (the cwd labels and the scratch directory names
were not queried and stay private, fail-closed). Then, as published:

```
go run scratch/august-repo-spend-paths.go -bench studio \
    -lo 2026-08-01 -hi 2026-08-31 -zone Australia/Sydney \
    -out <out> -csv august-2026-studio-bench-day-repo-bypath.csv
go run scratch/august-repo-spend-paths.go -bench studio -tilde \
    -lo 2026-08-01 -hi 2026-08-31 -zone Australia/Sydney \
    -out <out> -csv august-2026-studio-bench-day-repo-bypath-TILDE.csv
go run scratch/august-verify-attrib.go -day 2026-08-30 -zone Australia/Sydney
```

The whole pass over 1.1 GB of transcripts takes about ten seconds on this bench. Read the
run's own output, not just the CSV: it prints the reconciliation, the unclassified names, the
ambiguity census, the touched-vs-won table, the directory-to-repo mapping actually used, and
what fraction of the answer rests on a mapping that was not remote-verified. On the Studio
that last figure is **4,493,901,345 tokens, 45.33%** — but almost all of it is the two
cwd-fallback labels for sessions standing outside any checkout (below), which are not repos
and were never going to be remote-verified. The open source rows resting on a heuristic
mapping are `fixed3d` (35.3M, a checkout under a port directory), `apt` (0.25M) and
`serialize.js` (0.06M). The macbook's figure for the same question was 131,601,040 tokens,
0.61%.

### Step 2 — the object store adds nothing

The macbook bench's GCS sessions mirror (bucket name held internally) holds **only that
bench's own mirrored sessions**. A full pass over the mirror measured 10,315,509,634 tokens
across 43,197 turns — a strict subset of the macbook's 21,408,017,663, not additional data.
It contains nothing from the Studio (`rowan-tools#46`). The Studio's `~/.claude/projects` was
the only place its half existed, and that is what was measured.

### Step 3 — merge

```
go run scratch/august-merge-benches.go \
    -published data.csv \
    -raw <out>/august-2026-studio-bench-day-repo-bypath.csv \
    -raw-tilde <out>/august-2026-studio-bench-day-repo-bypath-TILDE.csv \
    -published-tilde-oss 4987293753 \
    -out merged.csv
```

The merge is a **union of rows on `(bench, day, repo)`**:

* Rows whose tuples do not collide are simply concatenated. None collided; the benches are
  disjoint by construction and the tool checks anyway.
* **An overlapping tuple must agree byte-exact.** Two benches producing different numbers for
  the same `(bench, day, repo)` is not something to average or reconcile — it means one of the
  runs is wrong, and the tool exits non-zero and prints both rows.
* Before writing, the same privacy filter used on the macbook side is applied to the raw
  Studio rows: open source rows verbatim, every private row collapsed into one `PRIVATE`
  aggregate line per `(bench, day)`. The tool refuses to write unless the filtered rows still
  sum to the raw grand total and turn count. They do: `9,914,497,601` / `47,567`, delta 0.
* The tilde arm is read only for the range figure and is never merged into the file.

---

## 4. Caveats

**The open source figures are floors, and here is the size of the edge on each bench.**

*macbook.* In the published (absolute-path) arm, 34,331 of 79,579 turns — 43.14% of turns,
9,028,682,297 tokens, 42.17% of the bench — had no absolute path reference in their tool calls
or results and were decided by the session's `cwd` fallback. In the tilde arm, where
`~/rowan-working/...` references are promoted to real ones, the fallback share drops to 26,344
turns (33.10%) and 6,705,453,475 tokens (31.32%), and the open source share rises from 17.8%
to 23.3%.

*studio.* 26,226 of 47,567 turns — 55.13%, 5,572,529,433 tokens, 56.21% of the bench — fell to
the fallback. In the tilde arm that drops to 21,815 turns (45.86%) and 4,251,200,414 tokens
(42.88%), and the open source share rises from 6.4% to 10.3%. The fallback bites harder here
because of how the Studio works: its unattended roles run from scratch working directories
outside any checkout, so a turn with no path reference has no repo to fall back to at all and
is labelled as standing outside the working root. Those two labels together hold 4.42 billion
tokens, 44.5% of the bench, all booked private.

**The true number is somewhere in the 14.2%–19.2% band and the low end is what gets
published.**

Other things a reader should know:

* **Winner-take-all.** A turn is attributed whole to one repo. macbook: 6,180 turns (7.77%,
  2,329,986,005 tokens) referenced two or more repos and only one was credited; 1,615 tied and
  fell to the fallback. studio: 3,851 turns (8.10%, 845,041,981 tokens) referenced two or
  more; 1,294 tied. Map iteration order never picks a winner.
* **Open source work in unresolved scratch directories, studio.** The Studio does its
  benchmarking runs, its packaging patrol and its bench provisioning in scratch directories
  that are not git checkouts and whose names match no repo. The rule leaves them unresolved
  and books them private, as designed. Measured, with the rule stated: eleven names carry
  `unresolved` attribution in the raw output; two are private repositories, one is a path
  fragment, and the
  remaining **eight scratch directories hold 41,569,026 tokens, 0.419% of the Studio bench**,
  almost all of it open source work. They are not reclassified here, because reclassifying by
  hand is exactly what the fail-closed rule exists to refuse.
* **Five studio days carry turns and no tokens** (2026-08-16, 17, 18, 19 and 27).
  Every usage record on those days is a `<synthetic>` refusal — the bench's usage allowance
  was exhausted, and each scheduled role that woke got a one-line refusal instead of a model.
  The turns are real and are counted; the zero is the honest cost of a run that did no work.
  Nothing from those days is missing or elsewhere.
* **One `(none)` reference edge, found by the independent check.** An elided path such as
  `<working-root>/...` in a tool result is counted by the main pass as a reference to a
  directory named `(none)`, which can then *tie* with a real repo and push the turn to the
  fallback; the independent check drops it, which is what the stated rule implies. Measured
  over the whole Studio month: it decides **exactly one turn, 441,180 tokens (0.0045% of the
  bench)**, which would otherwise be `nova-tools`. The published tool is left as it ran on
  both benches so the two halves follow the identical rule; the fix is queued for the next
  month's run, and the macbook data was not re-measured for this class.
* **A sensitivity that did not matter.** Counting distinct files instead of every occurrence
  would have changed the winner in 32 macbook turns (0.04%, 6,970,222 tokens) and 22 studio
  turns (0.05%, 5,339,671 tokens). Measured and reported, not applied.
* **These are list-price tokens, not dollars.** No conversion to money is made here.
  Converting requires per-model rates, and one model id in this window (`<synthetic>`) is not
  in the rate table at all — anything priced for it would be a guess, not a measurement.
* **Cache reads dominate the raw counts.** `cache_read` is the bulk of every row. That is
  normal for long agentic sessions and it is left visible rather than netted out, because
  netting it out is an editorial choice and the columns let anyone make it themselves.
* **The Studio's month is 22 days.** Its first row is 2026-08-10, the day the bench was
  provisioned. Nothing is missing before that; the bench did not exist.
