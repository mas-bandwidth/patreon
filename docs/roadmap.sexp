; The patreon roadmap: the data ROADMAP.md is generated from.
; Same shape as nova-sprint's docs/roadmap.sexp (internal/roadmapdoc's package comment):
; data, never evaluated. Edit this file, then regenerate ROADMAP.md with nova-sprint's
; generator, from a nova-sprint checkout:
;   go build -o /tmp/roadmap ./tools/roadmap && (cd <this repo> && /tmp/roadmap --file docs/roadmap.sexp --out ROADMAP.md)
; Never edit ROADMAP.md by hand. Created 2026-10-14 when this repository's open GitHub issues and
; unfinished pull requests were moved into it.
(roadmap "v1"
 :title "patreon roadmap"
 :text "Planned work for patreon. The items below were open GitHub issues and pull requests, moved here on 2026-10-14 when roadmaps became the standard for future work (Glenn, 2026-10-14). Each item names its issue or pull request; the branches of moved pull requests are kept to resume from. Landed work is in git."

 :groups
 (
  (group "github-issues" :title "Moved from GitHub issues on 2026-10-14"
   :text "Open GitHub issues of mas-bandwidth/patreon, recorded here on 2026-10-14 when roadmaps became the standard for future work (Glenn, 2026-10-14). Each item names its issue; the issue is closed with a pointer here.")
 )

 :items
 (
  (item "issue-6" :group "github-issues" :area "patreon"
   :title "token-spend: an elided path counts as a reference to (none) and can tie a real repo into the fallback"
   :text "Found by the independent check on the studio bench's busiest day (tools/august-verify-attrib.go vs tools/august-repo-spend-paths.go, 2026-08-30). In the main pass, countWith records a reference whenever the working-root prefix is followed by path characters, and resolver.repoFor trims trailing dots and maps an all-dots name to the repo key (none)."
   :date "2026-10-14"
   :origin "mas-bandwidth/patreon issue #6, opened 2026-09-02 by rowan-claude; moved here 2026-10-14"
   :kept (:github-issue "6" :labels ""))
 ))
