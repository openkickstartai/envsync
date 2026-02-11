package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCreatesDir(t *testing.T) {
	tmp := t.TempDir()
	old, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(old)

	if err := cmdInit(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(syncDir, "key")); err != nil {
		t.Error("key file not created")
	}
}

func TestPushPullRoundtrip(t *testing.T) {
	tmp := t.TempDir()
	old, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(old)

	cmdInit()
	os.WriteFile(".env", []byte("SECRET=hello\nDB=postgres://localhost\n"), 0644)

	if err := cmdPush(".env", "test"); err != nil {
		t.Fatal(err)
	}
	os.Remove(".env")

	if err := cmdPull("test"); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(".env")
	if string(data) != "SECRET=hello\nDB=postgres://localhost\n" {
		t.Errorf("roundtrip failed: got %q", string(data))
	}
}

func TestXorEncrypt(t *testing.T) {
	key := []byte("secretkey1234567")
	plain := []byte("hello world")
	enc := xorEncrypt(plain, key)
	if string(enc) == string(plain) {
		t.Error("encryption should change data")
	}
	dec := xorEncrypt(enc, key)
	if string(dec) != string(plain) {
		t.Error("decryption should restore original")
	}
}
