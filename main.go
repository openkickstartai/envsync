package main

import (
	"fmt"
	"os"
)

var version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: envsync <init|push|pull|diff|add-key|version>")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "init":
		if err := cmdInit(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "push":
		env := flagEnv()
		if len(os.Args) < 3 {
			fmt.Println("usage: envsync push <file> [--env <name>]")
			os.Exit(1)
		}
		if err := cmdPush(os.Args[2], env); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "pull":
		env := flagEnv()
		if err := cmdPull(env); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "diff":
		env := flagEnv()
		if err := cmdDiff(env); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "version":
		fmt.Printf("envsync %s\n", version)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func flagEnv() string {
	for i, a := range os.Args {
		if a == "--env" && i+1 < len(os.Args) {
			return os.Args[i+1]
		}
	}
	return "default"
}
