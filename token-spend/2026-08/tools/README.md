# The tools that produced `../data.csv`

`august-repo-spend-paths.go` is the analysis pass: it walks a bench's `~/.claude/projects`
transcript tree, folds turns by assistant message id, converts UTC timestamps to local days,
attributes each turn to a GitHub repository by the paths appearing in its tool calls and tool
results, folds working directories onto repos through their git `origin` remotes, and writes
the `(bench, day, repo)` CSV — refusing to emit anything unless its grand total and turn count
reconcile exactly against the shared accounting scanner. `august-verify-attrib.go` is a
deliberately independent recomputation of one day, sharing no code with the pass it checks:
its own walk, record struct, day conversion, message-id fold, path scan and remote read. Both
files are published as they ran, with two changes: the working-directory root is now taken
from a `-working` flag (default `$HOME/rowan-working`) instead of a hard-coded absolute path,
and the name table in `seedRepos` carries the public repo names only — the private names each
bench also works on were removed, which cannot misclassify anything because a name that table
does not know is reported as unresolved and classified private fail-closed. **The sanitized
copies were rebuilt and re-run against the same transcript store and reproduced `data.csv`'s
source CSV byte-for-byte.** Both import `github.com/mas-bandwidth/rowan-tools/modules/spend`
and `.../modules/usage`, which live in the private workshop repo the Studio already has; drop
these files into `rowan-tools/scratch/` and `go run` them from the repo root. Running them
needs a `repo-visibility.json` in the output directory, built from real
`gh repo view <name> --json isPrivate` answers — it is not shipped here, because on this bench
it is a list of repository names.
