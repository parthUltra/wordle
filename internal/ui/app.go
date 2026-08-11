package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/parthUltra/wordle/internal/game"
	"github.com/parthUltra/wordle/internal/store"
	"github.com/parthUltra/wordle/internal/words"
)

type screen int

const (
	screenMenu screen = iota
	screenPlay
	screenSettings
	screenStats
	screenSeed
)

var menuItems = []string{"daily", "random", "seeded", "settings", "stats", "quit"}

type Model struct {
	lists       *words.Lists
	st          *store.Store
	screen      screen
	menuIdx     int
	settingsIdx int
	session     *game.Session
	mode        string
	seedInput   string
	lastAnswer  string
	width       int
	height      int
	recorded    bool
	help        bool
}

func New(lists *words.Lists, st *store.Store, cmd string) Model {
	m := Model{lists: lists, st: st, width: 80, height: 24}
	switch cmd {
	case "daily":
		m.openDaily()
	case "random":
		m.start("random", lists.Random(0))
	case "stats":
		m.screen = screenStats
	}
	return m
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if m.help {
			m.help = false
			return m, nil
		}
		switch m.screen {
		case screenMenu:
			return m.updateMenu(msg)
		case screenPlay:
			return m.updatePlay(msg)
		case screenSettings:
			return m.updateSettings(msg)
		case screenStats:
			if msg.String() == "esc" || msg.String() == "q" || msg.String() == "enter" {
				m.screen = screenMenu
			}
			return m, nil
		case screenSeed:
			return m.updateSeed(msg)
		}
	}
	return m, nil
}

func (m Model) updateMenu(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.menuIdx > 0 {
			m.menuIdx--
		}
	case "down", "j":
		if m.menuIdx < len(menuItems)-1 {
			m.menuIdx++
		}
	case "q":
		return m, tea.Quit
	case "enter":
		switch menuItems[m.menuIdx] {
		case "daily":
			m.openDaily()
		case "random":
			m.start("random", m.lists.Random(0))
		case "seeded":
			m.seedInput = ""
			m.screen = screenSeed
		case "settings":
			m.screen = screenSettings
		case "stats":
			m.screen = screenStats
		case "quit":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) updateSeed(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.screen = screenMenu
	case "enter":
		n, _ := strconv.ParseUint(m.seedInput, 10, 64)
		m.start("seed "+m.seedInput, m.lists.BySeed(n))
	case "backspace":
		if m.seedInput != "" {
			m.seedInput = m.seedInput[:len(m.seedInput)-1]
		}
	default:
		if len(msg.Text) == 1 && msg.Text[0] >= '0' && msg.Text[0] <= '9' {
			m.seedInput += msg.Text
		}
	}
	return m, nil
}

func (m Model) updateSettings(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		_ = m.st.SaveConfig()
		m.screen = screenMenu
	case "up", "k":
		if m.settingsIdx > 0 {
			m.settingsIdx--
		}
	case "down", "j":
		if m.settingsIdx < 3 {
			m.settingsIdx++
		}
	case "enter", "left", "right", "h", "l", "space":
		switch m.settingsIdx {
		case 0:
			m.cyclePalette(msg.String() == "left" || msg.String() == "h")
		case 1:
			m.st.Config.Hard = !m.st.Config.Hard
		case 2:
			delta := 1
			if msg.String() == "left" || msg.String() == "h" {
				delta = -1
			}
			m.st.Config.MaxGuesses += delta
			if m.st.Config.MaxGuesses < 4 {
				m.st.Config.MaxGuesses = 4
			}
			if m.st.Config.MaxGuesses > 8 {
				m.st.Config.MaxGuesses = 8
			}
		case 3:
			_ = m.st.SaveConfig()
			m.screen = screenMenu
		}
	}
	return m, nil
}

func (m *Model) cyclePalette(back bool) {
	cur := 0
	for i, p := range store.Palettes {
		if p == m.st.Config.Palette {
			cur = i
			break
		}
	}
	if back {
		cur = (cur - 1 + len(store.Palettes)) % len(store.Palettes)
	} else {
		cur = (cur + 1) % len(store.Palettes)
	}
	m.st.Config.Palette = store.Palettes[cur]
}

