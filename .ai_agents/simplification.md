# beads simplification plan (v1)

## Goal
- Keep core: git-backed JSONL collaboration + local fast DB + dependency-aware ready queue.
- Cut moving parts; reduce duplicate code paths.

## Keep intact
- Hash IDs + distributed merge model (`internal/types/id_generator.go`, `docs/COLLISION_MATH.md`).
- `bd ready` blocker semantics (`internal/storage/sqlite/ready.go`).
- JSONL sync artifact (`docs/SYNC.md`).

## Simplification opportunities (priority)
1. Single sync engine
- Merge `cmd/bd/sync.go` + `cmd/bd/daemon_sync.go` into one `SyncService`.
- CLI + daemon call same state machine.

2. Single command execution layer
- Commands call app service interface; transport hidden (RPC/direct).
- Remove per-command daemon/direct branching in files like `cmd/bd/create.go`, `cmd/bd/ready.go`.

3. Central lock manager
- One package for lock ordering/lifetime.
- Cover `.sync.lock`, `.jsonl.lock`, daemon lock, exclusive lock.

4. Adaptive blocked cache
- Small repos: compute blocked via direct query/CTE.
- Large repos: enable `blocked_issues_cache`.
- Cut invalidation complexity for common case.

5. De-scope Dolt from core runtime
- Move Dolt-heavy flows to optional mode/plugin/build tag.
- Keep SQLite+JSONL as default core path.

6. Reduce watcher complexity
- Replace heavy watcher/event plumbing with:
  - pull/import freshness check at command start
  - periodic sync tick in daemon
  - explicit `bd sync` for force

7. Narrow dependency semantics
- Split edge semantics:
  - workflow edges (affect ready): `blocks`, `parent-child`, `conditional-blocks`, `waits-for`
  - knowledge edges (no ready impact): `relates-to`, `duplicates`, `supersedes`, `replies-to`
- Keep one table if desired; simplify invalidation + cycle rules by class.

## Q&A from review

### 5) What is Dolt?
- Dolt = SQL database with git-like version control (commits/branches/merges on tables/rows).
- In beads: optional backend alternative to SQLite (`internal/storage/dolt/*`, `docs/DOLT.md`).
- Why complexity grows:
  - extra backend code paths
  - extra sync modes (`dolt-native`, `belt-and-suspenders`)
  - backend-specific locks/daemon constraints
- Simplification proposal not “remove Dolt”; “isolate Dolt from default runtime path”.

### 5b) Dolt vs SQLite+JSONL: alternative or additional?
- SQLite vs Dolt: alternative DB backend.
- JSONL usage depends on sync mode (`docs/CONFIG.md`):
  - `git-portable` (default): JSONL is sync path.
  - `dolt-native`: Dolt remote is sync path; JSONL not used for sync (still import/export capable).
  - `belt-and-suspenders`: both Dolt remote + JSONL.
- So: Dolt can replace SQLite; JSONL can be primary, secondary, or bypassed for sync depending mode.

### 6) Explain “reduce watcher complexity”
- Current daemon path has fsnotify watcher + polling fallback + debounce + dropped-event checks (`cmd/bd/daemon_watcher.go`, `cmd/bd/daemon_event_loop.go`).
- Proposal:
  - on each command start: cheap staleness check; auto-import if JSONL newer
  - daemon: periodic remote sync only
  - keep manual override: `bd sync`
- Tradeoff:
  - less instant propagation between concurrent shells
  - much simpler failure surface/race model

### 6b) Is this specifically because of sync branch worktree (`beads-sync`)?
- Not primarily.
- Watcher/event complexity exists even without sync branch (cross-platform FS events, debounce/race recovery).
- Sync-branch worktree adds extra edge cases/path handling (`docs/PROTECTED_BRANCHES.md`, `internal/git/worktree.go`), so it increases complexity, but not the root cause.

### 7) Explain “narrow dependency type semantics”
- Today many edge types share one dependency system; only some should gate execution.
- Proposal: formal edge classes.
  - execution-critical edges drive `ready`, blocker cache, strict cycle checks
  - informational edges only for graph nav/audit/threading
