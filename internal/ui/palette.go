package ui

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/parthUltra/wordle/internal/game"
	"github.com/parthUltra/wordle/internal/store"
)

// NYT Wordle dark-mode tokens:
// https://www.nytimes.com/games/wordle (dark theme CSS)
// empty border #3a3a3c, typed #565758, correct #538d4e, present #b59f3b

type theme struct {
	page      string
	empty     string
	fg        string
	muted     string
	line      string
	typed     string
	correct   string
	present   string
	absent    string
	keyUnused string
	keyFg     string
	onTile    string
	ascii     bool
}

func themeFor(p store.Palette) theme {
	switch p.Canonical() {
	case store.PalettePaper:
		// NYT Wordle light
		return theme{
			page: "#ffffff", empty: "#ffffff", fg: "#1a1a1b", muted: "#787c7e",
			line: "#d3d6da", typed: "#878a8c",
			correct: "#6aaa64", present: "#c9b458", absent: "#787c7e",
			keyUnused: "#d3d6da", keyFg: "#1a1a1b", onTile: "#ffffff",
		}
	case store.PaletteHighContrast:
		return theme{
			page: "#000000", empty: "#000000", fg: "#ffffff", muted: "#c0c0c0",
			line: "#ffffff", typed: "#ffff00",
			correct: "#f5793a", present: "#85c0f9", absent: "#3a3a3c",
			keyUnused: "#5c5c5c", keyFg: "#ffffff", onTile: "#000000",
		}
	case store.PaletteColorblind:
		return theme{
			page: "#121213", empty: "#121213", fg: "#ffffff", muted: "#818384",
			line: "#3a3a3c", typed: "#85c0f9",
			correct: "#2e86ab", present: "#f6ae2d", absent: "#3a3a3c",
			keyUnused: "#565758", keyFg: "#ffffff", onTile: "#ffffff",
		}
	case store.PalettePhosphor:
		return theme{
			page: "#07110c", empty: "#07110c", fg: "#d4f5dc", muted: "#6b9a78",
			line: "#1c3a28", typed: "#3ddc84",
			correct: "#2fbf71", present: "#e6a817", absent: "#14301f",
			keyUnused: "#0f2418", keyFg: "#d4f5dc", onTile: "#07110c",
		}
	case store.PaletteClay:
		return theme{
			page: "#1a1210", empty: "#1a1210", fg: "#f4e4d4", muted: "#a08070",
			line: "#4a3028", typed: "#d4a574",
			correct: "#c14a2a", present: "#d4a017", absent: "#2e201c",
			keyUnused: "#2a1c18", keyFg: "#f4e4d4", onTile: "#fff5eb",
		}
	case store.PaletteASCII:
		return theme{ascii: true}
	default:
		// NYT Wordle dark
		return theme{
			page: "#121213", empty: "#121213", fg: "#ffffff", muted: "#818384",
			line: "#3a3a3c", typed: "#565758",
			correct: "#538d4e", present: "#b59f3b", absent: "#3a3a3c",
			keyUnused: "#818384", keyFg: "#ffffff", onTile: "#ffffff",
		}
	}
}

func (th theme) col(hex string) color.Color {
	return lipgloss.Color(hex)
}

func (th theme) tile(t game.Tile, letter string, filled, compact bool) string {
	if letter == "" {
		letter = " "
	}
	if th.ascii {
		return asciiTile(t, letter, filled, compact)
	}

	bg, fg, border := th.empty, th.onTile, th.line
	switch t {
	case game.Correct:
		bg, border, fg = th.correct, th.correct, th.onTile
	case game.Present:
		bg, border, fg = th.present, th.present, th.onTile
	case game.Absent:
		bg, border, fg = th.absent, th.absent, th.onTile
	default:
		if filled {
			border, fg = th.typed, th.fg
		} else {
			fg = th.empty
		}
	}

	s := lipgloss.NewStyle().
		Foreground(th.col(fg)).
		Background(th.col(bg)).
		Bold(true).
		Width(5).
		Align(lipgloss.Center).
		BorderForeground(th.col(border)).
		MarginRight(1)

	if compact {
		s = s.Border(lipgloss.NormalBorder(), false, true, false, true) // left+right only
	} else {
		s = s.Border(lipgloss.ThickBorder()).Height(1)
	}
	return s.Render(letter)
}

