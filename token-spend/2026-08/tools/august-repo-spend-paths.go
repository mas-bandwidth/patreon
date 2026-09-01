//go:build ignore

// august-repo-spend-paths — THROWAWAY. August 2026 token spend per GITHUB REPO for THIS
// bench, attributed by the paths each TURN actually touched.
//
// WHY THIS EXISTS. The first run attributed a turn to its SESSION's working directory. The
// estate runs almost everything from orchestrator sessions whose cwd is ONE private working
// directory, so that run booked a month of schema and serialize work as 99.996% private --
// an artifact of where the shell was standing, not of what was worked on.
//
// THE RULE, deterministic so it can be reproduced:
//
//  1. A TURN is one folded assistant message (message id, largest usage seen), exactly as
//     usage.FoldTurns defines it. That definition is not re-litigated here.
//  2. A turn's SEGMENT is every transcript record since the previous usage-carrying record
//     (exclusive) through the usage-carrying record itself (inclusive). That is the tool
//     RESULTS the turn read plus the tool CALLS the turn made, and every record in the file
//     belongs to exactly one segment, so nothing is counted twice or dropped.
//  3. Inside a segment, every occurrence of the literal <working-root>/<dir>
//     is one REFERENCE to <dir>.
//  4. <dir> is folded to its GITHUB REPO first, by reading that checkout's origin remote --
//     the owner's ruling is that the github repo is the key and the directory is only how it
//     is recovered. So apt-schema is apt (its remote says so) and NOT schema, which the old
//     name heuristic would have got wrong.
//  5. The turn attributes WHOLE to the repo with the most references. NO references, or a
//     TIE for the most, falls back to the session's own cwd repo.
//
// Folding to repos BEFORE the argmax is deliberate: three references to schema-bench and two
// to schema-gowrite are five references to schema, not a two-way loss to something else.
//
// Attribution moves tokens between buckets. It must not create or destroy any, so the run
// re-derives the same window through modules/spend and REFUSES to report unless the two
// grand totals agree exactly.
//
// Usage:
//
//	go run scratch/august-repo-spend-paths.go -bench macbook -out /path/to/dir
package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mas-bandwidth/rowan-tools/modules/spend"
	"github.com/mas-bandwidth/rowan-tools/modules/usage"
)

// workingPrefix is the literal a path reference must start with. Absolute only, per the
// ruling: a tilde form is measured separately below as a sensitivity check, never mixed into
// the attribution, so the stated rule stays the rule that ran.
//
// It is derived from -working (default $HOME/rowan-working) rather than hard coded, so the
// same rule runs unchanged on any bench.
var workingPrefix string

// workingRoot is workingPrefix without its trailing separator: the directory the checkouts
// live in, which is also where the remotes are read from.
var workingRoot string

// tildePrefix is the same location written the other way -- overwhelmingly inside Bash
// commands, which is real work on a real checkout. The stated rule counts ABSOLUTE paths
// only, so by default these are tallied and reported but never attributed; -tilde runs the
// same measurement with them promoted to references, as a sensitivity arm.
var tildePrefix string

// countTilde promotes tilde references to real ones. Set once from the flag before the pass.
var countTilde bool

// masOwner is the github owner whose repos print bare. Anything under a different owner
// keeps its owner in the name, because a fork or a vendored upstream is not the same repo.
const masOwner = "mas-bandwidth"

