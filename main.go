package main

import (
	"fmt"
	"hash/fnv"
	"os"
	"strings"
)

// Build metadata, injected by GoReleaser via
// -ldflags "-X main.version=… -X main.commit=… -X main.date=…".
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Printf("autoport %s (%s, %s)\n", version, commit, date)
		return
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <scope-1> <scope-2> ... <scope-n>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\t Assigns a random port based on the given scopes that is unlikely to collide\n")
		os.Exit(1)
	}

	scopes := os.Args[1:]
	port := portFor(scopes)

	fmt.Fprintf(os.Stderr, "autoport(%s): %d\n", strings.Join(scopes, "::"), port)
	fmt.Println(port)
}

const (
	minPort   = 10000
	portRange = 65536 - minPort
)

// portFor maps scopes to a stable port in [minPort, 65535]. The scopes are
// joined with "::" and hashed with FNV-1a, so the same scopes always give the
// same port.
func portFor(scopes []string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(strings.Join(scopes, "::")))
	return minPort + h.Sum32()%portRange
}
