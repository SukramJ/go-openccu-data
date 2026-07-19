#!/bin/bash
# Regenerate the embedded data snapshot from an openccu-data checkout.
#
# Usage: script/regenerate.sh <path-to-openccu-data-checkout>
#
# Copies the distributed artifacts into data/ and stamps
# SnapshotVersion in openccudata.go from the checkout's const.py.
# Run by the regenerate-on-data-release workflow on every openccu-data
# release; usable locally for development snapshots.
set -euo pipefail

SRC="${1:?usage: script/regenerate.sh <openccu-data-checkout>}"
DATA_SRC="$SRC/openccu_data/data"
[ -d "$DATA_SRC" ] || { echo "not an openccu-data checkout: $SRC" >&2; exit 1; }

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

rm -rf "$ROOT/data"
mkdir -p "$ROOT/data/profiles" "$ROOT/data/translation_custom"
cp "$DATA_SRC/translation_extract.json.gz" "$ROOT/data/"
cp "$DATA_SRC/easymode_extract.json.gz" "$ROOT/data/"
cp "$DATA_SRC/device_semantics.json" "$ROOT/data/"
cp "$DATA_SRC"/profiles/*.json.gz "$ROOT/data/profiles/"
cp "$DATA_SRC/profiles/_receiver_type_aliases.json" "$ROOT/data/profiles/"
cp "$DATA_SRC"/translation_custom/*.json "$ROOT/data/translation_custom/"

VERSION="$(sed -n 's/^VERSION: Final = "\(.*\)"$/\1/p' "$SRC/openccu_data/const.py")"
[ -n "$VERSION" ] || { echo "could not read VERSION from const.py" >&2; exit 1; }
sed -i.bak "s/^const SnapshotVersion = \".*\"$/const SnapshotVersion = \"$VERSION\"/" "$ROOT/openccudata.go"
rm -f "$ROOT/openccudata.go.bak"

echo "snapshot regenerated from openccu-data $VERSION"