func (m Model) updatePlay(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if m.session == nil {
		m.screen = screenMenu
		return m, nil
	}
	key := msg.String()
	if m.session.Status != game.Playing {
		switch key {
		case "n":
			if m.mode == "daily" {
				m.session.Message = "Daily is done — come back tomorrow"
				break
			}
			m.start("random", m.lists.Random(0))
		case "r":
			if m.mode == "daily" {
				break
			}
			m.start(m.mode, m.lastAnswer)
		case "c":
			return m, tea.SetClipboard(m.session.Share(m.shareTitle()))
		case "esc":
			m.screen = screenMenu
		}
		return m, nil
	}
	switch key {
	case "esc":
		m.persistDaily()
		m.screen = screenMenu
	case "?":
		m.help = true
	case "tab":
		_, _ = m.session.UseHint()
		m.afterMove()
	case "ctrl+g":
		m.session.GiveUp()
		m.afterMove()
	case "enter":
		_ = m.session.Submit()
		m.afterMove()
	case "backspace":
		m.session.Backspace()
	default:
		if len(msg.Text) == 1 {
			m.session.Type(rune(msg.Text[0]))
		}
	}
	return m, nil
}

func (m *Model) openDaily() {
	if d := m.st.Today(today()); d != nil {
		m.session = game.Restore(m.lists.Allowed, m.lists.Answers, game.Snapshot{
			Answer:        d.Answer,
			Guesses:       d.Guesses,
			Hard:          d.Hard,
			MaxGuesses:    d.MaxGuesses,
			HintPenalties: d.HintPenalties,
			HintLog:       d.HintLog,
			Won:           d.Won,
			Lost:          d.Lost,
		})
		m.mode = "daily"
		m.lastAnswer = d.Answer
		m.recorded = d.Done()
		m.screen = screenPlay
		m.help = false
		if d.Done() {
			m.session.Message = "Today's puzzle"
		}
		return
	}
	m.start("daily", m.lists.Daily(time.Now()))
	m.persistDaily()
}

func (m *Model) start(mode, answer string) {
	m.mode = mode
	m.lastAnswer = answer
	m.recorded = false
	m.session = game.NewSession(answer, m.lists.Allowed, m.lists.Answers, game.Options{
		MaxGuesses: m.st.Config.MaxGuesses,
		Hard:       m.st.Config.Hard,
	})
	m.screen = screenPlay
	m.help = false
}

func (m *Model) afterMove() {
	m.persistDaily()
	m.maybeRecord()
}

func (m *Model) persistDaily() {
	if m.session == nil || m.mode != "daily" {
		return
	}
	snap := m.session.Snapshot()
	m.st.Daily = store.Daily{
		Date:          today(),
		Answer:        snap.Answer,
		Guesses:       snap.Guesses,
		Hard:          snap.Hard,
		MaxGuesses:    snap.MaxGuesses,
		HintPenalties: snap.HintPenalties,
		HintLog:       snap.HintLog,
		Won:           snap.Won,
		Lost:          snap.Lost,
	}
	_ = m.st.SaveDaily()
}

func (m *Model) maybeRecord() {
	if m.session == nil || m.recorded || m.session.Status == game.Playing {
		return
	}
	daily := ""
	if m.mode == "daily" {
		daily = today()
	}
	m.st.Stats.Record(m.session.Status == game.Won, len(m.session.Guesses), daily)
	_ = m.st.SaveStats()
	m.recorded = true
}

func today() string { return time.Now().Format("2006-01-02") }

func (m Model) shareTitle() string {
	if m.mode == "daily" {
		return "Wordle " + today()
	}
	return "Wordle"
}

