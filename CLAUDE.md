# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```sh
go test ./...                  # run all tests
go test -run TestName ./...    # run a single test
go vet ./...                   # static checks (CI runs this)
go run mvdan.cc/gofumpt@latest -l .   # CI fails if any file is not gofumpt-formatted
script/regenerate.sh <path-to-openccu-data-checkout>   # regenerate data/ snapshot
```

There are no external dependencies; the module is a single package at the repo root.

## What this repository is

A **data artifact, not a lookup framework** — this constraint is deliberate and should be preserved. It ships the [openccu-data](https://github.com/SukramJ/openccu-data) metadata extracts (translations, easymodes, MASTER-paramset profiles, curated overlays, device semantics, the WebUI device images under `device_images/250/`) as an embedded Go module. Consumers (primarily [openccu-loom](https://github.com/SukramJ/openccu-loom), `internal/ccudata`) decode the archives themselves and layer lookup semantics on top. Resist adding parsing/lookup logic here; the only typed accessor is `DoorbellModels()` for the small curated `device_semantics.json`.

## Architecture

- `openccudata.go` — the entire API: `ReadFile`/`ReadDir` over an `embed.FS`, plus `SnapshotVersion` and `DoorbellModels()`. The embed uses `//go:embed all:data` — the `all:` prefix is required, otherwise underscore-prefixed files (`profiles/_receiver_type_aliases.json`) are silently dropped.
- `data/` — **generated, never hand-edited.** `script/regenerate.sh` deletes and recreates the whole tree from an openccu-data checkout and stamps `const SnapshotVersion` in `openccudata.go` from the checkout's `const.py`.
- `.github/workflows/regenerate-on-data-release.yml` — openccu-data's release workflow fires `repository_dispatch` (event `data-release`); this regenerates the snapshot and opens a version-bump PR. Manual fallback: `workflow_dispatch` with the release tag. It pushes with `RELEASE_PAT` (not `GITHUB_TOKEN`) so the PR still triggers CI.

## Versioning

Two independent version schemes: module tags are SemVer (`v0.x`), while `SnapshotVersion` is the upstream openccu-data CalVer release (e.g. `2026.7.0`) the embedded data was taken from. After a regeneration PR merges, a new module tag is expected; pushing a tag triggers `release-on-tag.yml`, which creates the GitHub Release with the embedded snapshot version in its notes.

## Licensing split

Code is MIT; the extracted data under `data/` inherits the eQ-3 HomeMatic Software License (free for private and non-commercial use only). Exception: `data/translation_custom/*.json` are hand-curated by the openccu-data maintainers and MIT. See `NOTICE.md` before moving or redistributing data files.

---

## Which model does which work

