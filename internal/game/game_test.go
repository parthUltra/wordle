package game

import (
	"strings"
	"testing"
)

func tilesOf(marks string) [5]Tile {
	var out [5]Tile
	for i, r := range marks {
		switch r {
		case 'C':
			out[i] = Correct
		case 'P':
			out[i] = Present
		case 'A':
			out[i] = Absent
		default:
			tPanic("bad mark " + string(r))
		}
	}
	return out
}

func tPanic(msg string) { panic(msg) }

func TestScoreExactMatch(t *testing.T) {
	got := Score("CRANE", "CRANE")
	want := tilesOf("CCCCC")
	if got != want {
		t.Fatalf("Score(CRANE, CRANE) = %v, want %v", got, want)
	}
}

func TestScoreDuplicateLetters(t *testing.T) {
	// ARRAY vs CRANE: R green, first A yellow, extra A/R/Y gray.
	got := Score("ARRAY", "CRANE")
	want := tilesOf("PCAAA")
	if got != want {
		t.Fatalf("Score(ARRAY, CRANE) = %v, want %v", got, want)
	}
}

func TestScoreGreenConsumesBeforeYellow(t *testing.T) {
	// ERASE vs SPEED: both E's yellow, S yellow, R/A gray.
	got := Score("ERASE", "SPEED")
	want := tilesOf("PAAPP")
	if got != want {
		t.Fatalf("Score(ERASE, SPEED) = %v, want %v", got, want)
	}
}

func TestScoreAllAbsent(t *testing.T) {
	got := Score("BLIMP", "CRANE")
	want := tilesOf("AAAAA")
	if got != want {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestHardModeRequiresGreenInPlace(t *testing.T) {
	s := newTestSession("CRANE")
	if err := s.mustGuess("TRACE"); err != nil {
		t.Fatal(err)
	}
	// T is absent, R green pos1, A green pos2, C present, E green pos4.
	err := s.SubmitWord("CRONY")
	if err == nil {
		t.Fatal("expected hard-mode error: E must stay in place")
	}
}

func TestHardModeRequiresYellowReuse(t *testing.T) {
	s := newTestSession("CRANE")
	if err := s.mustGuess("PILOT"); err != nil {
		t.Fatal(err)
	}
	// no overlap — not a useful setup. Use a yellow.
	s = newTestSession("CRANE")
	if err := s.mustGuess("BEACH"); err != nil {
		t.Fatal(err)
	}
	// A yellow, E yellow, C yellow (BEACH vs CRANE)
	err := s.SubmitWord("BLIMP")
	if err == nil {
		t.Fatal("expected hard-mode error: must reuse A/C/E")
	}
}

func TestHardModeAllowsGrayLetters(t *testing.T) {
	s := newTestSession("CRANE")
	if err := s.mustGuess("BLIMP"); err != nil {
		t.Fatal(err)
	}
	// All gray. Official hard mode still allows those letters.
	if err := s.SubmitWord("BLUNT"); err != nil {
		t.Fatalf("gray reuse should be allowed: %v", err)
	}
}

func TestHardModeAcceptsValidFollowUp(t *testing.T) {
	s := newTestSession("CRANE")
	if err := s.mustGuess("TRACE"); err != nil {
		t.Fatal(err)
	}
	if err := s.SubmitWord("CRANE"); err != nil {
		t.Fatalf("valid hard follow-up: %v", err)
	}
	if s.Status != Won {
		t.Fatalf("status = %v, want Won", s.Status)
	}
}

func TestRejectsUnknownWord(t *testing.T) {
	s := newTestSession("CRANE")
	if err := s.SubmitWord("ZZZZZ"); err == nil {
		t.Fatal("expected unknown-word error")
	}
}

func TestRejectsWrongLength(t *testing.T) {
	s := newTestSession("CRANE")
	if err := s.SubmitWord("HI"); err == nil {
		t.Fatal("expected length error")
	}
}

func TestWinAndLoss(t *testing.T) {
	s := newTestSession("CRANE")
	s.MaxGuesses = 2
	if err := s.mustGuess("TRACE"); err != nil {
		t.Fatal(err)
	}
	if s.Status != Playing {
		t.Fatalf("status = %v", s.Status)
	}
	if err := s.mustGuess("CRANE"); err != nil {
		t.Fatal(err)
	}
	if s.Status != Won {
		t.Fatalf("status = %v, want Won", s.Status)
	}

	s = newTestSession("CRANE")
	s.MaxGuesses = 1
	if err := s.mustGuess("TRACE"); err != nil {
		t.Fatal(err)
	}
	if s.Status != Lost {
		t.Fatalf("status = %v, want Lost", s.Status)
	}
}

func TestHintPrefersRarerLetter(t *testing.T) {
	s := newTestSession("CRANE")
	text, err := s.UseHint()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, "C") {
		t.Fatalf("expected rarest unused letter C in %q", text)
	}
}

func TestHintContainsUnknownLetter(t *testing.T) {
	s := newTestSession("CRANE")
	if err := s.mustGuess("BLIMP"); err != nil {
		t.Fatal(err)
	}
	text, err := s.UseHint()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.ToLower(text), "contains") {
		t.Fatalf("hint %q should say contains", text)
	}
	found := false
	for _, c := range "CRANE" {
		if strings.Contains(text, string(c)) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("hint %q should name a letter from CRANE", text)
	}
	if s.HintsUsed != 1 {
		t.Fatalf("HintsUsed = %d", s.HintsUsed)
	}
	if s.HintPenalties != 0 {
		t.Fatal("first hint should be free")
	}
}

