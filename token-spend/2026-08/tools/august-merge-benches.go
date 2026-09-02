//go:build ignore

// august-merge-benches — THROWAWAY. Merges one bench's RAW (bench, day, repo) CSV from
// august-repo-spend-paths.go with the already-published data.csv of another bench, applying
// the published privacy filter to the raw side first, and refusing to write anything unless
// the numbers conserve.
//
// THE RULES, verbatim from the published README (section 3, step 3):
//
//  1. The merge is a UNION of rows on (bench, day, repo).
//  2. An overlapping tuple must agree byte-exact. Disagreement is a finding, not something to
//     average; the tool exits non-zero and prints both rows.
//  3. Before publishing: open source rows verbatim; every private row collapsed into ONE
//     aggregate line per (bench, day) whose repo and visibility are literally PRIVATE, whose
//     attribution is `aggregate`, whose token columns and turns are the sums of the rows it
//     replaces, and whose sessions_touching is blank.
//  4. The published file must still sum to the unfiltered grand total.
//
// The raw CSV has 12 columns; the published one has 14 (two object-store mirror columns). The
// raw side gets `0,none` for those: the mirror holds only the other bench's sessions
// (README section 3, step 2), so on this bench nothing was mirrored, and that is a statement
// about the mirror rather than about the measurement.
//
// Usage:
//
//	go run scratch/august-merge-benches.go -published data.csv -raw studio.csv \
//	    -raw-tilde studio-TILDE.csv -published-tilde-oss 4987293753 -out merged.csv
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const header = "bench,day,repo,visibility,attribution,input,output,cache_creation,cache_read,total,turns,sessions_touching,gcs_mirrored_total,gcs_coverage"

type row struct {
	bench, day, repo, vis, attr string
	in, out, cc, cr, total      int64
	turns                       int
	sessions                    string // blank on an aggregate line
	gcsTotal, gcsCov            string
}

func (r row) tuple() string { return r.bench + "\x00" + r.day + "\x00" + r.repo }

func (r row) fields() []string {
	return []string{r.bench, r.day, r.repo, r.vis, r.attr, i64(r.in), i64(r.out), i64(r.cc),
		i64(r.cr), i64(r.total), strconv.Itoa(r.turns), r.sessions, r.gcsTotal, r.gcsCov}
}