func (m Model) View() tea.View {
	th := themeFor(m.st.Config.Palette)
	var body string
	switch {
	case m.help:
		body = m.viewHelp(th)
	case m.screen == screenMenu:
		body = m.viewMenu(th)
	case m.screen == screenPlay:
		body = m.viewPlay(th)
	case m.screen == screenSettings:
		body = m.viewSettings(th)
	case m.screen == screenStats:
		body = m.viewStats(th)
	case m.screen == screenSeed:
		body = m.viewSeed(th)
	}
	if m.width > 0 && m.height > 0 {
		opts := []lipgloss.WhitespaceOption{}
		if !th.ascii && th.page != "" {
			opts = append(opts, lipgloss.WithWhitespaceStyle(
				lipgloss.NewStyle().Background(th.col(th.page)),
			))
		}
		body = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, body, opts...)
	}
	v := tea.NewView(body)
	v.AltScreen = true
	v.WindowTitle = "Wordle"
	if !th.ascii && th.page != "" {
		v.BackgroundColor = th.col(th.page)
		v.ForegroundColor = th.col(th.fg)
	}
	return v
}

func (m Model) compactTiles() bool {
	rows := 6
	if m.session != nil {
		rows = m.session.MaxGuesses
	}
	// 3-line boxes need ~ header 3 + rows*3 + keyboard 5 + footer 3
	return m.height > 0 && m.height < rows*3+12
}

func (m Model) viewMenu(th theme) string {
	var b strings.Builder
	b.WriteString(th.header(fmt.Sprintf("streak %d  ·  %d%% won  ·  %s",
		m.st.Stats.Streak, m.st.Stats.WinPct(), m.st.Config.Palette.Label()), 36) + "\n\n")
	for i, item := range menuItems {
		label := strings.ToUpper(item)
		if i == m.menuIdx {
			b.WriteString(th.key("  "+label+"  ", game.Correct, true) + "\n\n")
		} else {
			b.WriteString(th.key("  "+label+"  ", game.Empty, true) + "\n\n")
		}
	}
	b.WriteString(th.dim("↑↓ select   enter play"))
	return b.String()
}

func praise(n int) string {
	switch n {
	case 1:
		return "Genius"
	case 2:
		return "Magnificent"
	case 3:
		return "Impressive"
	case 4:
		return "Splendid"
	case 5:
		return "Great"
	default:
		return "Phew"
	}
}

func (m Model) viewSeed(th theme) string {
	return th.header("seeded puzzle", 36) + "\n\n" +
		th.text("number  "+m.seedInput+"_") + "\n\n" +
		th.dim("enter play   esc back")
}

func (m Model) viewSettings(th theme) string {
	rows := []string{
		fmt.Sprintf("palette     < %s >", m.st.Config.Palette.Label()),
		fmt.Sprintf("hard mode   %s", onOff(m.st.Config.Hard)),
		fmt.Sprintf("guesses     %d", m.st.Config.MaxGuesses),
		"back",
	}
	var b strings.Builder
	b.WriteString(th.header("settings", 36) + "\n\n")
	for i, row := range rows {
		cur := "  "
		if i == m.settingsIdx {
			cur = "> "
		}
		line := cur + row
		if i == m.settingsIdx {
			b.WriteString(th.title(line) + "\n")
		} else {
			b.WriteString(th.text(line) + "\n")
		}
	}
	b.WriteString("\n" + th.dim("← → cycle · enter toggle · esc save"))
	return b.String()
}

func onOff(v bool) string {
	if v {
		return "on"
	}
	return "off"
}

func (m Model) viewStats(th theme) string {
	s := m.st.Stats
	var b strings.Builder
	b.WriteString(th.header("statistics", 40) + "\n\n")

	cells := []string{
		th.statCell(fmt.Sprintf("%d", s.Played), "Played"),
		th.statCell(fmt.Sprintf("%d", s.WinPct()), "Win %"),
		th.statCell(fmt.Sprintf("%d", s.Streak), "Streak"),
		th.statCell(fmt.Sprintf("%d", s.MaxStreak), "Max"),
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, cells...) + "\n\n")
	b.WriteString(th.title("Guess distribution") + "\n\n")

	if s.Played == 0 {
		b.WriteString(th.dim("No games yet. Play a daily or random puzzle.") + "\n")
	} else {
		limit := m.st.Config.MaxGuesses
		if limit < 6 {
			limit = 6
		}
		if limit > 8 {
			limit = 8
		}
		peak := 1
		for i := 1; i <= limit; i++ {
			if s.Dist[i] > peak {
				peak = s.Dist[i]
			}
		}
		for i := 1; i <= limit; i++ {
			hit := s.LastWon && s.LastGuesses == i
			b.WriteString(th.distBar(i, s.Dist[i], peak, hit) + "\n")
		}
	}
	b.WriteString("\n" + th.dim("esc back"))
	return b.String()
}