- Benefits:
  - fewer cache rebuild triggers
  - fewer surprising “why blocked?” cases
  - clearer mental model for users + agents

## Suggested rollout
1. Start with #2 (single command service) + #1 (single sync engine). Biggest duplication cut.
2. Then #3 (lock manager) while tests protect race behavior.
3. Then #6 and #7 behind feature flags.
4. Keep #5 as packaging/boundary refactor; no behavior change first pass.

## External analysis review (fact-check)

### Confirmed
- `createInRig()` is heavy duplication of main create path (`cmd/bd/create.go`).
- Routing boilerplate repeated in many commands (`needsRouting` + `resolveAndGetIssueWithRouting`).
- Routing docs in website advertise `bd routes *` commands not present in CLI (`website/docs/multi-agent/routing.md`).
- Daemon surface area is large; simplification potential is real (`cmd/bd/daemon*.go`).
- Dolt adds substantial optional surface (`internal/storage/dolt/*`, `cmd/bd/dolt*.go`).

### Partly true / needs nuance
- Formula dead-code claim overstated:
  - `ApplyAdvice` called from cook path.
  - `ApplyLoops`/`ApplyBranches` called by `ApplyControlFlow`, which cook uses.
  - Better framing: advanced formula features are under-used by current shipped formula, not fully dead.
- Issue field deadness mixed:
  - `Crystallizes` is active in sqlite/dolt paths.
  - `WorkType`/`QualityScore` appear schema-level but weakly wired in sqlite runtime.
  - `Creator`/`Validations` affect content hash + merge semantics even if sqlite persistence is incomplete.
- Event-system “duplication” also nuanced:
  - event-beads (`type=event`) are domain objects.
  - events table is audit trail.
  - overlap exists, but not automatically redundant.

### Action from this review
- Prioritize low-risk dedup + doc correctness first.
- Treat formula and event-system deletions as later, evidence-driven work.

## Dolt decision options

### Option A (recommended): SQLite core, Dolt optional module
- Keep default path: SQLite + JSONL.
- Move Dolt behind plugin/build tag/module boundary.
- Fast simplification without breaking default users.

### Option B: remove Dolt now
- Maximum simplification quickly.
- Breaks Dolt users + migration/federation workflows.
- Requires explicit deprecation + migration plan.

### Option C: Dolt-only
- Not recommended currently.
- Conflicts with current default workflows/docs; daemon behavior differs; higher migration risk.

## Execution plan (approved direction): remove Dolt entirely

## Scope decision
- Hard-remove Dolt backend/federation from this fork.
- Keep only: `sqlite` backend + `jsonl` no-db mode.
- Accept breaking Dolt users.

## Phase 0: Safety + branch prep
1. Create migration note doc: `docs/MIGRATION_DOLT_REMOVED.md`.
2. Add startup guard:
- If metadata backend == `dolt`, fail fast with explicit recovery steps.
- Recovery steps: pin to older commit to export data, then re-init sqlite.
3. Add test for fail-fast guard.

## Phase 1: Backend/config contraction
Files:
- `internal/configfile/configfile.go`
- `internal/config/sync.go`
- `cmd/bd/backend.go`
- `cmd/bd/init.go`
- `cmd/bd/main.go`
- `cmd/bd/config.go`
- `cmd/bd/sync_mode_cmd.go`

Changes:
- Remove `BackendDolt` constants/branches.
- Remove Dolt server-mode config fields/logic.
- Sync modes: keep only `git-portable`, `realtime`.
- Remove `--dolt-auto-commit` flag + plumbing.
- `bd backend list/show`: only sqlite/jsonl.
- `bd init --backend`: reject `dolt` with clear error.

## Phase 2: Delete Dolt/federation command surface
Delete files:
- `cmd/bd/dolt*.go`
- `cmd/bd/federation*.go`
- `cmd/bd/migrate_dolt*.go`
- Dolt/federation-specific doctor checks:
  - `cmd/bd/doctor/dolt.go`
  - `cmd/bd/doctor/federation.go`
  - `cmd/bd/doctor/federation_nocgo.go`
  - `cmd/bd/doctor/perf_dolt.go`
  - dolt migration-validation files if solely Dolt-focused