func main() {
	bench := flag.String("bench", "macbook", "bench name for the uniqueness tuple")
	zone := flag.String("zone", "Australia/Sydney", "zone the local days are counted in")
	lo := flag.String("lo", "2026-08-01", "inclusive first local day")
	hi := flag.String("hi", "2026-08-31", "inclusive last local day")
	out := flag.String("out", ".", "directory for the CSV and the visibility cache")
	expect := flag.Int64("expect", 0, "grand total the cwd run measured; 0 to skip that check")
	tilde := flag.Bool("tilde", false, "SENSITIVITY ARM: also count ~/rowan-working references")
	name := flag.String("csv", "august-2026-macbook-bench-day-repo-bypath.csv", "CSV file name")
	work := flag.String("working", "", "directory holding the checkouts (default $HOME/rowan-working)")
	flag.Parse()
	countTilde = *tilde

	home, err := os.UserHomeDir()
	must(err)
	if *work == "" {
		*work = filepath.Join(home, "rowan-working")
	}
	workingRoot = strings.TrimSuffix(*work, "/")
	workingPrefix = workingRoot + "/"
	tildePrefix = "~/" + filepath.Base(workingRoot) + "/"
	loc, err := time.LoadLocation(*zone)
	must(err)
	root := spend.ProjectsPath(home)

	// ---- BASELINE. The same window through the shared scanner, whose fail-closed
	// accounting is the thing being reused. Its grand total is what the new attribution has
	// to preserve.
	base, err := spend.Scan(os.DirFS(root), spend.ScanOptions{Lo: *lo, Hi: *hi, Loc: loc, Root: root})
	must(err)
	var baseTotal int64
	baseTurns := 0
	for _, r := range base.Rows {
		baseTotal += r.Tokens.Total()
		baseTurns += r.Turns
	}
	fmt.Printf("BASELINE (modules/spend) root=%s zone=%s window=%s..%s\n", root, *zone, *lo, *hi)
	fmt.Printf("BASELINE files_read=%d skipped=%d walked=%v complete=%v rows=%d turns=%d total=%d\n",
		base.Files, base.Skipped, base.Walked, base.Complete(), len(base.Rows), baseTurns, baseTotal)
	for _, n := range base.Notes {
		fmt.Printf("BASELINE NOTE [%v] %s\n", n.Level, n.Msg)
	}
	if !base.Complete() {
		fmt.Println("BASELINE INCOMPLETE -- every total below is a FLOOR, not a total")
	}

	// ---- RESOLVER. Directory -> github repo, remote-verified where the checkout survives.
	res := newResolver(workingRoot)

	// ---- THE PATH PASS.
	pass, err := scanPaths(root, *lo, *hi, loc)
	must(err)

	fmt.Printf("\nPATH PASS files_read=%d skipped=%d malformed_usage_records=%d turns=%d total=%d\n",
		pass.files, pass.skipped, pass.bad, pass.turns, pass.total)

	// ---- RECONCILIATION. Refuse to report a table that invented or lost tokens.
	fmt.Printf("\nRECONCILE path-pass total=%d vs modules/spend total=%d delta=%d\n",
		pass.total, baseTotal, pass.total-baseTotal)
	fmt.Printf("RECONCILE path-pass turns=%d vs modules/spend turns=%d delta=%d\n",
		pass.turns, baseTurns, pass.turns-baseTurns)
	if pass.total != baseTotal || pass.turns != baseTurns {
		fmt.Fprintln(os.Stderr, "FATAL: the two passes disagree -- that is a parser bug, "+
			"not an attribution difference. Refusing to print a table.")
		os.Exit(1)
	}
	if *expect != 0 {
		fmt.Printf("RECONCILE against the cwd run's published total %d: delta=%d\n",
			*expect, pass.total-*expect)
		if pass.total != *expect {
			fmt.Fprintln(os.Stderr, "FATAL: this window no longer sums to the cwd run's "+
				"grand total. Attribution must move tokens, never change how many there are.")
			os.Exit(1)
		}
	}

	// ---- ATTRIBUTION.
	type key struct{ day, repo string }
	agg := map[key]*usage.Tokens{}
	turnsBy := map[key]int{}
	sessBy := map[key]map[string]bool{}
	repoDirs := map[string]map[string]bool{}

	var multiTok, fallbackTok, tieTok, tildeOnlyTok int64
	var multiTurns, fallbackTurns, tieTurns, tildeOnlyTurns int
	var distinctDisagreeTok int64
	var distinctDisagreeTurns int

	// TOUCH DIAGNOSTIC. The winner-take-all rule hides how often a repo was worked on but
	// out-referenced by whatever else the same turn mentioned. These three say so: how many
	// references a repo drew in total, how many turns mentioned it at all, and how many
	// tokens rode those turns. A repo with a large touched figure and a small won figure is
	// losing every argmax to a louder neighbour, which is the exact failure the cwd run had.
	refTotal := map[string]int{}
	touchTurns := map[string]int{}
	touchTok := map[string]int64{}

	for _, t := range pass.rows {
		// Fold directory references onto github repos BEFORE choosing a winner.
		byRepo := map[string]int{}
		byRepoDistinct := map[string]int{}
		for dir, n := range t.refs {
			r := res.repoFor(dir)
			byRepo[r] += n
			byRepoDistinct[r] += t.distinct[dir]
			if repoDirs[r] == nil {
				repoDirs[r] = map[string]bool{}
			}
			repoDirs[r][dir] = true
		}
		for r, n := range byRepo {
			refTotal[r] += n
			touchTurns[r]++
			touchTok[r] += t.tok.Total()
		}
		cwdRepo := res.repoFor(firstSeg(strings.TrimPrefix(t.cwd, workingPrefix)))
		if !strings.HasPrefix(t.cwd, workingPrefix) {
			// A session that ran outside the working root has no recoverable repo, so its
			// fallback is the directory itself, which classifies PRIVATE.
			cwdRepo = outsideRepo(t.cwd, home)
		}

		repo, decided := argmax(byRepo)
		if !decided {
			repo = cwdRepo
			fallbackTurns++
			fallbackTok += t.tok.Total()
			if len(byRepo) > 0 {
				tieTurns++
				tieTok += t.tok.Total()
			}
			if t.tilde > 0 && len(t.refs) == 0 {
				tildeOnlyTurns++
				tildeOnlyTok += t.tok.Total()
			}
		}
		if len(byRepo) > 1 {
			multiTurns++
			multiTok += t.tok.Total()
		}
		// SENSITIVITY: would counting DISTINCT files instead of every occurrence have
		// changed this turn's winner? Reported, never applied.
		if dr, ok := argmax(byRepoDistinct); ok && decided && dr != repo {
			distinctDisagreeTurns++
			distinctDisagreeTok += t.tok.Total()
		}

		k := key{t.day, repo}
		if agg[k] == nil {
			agg[k] = &usage.Tokens{}
			sessBy[k] = map[string]bool{}
		}
		*agg[k] = agg[k].Add(t.tok)
		turnsBy[k]++
		sessBy[k][sessionOf(t.path)] = true
	}

	// ---- VISIBILITY. Cached answers only; an unknown name is PRIVATE, never assumed open.
	vis := loadVisibility(*out)
	unknown := map[string]bool{}
	for k := range agg {
		if _, ok := vis[k.repo]; !ok {
			unknown[k.repo] = true
		}
	}
	if len(unknown) > 0 {
		names := sortedKeys(unknown)
		fmt.Printf("\nUNCLASSIFIED REPOS (not in %s) -- bucketed PRIVATE fail-closed:\n  %s\n",
			filepath.Join(*out, visFile), strings.Join(names, ", "))
	}

	// ---- CSV on the (bench, day, repo) tuple.
	must(os.MkdirAll(*out, 0o755))
	csvPath := filepath.Join(*out, *name)
	f, err := os.Create(csvPath)
	must(err)
	w := csv.NewWriter(f)
	must(w.Write([]string{"bench", "day", "repo", "visibility", "attribution",
		"input", "output", "cache_creation", "cache_read", "total", "turns", "sessions_touching"}))

	keys := make([]key, 0, len(agg))
	for k := range agg {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].day != keys[j].day {
			return keys[i].day < keys[j].day
		}
		return keys[i].repo < keys[j].repo
	})
	for _, k := range keys {
		t := *agg[k]
		must(w.Write([]string{*bench, k.day, k.repo, visibility(vis, k.repo), res.provenance(k.repo),
			i64(t.Input), i64(t.Output), i64(t.CacheCreation), i64(t.CacheRead), i64(t.Total()),
			strconv.Itoa(turnsBy[k]), strconv.Itoa(len(sessBy[k]))}))
	}
	w.Flush()
	must(w.Error())
	must(f.Close())
	fmt.Printf("\nCSV %s (%d rows)\n", csvPath, len(keys))

	// ---- MONTH ROLL-UP.
	type roll struct {
		tok   usage.Tokens
		turns int
		days  map[string]bool
		sess  map[string]bool
	}
	byRepo := map[string]*roll{}
	for _, k := range keys {
		r := byRepo[k.repo]
		if r == nil {
			r = &roll{days: map[string]bool{}, sess: map[string]bool{}}
			byRepo[k.repo] = r
		}
		r.tok = r.tok.Add(*agg[k])
		r.turns += turnsBy[k]
		r.days[k.day] = true
		for s := range sessBy[k] {
			r.sess[s] = true
		}
	}
	var grand int64
	for _, r := range byRepo {
		grand += r.tok.Total()
	}
	names := make([]string, 0, len(byRepo))
	for n := range byRepo {
		names = append(names, n)
	}
	sort.Slice(names, func(i, j int) bool {
		if byRepo[names[i]].tok.Total() != byRepo[names[j]].tok.Total() {
			return byRepo[names[i]].tok.Total() > byRepo[names[j]].tok.Total()
		}
		return names[i] < names[j]
	})

	fmt.Printf("\nREPO ROLLUP (grand total %d tokens)\n", grand)
	fmt.Printf("%-26s %-8s %-10s %14s %12s %12s %14s %14s %5s %7s %8s\n",
		"repo", "vis", "mapping", "total", "input", "output", "cache_create", "cache_read",
		"days", "turns", "share")
	var privTok usage.Tokens
	privTurns, privDays, privSess := 0, map[string]bool{}, map[string]bool{}
	privNames := []string{}
	for _, n := range names {
		r := byRepo[n]
		if visibility(vis, n) != "OSS" {
			privTok = privTok.Add(r.tok)
			privTurns += r.turns
			for d := range r.days {
				privDays[d] = true
			}
			for s := range r.sess {
				privSess[s] = true
			}
			privNames = append(privNames, n)
			continue
		}
		fmt.Printf("%-26s %-8s %-10s %14d %12d %12d %14d %14d %5d %7d %7.3f%%\n",
			n, "OSS", res.provenance(n), r.tok.Total(), r.tok.Input, r.tok.Output,
			r.tok.CacheCreation, r.tok.CacheRead, len(r.days), r.turns,
			100*float64(r.tok.Total())/float64(grand))
	}
	fmt.Printf("%-26s %-8s %-10s %14d %12d %12d %14d %14d %5d %7d %7.3f%%\n",
		"(all private, one line)", "PRIVATE", "-", privTok.Total(), privTok.Input,
		privTok.Output, privTok.CacheCreation, privTok.CacheRead, len(privDays), privTurns,
		100*float64(privTok.Total())/float64(grand))
	fmt.Printf("  private constituents: %s\n", strings.Join(privNames, ", "))

	// ---- AMBIGUITY.
	fmt.Printf("\nAMBIGUITY (denominator = %d turns, %d tokens)\n", pass.turns, grand)
	amb := func(label string, turns int, tok int64) {
		fmt.Printf("  %-42s %8d turns (%6.2f%%)  %14d tokens (%6.2f%%)\n", label, turns,
			100*float64(turns)/float64(pass.turns), tok, 100*float64(tok)/float64(grand))
	}
	amb("turns referencing 2+ repos", multiTurns, multiTok)
	amb("turns decided by cwd FALLBACK (all)", fallbackTurns, fallbackTok)
	amb("  of those, ties between repos", tieTurns, tieTok)
	amb("  of those, tilde-path-only turns", tildeOnlyTurns, tildeOnlyTok)
	amb("SENSITIVITY: distinct-file rule would differ", distinctDisagreeTurns, distinctDisagreeTok)

	// ---- TOUCHED vs WON.
	fmt.Printf("\nTOUCHED vs WON (a repo referenced in a turn the argmax gave to someone else)\n")
	fmt.Printf("%-26s %10s %10s %16s %16s %7s\n",
		"repo", "refs", "turns_ref", "tokens_in_those", "tokens_won", "won/ref")
	touched := make([]string, 0, len(touchTok))
	for n := range touchTok {
		touched = append(touched, n)
	}
	sort.Slice(touched, func(i, j int) bool { return touchTok[touched[i]] > touchTok[touched[j]] })
	for i, n := range touched {
		if i >= 25 {
			break
		}
		var won int64
		if r := byRepo[n]; r != nil {
			won = r.tok.Total()
		}
		fmt.Printf("%-26s %10d %10d %16d %16d %6.2f%%\n", n, refTotal[n], touchTurns[n],
			touchTok[n], won, 100*float64(won)/float64(touchTok[n]))
	}

	// ---- MAPPING PROVENANCE.
	fmt.Printf("\nDIRECTORY -> REPO MAPPING ACTUALLY USED\n")
	for _, n := range names {
		dirs := sortedKeys(repoDirs[n])
		fmt.Printf("  %-24s [%s] <- %s\n", n, res.provenance(n), strings.Join(dirs, ", "))
	}
	fmt.Printf("\nHEURISTIC (checkout gone from disk, folded onto a verified sibling or the "+
		"name table):\n  %s\n", strings.Join(res.heuristicDirs(), ", "))
	fmt.Printf("\nUNRESOLVED directory names (no remote, no verified sibling, no known "+
		"repo -- almost all truncated path fragments):\n  %s\n",
		strings.Join(orNone(res.unresolvedDirs()), ", "))

	// How much of the answer rests on a mapping that was not read off a remote?
	var heurTok int64
	heurRepos := []string{}
	for _, n := range names {
		if res.provenance(n) != "remote" {
			heurTok += byRepo[n].tok.Total()
			heurRepos = append(heurRepos, fmt.Sprintf("%s(%s)", n, res.provenance(n)))
		}
	}
	fmt.Printf("\nMAPPING CONFIDENCE: %d tokens (%.4f%%) sit in repos whose mapping was NOT "+
		"remote-verified\n  %s\n", heurTok, 100*float64(heurTok)/float64(grand),
		strings.Join(orNone(heurRepos), ", "))
}