func main() {
	published := flag.String("published", "", "the other bench's already-published data.csv (14 columns)")
	raw := flag.String("raw", "", "this bench's RAW csv from august-repo-spend-paths.go (12 columns)")
	rawTilde := flag.String("raw-tilde", "", "this bench's RAW tilde-arm csv (12 columns), for the range figure only")
	pubTildeOSS := flag.Int64("published-tilde-oss", 0, "the other bench's tilde-arm open source total, from its README")
	out := flag.String("out", "", "merged, filtered CSV to write")
	flag.Parse()
	if *published == "" || *raw == "" || *out == "" {
		fmt.Fprintln(os.Stderr, "need -published, -raw and -out")
		os.Exit(2)
	}

	pub := readCSV(*published, 14)
	rawRows := readCSV(*raw, 12)

	// ---- The published side is taken as-is, but its invariants are re-checked rather than
	// trusted: every PRIVATE line is an aggregate with a blank sessions column, and every
	// total is the sum of its parts.
	var pubTotal int64
	pubTurns := 0
	for _, r := range pub {
		if r.vis == "PRIVATE" && (r.repo != "PRIVATE" || r.attr != "aggregate" || r.sessions != "") {
			fatal("published file carries an un-collapsed private row: %v", r.fields())
		}
		if r.total != r.in+r.out+r.cc+r.cr {
			fatal("published row total does not equal its parts: %v", r.fields())
		}
		pubTotal += r.total
		pubTurns += r.turns
	}

	// ---- The raw side: open source rows verbatim, private rows collapsed per (bench, day).
	var rawTotal int64
	rawTurns := 0
	filtered := []row{}
	agg := map[string]*row{}
	aggOrder := []string{}
	rawSeen := map[string]bool{}
	for _, r := range rawRows {
		// The analysis pass emits one row per tuple. A repeated tuple here is an input fed
		// twice, and a repeated PRIVATE row would be absorbed into the aggregate without any
		// later check noticing, so it is refused at the door.
		if rawSeen[r.tuple()] {
			fatal("raw input carries a duplicate tuple: %v", r.fields())
		}
		rawSeen[r.tuple()] = true
		rawTotal += r.total
		rawTurns += r.turns
		if r.total != r.in+r.out+r.cc+r.cr {
			fatal("raw row total does not equal its parts: %v", r.fields())
		}
		if r.vis == "OSS" {
			r.gcsTotal, r.gcsCov = "0", "none"
			filtered = append(filtered, r)
			continue
		}
		k := r.bench + "\x00" + r.day
		a := agg[k]
		if a == nil {
			a = &row{bench: r.bench, day: r.day, repo: "PRIVATE", vis: "PRIVATE", attr: "aggregate",
				gcsTotal: "0", gcsCov: "none"}
			agg[k] = a
			aggOrder = append(aggOrder, k)
		}
		a.in += r.in
		a.out += r.out
		a.cc += r.cc
		a.cr += r.cr
		a.total += r.total
		a.turns += r.turns
	}
	for _, k := range aggOrder {
		filtered = append(filtered, *agg[k])
	}

	// ---- Conservation on the raw side: the filter moves rows, never tokens.
	var fTotal int64
	fTurns := 0
	for _, r := range filtered {
		fTotal += r.total
		fTurns += r.turns
	}
	fmt.Printf("RAW %s: rows=%d total=%d turns=%d\n", *raw, len(rawRows), rawTotal, rawTurns)
	fmt.Printf("FILTERED: rows=%d (%d open source + %d aggregate) total=%d turns=%d delta_total=%d delta_turns=%d\n",
		len(filtered), len(filtered)-len(aggOrder), len(aggOrder), fTotal, fTurns, fTotal-rawTotal, fTurns-rawTurns)
	if fTotal != rawTotal || fTurns != rawTurns {
		fatal("the privacy filter changed the grand total. Refusing to write.")
	}

	// ---- Union on the tuple. ANY collision is fatal: two rows that differ mean one run is
	// wrong; two rows that are identical mean an input was fed twice, and dropping one would
	// silently lose its tokens from the conservation the file has to keep.
	seen := map[string]row{}
	all := []row{}
	add := func(r row, src string) {
		if prev, ok := seen[r.tuple()]; ok {
			fatal("tuple collision (%s): %s/%s/%s\n  %v\n  %v", src, r.bench, r.day, r.repo, prev.fields(), r.fields())
		}
		seen[r.tuple()] = r
		all = append(all, r)
	}
	for _, r := range pub {
		add(r, "published")
	}
	for _, r := range filtered {
		add(r, "raw")
	}

	// bench, then day, then open source rows by repo, then the aggregate line last -- the
	// order the published file already uses.
	sort.SliceStable(all, func(i, j int) bool {
		a, b := all[i], all[j]
		if a.bench != b.bench {
			return a.bench < b.bench
		}
		if a.day != b.day {
			return a.day < b.day
		}
		ap, bp := a.repo == "PRIVATE", b.repo == "PRIVATE"
		if ap != bp {
			return !ap
		}
		return a.repo < b.repo
	})

	// Written to a temporary path and renamed into place only after the re-read below passes,
	// so a refusal never leaves a rejected file behind -- or clobbers -published when the
	// caller names the same path for both.
	tmp := *out + ".tmp"
	f, err := os.Create(tmp)
	must(err)
	w := csv.NewWriter(f)
	must(w.Write(strings.Split(header, ",")))
	for _, r := range all {
		must(w.Write(r.fields()))
	}
	w.Flush()
	must(w.Error())
	must(f.Close())

	// ---- Re-read what was written and sum it: the published file is checked from disk, not
	// from the slice that produced it.
	back := readCSV(tmp, 14)
	type tot struct {
		total, oss int64
		turns      int
		rows       int
	}
	byBench := map[string]*tot{}
	var grand tot
	repoTok := map[string]int64{}
	repoTurns := map[string]int{}
	repoDays := map[string]map[string]bool{}
	repoBench := map[string]map[string]bool{}
	for _, r := range back {
		b := byBench[r.bench]
		if b == nil {
			b = &tot{}
			byBench[r.bench] = b
		}
		b.total += r.total
		b.turns += r.turns
		b.rows++
		grand.total += r.total
		grand.turns += r.turns
		grand.rows++
		if r.vis == "OSS" {
			b.oss += r.oss()
			grand.oss += r.oss()
			repoTok[r.repo] += r.total
			repoTurns[r.repo] += r.turns
			if repoDays[r.repo] == nil {
				repoDays[r.repo] = map[string]bool{}
				repoBench[r.repo] = map[string]bool{}
			}
			repoDays[r.repo][r.day] = true
			repoBench[r.repo][r.bench] = true
		}
	}
	fmt.Printf("\nWRITTEN %s: %d rows\n", *out, len(back))
	// ---- The written file must sum to published + raw, exactly. This is checked on the
	// re-read from disk, after the union, so nothing between the inputs and the file can
	// lose a token without the run going red.
	if grand.total != pubTotal+rawTotal || grand.turns != pubTurns+rawTurns {
		_ = os.Remove(tmp)
		fatal("written file sums to %d tokens / %d turns; inputs sum to %d / %d. Tokens were lost or invented.",
			grand.total, grand.turns, pubTotal+rawTotal, pubTurns+rawTurns)
	}
	must(os.Rename(tmp, *out))
	benches := []string{}
	for b := range byBench {
		benches = append(benches, b)
	}
	sort.Strings(benches)
	for _, b := range benches {
		t := byBench[b]
		fmt.Printf("  %-8s rows=%4d total=%d turns=%d open_source=%d (%.3f%%) private=%d (%.3f%%)\n",
			b, t.rows, t.total, t.turns, t.oss, pct(t.oss, t.total), t.total-t.oss, pct(t.total-t.oss, t.total))
	}
	fmt.Printf("  %-8s rows=%4d total=%d turns=%d open_source=%d (%.3f%%) private=%d (%.3f%%)\n",
		"BOTH", grand.rows, grand.total, grand.turns, grand.oss, pct(grand.oss, grand.total),
		grand.total-grand.oss, pct(grand.total-grand.oss, grand.total))

	// ---- The tilde arm: a range figure, never merged into the file.
	if *rawTilde != "" {
		var tOSS, tTotal int64
		for _, r := range readCSV(*rawTilde, 12) {
			tTotal += r.total
			if r.vis == "OSS" {
				tOSS += r.total
			}
		}
		fmt.Printf("\nTILDE ARM (sensitivity, not in the file): raw bench open_source=%d of %d (%.3f%%)\n",
			tOSS, tTotal, pct(tOSS, tTotal))
		if tTotal != rawTotal {
			fatal("tilde arm grand total %d != absolute arm %d -- the arms must conserve", tTotal, rawTotal)
		}
		if *pubTildeOSS > 0 {
			both := tOSS + *pubTildeOSS
			fmt.Printf("TILDE ARM both benches: open_source=%d of %d (%.3f%%)\n", both, grand.total, pct(both, grand.total))
		}
	}

	// ---- Combined per-repo rollup, as a markdown table for the README.
	names := []string{}
	for n := range repoTok {
		names = append(names, n)
	}
	sort.Slice(names, func(i, j int) bool {
		if repoTok[names[i]] != repoTok[names[j]] {
			return repoTok[names[i]] > repoTok[names[j]]
		}
		return names[i] < names[j]
	})
	fmt.Printf("\n| Repository | Tokens | Share of month | Days | Turns | Benches |\n|---|---:|---:|---:|---:|---|\n")
	for _, n := range names {
		bs := []string{}
		for b := range repoBench[n] {
			bs = append(bs, b)
		}
		sort.Strings(bs)
		fmt.Printf("| %s | %s | %.3f%% | %d | %s | %s |\n", n, commas(repoTok[n]), pct(repoTok[n], grand.total),
			len(repoDays[n]), commas(int64(repoTurns[n])), strings.Join(bs, ", "))
	}
	fmt.Printf("| **all open source** | **%s** | **%.3f%%** | | | |\n", commas(grand.oss), pct(grand.oss, grand.total))
	fmt.Printf("| **everything else (aggregate)** | **%s** | **%.3f%%** | | | |\n",
		commas(grand.total-grand.oss), pct(grand.total-grand.oss, grand.total))
}

