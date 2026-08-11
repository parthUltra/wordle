package main

import "testing"

func TestParseCmd(t *testing.T) {
	got, err := parseCmd(nil)
	if err != nil || got != "" {
		t.Fatalf("empty args: %q %v", got, err)
	}
	for _, cmd := range []string{"daily", "random", "stats", "help"} {
		got, err = parseCmd([]string{cmd})
		if err != nil || got != cmd {
			t.Fatalf("%s: %q %v", cmd, got, err)
		}
	}
	got, err = parseCmd([]string{"--help"})
	if err != nil || got != "help" {
		t.Fatalf("help flag: %q %v", got, err)
	}
	if _, err = parseCmd([]string{"nope"}); err == nil {
		t.Fatal("expected error")
	}
}