// ---------------------------------------------------------------- the path pass

// turnRow is one folded turn with the places it touched.
type turnRow struct {
	path     string
	day      string
	cwd      string
	tok      usage.Tokens
	refs     map[string]int // dir -> occurrences
	distinct map[string]int // dir -> distinct full paths
	tilde    int            // ~/rowan-working references, counted never attributed
}

type passResult struct {
	rows                       []turnRow
	files, skipped, bad, turns int
	total                      int64
}

// scanPaths walks the transcript tree and folds each file into turns carrying their path
// references. It mirrors modules/spend's prefilter, window and folding EXACTLY -- that is
// what makes the reconciliation above a real check rather than a tautology.
func scanPaths(root, lo, hi string, loc *time.Location) (passResult, error) {
	var pr passResult
	loMidnight, err := spend.ParseDay(lo)
	if err != nil {
		return pr, err
	}
	cutoff := time.Date(loMidnight.Year(), loMidnight.Month(), loMidnight.Day(), 0, 0, 0, 0,
		loc).AddDate(0, 0, -1)

	var files []string
	err = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if p == root {
				return err
			}
			pr.skipped++
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(p, spend.TranscriptExt) {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return pr, err
	}
	sort.Strings(files)

	for _, p := range files {
		info, err := os.Stat(p)
		if err != nil {
			pr.skipped++
			continue
		}
		if info.ModTime().Before(cutoff) {
			continue
		}
		rel, _ := filepath.Rel(root, p)
		rel = filepath.ToSlash(rel)
		rows, bad, err := scanFile(p, rel, lo, hi, loc)
		if err != nil {
			pr.skipped++
			continue
		}
		pr.files++
		pr.bad += bad
		for _, r := range rows {
			pr.rows = append(pr.rows, r)
			pr.turns++
			pr.total += r.tok.Total()
		}
	}
	return pr, nil
}

