//go:build ignore

// august-verify-attrib — THROWAWAY. Independent recomputation of ONE day's PATH-ATTRIBUTED
// spend, deliberately sharing no code with the main pass.
//
// A separate check on the PARSER (are the tokens right?) lives beside this one in the
// workshop repo. THIS file checks the ATTRIBUTION (did the tokens land on the right repo?),
// which is the half the re-run changed, so it is pointed at a heavy schema day rather than a
// quiet one.
//
// Everything is re-derived from scratch and on purpose: its own walk, its own record struct,
// its own UTC->local conversion, its own message-id fold, its own path scan, its own git
// remote read. It does NOT import modules/spend or modules/usage. If the two agree on a busy
// day, the agreement means something; if they were built from the same parts, it would not.
//
// It also prints the WHOLE day's turn-level decision census -- how many turns each repo won,
// by how large a margin -- so a reviewer can see the rule operating rather than trust it.
//
// Usage: go run scratch/august-verify-attrib.go -day 2026-08-31
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// workingPrefix is the literal a path reference must start with, derived from -working
// (default $HOME/rowan-working) so the same rule runs unchanged on any bench.
var workingPrefix string

func main() {
	day := flag.String("day", "2026-08-31", "local day to recompute")
	zone := flag.String("zone", "Australia/Sydney", "local zone")
	top := flag.Int("top", 15, "how many repos to print")
	// rootFlag lets this independent recomputation be aimed at a MIRRORED transcript tree
	// (an object-store pull) as well as this bench's own. The walk, the fold, the day
	// conversion, the path scan and the remote read are all this file's own, sharing nothing
	// with the main pass, which is the whole point of it.
	rootFlag := flag.String("root", "", "transcript tree to scan (default: this bench's ~/.claude/projects)")
	work := flag.String("working", "", "directory holding the checkouts (default $HOME/rowan-working)")
	flag.Parse()

	home, _ := os.UserHomeDir()
	root := filepath.Join(home, ".claude", "projects")
	if *rootFlag != "" {
		root = *rootFlag
	}
	loc, err := time.LoadLocation(*zone)
	if err != nil {
		panic(err)
	}

	if *work == "" {
		*work = filepath.Join(home, "rowan-working")
	}
	workingPrefix = strings.TrimSuffix(*work, "/") + "/"

	remotes := readRemotes(*work)

	type tok struct{ in, cc, cr, out int64 }
	type turn struct {
		tok   tok
		total int64
		day   string
		cwd   string
		refs  map[string]int
	}

	perRepo := map[string]*tok{}
	turnsPerRepo := map[string]int{}
	var grand int64
	var turnsAll, filesSeen, filesRead, badLines int
	var fallback, multi int

	err = filepath.WalkDir(root, func(p string, d fs.DirEntry, werr error) error {
		if werr != nil || d.IsDir() || !strings.HasSuffix(p, ".jsonl") {
			return nil
		}
		filesSeen++
		f, oerr := os.Open(p)
		if oerr != nil {
			return nil
		}
		defer f.Close()
		filesRead++

		order := []string{}
		byID := map[string]*turn{}
		seg := map[string]int{}
		cwd := ""

		br := bufio.NewReader(f)
		for lineno := 0; ; lineno++ {
			line, rerr := br.ReadString('\n')
			if line != "" {
				scanToolPaths(line, seg)
				if cwd == "" {
					var c struct {
						Cwd string `json:"cwd"`
					}
					if json.Unmarshal([]byte(line), &c) == nil {
						cwd = c.Cwd
					}
				}
				if strings.Contains(line, `"usage"`) {
					var r struct {
						Timestamp string `json:"timestamp"`
						Message   *struct {
							ID    string `json:"id"`
							Usage *struct {
								In  int64 `json:"input_tokens"`
								CC  int64 `json:"cache_creation_input_tokens"`
								CR  int64 `json:"cache_read_input_tokens"`
								Out int64 `json:"output_tokens"`
							} `json:"usage"`
						} `json:"message"`
					}
					if json.Unmarshal([]byte(line), &r) != nil {
						badLines++
					} else if r.Message != nil && r.Message.Usage != nil && r.Timestamp != "" {
						t, perr := time.Parse(time.RFC3339Nano, r.Timestamp)
						if perr != nil {
							badLines++
						} else {
							u := r.Message.Usage
							id := r.Message.ID
							if id == "" {
								id = fmt.Sprintf("noid:%s:%d", p, lineno)
							}
							tt := u.In + u.CC + u.CR + u.Out
							cur, seen := byID[id]
							if !seen {
								cur = &turn{day: t.In(loc).Format("2006-01-02"), cwd: cwd,
									refs: map[string]int{}}
								byID[id] = cur
								order = append(order, id)
							}
							if !seen || tt > cur.total {
								cur.total = tt
								cur.tok = tok{u.In, u.CC, u.CR, u.Out}
							}
							for k, v := range seg {
								cur.refs[k] += v
							}
							seg = map[string]int{}
						}
					}
				}
			}
			if rerr != nil {
				if rerr != io.EOF {
					badLines++
				}
				break
			}
		}

		for _, id := range order {
			t := byID[id]
			if t.day != *day {
				continue
			}
			turnsAll++
			grand += t.total

			byRepo := map[string]int{}
			for dir, n := range t.refs {
				byRepo[foldRepo(dir, remotes)] += n
			}
			if len(byRepo) > 1 {
				multi++
			}
			repo, ok := pick(byRepo)
			if !ok {
				fallback++
				c := t.cwd
				if c == "" {
					c = cwd
				}
				repo = foldRepo(firstSeg(strings.TrimPrefix(c, workingPrefix)), remotes)
				if !strings.HasPrefix(c, workingPrefix) {
					repo = "(cwd:" + firstSeg(strings.TrimPrefix(c, home+"/")) + ")"
				}
			}
			if perRepo[repo] == nil {
				perRepo[repo] = &tok{}
			}
			perRepo[repo].in += t.tok.in
			perRepo[repo].cc += t.tok.cc
			perRepo[repo].cr += t.tok.cr
			perRepo[repo].out += t.tok.out
			turnsPerRepo[repo]++
		}
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "walk:", err)
		os.Exit(1)
	}

	fmt.Printf("INDEPENDENT ATTRIBUTION CHECK day=%s zone=%s\n", *day, *zone)
	fmt.Printf("  jsonl files seen=%d read=%d ; unparseable/undated usage records=%d\n",
		filesSeen, filesRead, badLines)
	fmt.Printf("  turns on the day=%d ; grand total=%d\n", turnsAll, grand)
	fmt.Printf("  turns referencing 2+ repos=%d ; turns decided by cwd fallback=%d\n\n",
		multi, fallback)

	names := make([]string, 0, len(perRepo))
	for n := range perRepo {
		names = append(names, n)
	}
	sort.Slice(names, func(i, j int) bool {
		a, b := perRepo[names[i]], perRepo[names[j]]
		return a.in+a.cc+a.cr+a.out > b.in+b.cc+b.cr+b.out
	})
	fmt.Printf("%-26s %14s %10s %8s\n", "repo", "total", "turns", "share")
	for i, n := range names {
		if i >= *top {
			break
		}
		t := perRepo[n]
		s := t.in + t.cc + t.cr + t.out
		fmt.Printf("%-26s %14d %10d %7.3f%%\n", n, s, turnsPerRepo[n],
			100*float64(s)/float64(grand))
	}
}

