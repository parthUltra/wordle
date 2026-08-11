package words

import (
	"testing"
	"time"
)

const sampleAnswers = "CRANE\nSLATE\nTRACE\n"
const sampleAllowed = "CRANE\nSLATE\nTRACE\nZZZZZ\nPILOT\n"

func TestParseAndPick(t *testing.T) {
	l, err := Parse(sampleAnswers, sampleAllowed)
	if err != nil {
		t.Fatal(err)
	}
	if len(l.Answers) != 3 {
		t.Fatalf("answers = %d", len(l.Answers))
	}
	if _, ok := l.Allowed["PILOT"]; !ok {
		t.Fatal("PILOT should be allowed")
	}
	if l.BySeed(1) != l.BySeed(1) {
		t.Fatal("seed not stable")
	}
	if l.BySeed(2) == "" {
		t.Fatal("empty pick")
	}
}

func TestDailyStable(t *testing.T) {
	l, err := Parse(sampleAnswers, sampleAllowed)
	if err != nil {
		t.Fatal(err)
	}
	d1 := time.Date(2026, 8, 11, 9, 0, 0, 0, time.Local)
	d2 := time.Date(2026, 8, 11, 21, 0, 0, 0, time.Local)
	d3 := time.Date(2026, 8, 12, 9, 0, 0, 0, time.Local)
	if l.Daily(d1) != l.Daily(d2) {
		t.Fatal("same local day should match")
	}
	if _, ok := l.Allowed[l.Daily(d3)]; !ok {
		t.Fatal("daily word not in list")
	}
}