// scanFile reads ONE transcript, segmenting path references onto the turn that closes each
// segment, folding by message id, and returning only the turns inside the window.
func scanFile(abs, rel, lo, hi string, loc *time.Location) ([]turnRow, int, error) {
	f, err := os.Open(abs)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	type acc struct {
		row   turnRow
		total int64
	}
	idx := map[string]int{}
	accs := []acc{}

	segRefs := map[string]int{}
	segPaths := map[string]map[string]bool{}
	segTilde := 0
	cwd := ""
	bad := 0

	br := bufio.NewReader(f)
	for lineno := 0; ; lineno++ {
		line, rerr := br.ReadString('\n')
		if line != "" {
			// Only TOOL CALLS AND TOOL RESULTS contribute references, and that restriction is
			// the whole measurement.
			//
			// FOUND BY THE DIAGNOSTIC BELOW, not predicted. Scanning the whole record instead
			// put ONE private orchestration checkout in 79,345 of 79,579 turns and handed it
				// 97.6% of every token it
			// was seen near. The reason is that EVERY record carries a "cwd" field, so a raw
			// scan re-imports the session's working directory as if it were a thing the turn
			// touched -- which is the exact bias the owner's ruling exists to remove, smuggled
			// back in through the metadata. Injected attachments (CLAUDE.md, the memory
			// index) did the same thing more quietly.
			//
			// A turn's own prose is excluded too: "the turn's tool calls/results" is the rule,
			// and reasoning that merely NAMES a path is not evidence the path was worked on.
			countToolRefs(line, segRefs, segPaths)
			segTilde += strings.Count(line, tildePrefix)

			if cwd == "" {
				if c := fieldString(line, `"cwd"`); c != "" {
					cwd = c
				}
			}
			if strings.Contains(line, `"usage"`) {
				var d rawRecord
				if json.Unmarshal([]byte(line), &d) != nil {
					bad++
				} else if d.Message != nil && d.Message.Usage != nil {
					ts := d.Timestamp
					if ts == "" {
						ts = d.Message.Timestamp
					}
					if ts == "" {
						bad++
					} else {
						u := d.Message.Usage
						tok := usage.Tokens{Input: u.Input, CacheCreation: u.CacheCreation,
							CacheRead: u.CacheRead, Output: u.Output}
						key := d.Message.ID
						if key == "" {
							key = fmt.Sprintf("no-id:line%d", lineno)
						}
						tot := tok.Total()
						if i, seen := idx[key]; seen {
							// Same fold as usage.FoldTurns for tokens: largest usage wins,
							// first timestamp kept. References UNION across the duplicate
							// records, because each streamed record's segment holds real
							// evidence and a later record must not erase an earlier one.
							if tot > accs[i].total {
								accs[i].total = tot
								accs[i].row.tok = tok
							}
							mergeRefs(accs[i].row.refs, segRefs)
							mergeDistinct(accs[i].row.distinct, segPaths)
							accs[i].row.tilde += segTilde
						} else {
							idx[key] = len(accs)
							r := turnRow{path: rel, cwd: cwd, tok: tok,
								refs: map[string]int{}, distinct: map[string]int{}, tilde: segTilde}
							r.day, _ = spend.LocalDay(ts, loc)
							mergeRefs(r.refs, segRefs)
							mergeDistinct(r.distinct, segPaths)
							accs = append(accs, acc{row: r, total: tot})
						}
						// The segment closes on a usage record, so nothing is double counted.
						segRefs = map[string]int{}
						segPaths = map[string]map[string]bool{}
						segTilde = 0
					}
				}
			}
		}
		if rerr != nil {
			if rerr == io.EOF {
				break
			}
			return nil, bad, rerr
		}
	}

	out := make([]turnRow, 0, len(accs))
	for _, a := range accs {
		if !spend.InWindow(a.row.day, lo, hi) {
			continue
		}
		if a.row.cwd == "" {
			a.row.cwd = cwd
		}
		out = append(out, a.row)
	}
	return out, bad, nil
}

