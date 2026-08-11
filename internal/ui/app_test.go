package ui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/parthUltra/wordle/internal/game"
	"github.com/parthUltra/wordle/internal/store"
	"github.com/parthUltra/wordle/internal/words"
)

func testModel(t *testing.T) Model {
	t.Helper()
	lists, err := words.Parse("CRANE\nSLATE\n", "CRANE\nSLATE\nTRACE\nBEACH\n")
	if err != nil {
		t.Fatal(err)
	}
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return New(lists, st, "")
}

func TestViewDoesNotPanic(t *testing.T) {
	m := testModel(t)
	v := m.View()
	if v.Content == "" && !strings.Contains(viewString(v), "WORDLE") {
		// NewView may store content internally; just ensure View returns.
		_ = v
	}
	m.start("random", "CRANE")
	_ = m.View()
	m.help = true
	_ = m.View()
}

func TestOpenDailyShowsFinishedPuzzle(t *testing.T) {
	m := testModel(t)
	m.st.Daily = store.Daily{
		Date: today(), Answer: "CRANE", Guesses: []string{"TRACE", "CRANE"},
		Won: true, MaxGuesses: 6,
	}
	m.openDaily()
	if m.session == nil || m.session.Status != game.Won || len(m.session.Guesses) != 2 {
		t.Fatalf("expected restored daily, got %+v", m.session)
	}
	if m.session.Guesses[1].Word != "CRANE" {
		t.Fatalf("guesses = %+v", m.session.Guesses)
	}
}

func TestTypingHIsNotAHint(t *testing.T) {
	m := testModel(t)
	m.start("random", "CRANE")
	updated, _ := m.Update(tea.KeyPressMsg{Code: 'h', Text: "h"})
	m = updated.(Model)
	if m.session.Input != "H" {
		t.Fatalf("expected to type H, got %q (hints=%d)", m.session.Input, m.session.HintsUsed)
	}
	if m.session.HintsUsed != 0 {
		t.Fatal("letter h should not spend a hint")
	}
}

func viewString(v tea.View) string {
	return v.Content
}