func (m Model) viewHelp(th theme) string {
	return th.header("how to play", 36) + "\n\n" + th.text(strings.TrimSpace(`
a–z     type a letter
enter   submit guess
bksp    delete
tab     hint (first free, then costs a guess)
ctrl+g  give up
esc     menu
?       this help

hard mode keeps greens in place and
requires yellows in later guesses.
grays are still allowed.

any key closes this overlay
`))
}

func (m Model) viewPlay(th theme) string {
	s := m.session
	meta := m.mode
	if s.Hard {
		meta += "  ·  hard"
	}
	var b strings.Builder
	b.WriteString(th.header(meta, 40) + "\n\n")
	b.WriteString(m.renderGrid(th) + "\n\n")
	b.WriteString(m.renderKeyboard(th) + "\n")
	if s.Message != "" {
		b.WriteString("\n" + th.title(s.Message) + "\n")
	}
	if s.Status == game.Won {
		b.WriteString("\n" + th.title(fmt.Sprintf("%s  ·  %d/%d", praise(len(s.Guesses)), len(s.Guesses), s.MaxGuesses)) + "\n")
		if m.mode == "daily" {
			b.WriteString(th.dim("c copy   esc menu"))
		} else {
			b.WriteString(th.dim("n new   r replay   c copy   esc menu"))
		}
	} else if s.Status == game.Lost {
		b.WriteString("\n" + th.title(s.Answer) + "\n")
		if m.mode == "daily" {
			b.WriteString(th.dim("c copy   esc menu"))
		} else {
			b.WriteString(th.dim("n new   r replay   c copy   esc menu"))
		}
	} else {
		b.WriteString("\n" + th.dim("tab hint   ? help   esc menu"))
	}
	return b.String()
}

func (m Model) renderGrid(th theme) string {
	s := m.session
	compact := m.compactTiles()
	rows := make([]string, 0, s.MaxGuesses)
	for i := 0; i < s.MaxGuesses; i++ {
		switch {
		case i < len(s.Guesses):
			g := s.Guesses[i]
			tiles := make([]string, 5)
			for c := 0; c < 5; c++ {
				tiles[c] = th.tile(g.Tiles[c], string(g.Word[c]), true, compact)
			}
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, tiles...))
		case i == len(s.Guesses) && s.Status == game.Playing:
			tiles := make([]string, 5)
			for c := 0; c < 5; c++ {
				ch := " "
				if c < len(s.Input) {
					ch = string(s.Input[c])
				}
				tiles[c] = th.tile(game.Empty, ch, c < len(s.Input), compact)
			}
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, tiles...))
		default:
			tiles := make([]string, 5)
			for c := 0; c < 5; c++ {
				tiles[c] = th.tile(game.Empty, " ", false, compact)
			}
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, tiles...))
		}
	}
	return lipgloss.JoinVertical(lipgloss.Center, rows...)
}

func (m Model) renderKeyboard(th theme) string {
	kb := m.session.Keyboard()
	row := func(letters string) string {
		keys := make([]string, 0, len(letters))
		for i := 0; i < len(letters); i++ {
			ch := letters[i]
			keys = append(keys, th.key(string(ch), kb[ch], false))
		}
		return lipgloss.JoinHorizontal(lipgloss.Top, keys...)
	}
	top := row("QWERTYUIOP")
	mid := "  " + row("ASDFGHJKL")
	bot := "     " + row("ZXCVBNM")
	return lipgloss.JoinVertical(lipgloss.Left, top, mid, bot)
}