// toolRecord is the part of a record that says what was TOUCHED, as opposed to where the
// shell happened to be standing.
type toolRecord struct {
	Message *struct {
		Content json.RawMessage `json:"content"`
	} `json:"message"`
	ToolUseResult json.RawMessage `json:"toolUseResult"`
}

type contentBlock struct {
	Type    string          `json:"type"`
	Input   json.RawMessage `json:"input"`   // tool_use: the call's arguments
	Content json.RawMessage `json:"content"` // tool_result: what came back
}

// countToolRefs decodes one record far enough to reach its tool calls and tool results, and
// counts path references in those alone.
//
// The two cheap gates come first because the decode is the expensive part of the whole run
// and most records are neither a call nor a result.
func countToolRefs(line string, refs map[string]int, paths map[string]map[string]bool) {
	if !strings.Contains(line, workingPrefix) && !(countTilde && strings.Contains(line, tildePrefix)) {
		return
	}
	if !strings.Contains(line, `"tool_use"`) && !strings.Contains(line, `"tool_result"`) &&
		!strings.Contains(line, `"toolUseResult"`) {
		return
	}
	var rec toolRecord
	if json.Unmarshal([]byte(line), &rec) != nil {
		return
	}
	if len(rec.ToolUseResult) > 0 {
		countRefs(string(rec.ToolUseResult), refs, paths)
	}
	if rec.Message == nil || len(rec.Message.Content) == 0 {
		return
	}
	var blocks []contentBlock
	if json.Unmarshal(rec.Message.Content, &blocks) != nil {
		return // a string-valued content has no tool blocks in it
	}
	for _, b := range blocks {
		switch b.Type {
		case "tool_use":
			countRefs(string(b.Input), refs, paths)
		case "tool_result":
			countRefs(string(b.Content), refs, paths)
		}
	}
}

