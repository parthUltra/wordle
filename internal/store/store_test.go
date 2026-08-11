package store

import "testing"

func TestPaletteCanonical(t *testing.T) {
	if Palette("nord").Canonical() != PalettePhosphor {
		t.Fatal("nord should map to phosphor")
	}
	if Palette("solarized").Canonical() != PalettePaper {
		t.Fatal("solarized should map to paper")
	}
	if Palette("").Canonical() != PaletteClassic {
		t.Fatal("empty should be classic")
	}
}

func TestRecordAndPersist(t *testing.T) {
	dir := t.TempDir()
	st, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Config.Hard = true
	st.Config.Palette = PaletteClay
	if err := st.SaveConfig(); err != nil {
		t.Fatal(err)
	}
	st.Stats.Record(true, 3, "2026-08-11")
	st.Stats.Record(false, 6, "")
	if err := st.SaveStats(); err != nil {
		t.Fatal(err)
	}

	st2, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !st2.Config.Hard || st2.Config.Palette != PaletteClay {
		t.Fatalf("config = %+v", st2.Config)
	}
	if st2.Stats.Played != 2 || st2.Stats.Wins != 1 || st2.Stats.Streak != 0 {
		t.Fatalf("stats = %+v", st2.Stats)
	}
	if st2.Stats.Dist[3] != 1 {
		t.Fatalf("dist = %v", st2.Stats.Dist)
	}
	if st2.Stats.WinPct() != 50 {
		t.Fatalf("win%% = %d", st2.Stats.WinPct())
	}
	if st2.Stats.LastWon || st2.Stats.LastGuesses != 6 {
		t.Fatalf("last = won:%v guesses:%d", st2.Stats.LastWon, st2.Stats.LastGuesses)
	}

	st2.Daily = Daily{Date: "2026-08-11", Answer: "CRANE", Guesses: []string{"BLIMP", "CRANE"}, Won: true}
	if err := st2.SaveDaily(); err != nil {
		t.Fatal(err)
	}
	st3, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if st3.Today("2026-08-11") == nil || !st3.Daily.Done() || st3.Daily.Answer != "CRANE" {
		t.Fatalf("daily = %+v", st3.Daily)
	}
	if st3.Today("2026-08-12") != nil {
		t.Fatal("wrong day should miss")
	}
}
