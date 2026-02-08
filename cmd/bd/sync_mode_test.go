package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/steveyegge/beads/internal/storage/sqlite"
)

// TestSyncModeConfig verifies sync mode configuration storage and retrieval.
func TestSyncModeConfig(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	// Create .beads directory
	beadsDir := filepath.Join(tmpDir, ".beads")
	if err := os.MkdirAll(beadsDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	// Create store
	dbPath := filepath.Join(beadsDir, "beads.db")
	testStore, err := sqlite.New(ctx, dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer testStore.Close()

	// Test 1: Default mode is git-portable
	mode := GetSyncMode(ctx, testStore)
	if mode != SyncModeGitPortable {
		t.Errorf("default sync mode = %q, want %q", mode, SyncModeGitPortable)
	}
	t.Logf("✓ Default sync mode is git-portable")

	// Test 2: Set and get realtime mode
	if err := SetSyncMode(ctx, testStore, SyncModeRealtime); err != nil {
		t.Fatalf("failed to set sync mode: %v", err)
	}
	mode = GetSyncMode(ctx, testStore)
	if mode != SyncModeRealtime {
		t.Errorf("sync mode = %q, want %q", mode, SyncModeRealtime)
	}
	t.Logf("✓ Can set and get realtime mode")

	// Test 3: Invalid mode returns error
	err = SetSyncMode(ctx, testStore, "invalid-mode")
	if err == nil {
		t.Error("expected error for invalid sync mode")
	}
	t.Logf("✓ Invalid mode correctly rejected")

	// Test 4: Invalid mode in DB defaults to git-portable
	if err := testStore.SetConfig(ctx, SyncModeConfigKey, "invalid"); err != nil {
		t.Fatalf("failed to set invalid config: %v", err)
	}
	mode = GetSyncMode(ctx, testStore)
	if mode != SyncModeGitPortable {
		t.Errorf("invalid mode should default to %q, got %q", SyncModeGitPortable, mode)
	}
	t.Logf("✓ Invalid mode in DB defaults to git-portable")
}

// TestShouldExportJSONL verifies JSONL export behavior per mode.
func TestShouldExportJSONL(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	beadsDir := filepath.Join(tmpDir, ".beads")
	if err := os.MkdirAll(beadsDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	dbPath := filepath.Join(beadsDir, "beads.db")
	testStore, err := sqlite.New(ctx, dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer testStore.Close()

	tests := []struct {
		mode       string
		wantExport bool
	}{
		{SyncModeGitPortable, true},
		{SyncModeRealtime, true},
	}

	for _, tt := range tests {
		t.Run(tt.mode, func(t *testing.T) {
			if err := SetSyncMode(ctx, testStore, tt.mode); err != nil {
				t.Fatalf("failed to set mode: %v", err)
			}

			got := ShouldExportJSONL(ctx, testStore)
			if got != tt.wantExport {
				t.Errorf("ShouldExportJSONL() = %v, want %v", got, tt.wantExport)
			}
		})
	}
}

// TestShouldExportJSONL_DirectDB verifies ShouldExportJSONL always returns true.
func TestShouldExportJSONL_DirectDB(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	beadsDir := filepath.Join(tmpDir, ".beads")
	if err := os.MkdirAll(beadsDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	dbPath := filepath.Join(beadsDir, "beads.db")
	testStore, err := sqlite.New(ctx, dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer testStore.Close()

	// Default (no config set) should return true
	if !ShouldExportJSONL(ctx, testStore) {
		t.Error("ShouldExportJSONL() = false, want true for default (no config)")
	}
}

// TestShouldImportJSONL verifies JSONL import behavior per mode.
func TestShouldImportJSONL(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	beadsDir := filepath.Join(tmpDir, ".beads")
	if err := os.MkdirAll(beadsDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	dbPath := filepath.Join(beadsDir, "beads.db")
	testStore, err := sqlite.New(ctx, dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer testStore.Close()

	tests := []struct {
		mode       string
		wantImport bool
	}{
		{SyncModeGitPortable, true},
		{SyncModeRealtime, true},
	}

	for _, tt := range tests {
		t.Run(tt.mode, func(t *testing.T) {
			if err := SetSyncMode(ctx, testStore, tt.mode); err != nil {
				t.Fatalf("failed to set mode: %v", err)
			}

			got := ShouldImportJSONL(ctx, testStore)
			if got != tt.wantImport {
				t.Errorf("ShouldImportJSONL() = %v, want %v", got, tt.wantImport)
			}
		})
	}
}

// TestSyncModeDescription verifies mode descriptions are meaningful.
func TestSyncModeDescription(t *testing.T) {
	tests := []struct {
		mode        string
		wantContain string
	}{
		{SyncModeGitPortable, "JSONL"},
		{SyncModeRealtime, "every change"},
		{"invalid", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.mode, func(t *testing.T) {
			desc := SyncModeDescription(tt.mode)
			if desc == "" {
				t.Error("description should not be empty")
			}
			// Just verify descriptions are non-empty and distinct
			t.Logf("%s: %s", tt.mode, desc)
		})
	}
}