func TestExtraHintCostsGuess(t *testing.T) {
	s := newTestSession("CRANE")
	s.MaxGuesses = 3
	if err := s.mustGuess("BLIMP"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.UseHint(); err != nil {
		t.Fatal(err)
	}
	if _, err := s.UseHint(); err != nil {
		t.Fatal(err)
	}
	if s.HintPenalties != 1 {
		t.Fatalf("penalties = %d, want 1", s.HintPenalties)
	}
	if s.GuessesUsed() != 2 {
		t.Fatalf("used = %d, want 2", s.GuessesUsed())
	}
}

func TestRemainingCandidatesShrinks(t *testing.T) {
	s := newTestSession("CRANE")
	before := s.RemainingAnswers()
	if err := s.mustGuess("CRANE"); err != nil {
		t.Fatal(err)
	}
	after := s.RemainingAnswers()
	if after != 1 {
		t.Fatalf("remaining after exact guess = %d", after)
	}
	if before <= after {
		t.Fatalf("remaining should shrink: %d -> %d", before, after)
	}
}

func TestShareGrid(t *testing.T) {
	s := newTestSession("CRANE")
	if err := s.mustGuess("BLIMP"); err != nil {
		t.Fatal(err)
	}
	if err := s.mustGuess("CRANE"); err != nil {
		t.Fatal(err)
	}
	out := s.Share("Wordle")
	if !strings.Contains(out, "⬛") || !strings.Contains(out, "🟩") {
		t.Fatalf("share missing tiles:\n%s", out)
	}
	if !strings.Contains(out, "2/6") {
		t.Fatalf("share missing score:\n%s", out)
	}
}

func TestRestoreCompletedGame(t *testing.T) {
	s := newTestSession("CRANE")
	if err := s.mustGuess("BLIMP"); err != nil {
		t.Fatal(err)
	}
	if err := s.mustGuess("CRANE"); err != nil {
		t.Fatal(err)
	}
	got := Restore(s.Allowed, s.Answers, s.Snapshot())
	if got.Status != Won || len(got.Guesses) != 2 || got.Guesses[0].Word != "BLIMP" {
		t.Fatalf("restore = status %v guesses %+v", got.Status, got.Guesses)
	}
}

func TestTypeAndBackspace(t *testing.T) {
	s := newTestSession("CRANE")
	s.Type('c')
	s.Type('r')
	s.Type('a')
	if s.Input != "CRA" {
		t.Fatalf("input = %q", s.Input)
	}
	s.Backspace()
	if s.Input != "CR" {
		t.Fatalf("input = %q", s.Input)
	}
	s.Type('1')
	if s.Input != "CR" {
		t.Fatalf("ignored non-letter: %q", s.Input)
	}
}

func newTestSession(answer string) *Session {
	allowed := map[string]struct{}{
		"CRANE": {}, "TRACE": {}, "CRONY": {}, "BEACH": {}, "BLIMP": {},
		"BLUNT": {}, "ARRAY": {}, "ERASE": {}, "SPEED": {}, "PILOT": {},
		"ZZZZZ": {}, // not included — used as invalid
	}
	delete(allowed, "ZZZZZ")
	answers := []string{"CRANE", "TRACE", "SPEED", "BEACH", "BLIMP", "ARRAY"}
	return NewSession(answer, allowed, answers, Options{MaxGuesses: 6, Hard: true})
}

func (s *Session) mustGuess(word string) error {
	s.Input = word
	return s.Submit()
}

func (s *Session) SubmitWord(word string) error {
	s.Input = word
	return s.Submit()
}
