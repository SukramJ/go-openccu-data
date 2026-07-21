// SPDX-License-Identifier: MIT
// Copyright (C) 2026 go-openccu-data authors.

// Package openccudata ships the OCCU/OpenCCU metadata extracts of the
// openccu-data project as a versioned Go artifact: the raw archives
// (translations, easymodes, MASTER-paramset profiles, curated
// overlays) plus thin typed accessors for the curated
// device-semantics classifications.
//
// The data snapshot under data/ is regenerated on every openccu-data
// release (see script/regenerate.sh and the repository_dispatch
// workflow); SnapshotVersion records the upstream release it was
// taken from. Consumers decode the archives themselves — this module
// deliberately stays a data artifact, not a lookup framework.
//
// Licensing is split: this module's code is MIT; the extracted data
// inherits the eQ-3 HomeMatic Software License (free for private and
// non-commercial use). See NOTICE.md.
package openccudata

import (
	"embed"
	"encoding/json"
	"strings"
	"sync"
)

// SnapshotVersion is the openccu-data release this data snapshot was
// generated from. Stamped by script/regenerate.sh.
const SnapshotVersion = "2026.7.1"

// files carries the embedded data snapshot. Exposed through ReadFile
// so the path layout stays an implementation detail.
//
// The all: prefix keeps underscore-prefixed files
// (profiles/_receiver_type_aliases.json) in the embed — plain
// path patterns silently skip them.
//
//go:embed all:data
var files embed.FS

// ReadFile returns one embedded artifact by its snapshot-relative
// name, e.g. "translation_extract.json.gz",
// "profiles/SWITCH_VIRTUAL_RECEIVER.json.gz",
// "translation_custom/device_models_de.json" or
// "device_semantics.json".
func ReadFile(name string) ([]byte, error) {
	return files.ReadFile("data/" + name)
}

// ReadDir lists one embedded snapshot directory ("." for the root),
// so consumers can enumerate profile archives and curated overlays
// without hardcoding the file set.
func ReadDir(name string) ([]string, error) {
	path := "data"
	if name != "" && name != "." {
		path += "/" + name
	}
	entries, err := files.ReadDir(path)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		out = append(out, e.Name())
	}
	return out, nil
}

// deviceSemantics mirrors data/device_semantics.json: curated device
// classifications shared across the stack. Keys starting with "_" are
// documentation.
type deviceSemantics struct {
	DoorbellModels []string `json:"doorbell_models"`
}

var (
	semanticsOnce sync.Once
	doorbellSet   map[string]struct{}
)

// DoorbellModels returns the curated set of device models whose
// press/ring channel is a doorbell rather than a generic button.
// Consumers map the ring press of these devices onto their platform's
// doorbell semantics (e.g. Home Assistant's standard `ring` event
// type). Returns an empty set when the embedded document is missing
// or malformed — callers then fall back to generic button semantics.
func DoorbellModels() map[string]struct{} {
	semanticsOnce.Do(func() {
		doorbellSet = map[string]struct{}{}
		raw, err := ReadFile("device_semantics.json")
		if err != nil {
			return
		}
		var doc deviceSemantics
		if err := json.Unmarshal(raw, &doc); err != nil {
			return
		}
		for _, m := range doc.DoorbellModels {
			if m = strings.TrimSpace(m); m != "" {
				doorbellSet[m] = struct{}{}
			}
		}
	})
	return doorbellSet
}
