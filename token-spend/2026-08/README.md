# August 2026 token spend, by repository

Preliminary data for the question "how much of the AI spend went into the open source?"

One bench's half of August 2026 is finished and is published here in full: the per-day,
per-repository numbers for every open source repo, the method that produced them, and the
tools that ran it. The other bench's half has not been measured yet. **This is a handover:
the Studio bench runs the same pipeline over its own transcript store and merges the result,
and the union is the month's number.**

Everything below is *tokens*, not dollars, and every open source figure is a **floor** — the
attribution rule is deliberately conservative and gives ambiguous work to the bucket that
makes the open source share smaller. The caveats section says exactly how much sits on that
edge.

---

## 1. Current state

Measured bench: `macbook`. Window `2026-08-01` .. `2026-08-31` inclusive, days counted in
`Australia/Sydney`.

| | |
|---|---|
| Transcript files read | 1,426 (0 skipped) |
| Folded turns | 79,579 |
| **Total tokens on this bench, August 2026** | **21,408,017,663** |
| Open source, absolute-path rule (**the published floor**) | **3,812,106,080 — 17.807%** |
| Open source, tilde-path sensitivity arm (upper end of the range) | 4,987,293,753 — 23.296% |
| Everything else (private work), one aggregate line per day | 17,595,911,583 — 82.193% |

So the honest statement for this bench is: **between 17.8% and 23.3% of August's tokens went
into open source repositories, and 17.8% is the number that survives the strictest reading.**

Two independent passes over the same window agree to the token: the attribution pass and the
shared accounting scanner both report 21,408,017,663 tokens and 79,579 turns, delta 0. The
run refuses to print a table if they ever disagree.

### Open source repositories, August 2026, this bench

| Repository | Tokens | Share of month | Days | Turns |
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

`schema` is not a mistake and not double counting: it is the largest single body of work in
the month by a wide margin, and it is open source. Repos printed with an owner prefix
(`microsoft/vcpkg`, `jedisct1/libhydrogen`, `gafferongames/vcpkg`) are upstreams and forks
worked on directly, not repos under this org.

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
  the file gives **21,408,017,663** and summing `turns` gives **79,579** — the bench's whole
  month, both matching section 1. The file is complete, not a sample.
* `attribution` records how the directory-to-repo mapping was obtained: `remote` (read off
  that checkout's git origin), `heuristic` (checkout gone from disk, folded onto a verified
  sibling or a known repo name), `unresolved` (none of the above; classified private
  fail-closed), `aggregate` (a collapsed private line).
* `gcs_mirrored_total` / `gcs_coverage` say how much of that row was also present in the
  object-store mirror of this bench's sessions (`full` / `partial` / `none`, `mixed` on an
  aggregate line). See section 3 — the mirror is a strict subset and changes no number here.

167 open source rows + 28 aggregate lines = 195 rows.

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
not line up.

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
orchestration checkout in 79,345 of 79,579 turns and handed it 97.6% of every token it was
seen near. The reason: *every* transcript record carries a `cwd` field, so a raw scan
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
`repo-visibility.json` beside the output. It is a **cache of real `gh` answers, not a guess**.
A repo name missing from the cache is listed in the run output and bucketed **private,
fail-closed**. That is why every figure here is a floor: an unclassified name can only ever
push the open source share down.

### Conservation

Attribution moves tokens between buckets. It must not create or destroy any. So the run
re-derives the same window through the shared accounting scanner and **refuses to print a
table** unless the two grand totals *and* the two turn counts agree exactly. They do:
`21,408,017,663` and `79,579` both ways.

There is a second, deliberately independent check (`tools/august-verify-attrib.go`) that
recomputes one day from scratch — its own walk, its own record struct, its own UTC-to-local
conversion, its own message-id fold, its own path scan, its own git remote read, sharing no
code with the main pass. If two programs built from the same parts agree, the agreement means
nothing; these were not.

---

## 3. How the Studio bench completes this

### Step 1 — run the pipeline on the Studio's own store

The two tools in `tools/` are the ones that produced this data. They import
`github.com/mas-bandwidth/rowan-tools/modules/{spend,usage}`; those modules live in the
private workshop repo, which the Studio has. Drop both files into `rowan-tools/scratch/` and
run them from the repo root.

