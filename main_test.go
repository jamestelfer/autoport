package main

import (
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"
)

func assertInRange(t *testing.T, scopes []string, port uint32) {
	t.Helper()
	if port < minPort || port > 65535 {
		t.Errorf("portFor(%q) = %d, want within [%d, 65535]", scopes, port, minPort)
	}
}

// Known values pin the hash. A change here moves every user's ports.
func TestPortForKnownValues(t *testing.T) {
	cases := []struct {
		scopes []string
		want   uint32
	}{
		{[]string{"foo", "bar"}, 39528},
	}
	for _, c := range cases {
		if got := portFor(c.scopes); got != c.want {
			t.Errorf("portFor(%q) = %d, want %d", c.scopes, got, c.want)
		}
	}
}

func TestPortForIsDeterministic(t *testing.T) {
	for _, scopes := range edgeCaseScopes() {
		first := portFor(scopes)
		for range 100 {
			if got := portFor(scopes); got != first {
				t.Fatalf("portFor(%q) changed between calls: %d then %d", scopes, first, got)
			}
		}
	}
}

func TestPortForEdgeCasesStayInRange(t *testing.T) {
	for _, scopes := range edgeCaseScopes() {
		assertInRange(t, scopes, portFor(scopes))
	}
}

// Random inputs of varied length and count exercise the whole hash space far
// more than hand-picked cases can.
func TestPortForRandomInputsStayInRange(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 100_000 {
		scopes := make([]string, 1+r.IntN(20))
		for i := range scopes {
			b := make([]byte, r.IntN(200))
			for j := range b {
				b[j] = byte(r.UintN(256))
			}
			scopes[i] = string(b)
		}
		assertInRange(t, scopes, portFor(scopes))
	}
}

// Different scopes spread across the range rather than clustering.
func TestPortForSpreadsDistinctInputs(t *testing.T) {
	const n = 10_000
	seen := make(map[uint32]bool, n)
	for i := range n {
		seen[portFor([]string{"project", strings.Repeat("x", i%50), string(rune('a' + i%26)), strconv.Itoa(i)})] = true
	}
	// With 55536 slots and 10000 inputs, about 9200 distinct ports are
	// expected. A badly broken mapping would give far fewer.
	if len(seen) < 8500 {
		t.Errorf("%d inputs produced only %d distinct ports", n, len(seen))
	}
}

func FuzzPortFor(f *testing.F) {
	for _, scopes := range edgeCaseScopes() {
		f.Add(strings.Join(scopes, "\x00"))
	}
	f.Fuzz(func(t *testing.T, joined string) {
		scopes := strings.Split(joined, "\x00")
		port := portFor(scopes)
		assertInRange(t, scopes, port)
		if again := portFor(scopes); again != port {
			t.Errorf("portFor(%q) changed between calls: %d then %d", scopes, port, again)
		}
	})
}

func edgeCaseScopes() [][]string {
	return [][]string{
		{""},
		{"a"},
		{"", ""},
		{"::"},
		{"a", "b"},
		{"b", "a"},
		{"日本語", "🚀"},
		{"\x00\xff"},
		{strings.Repeat("long", 100_000)},
		strings.Split(strings.Repeat("s,", 10_000), ","),
	}
}
