package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

const syncDir = ".envsync"

func cmdInit() error {
	if err := os.MkdirAll(syncDir, 0755); err != nil {
		return err
	}
	// Generate a simple symmetric key for demo
	key := make([]byte, 32)
	rand.Read(key)
	keyPath := filepath.Join(syncDir, "key")
	if err := os.WriteFile(keyPath, []byte(hex.EncodeToString(key)), 0600); err != nil {
		return err
	}
	// Create .gitignore entry
	gi := filepath.Join(syncDir, ".gitignore")
	os.WriteFile(gi, []byte("key\n"), 0644)
	fmt.Println("Initialized .envsync/ directory")
	fmt.Println("Key generated at .envsync/key (do NOT commit this)")
	return nil
}

func cmdPush(file, env string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", file, err)
	}
	keyHex, err := os.ReadFile(filepath.Join(syncDir, "key"))
	if err != nil {
		return fmt.Errorf("no key found, run envsync init first")
	}
	key, _ := hex.DecodeString(string(keyHex))
	encrypted := xorEncrypt(data, key)
	outPath := filepath.Join(syncDir, env+".enc")
	if err := os.WriteFile(outPath, encrypted, 0644); err != nil {
		return err
	}
	fmt.Printf("Encrypted %s -> %s (%d bytes)\n", file, outPath, len(encrypted))
	return nil
}

func cmdPull(env string) error {
	encPath := filepath.Join(syncDir, env+".enc")
	data, err := os.ReadFile(encPath)
	if err != nil {
		return fmt.Errorf("no encrypted env found for '%s'", env)
	}
	keyHex, err := os.ReadFile(filepath.Join(syncDir, "key"))
	if err != nil {
		return fmt.Errorf("no key found")
	}
	key, _ := hex.DecodeString(string(keyHex))
	decrypted := xorEncrypt(data, key) // XOR is symmetric
	if err := os.WriteFile(".env", decrypted, 0600); err != nil {
		return err
	}
	fmt.Printf("Decrypted %s -> .env\n", encPath)
	return nil
}

func cmdDiff(env string) error {
	local, err := os.ReadFile(".env")
	if err != nil {
		return fmt.Errorf(".env not found")
	}
	encPath := filepath.Join(syncDir, env+".enc")
	encData, err := os.ReadFile(encPath)
	if err != nil {
		return fmt.Errorf("no encrypted env for '%s'", env)
	}
	keyHex, _ := os.ReadFile(filepath.Join(syncDir, "key"))
	key, _ := hex.DecodeString(string(keyHex))
	remote := xorEncrypt(encData, key)
	if string(local) == string(remote) {
		fmt.Println("No differences")
	} else {
		fmt.Println("Local .env differs from encrypted version")
		fmt.Printf("Local: %d bytes, Remote: %d bytes\n", len(local), len(remote))
	}
	return nil
}

// Simple XOR encryption (placeholder for age in production)
func xorEncrypt(data, key []byte) []byte {
	out := make([]byte, len(data))
	for i, b := range data {
		out[i] = b ^ key[i%len(key)]
	}
	return out
}