func asciiTile(t game.Tile, letter string, filled, compact bool) string {
	if compact {
		switch {
		case t == game.Correct:
			return "[" + letter + "]"
		case t == game.Present:
			return "(" + letter + ")"
		case filled:
			return "|" + letter + "|"
		default:
			return "| |"
		}
	}
	inner := " " + letter + " "
	switch {
	case t == game.Correct:
		return "+---+\n|" + inner + "|\n+---+"
	case t == game.Present:
		return "+---+\n(" + inner + ")\n+---+"
	case filled:
		return "+---+\n|" + inner + "|\n+---+"
	default:
		return "+---+\n|   |\n+---+"
	}
}

func (th theme) key(label string, t game.Tile, wide bool) string {
	if th.ascii {
		switch t {
		case game.Correct:
			return "[" + label + "]"
		case game.Present:
			return "(" + label + ")"
		default:
			return " " + label + " "
		}
	}
	bg, fg := th.keyUnused, th.keyFg
	switch t {
	case game.Correct:
		bg, fg = th.correct, th.onTile
	case game.Present:
		bg, fg = th.present, th.onTile
	case game.Absent:
		bg, fg = th.absent, th.onTile
	}
	pad := 1
	if wide {
		pad = 2
	}
	return lipgloss.NewStyle().
		Foreground(th.col(fg)).
		Background(th.col(bg)).
		Bold(true).
		Padding(0, pad).
		MarginRight(1).
		Align(lipgloss.Center).
		Render(label)
}

func (th theme) text(s string) string {
	if th.ascii {
		return s
	}
	return lipgloss.NewStyle().Foreground(th.col(th.fg)).Render(s)
}

func (th theme) dim(s string) string {
	if th.ascii {
		return s
	}
	return lipgloss.NewStyle().Foreground(th.col(th.muted)).Render(s)
}

func (th theme) title(s string) string {
	if th.ascii {
		return s
	}
	return lipgloss.NewStyle().Bold(true).Foreground(th.col(th.fg)).Render(s)
}

func (th theme) statCell(value, label string) string {
	num := th.title(value)
	if !th.ascii {
		num = lipgloss.NewStyle().
			Bold(true).
			Foreground(th.col(th.fg)).
			Align(lipgloss.Center).
			Width(10).
			Render(value)
	}
	cap := th.dim(label)
	if !th.ascii {
		cap = lipgloss.NewStyle().
			Foreground(th.col(th.muted)).
			Align(lipgloss.Center).
			Width(10).
			Render(label)
	}
	return lipgloss.JoinVertical(lipgloss.Center, num, cap)
}

func (th theme) distBar(guess, count, peak int, highlight bool) string {
	label := fmt.Sprintf("%d", guess)
	width := 1
	if peak > 0 && count > 0 {
		width = 1 + (count*18)/peak
	}
	bar := strings.Repeat("█", width)
	text := fmt.Sprintf(" %d ", count)
	if len(text) < width {
		bar = strings.Repeat("█", width-len(text)) + text
	} else {
		bar = text
	}
	if th.ascii {
		return fmt.Sprintf("%s  %s", label, strings.TrimSpace(bar))
	}
	bg, fg := th.absent, th.fg
	if highlight {
		bg, fg = th.correct, th.onTile
	}
	body := lipgloss.NewStyle().
		Foreground(th.col(fg)).
		Background(th.col(bg)).
		Bold(true).
		Render(bar)
	return lipgloss.JoinHorizontal(lipgloss.Center,
		lipgloss.NewStyle().Foreground(th.col(th.fg)).Width(2).Render(label),
		body,
	)
}

func (th theme) header(subtitle string, width int) string {
	name := "Wordle"
	if !th.ascii {
		name = lipgloss.NewStyle().Bold(true).Foreground(th.col(th.fg)).Render("Wordle")
	}
	n := 28
	if width > 28 {
		n = width
		if n > 56 {
			n = 56
		}
	}
	return name + "\n" + th.dim(subtitle) + "\n" + th.dim(strings.Repeat("─", n))
}
