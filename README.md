# go-openccu-data

[openccu-data](https://github.com/SukramJ/openccu-data) as a versioned
Go artifact: the OCCU/OpenCCU metadata extracts (translations,
easymodes, MASTER-paramset profiles, curated overlays, curated device
semantics) embedded into a Go module with thin typed accessors.

The snapshot under `data/` is regenerated automatically on every
openccu-data release (`repository_dispatch` → regenerate workflow →
version-bump PR); `openccudata.SnapshotVersion` records the upstream
release it was taken from. Consumers decode the archives themselves —
this module deliberately stays a **data artifact**, not a lookup
framework.

## Usage

```go
import openccudata "github.com/SukramJ/go-openccu-data"

raw, err := openccudata.ReadFile("translation_extract.json.gz")
models := openccudata.DoorbellModels() // curated device semantics
```

Primary consumer:
[openccu-loom](https://github.com/SukramJ/openccu-loom)
(`internal/ccudata` decodes the archives and layers its lookup
semantics on top).

## Versioning

Module tags are independent SemVer (`v0.x`); the embedded data stand
is identified by `SnapshotVersion` (the openccu-data CalVer release).

## Regeneration

```sh
script/regenerate.sh <path-to-openccu-data-checkout>
```

CI regeneration is triggered by openccu-data's release workflow via
`repository_dispatch` (event `data-release`); manual fallback:
`workflow_dispatch` with the release tag.

## Licensing

Code is MIT. The extracted data inherits the eQ-3 HomeMatic Software
License (free for private and non-commercial use) — see
[`NOTICE.md`](./NOTICE.md).
