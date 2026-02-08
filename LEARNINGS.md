# Learnings

- 2026-02-08: Graph/workflow internals: blocking computed via blocked_issues_cache (full rebuild on dependency/status changes). Graph links (replies_to, relates_to, duplicates, supersedes) stored as dependency edges (Decision 004) with optional thread_id for threading.
- 2026-02-08: Dolt-removal simplification safe path: remove CLI/flag surfaces first, then delete dead `ShouldImport/ShouldExport` branches (now always true), then validate with `go test -short ./...` + `go build ./cmd/bd`.
- 2026-02-08: After removing Dolt runtime paths, simplify filesystem watcher to JSONL + `.beads` dir only; no need to probe `.beads/dolt/.dolt/*`.
- 2026-02-08: When init command validation uses `os.Exit`, avoid invalid-backend tests through `rootCmd.Execute()`; prefer positive-path assertions (e.g., sqlite BEADS_DIR path placement).
