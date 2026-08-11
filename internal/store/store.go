package store

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Palette string

const (
	PaletteClassic      Palette = "classic"
	PalettePaper        Palette = "paper"
	PaletteHighContrast Palette = "high-contrast"
	PaletteColorblind   Palette = "colorblind"
	PalettePhosphor     Palette = "phosphor"
	PaletteClay         Palette = "clay"
	PaletteASCII        Palette = "ascii"
)

var Palettes = []Palette{
	PaletteClassic,
	PalettePaper,
	PaletteHighContrast,
	PaletteColorblind,
	PalettePhosphor,
	PaletteClay,
	PaletteASCII,
}

func (p Palette) Label() string {
	switch p {
	case PalettePaper:
		return "paper"
	case PaletteHighContrast:
		return "high contrast"
	case PaletteColorblind:
		return "colorblind"
	case PalettePhosphor:
		return "phosphor"
	case PaletteClay:
		return "clay"
	case PaletteASCII:
		return "ascii"
	default:
		return "classic"
	}
}

func (p Palette) Canonical() Palette {
	switch p {
	case "nord":
		return PalettePhosphor
	case "solarized":
		return PalettePaper
	case PalettePaper, PaletteHighContrast, PaletteColorblind, PalettePhosphor, PaletteClay, PaletteASCII:
		return p
	default:
		return PaletteClassic
	}
}

type Config struct {
	Palette    Palette `json:"palette"`
	Hard       bool    `json:"hard"`
	MaxGuesses int     `json:"max_guesses"`
}

func DefaultConfig() Config {
	return Config{Palette: PaletteClassic, MaxGuesses: 6}
}

type Stats struct {
	Played      int    `json:"played"`
	Wins        int    `json:"wins"`
	Streak      int    `json:"streak"`
	MaxStreak   int    `json:"max_streak"`
	Dist        [9]int `json:"dist"` // Dist[n] = wins in n guesses
	LastDaily   string `json:"last_daily"`
	LastGuesses int    `json:"last_guesses"`
	LastWon     bool   `json:"last_won"`
}

func (s *Stats) Record(won bool, guesses int, dailyDate string) {
	s.Played++
	if won {
		s.Wins++
		s.Streak++
		if s.Streak > s.MaxStreak {
			s.MaxStreak = s.Streak
		}
		if guesses >= 1 && guesses < len(s.Dist) {
			s.Dist[guesses]++
		}
	} else {
		s.Streak = 0
	}
	s.LastGuesses = guesses
	s.LastWon = won
	if dailyDate != "" {
		s.LastDaily = dailyDate
	}
}

func (s Stats) WinPct() int {
	if s.Played == 0 {
		return 0
	}
	return (s.Wins * 100) / s.Played
}

type Daily struct {
	Date          string   `json:"date"`
	Answer        string   `json:"answer"`
	Guesses       []string `json:"guesses"`
	Hard          bool     `json:"hard"`
	MaxGuesses    int      `json:"max_guesses"`
	HintPenalties int      `json:"hint_penalties"`
	HintLog       []string `json:"hint_log"`
	Won           bool     `json:"won"`
	Lost          bool     `json:"lost"`
}

func (d Daily) Done() bool { return d.Won || d.Lost }

type Store struct {
	Dir    string
	Config Config
	Stats  Stats
	Daily  Daily
}

func Open(dir string) (*Store, error) {
	if dir == "" {
		base, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}
		dir = filepath.Join(base, "wordle")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	st := &Store{Dir: dir, Config: DefaultConfig()}
	_ = st.loadJSON("config.json", &st.Config)
	_ = st.loadJSON("stats.json", &st.Stats)
	_ = st.loadJSON("daily.json", &st.Daily)
	if st.Config.MaxGuesses <= 0 {
		st.Config.MaxGuesses = 6
	}
	st.Config.Palette = st.Config.Palette.Canonical()
	return st, nil
}

func (st *Store) SaveConfig() error {
	return st.saveJSON("config.json", st.Config)
}

func (st *Store) SaveStats() error {
	return st.saveJSON("stats.json", st.Stats)
}

func (st *Store) SaveDaily() error {
	return st.saveJSON("daily.json", st.Daily)
}

func (st *Store) Today(day string) *Daily {
	if st.Daily.Date == day && st.Daily.Answer != "" {
		return &st.Daily
	}
	return nil
}

func (st *Store) loadJSON(name string, v any) error {
	b, err := os.ReadFile(filepath.Join(st.Dir, name))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func (st *Store) saveJSON(name string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(st.Dir, name), b, 0o644)
}
