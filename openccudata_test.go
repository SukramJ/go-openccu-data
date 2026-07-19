// SPDX-License-Identifier: MIT
// Copyright (C) 2026 go-openccu-data authors.

package openccudata

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"testing"
)

// TestSnapshotCarriesCoreArchives pins the artifact set consumers rely
// on: the two extracts, the curated overlays, the profile archives and
// the device-semantics document.
func TestSnapshotCarriesCoreArchives(t *testing.T) {
	for _, name := range []string{
		"translation_extract.json.gz",
		"easymode_extract.json.gz",
		"device_semantics.json",
		"profiles/_receiver_type_aliases.json",
	} {
		if _, err := ReadFile(name); err != nil {
			t.Errorf("ReadFile(%q): %v", name, err)
		}
	}
	profiles, err := ReadDir("profiles")
	if err != nil || len(profiles) < 2 {
		t.Fatalf("ReadDir(profiles) = %v, err=%v — want several archives", profiles, err)
	}
	custom, err := ReadDir("translation_custom")
	if err != nil || len(custom) == 0 {
		t.Fatalf("ReadDir(translation_custom) = %v, err=%v — want curated overlays", custom, err)
	}
}

// TestTranslationExtractDecodes proves the gz archive is a valid
// gzip-compressed JSON document.
func TestTranslationExtractDecodes(t *testing.T) {
	raw, err := ReadFile("translation_extract.json.gz")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("gzip: %v", err)
	}
	blob, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var doc map[string]any
	if err := json.Unmarshal(blob, &doc); err != nil {
		t.Fatalf("json: %v", err)
	}
	if len(doc) == 0 {
		t.Fatal("translation extract decoded empty")
	}
}

// TestDoorbellModelsCuratedSet pins the curated doorbell
// classification.
func TestDoorbellModelsCuratedSet(t *testing.T) {
	models := DoorbellModels()
	for _, want := range []string{"HM-Sen-DB-PCB", "HmIP-DBB", "HmIP-DSD-PCB"} {
		if _, ok := models[want]; !ok {
			t.Errorf("DoorbellModels missing %s", want)
		}
	}
	if _, ok := models["HmIP-PS"]; ok {
		t.Error("HmIP-PS must not classify as doorbell")
	}
}

// TestSnapshotVersionStamped guards the regenerate script's stamp.
func TestSnapshotVersionStamped(t *testing.T) {
	if SnapshotVersion == "" {
		t.Fatal("SnapshotVersion empty")
	}
}