// countRefs finds every <working-root>/<dir> in one chunk of JSON. A hand-rolled
// index scan rather than a regexp: this runs over every byte of a multi-gigabyte tree.
func countRefs(line string, refs map[string]int, paths map[string]map[string]bool) {
	countWith(line, workingPrefix, refs, paths)
	if countTilde {
		countWith(line, tildePrefix, refs, paths)
	}
}

func countWith(line, prefix string, refs map[string]int, paths map[string]map[string]bool) {
	i := 0
	for {
		j := strings.Index(line[i:], prefix)
		if j < 0 {
			return
		}
		start := i + j + len(prefix)
		dir := segmentAt(line, start)
		if dir != "" {
			refs[dir]++
			full := line[start : start+len(dir)]
			if k := pathEnd(line, start); k > start {
				full = line[start:k]
			}
			if paths[dir] == nil {
				paths[dir] = map[string]bool{}
			}
			paths[dir][full] = true
		}
		i = start
	}
}

// segmentAt reads the directory name beginning at start: the run of characters a checkout
// name may contain, stopped by a separator or by anything a path cannot hold.
func segmentAt(s string, start int) string {
	k := start
	for k < len(s) && pathChar(s[k]) {
		k++
	}
	return s[start:k]
}

// pathEnd extends past the directory to the end of the whole path reference, so the
// distinct-file sensitivity counts files rather than directories.
func pathEnd(s string, start int) int {
	k := start
	for k < len(s) && (pathChar(s[k]) || s[k] == '/') {
		k++
	}
	return k
}

// pathChar is what may appear in one path segment. Deliberately narrow: a JSON escape, a
// quote or a space ends the reference, which is what keeps a sentence quoting a path from
// swallowing the rest of the line.
func pathChar(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' ||
		c == '.' || c == '-' || c == '_' || c == '+'
}

func mergeRefs(dst, src map[string]int) {
	for k, v := range src {
		dst[k] += v
	}
}

func mergeDistinct(dst map[string]int, src map[string]map[string]bool) {
	for k, v := range src {
		dst[k] += len(v)
	}
}

// argmax returns the single highest-scoring key. A TIE IS NOT A DECISION: it returns false,
// and the caller falls back, rather than letting Go's map order pick a winner and make the
// whole report unreproducible.
func argmax(m map[string]int) (string, bool) {
	best, bestN, ties := "", 0, 0
	for k, n := range m {
		if n > bestN {
			best, bestN, ties = k, n, 1
			continue
		}
		if n == bestN {
			ties++
		}
	}
	if bestN == 0 || ties != 1 {
		return "", false
	}
	return best, true
}

// rawRecord is the part of a transcript record this tool decodes.
type rawRecord struct {
	Timestamp string `json:"timestamp"`
	Message   *struct {
		ID        string `json:"id"`
		Timestamp string `json:"timestamp"`
		Usage     *struct {
			Input         int64 `json:"input_tokens"`
			CacheCreation int64 `json:"cache_creation_input_tokens"`
			CacheRead     int64 `json:"cache_read_input_tokens"`
			Output        int64 `json:"output_tokens"`
		} `json:"usage"`
	} `json:"message"`
}

// fieldString pulls a top-level string field out of a record without decoding it, because
// decoding every line of a gigabyte of transcripts to find one field costs more than the
// whole rest of the pass.
func fieldString(line, quotedKey string) string {
	i := strings.Index(line, quotedKey+":")
	if i < 0 {
		return ""
	}
	j := strings.Index(line[i+len(quotedKey)+1:], `"`)
	if j < 0 {
		return ""
	}
	start := i + len(quotedKey) + 1 + j + 1
	k := strings.Index(line[start:], `"`)
	if k < 0 {
		return ""
	}
	return line[start : start+k]
}

// ---------------------------------------------------------------- directory -> repo