// scanToolPaths counts working-root references in this record's TOOL CALLS AND RESULTS
// only. The cwd field is on every record, so counting the whole line would charge every turn
// to whatever directory the shell was standing in -- the bias this whole re-run exists to
// remove.
func scanToolPaths(line string, seg map[string]int) {
	if !strings.Contains(line, workingPrefix) {
		return
	}
	var rec struct {
		Message *struct {
			Content json.RawMessage `json:"content"`
		} `json:"message"`
		ToolUseResult json.RawMessage `json:"toolUseResult"`
	}
	if json.Unmarshal([]byte(line), &rec) != nil {
		return
	}
	if len(rec.ToolUseResult) > 0 {
		harvest(string(rec.ToolUseResult), seg)
	}
	if rec.Message == nil || len(rec.Message.Content) == 0 {
		return
	}
	var blocks []struct {
		Type    string          `json:"type"`
		Input   json.RawMessage `json:"input"`
		Content json.RawMessage `json:"content"`
	}
	if json.Unmarshal(rec.Message.Content, &blocks) != nil {
		return
	}
	for _, b := range blocks {
		if b.Type == "tool_use" {
			harvest(string(b.Input), seg)
		} else if b.Type == "tool_result" {
			harvest(string(b.Content), seg)
		}
	}
}