Set by the owner for the go-fabric family on 2026-10-10 (go-fabric's
`CLAUDE.md`); the review split of 2026-11-08 (go-fabric #456) likewise.
Adopted here on 2026-11-08. A statement about openccu-data's extracts that
lands in this repo is cited later as fact — a comment that names a caller
or a consumer, a docstring that states what the other side does, a
release note — so the work that produces such statements gets
the strongest model, and the work a guard checks may use a cheaper one. The
agents in [Sub-Agent Delegation](#sub-agent-delegation) carry a default
model each; the rows below decide when a call overrides it (`model:` on the
Agent tool). Implementation defaults to Opus.

| Work | Model |
| --- | --- |
| Analysis, planning, reviews of design, audits; anything written into `NOTICE.md` or `README.md`, or a comment that names an openccu-data artefact, its CalVer or its licence; the regeneration script and workflow (`script/regenerate.sh`, `regenerate-on-data-release.yml` with `RELEASE_PAT`), the licence boundary in `NOTICE.md`; *which* guard gets built and its bite line; rebases and merge conflicts where `main` may have overtaken the change | **Fable** (the default for this repo; also the `hunt` agent) |
| Implementation of a finished plan with guards behind it: a test on the embedded data once the pattern exists, the release note for a finished diff | **Opus** (the `impl` default; also what fast mode runs on) |
| **Review of a slice's pull request, split (owner decision 2026-11-08, go-fabric #456).** The code review — correctness, concurrency, error paths, tests, the fixed report shape — is **Opus**. On top of it **Fable** reads, bounded: every new comment or docstring that names an openccu-data artefact, its CalVer or its licence against openccu-data at `SnapshotVersion`, every line that enters `NOTICE.md` or `README.md`, the paths of the first row, and the brief's decisions against the reference — the behavioural gaps a code review does not see. The load-bearing claims of the Opus review are verified before they reach the implementer. | **Opus**, then **Fable** |
| Work a test or generator checks immediately: `script/regenerate.sh`, `gofumpt`; search and first reading of a tree whose result a stronger model or the owner reads again | **Sonnet** (the `guard` default; `impl` with `model: "sonnet"`) |
| Mechanical inventories whose answer is a list the caller can re-grep: which files reference X, which links are dead, which entries lack a counterpart | **Haiku** (`sweep` default) |

Sonnet never writes a claim about an openccu-data artefact, its CalVer or its licence into the repo, and
never touches the paths named in the first row. Haiku returns lists, never
verdicts. When in doubt about which row a task is in, it is in the first.

---

## Sub-Agent Delegation

As go-fabric's and openccu-loom's rule (owner, 2026-10-10). The main
conversation owns planning, contract and wire decisions, guard
specifications, sub-agent steering, and final verification. It does not own
typing.

**Delegate what a gate can prove; keep what only a reader can accept.**
Delegate freely: table cases, fakes, scaffolding, fixture updates, doc and
link cleanups, version bumps across the changelog set, inventories and grep
sweeps. Never delegate: the regeneration script, the workflows, the licence boundary, and *which* guard gets
built. Locating is delegable; reading is not.

| Agent | Default model | Use for |
|---|---|---|
| `impl` | Opus | scoped implementation with a stated acceptance command |
| `guard` | Sonnet | tests from a caller-written guard spec, plus the bite proof |
| `sweep` | Haiku | read-only inventories and grep sweeps, high fan-out |
| `hunt` | Fable | adversarial read-only defect hunt, ranked candidates |

Each default is one row of
[Which model does which work](#which-model-does-which-work); a call that
falls into another row passes `model:` on the Agent tool — `impl` with
`model: "sonnet"` for work a generator or formatter checks immediately.

**Size CPU-bound fan-out from the host, never from a constant** — the owner
works from a 4-core box and a 14-core box:

```sh
cores=$(nproc); agents=$(( (cores - 1) / share ))   # share = cores per agent
```

Check the 1-minute load first; if it is at or above the core count, run one
CPU-bound agent and queue the rest. Eight concurrent agents is a practical
ceiling whatever the core count. Every CPU-bound agent must **pin** its share
in the command it runs (`GOMAXPROCS=<share> go test -p <share> ./internal/x/...`, `golangci-lint run --concurrency=<share>`) — handing out a number without
pinning it changes nothing. Read-only agents cost no core worth counting;
6–8 is a practical ceiling there.

Two hard rules: **no sub-agent runs the full gate (`go test ./...`, `go vet ./...`, `gofumpt -l .`) or a
repo-wide lint** — the full gate runs once, in the main conversation, at the
end of the slice — and **no two writing agents in the same package**
(disjoint file sets, or `isolation: "worktree"`).

Every brief carries five things: the files and public surface plus one style
reference; the acceptance command the agent runs itself; the invariants it
touches; **stop conditions** (new dependency, API or wire change, credential
handling, an ADR-shaped decision — a stop is a successful outcome); and a
report format ≤ 250 words, large artefacts to the scratchpad as a path.

**A report is a claim.** Read the diff, run the gate yourself. For `guard`,
the acceptance artefact is not a green test but the **bite proof**: the named
production line removed, the test observed red with its message, the line
restored, green again.

---

## Token use

Set by the owner on 2026-10-14 ("Wichtig ist mir eine effiziente Nutzung
der Token"), as go-fabric's rule of the same date. A slice costs on
the order of half a million to a million tokens per implementing agent
before the owner's review; the rules below are what keeps that from
doubling.

- **Read slices, not files.** `sed -n a,bp`, `grep -n` with a tight
  pattern and `head`; a whole file only when the task is the file. Never
  re-read what is already in context; cite the earlier read.
- **Delegate first reads.** A source tree's inventory (openccu-data, openccu-loom's `internal/ccudata`, this repo)
  is a `sweep` or Sonnet read written to the scratchpad with `file:line`
  for every fact; the owner spot-checks the load-bearing lines and verifies
  every fact about openccu-data's extracts before it enters a brief or
  `NOTICE.md` or `README.md` — the agent's result is a map, never a citation.
- **Logs by their summary.** The gate (`go test ./...`, `go vet ./...`, `gofumpt -l .`) runs in the
  background with a summary file (one line per gate, `FAIL`/`ok` lines
  only); read the summary, and open a log only on failure, with `grep` for
  the failure markers and the error strings — never `cat`. The same for CI:
  the job log is fetched once to a file and grepped.
- **One turn, every independent call.** Reads, checks and background
  starts that do not depend on each other go out together; a wait is a
  `Monitor` or an `until` loop in the background, not polling by hand.
- **Short artefacts.** A brief names files, decisions and the acceptance
  command, not prose; a report follows the fixed shape of
  [Sub-Agent Delegation](#sub-agent-delegation) (≤ 250 words, large
  artefacts as a scratchpad path); status messages to the owner are short,
  outcome first, in German.
- **Skip what cannot have changed.** No regeneration for a comment change;
  no second gate run for a notes-only commit — say what was skipped and
  why.
- **Edit by script.** Several replacements in several files are one
  scripted call with asserted anchors, not a chain of interactive edits;
  a generated file (`data/`) is regenerated, never patched.
- **Agents read what the brief names.** An implementing agent gets the
  file list and the pinned sources in the brief, reads those, and reports
  in the fixed shape; it does not survey the tree. One agent per brief;
  two agents only on disjoint paths, at most one Opus agent at a time.
- **Push, pull request and tag are the session's own calls.** An
  implementing agent works locally — branch, commits, guards — and
  reports; the owner's session reads the diff, runs the gate, pushes,
  opens the pull request and tags the merge commit; the owner merges. A
  delegated push, pull request or tag is refused by the auto-mode
  classifier ("Create Public Surface", 2026-11-07), and the review in
  [Which model does which work](#which-model-does-which-work) sits
  between the agent's report and the push anyway.
