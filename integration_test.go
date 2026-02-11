package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFullRoundtrip(t *testing.T) {
	tmp := t.TempDir()
	old, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(old)

	// Init
	if err := cmdInit(); err != nil { t.Fatal(err) }

	// Create .env with various content
	envContent := "DB_HOST=localhost\nDB_PASS=s3cr3t!@#$\nEMPTY=\nSPACES=hello world\n"
	os.WriteFile(".env", []byte(envContent), 0644)

	// Push
	if err := cmdPush(".env", "test"); err != nil { t.Fatal(err) }

	// Verify encrypted file exists
	encPath := filepath.Join(syncDir, "test.enc")
	if _, err := os.Stat(encPath); err != nil { t.Fatal("encrypted file not created") }

	// Remove original
	os.Remove(".env")

	// Pull
	if err := cmdPull("test"); err != nil { t.Fatal(err) }

	// Verify roundtrip
	data, _ := os.ReadFile(".env")
	if string(data) != envContent {
		t.Errorf("roundtrip failed:\ngot:  %q\nwant: %q", string(data), envContent)
	}
}

func TestPushMissingKey(t *testing.T) {
	tmp := t.TempDir()
	old, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(old)

	os.WriteFile(".env", []byte("X=1"), 0644)
	err := cmdPush(".env", "test")
	if err == nil { t.Error("should fail without init") }
}

func TestDiffIdentical(t *testing.T) {
	tmp := t.TempDir()
	old, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(old)

	cmdInit()
	os.WriteFile(".env", []byte("A=1\n"), 0644)
	cmdPush(".env", "test")
	// Diff should show no differences
	err := cmdDiff("test")
	if err != nil { t.Error(err) }
}