Refactor:
- Strip references from:
  - `cmd/bd/doctor/*.go`
  - `cmd/bd/daemon_*`
  - `cmd/bd/sync.go`
  - `cmd/bd/hook.go`
  - `cmd/bd/activity_watcher.go`

## Phase 3: Remove storage implementation
Delete:
- `internal/storage/dolt/` (entire package).

Refactor:
- `internal/storage/factory/factory.go`
  - Remove Dolt options/registry expectations tied to Dolt.
  - Keep sqlite construction path minimal.
- `internal/storage/versioned.go`
  - Remove federation/versioned interfaces if now unused.
  - Or keep only pieces used by sqlite path.

## Phase 4: Sync path simplification
Files:
- `cmd/bd/sync.go`
- `cmd/bd/sync_export.go`
- `cmd/bd/daemon_sync_branch.go`
- `internal/rpc/server_export_import_auto.go`

Changes:
- Remove dolt-native / belt-and-suspenders branches.
- Remove Dolt remote push/pull/commit paths.
- Keep pure git+JSONL sync behavior.

## Phase 5: Dependency cleanup
Files:
- `go.mod`, `go.sum`

Actions:
- Remove direct Dolt deps (`github.com/dolthub/driver`, etc).
- Remove replace for `go-icu-regex` if unused.
- `go mod tidy`.
- Confirm no cgo-only Dolt stubs remain.

## Phase 6: Tests and docs cleanup
Delete/update tests:
- Remove Dolt-only tests under `cmd/bd/*dolt*`, `internal/storage/dolt/*`, related doctor tests.
- Update sync-mode tests to 2 modes only.

Add regression tests:
1. `bd init --backend dolt` => explicit unsupported error.
2. `bd sync mode set dolt-native` => invalid mode error.
3. `bd backend list` no `dolt`.
4. Workspace with backend=`dolt` fails with migration hint.

Docs:
- Delete/replace:
  - `docs/DOLT.md`
  - `docs/DOLT-BACKEND.md`
- Update references in:
  - `README.md`
  - `docs/CONFIG.md`
  - `docs/CLI_REFERENCE.md`
  - `docs/QUICKSTART.md`
  - `docs/GIT_INTEGRATION.md`
  - `docs/INSTALLING.md`
  - `docs/EXTENDING.md`
- Fix fictional routes page:
  - `website/docs/multi-agent/routing.md` (remove non-existent `bd routes` commands).

## Phase 7: Validation gates
Run:
1. `go test -short ./...`
2. `go test ./cmd/bd/...`
3. `go test ./internal/...`
4. `golangci-lint run ./...`
5. `go build ./cmd/bd`

Manual smoke:
1. `bd init`
2. `bd create "x"`
3. `bd ready`
4. `bd sync`
5. `bd backend list`
6. `bd sync mode show`

## Delivery strategy
- PR-A: Phase 0-2 (command/config/API surface removal).
- PR-B: Phase 3-5 (storage + deps cleanup).
- PR-C: Phase 6-7 (tests/docs finalization).

## Risk controls
- Keep fail-fast message for old Dolt workspaces until docs and migration note merged.
- Do not silently auto-convert Dolt metadata to sqlite.
- Prefer explicit error + operator action.

## Verification snapshot (2026-02-08)

### Completed in this pass
- Removed CLI surfaces:
  - `cmd/bd/dolt_nocgo.go`
  - `cmd/bd/federation_nocgo.go`
  - `cmd/bd/migrate_dolt_nocgo.go`
  - `cmd/bd/migrate_dolt_cmd_nocgo.go`
  - `cmd/bd/diff.go`
  - `cmd/bd/history.go`
- Removed migrate flags `--to-dolt` and `--to-sqlite` (`cmd/bd/migrate.go`).
- Config validation now only accepts:
  - `sync.mode`: `git-portable|realtime`
  - `conflict.strategy`: `newest|manual|ours|theirs`
- Removed daemon federation CLI flags; trimmed daemon function signatures:
  - no `--federation`, `--federation-port`, `--remotesapi-port`
  - removed corresponding plumbing in `cmd/bd/daemon*.go`.