```
# 1. build the visibility cache from real gh answers, into the output directory
#    (one "<name>": true|false entry per repo name; true == public)
gh repo view <name> --json isPrivate

# 2. the attribution pass over this bench's own ~/.claude/projects
go run scratch/august-repo-spend-paths.go \
    -bench studio \
    -lo 2026-08-01 -hi 2026-08-31 \
    -zone Australia/Sydney \
    -out /path/to/out \
    -csv august-2026-studio-bench-day-repo-bypath.csv

# 3. the sensitivity arm (also counts ~/rowan-working references, which appear
#    overwhelmingly inside shell commands). Reported as a range, never merged into the floor.
go run scratch/august-repo-spend-paths.go -bench studio -tilde -out /path/to/out \
    -csv august-2026-studio-bench-day-repo-bypath-TILDE.csv

# 4. the independent recomputation, pointed at the busiest day in the CSV
go run scratch/august-verify-attrib.go -day <a busy day> -zone Australia/Sydney
```

`-bench studio` (or whatever that bench is called — just not `macbook`) is what keeps the
tuples disjoint. `-working` overrides the checkout directory if it is not `$HOME/rowan-working`.

Read the run's own output, not just the CSV. It prints the reconciliation, the list of
unclassified repo names, the ambiguity census, a touched-vs-won table (a repo referenced in
turns the argmax gave to someone else), the directory-to-repo mapping actually used, and what
fraction of the answer rests on a mapping that was not remote-verified. On this bench that
last figure is 131,601,040 tokens, 0.61%.

### Step 2 — the object store adds nothing, do not spend time on it

`gs://rowan-sydney/sessions` holds **only this bench's own mirrored sessions**. A full pass
over the mirror measured 10,315,509,634 tokens across 43,197 turns — a strict subset of the
21,408,017,663 already published here, not additional data. It contains nothing from any
other bench. See `rowan-tools#46`. Measure the Studio's own `~/.claude/projects`; that is the
only place the missing half exists.

The `gcs_coverage` column in `data.csv` records this: rows from the second half of the month
read `none` because mirroring had not caught up, and that is a statement about the mirror, not
about the measurement.

### Step 3 — merge

The merge is a **union of rows on `(bench, day, repo)`**.

* Rows whose tuples do not collide are simply concatenated.
* **An overlapping tuple must agree byte-exact.** Two benches producing different numbers for
  the same `(bench, day, repo)` is not something to average or reconcile — it means one of the
  runs is wrong, and it is a finding to report, not to smooth over.
* Sum `total` across the union for the month's grand total, sum the `OSS` rows for the open
  source total, and the ratio is the month's open source share. Same arithmetic as section 1,
  more rows.
* Before publishing the merged file, apply the same privacy filter used here: open source rows
  verbatim, every private row collapsed into one `PRIVATE` aggregate line per `(bench, day)`.
  Then check that the published file still sums to the unfiltered grand total.

---

## 4. Caveats

**The open source figures are floors, and here is the size of the edge.** In the published
(absolute-path) arm, 34,331 of 79,579 turns — 43.14% of turns, 9,028,682,297 tokens, 42.17%
of the month — had no absolute path reference in their tool calls or results and were decided
by the session's `cwd` fallback. That fallback is the biased signal this method exists to
avoid, and it points at private orchestration directories far more often than at anything
else. In the tilde arm, where `~/rowan-working/...` references are promoted to real ones, the
fallback share drops to 26,344 turns (33.10%) and 6,705,453,475 tokens (31.32%) — and the open
source share rises from 17.8% to 23.3%. **The true number is somewhere in that band and the
low end is what gets published.**

Other things a reader should know:

* **Winner-take-all.** A turn is attributed whole to one repo. 6,180 turns (7.77%,
  2,329,986,005 tokens) referenced two or more repos and only one of them was credited. 1,615
  turns tied and fell to the fallback rather than letting map iteration order pick a winner.
* **A sensitivity that did not matter.** Counting distinct files instead of every occurrence
  would have changed the winner in 32 turns (0.04%, 6,970,222 tokens). Measured and reported,
  not applied.
* **These are list-price tokens, not dollars.** No conversion to money is made here.
  Converting requires per-model rates, and one model id in this window (`<synthetic>`) is not
  in the rate table at all — anything priced for it would be a guess, not a measurement.
* **Cache reads dominate the raw counts.** `cache_read` is the bulk of every row. That is
  normal for long agentic sessions and it is left visible rather than netted out, because
  netting it out is an editorial choice and the columns let anyone make it themselves.
* **Half the month is missing.** Every number here is one bench. Nothing about the month as a
  whole should be stated until the Studio's half is merged.