type resolver struct {
	root     string            // the directory the checkouts live in
	verified map[string]string // directory name -> repo, read from its own origin remote
	vdirs    []string          // verified directory names, longest first
	cache    map[string]string // dir -> repo
	how      map[string]string // dir -> "remote" | "heuristic" | "unresolved"
	prov     map[string]string // repo -> best provenance seen
}

// newResolver reads the origin remote of EVERY checkout still under the working root up
// front. That set is the ground truth the heuristic is then anchored to, so a directory that
// has since been deleted is recovered through a sibling that still exists rather than
// through a guess about what its name meant.
func newResolver(root string) *resolver {
	r := &resolver{root: root, verified: map[string]string{}, cache: map[string]string{},
		how: map[string]string{}, prov: map[string]string{}}
	ents, _ := os.ReadDir(root)
	for _, e := range ents {
		if repo, ok := r.fromRemote(e.Name()); ok {
			r.verified[e.Name()] = repo
			r.vdirs = append(r.vdirs, e.Name())
		}
	}
	// Longest first, so a longer checkout name beats the shorter verified name it extends:
	// schema-r198cold resolves through schema, not the other way round.
	sort.Slice(r.vdirs, func(i, j int) bool { return len(r.vdirs[i]) > len(r.vdirs[j]) })
	return r
}

// repoFor maps one working directory name to its github repo.
//
// The REMOTE is the authority. Where the checkout is gone, the name is folded onto a
// checkout that IS verified -- schema-mm-pr6 and scratch-schema-x both fold through
// schema -- and only if that fails does it fall back to the bare name table.
// A name that answers to none of the three is left as itself and REPORTED as unresolved,
// never quietly bundled into a repo it might not belong to.
func (r *resolver) repoFor(dir string) string {
	// A path written at the end of a sentence carries the full stop into the segment, and an
	// elided path leaves only dots. Neither is a directory.
	dir = strings.TrimRight(dir, ".")
	if dir == "" {
		return "(none)"
	}
	if v, ok := r.cache[dir]; ok {
		return v
	}
	repo, how := r.resolve(dir)
	r.cache[dir] = repo
	r.how[dir] = how
	// A repo counts as remote-verified when ANY directory folded onto it was verified.
	if r.prov[repo] != "remote" {
		r.prov[repo] = how
	}
	return repo
}

func (r *resolver) resolve(dir string) (string, string) {
	if repo, ok := r.verified[dir]; ok {
		return repo, "remote"
	}
	// Fold onto the longest verified checkout this name extends on a '-' boundary, with and
	// without the scratch- prefix the estate uses for throwaway clones.
	for _, cand := range []string{dir, strings.TrimPrefix(dir, "scratch-")} {
		for _, v := range r.vdirs {
			if strings.HasPrefix(cand, v+"-") {
				return r.verified[v], "heuristic"
			}
		}
	}
	if n := upstreamOf(dir); n != dir {
		return n, "heuristic"
	}
	if seedRepos[dir] {
		return dir, "heuristic"
	}
	return dir, "unresolved"
}

// unresolvedDirs names every referenced directory that answered to no remote, no verified
// sibling and no known repo name. These are overwhelmingly truncated path fragments, and
// their size is reported rather than assumed small.
func (r *resolver) unresolvedDirs() []string {
	out := []string{}
	for d, h := range r.how {
		if h == "unresolved" {
			out = append(out, d)
		}
	}
	sort.Strings(out)
	return out
}

func (r *resolver) provenance(repo string) string {
	if p, ok := r.prov[repo]; ok {
		return p
	}
	return "heuristic"
}

func (r *resolver) heuristicDirs() []string {
	out := []string{}
	for d, h := range r.how {
		if h == "heuristic" {
			out = append(out, d)
		}
	}
	sort.Strings(out)
	if len(out) == 0 {
		return []string{"(none -- every directory referenced still exists on disk)"}
	}
	return out
}

// fromRemote reads a checkout's origin URL. Handles a worktree, whose .git is a FILE
// pointing into the parent's .git/worktrees/<name>, by walking back to the parent config.
func (r *resolver) fromRemote(dir string) (string, bool) {
	base := filepath.Join(r.root, dir)
	gitPath := filepath.Join(base, ".git")
	info, err := os.Stat(gitPath)
	if err != nil {
		return "", false
	}
	cfg := filepath.Join(gitPath, "config")
	if !info.IsDir() {
		b, err := os.ReadFile(gitPath)
		if err != nil {
			return "", false
		}
		gd := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(string(b)), "gitdir:"))
		if i := strings.Index(gd, "/worktrees/"); i >= 0 {
			gd = gd[:i]
		}
		cfg = filepath.Join(gd, "config")
	}
	b, err := os.ReadFile(cfg)
	if err != nil {
		return "", false
	}
	url := originURL(string(b))
	if url == "" {
		return "", false
	}
	return repoFromURL(url), true
}