- Removed dead Dolt-mode branches from runtime paths:
  - `cmd/bd/sync_export.go`
  - `cmd/bd/daemon_sync.go`
  - `cmd/bd/daemon_sync_branch.go`
  - `cmd/bd/staleness.go`
  - `cmd/bd/autoimport.go`
  - `cmd/bd/autoflush.go`
  - `cmd/bd/direct_mode.go`
- Removed dead `prefer-dolt` config parsing path (`cmd/bd/autoimport.go`).
- Removed hook-time Dolt skip branches (`cmd/bd/init_git_hooks.go`).
- Deleted docs:
  - `docs/DOLT.md`
  - `docs/DOLT-BACKEND.md`

### Validation results
- `go test -short ./...` ✅
- `go build ./cmd/bd` ✅

### Remaining Dolt/federation surfaces (next pass)
- Runtime/behavioral:
  - `cmd/bd/activity_watcher.go` still probes `.beads/dolt/.dolt/*` paths.
  - `internal/configfile/configfile.go` still retains `BackendDolt` backward-compat capability mapping.
  - `cmd/bd/dolt_server_nocgo.go` and doctor migration-validation stubs still exist as explicit error adapters.
- Docs:
  - `docs/CONFIG.md`, `docs/QUICKSTART.md`, `docs/GIT_INTEGRATION.md`, `docs/CLI_REFERENCE.md`, `website/docs/getting-started/quickstart.md` still reference Dolt modes/flags.
- Tests/comments:
  - many tests still mention `dolt-native`/`belt-and-suspenders` as invalid legacy values (mostly OK, but noisy).

## Verification snapshot (2026-02-08, pass 2)

### Completed in this pass
- Simplified watcher paths to JSONL + `.beads` directory only:
  - `cmd/bd/activity_watcher.go`
- Removed unused Dolt server stub:
  - `cmd/bd/dolt_server_nocgo.go`
- Simplified backend capability logic:
  - removed `BackendDolt` constant, treat unknown/legacy backend strings as conservative single-process:
  - `internal/configfile/configfile.go`
  - `internal/configfile/configfile_test.go`
- Removed stale Dolt wording in runtime comments:
  - `cmd/bd/create.go`
  - `internal/beads/beads.go`
- Updated user-facing docs to SQLite+JSONL only:
  - `docs/CONFIG.md`
  - `docs/QUICKSTART.md`
  - `docs/GIT_INTEGRATION.md`
  - `docs/CLI_REFERENCE.md`
  - `website/docs/getting-started/quickstart.md`

### Validation results
- `go test -short ./...` ✅
- `go build ./cmd/bd` ✅

### Remaining cleanup candidates
- Historical/reference text only:
  - release notes in `cmd/bd/info.go`
  - legacy/compat tests mentioning invalid Dolt modes.
- Doctor compatibility stubs still mention Dolt removal guidance:
  - `cmd/bd/doctor/migration_validation_nocgo.go`
  - lock cleanup checks for `dolt.bootstrap.lock` in doctor.

## Verification snapshot (2026-02-08, pass 3)

### Completed in this pass
- Removed remaining Dolt-branded wording/tokens from code/docs/tests under:
  - `cmd/`
  - `internal/`
  - `docs/`
  - `website/`
- Updated compatibility test and comments to neutral "legacy"/generic invalid-mode wording.
- Reworked `init` BEADS_DIR backend test to sqlite path verification:
  - `cmd/bd/init_test.go`
- Simplified doctor legacy migration stubs and lock checks:
  - switched bootstrap lock checks to `bootstrap.lock`
  - removed explicit Dolt naming from migration-validation structures/messages
  - files: `cmd/bd/doctor/locks.go`, `cmd/bd/doctor/fix/locks.go`, `cmd/bd/doctor/migration_validation_nocgo.go` (+ tests)
- Cleaned historical "what's new" strings to remove Dolt naming while preserving chronology:
  - `cmd/bd/info.go`

### Validation results
- `go test -short ./cmd/bd/...` ✅
- `go test -short ./internal/...` ✅
- `go build ./cmd/bd` ✅

### Remaining cleanup candidates
- Optional only:
  - rename legacy helper symbols/files (`migration_validation_nocgo`) for naming consistency.