func harvest(s string, seg map[string]int) {
	pre := workingPrefix
	for i := 0; ; {
		j := strings.Index(s[i:], pre)
		if j < 0 {
			return
		}
		k := i + j + len(pre)
		e := k
		for e < len(s) {
			c := s[e]
			ok := c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' ||
				c == '.' || c == '-' || c == '_' || c == '+'
			if !ok {
				break
			}
			e++
		}
		if name := strings.TrimRight(s[k:e], "."); name != "" {
			seg[name]++
		}
		i = k
	}
}

// pick returns the single most-referenced repo. A tie is not a decision.
func pick(m map[string]int) (string, bool) {
	best, n, ties := "", 0, 0
	for k, v := range m {
		if v > n {
			best, n, ties = k, v, 1
		} else if v == n {
			ties++
		}
	}
	return best, n > 0 && ties == 1
}

// readRemotes maps every checkout under the working root to its github repo, independently
// of the main tool's resolver.
func readRemotes(base string) map[string]string {
	out := map[string]string{}
	ents, _ := os.ReadDir(base)
	for _, e := range ents {
		g := filepath.Join(base, e.Name(), ".git")
		st, err := os.Stat(g)
		if err != nil {
			continue
		}
		cfgPath := filepath.Join(g, "config")
		if !st.IsDir() {
			b, rerr := os.ReadFile(g)
			if rerr != nil {
				continue
			}
			gd := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(string(b)), "gitdir:"))
			if i := strings.Index(gd, "/worktrees/"); i >= 0 {
				gd = gd[:i]
			}
			cfgPath = filepath.Join(gd, "config")
		}
		b, rerr := os.ReadFile(cfgPath)
		if rerr != nil {
			continue
		}
		in := false
		for _, ln := range strings.Split(string(b), "\n") {
			t := strings.TrimSpace(ln)
			if strings.HasPrefix(t, "[") {
				in = t == `[remote "origin"]`
				continue
			}
			if in && strings.HasPrefix(t, "url") {
				if i := strings.Index(t, "="); i >= 0 {
					out[e.Name()] = ghName(strings.TrimSpace(t[i+1:]))
				}
				break
			}
		}
	}
	return out
}

func ghName(u string) string {
	u = strings.ReplaceAll(strings.TrimSuffix(strings.TrimSpace(u), ".git"), ":", "/")
	parts := []string{}
	for _, p := range strings.Split(u, "/") {
		if p != "" {
			parts = append(parts, p)
		}
	}
	if len(parts) < 2 {
		return u
	}
	if parts[len(parts)-2] == "mas-bandwidth" {
		return parts[len(parts)-1]
	}
	return parts[len(parts)-2] + "/" + parts[len(parts)-1]
}

// foldRepo maps a directory name to its repo: its own remote first, else the longest
// still-existing checkout it extends on a '-' boundary, else itself.
func foldRepo(dir string, remotes map[string]string) string {
	dir = strings.TrimRight(dir, ".")
	if dir == "" {
		return "(none)"
	}
	if r, ok := remotes[dir]; ok {
		return r
	}
	cands := []string{dir}
	if s := strings.TrimPrefix(dir, "scratch-"); s != dir {
		cands = append(cands, s)
	}
	best, bestLen := "", 0
	for _, c := range cands {
		for name, repo := range remotes {
			if strings.HasPrefix(c, name+"-") && len(name) > bestLen {
				best, bestLen = repo, len(name)
			}
		}
	}
	if best != "" {
		return best
	}
	return dir
}

func firstSeg(p string) string {
	if i := strings.Index(p, "/"); i >= 0 {
		return p[:i]
	}
	return p
}