// originURL finds the url of the origin remote in a git config.
func originURL(cfg string) string {
	in := false
	for _, ln := range strings.Split(cfg, "\n") {
		t := strings.TrimSpace(ln)
		if strings.HasPrefix(t, "[") {
			in = t == `[remote "origin"]`
			continue
		}
		if in && strings.HasPrefix(t, "url") {
			if i := strings.Index(t, "="); i >= 0 {
				return strings.TrimSpace(t[i+1:])
			}
		}
	}
	return ""
}

// repoFromURL reduces every remote form the estate uses -- https, scp-style, an ssh alias, a
// 443 alias, with or without the .git suffix -- to owner/name, then to the bare name when
// the owner is the estate's own.
func repoFromURL(u string) string {
	u = strings.TrimSuffix(strings.TrimSpace(u), ".git")
	u = strings.ReplaceAll(u, ":", "/")
	parts := []string{}
	for _, p := range strings.Split(u, "/") {
		if p != "" {
			parts = append(parts, p)
		}
	}
	if len(parts) < 2 {
		return u
	}
	owner, name := parts[len(parts)-2], parts[len(parts)-1]
	if owner == masOwner {
		return name
	}
	return owner + "/" + name
}

// outsideRepo names the fallback for a session whose cwd is not under the working root. There
// is no github repo to recover, so the directory itself is the label and it classifies
// PRIVATE.
func outsideRepo(cwd, home string) string {
	rel := strings.TrimPrefix(cwd, home+"/")
	if rel == cwd {
		return "(outside home)"
	}
	return "(cwd:" + firstSeg(rel) + ")"
}

func firstSeg(p string) string {
	if i := strings.Index(p, "/"); i >= 0 {
		return p[:i]
	}
	return p
}

// upstreamOf is the NAME heuristic, used only where the checkout is gone. Longest matching
// known repo name on a '-' boundary.
func upstreamOf(base string) string {
	base = strings.TrimPrefix(base, "scratch-")
	if strings.HasPrefix(base, "serializego") {
		return "serialize.go"
	}
	best := ""
	for name := range seedRepos {
		if base == name {
			return name
		}
		if strings.HasPrefix(base, name+"-") && len(name) > len(best) {
			best = name
		}
	}
	if best != "" {
		return best
	}
	return base
}

// seedRepos is only the NAME table the heuristic matches against, used ONLY where the
// checkout is gone from disk. Visibility is NOT decided here -- it is decided by the cached
// gh answers below.
//
// THIS PUBLISHED COPY CARRIES THE PUBLIC REPO NAMES ONLY. The private names each bench also
// works on were removed before publication; add your own here. Removing a name cannot
// misclassify anything, because a name this table does not know is reported as unresolved
// and classified PRIVATE fail-closed.
var seedRepos = map[string]bool{
	"schema": true, "serialize": true, "serialize.go": true, "serialize.c": true,
	"serialize.modern": true, "serialize.rs": true, "serialize.elixir": true,
	"serialize.java": true, "serialize.dart": true, "serialize.js": true,
	"serialize.cs": true, "apt": true, "nova": true, "nova-tools": true,
	"reliable": true, "reliable.rs": true, "reliable.cs": true, "reliable.go": true,
	"netcode": true, "netcode.rs": true, "netcode.cs": true, "netcode.go": true,
	"fixed": true, "fixed3d": true, "fast3d": true, "table": true, "proton": true,
	"hydrogen": true, "yojimbo": true, "xdp": true, "udp": true, "matchmaker": true,
	"analytics": true, "cubes": true, "map": true, "next": true, "box3d": true,
	"patreon": true, "gafferongames": true, "gafferongames-new": true,
	"awesome-persistent-ai": true, "open-ledger": true, "homebrew-tap": true,
	"rowan-tools": true,
}

// ---------------------------------------------------------------- visibility

const visFile = "repo-visibility.json"

// loadVisibility reads the cached `gh repo view --json isPrivate` answers. A CACHE, not a
// guess: the file is written by the operator from real gh output, and a name missing from it
// is reported and bucketed PRIVATE rather than assumed open source.
func loadVisibility(dir string) map[string]bool {
	out := map[string]bool{}
	b, err := os.ReadFile(filepath.Join(dir, visFile))
	if err != nil {
		return out
	}
	_ = json.Unmarshal(b, &out)
	return out
}

func visibility(vis map[string]bool, repo string) string {
	if pub, ok := vis[repo]; ok && pub {
		return "OSS"
	}
	return "PRIVATE"
}

// ---------------------------------------------------------------- helpers

// sessionOf reduces a transcript path to its TOP-LEVEL session id, so a session's subagent
// transcripts do not each count as another session.
func sessionOf(p string) string {
	parts := strings.Split(p, "/")
	if len(parts) < 2 {
		return p
	}
	if len(parts) >= 3 {
		return parts[0] + "/" + parts[1]
	}
	return parts[0] + "/" + strings.TrimSuffix(parts[1], ".jsonl")
}

func orNone(s []string) []string {
	if len(s) == 0 {
		return []string{"(none)"}
	}
	return s
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func i64(n int64) string { return strconv.FormatInt(n, 10) }

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "FATAL:", err)
		os.Exit(1)
	}
}
