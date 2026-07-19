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

A **data artifact, not a lookup framework** — this constraint is deliberate and should be preserved. It ships the [openccu-data](https://github.com/SukramJ/openccu-data) metadata extracts (translations, easymodes, MASTER-paramset profiles, curated overlays, device semantics) as an embedded Go module. Consumers (primarily [openccu-loom](https://github.com/SukramJ/openccu-loom), `internal/ccudata`) decode the archives themselves and layer lookup semantics on top. Resist adding parsing/lookup logic here; the only typed accessor is `DoorbellModels()` for the small curated `device_semantics.json`.

## Architecture

- `openccudata.go` — the entire API: `ReadFile`/`ReadDir` over an `embed.FS`, plus `SnapshotVersion` and `DoorbellModels()`. The embed uses `//go:embed all:data` — the `all:` prefix is required, otherwise underscore-prefixed files (`profiles/_receiver_type_aliases.json`) are silently dropped.
- `data/` — **generated, never hand-edited.** `script/regenerate.sh` deletes and recreates the whole tree from an openccu-data checkout and stamps `const SnapshotVersion` in `openccudata.go` from the checkout's `const.py`.
- `.github/workflows/regenerate-on-data-release.yml` — openccu-data's release workflow fires `repository_dispatch` (event `data-release`); this regenerates the snapshot and opens a version-bump PR. Manual fallback: `workflow_dispatch` with the release tag. It pushes with `RELEASE_PAT` (not `GITHUB_TOKEN`) so the PR still triggers CI.

## Versioning

Two independent version schemes: module tags are SemVer (`v0.x`), while `SnapshotVersion` is the upstream openccu-data CalVer release (e.g. `2026.7.0`) the embedded data was taken from. After a regeneration PR merges, a new module tag is expected; pushing a tag triggers `release-on-tag.yml`, which creates the GitHub Release with the embedded snapshot version in its notes.

## Licensing split

Code is MIT; the extracted data under `data/` inherits the eQ-3 HomeMatic Software License (free for private and non-commercial use only). Exception: `data/translation_custom/*.json` are hand-curated by the openccu-data maintainers and MIT. See `NOTICE.md` before moving or redistributing data files.