func (r row) oss() int64 {
	if r.vis == "OSS" {
		return r.total
	}
	return 0
}

func readCSV(path string, cols int) []row {
	f, err := os.Open(path)
	must(err)
	defer f.Close()
	rd := csv.NewReader(f)
	recs, err := rd.ReadAll()
	must(err)
	if len(recs) == 0 {
		fatal("%s is empty", path)
	}
	want := strings.Split(header, ",")[:cols]
	if strings.Join(recs[0], ",") != strings.Join(want, ",") {
		fatal("%s: header is not the expected %d columns:\n  got  %v\n  want %v", path, cols, recs[0], want)
	}
	out := []row{}
	for i, rec := range recs[1:] {
		if len(rec) != cols {
			fatal("%s line %d: %d fields, want %d", path, i+2, len(rec), cols)
		}
		r := row{bench: rec[0], day: rec[1], repo: rec[2], vis: rec[3], attr: rec[4],
			in: pi(rec[5]), out: pi(rec[6]), cc: pi(rec[7]), cr: pi(rec[8]), total: pi(rec[9]),
			turns: int(pi(rec[10])), sessions: rec[11]}
		if cols == 14 {
			r.gcsTotal, r.gcsCov = rec[12], rec[13]
		}
		if r.vis != "OSS" && r.vis != "PRIVATE" {
			fatal("%s line %d: visibility %q is neither OSS nor PRIVATE", path, i+2, r.vis)
		}
		out = append(out, r)
	}
	return out
}

func pct(a, b int64) float64 { return 100 * float64(a) / float64(b) }

func commas(n int64) string {
	s := strconv.FormatInt(n, 10)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

func pi(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	must(err)
	return n
}

func i64(n int64) string { return strconv.FormatInt(n, 10) }

func must(err error) {
	if err != nil {
		fatal("%v", err)
	}
}

func fatal(f string, a ...any) {
	fmt.Fprintf(os.Stderr, "FATAL: "+f+"\n", a...)
	os.Exit(1)
}
