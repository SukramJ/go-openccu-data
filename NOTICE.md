# NOTICE

This repository ships two distinct kinds of files that are governed by
different licenses. Please respect both when redistributing.

## 1. Original code (MIT License)

All Go source code in this repository — `openccudata.go`, the tests, the
`script/` helpers and the build/CI configuration — is original work and
licensed under the [MIT License](./LICENSE).

## 2. Extracted data artifacts (eQ-3 Homematic Software License)

The files under `data/` are copied unchanged from the
[openccu-data](https://github.com/SukramJ/openccu-data) release named by
`SnapshotVersion`. They are derivative works generated from
[OpenCCU-Base](https://github.com/homematicip/OpenCCU-Base) maintained by
eQ-3 AG, and from compatible distributions such as
[OpenCCU](https://github.com/OpenCCU/OpenCCU).

Specifically these files:

- `data/easymode_extract.json.gz`
- `data/translation_extract.json.gz`
- `data/device_semantics.json`
- `data/translation_custom/*.json` _(curated additions, MIT)_
- `data/profiles/*.json.gz` (+ `_receiver_type_aliases.json`)
- `data/device_images/250/**/*.png`

are obtained by parsing TCL configuration and JavaScript translation files
under `www/` of OpenCCU-Base; the device images are unmodified byte-for-byte
copies of `www/config/img/devices/250/` of OpenCCU-Base. They retain the
licensing of the original upstream sources. Per OpenCCU-Base's
`licenses/licenses.md`, files are published under the "Homematic Software
License" version 2.0 (HMSL 2.0) unless stated otherwise; `www/` is not listed
among the exceptions there. Refer to OpenCCU-Base's `licenses/HMSL2.txt` for
the full terms — in short: free for private and non-commercial use; commercial
redistribution requires permission from eQ-3.

The `translation_custom/` files are the exception inside the data tree: they
contain hand-curated translation overrides authored by the openccu-data
maintainers and are released under the MIT License together with the rest of
the code.

## Trademarks

"Homematic" and "HomematicIP" are trademarks of eQ-3 AG. This project is not
affiliated with or endorsed by eQ-3 AG.
