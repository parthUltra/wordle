package game

import (
	"fmt"
	"strings"
	"unicode"
)

type Tile uint8

const (
	Empty Tile = iota
	Absent
	Present
	Correct
)

type Status uint8

const (
	Playing Status = iota
	Won
	Lost
)

type Guess struct {
	Word  string
	Tiles [5]Tile
}

type Options struct {
	MaxGuesses int
	Hard       bool
}

type Session struct {
	Answer        string
	Allowed       map[string]struct{}
	Answers       []string
	MaxGuesses    int
	Hard          bool
	Guesses       []Guess
	Input         string
	Status        Status
	HintsUsed     int
	HintPenalties int
	HintLog       []string
	Message       string
}

type Snapshot struct {
	Answer        string
	Guesses       []string
	Hard          bool
	MaxGuesses    int
	HintPenalties int
	HintLog       []string
	Won           bool
	Lost          bool
}

func (s *Session) Snapshot() Snapshot {
	words := make([]string, len(s.Guesses))
	for i, g := range s.Guesses {
		words[i] = g.Word
	}
	return Snapshot{
		Answer:        s.Answer,
		Guesses:       words,
		Hard:          s.Hard,
		MaxGuesses:    s.MaxGuesses,
		HintPenalties: s.HintPenalties,
		HintLog:       append([]string(nil), s.HintLog...),
		Won:           s.Status == Won,
		Lost:          s.Status == Lost,
	}
}

func Restore(allowed map[string]struct{}, answers []string, snap Snapshot) *Session {
	s := NewSession(snap.Answer, allowed, answers, Options{MaxGuesses: snap.MaxGuesses, Hard: false})
	for _, w := range snap.Guesses {
		s.Allowed[strings.ToUpper(w)] = struct{}{}
		s.Input = w
		_ = s.Submit()
	}
	s.Hard = snap.Hard
	s.HintPenalties = snap.HintPenalties
	s.HintLog = append([]string(nil), snap.HintLog...)
	s.HintsUsed = len(s.HintLog)
	if snap.Won {
		s.Status = Won
	} else if snap.Lost {
		s.Status = Lost
	}
	s.Input = ""
	return s
}

func NewSession(answer string, allowed map[string]struct{}, answers []string, opt Options) *Session {
	answer = strings.ToUpper(answer)
	if opt.MaxGuesses <= 0 {
		opt.MaxGuesses = 6
	}
	if allowed == nil {
		allowed = map[string]struct{}{}
	}
	allowed[answer] = struct{}{}
	return &Session{
		Answer:     answer,
		Allowed:    allowed,
		Answers:    answers,
		MaxGuesses: opt.MaxGuesses,
		Hard:       opt.Hard,
		Status:     Playing,
	}
}

// Score implements official Wordle marking: greens first, then yellows,
// each answer letter consumed at most once.
func Score(guess, answer string) [5]Tile {
	guess = strings.ToUpper(guess)
	answer = strings.ToUpper(answer)
	var tiles [5]Tile
	remain := []byte(answer)
	for i := 0; i < 5; i++ {
		if guess[i] == answer[i] {
			tiles[i] = Correct
			remain[i] = 0
		}
	}
	for i := 0; i < 5; i++ {
		if tiles[i] == Correct {
			continue
		}
		for j := 0; j < 5; j++ {
			if remain[j] == guess[i] {
				tiles[i] = Present
				remain[j] = 0
				break
			}
		}
		if tiles[i] == Empty {
			tiles[i] = Absent
		}
	}
	return tiles
}

func (s *Session) Type(r rune) {
	if s.Status != Playing || len(s.Input) >= 5 {
		return
	}
	if !unicode.IsLetter(r) {
		return
	}
	s.Input += strings.ToUpper(string(r))
	s.Message = ""
}

func (s *Session) Backspace() {
	if s.Status != Playing || s.Input == "" {
		return
	}
	s.Input = s.Input[:len(s.Input)-1]
	s.Message = ""
}

func (s *Session) Submit() error {
	if s.Status != Playing {
		return fmt.Errorf("game over")
	}
	word := strings.ToUpper(s.Input)
	if len(word) != 5 {
		s.Message = "Need 5 letters"
		return fmt.Errorf("need 5 letters")
	}
	if _, ok := s.Allowed[word]; !ok {
		s.Message = "Not in word list"
		return fmt.Errorf("not in word list")
	}
	if s.Hard {
		if err := s.hardError(word); err != nil {
			s.Message = err.Error()
			return err
		}
	}
	g := Guess{Word: word, Tiles: Score(word, s.Answer)}
	s.Guesses = append(s.Guesses, g)
	s.Input = ""
	s.Message = ""
	if word == s.Answer {
		s.Status = Won
		return nil
	}
	if s.GuessesUsed() >= s.MaxGuesses {
		s.Status = Lost
	}
	return nil
}

func (s *Session) GuessesUsed() int {
	return len(s.Guesses) + s.HintPenalties
}

func (s *Session) RemainingGuesses() int {
	n := s.MaxGuesses - s.GuessesUsed()
	if n < 0 {
		return 0
	}
	return n
}

func (s *Session) GiveUp() {
	if s.Status != Playing {
		return
	}
	s.Status = Lost
	s.Message = "The word was " + s.Answer
}

func (s *Session) Keyboard() map[byte]Tile {
	best := make(map[byte]Tile, 26)
	for _, g := range s.Guesses {
		for i := 0; i < 5; i++ {
			ch := g.Word[i]
			if g.Tiles[i] > best[ch] {
				best[ch] = g.Tiles[i]
			}
		}
	}
	return best
}

func (s *Session) UseHint() (string, error) {
	if s.Status != Playing {
		return "", fmt.Errorf("game over")
	}
	letter, ok := s.unknownAnswerLetter()
	if !ok {
		s.Message = "No unused letters to hint"
		return "", fmt.Errorf("no unused letters")
	}
	if s.HintsUsed > 0 {
		if s.RemainingGuesses() <= 0 {
			return "", fmt.Errorf("no guesses left")
		}
		s.HintPenalties++
		if s.GuessesUsed() >= s.MaxGuesses {
			s.Status = Lost
		}
	}
	s.HintsUsed++
	text := fmt.Sprintf("Contains %c · %d answers left", letter, s.RemainingAnswers())
	s.HintLog = append(s.HintLog, text)
	s.Message = text
	return text, nil
}

func (s *Session) unknownAnswerLetter() (byte, bool) {
	known := map[byte]bool{}
	for _, g := range s.Guesses {
		for i := 0; i < 5; i++ {
			if g.Tiles[i] == Present || g.Tiles[i] == Correct {
				known[g.Word[i]] = true
			}
		}
	}
	for _, h := range s.HintLog {
		if len(h) >= 10 && strings.HasPrefix(h, "Contains ") {
			known[h[9]] = true
		}
	}
	const rarest = "QJZXVWKFBHGPMYDCLUOTNRSAIE"
	best, bestRank := byte(0), len(rarest)+1
	for i := 0; i < 5; i++ {
		ch := s.Answer[i]
		if known[ch] {
			continue
		}
		rank := strings.IndexByte(rarest, ch)
		if rank < 0 {
			rank = len(rarest)
		}
		if rank < bestRank {
			best, bestRank = ch, rank
		}
	}
	if best == 0 {
		return 0, false
	}
	return best, true
}

func (s *Session) RemainingAnswers() int {
	n := 0
	for _, w := range s.Answers {
		if MatchesHistory(w, s.Guesses) {
			n++
		}
	}
	return n
}

func MatchesHistory(candidate string, guesses []Guess) bool {
	candidate = strings.ToUpper(candidate)
	if len(candidate) != 5 {
		return false
	}
	for _, g := range guesses {
		if Score(g.Word, candidate) != g.Tiles {
			return false
		}
	}
	return true
}

func (s *Session) Share(title string) string {
	var b strings.Builder
	score := "X"
	if s.Status == Won {
		score = fmt.Sprintf("%d", len(s.Guesses))
	}
	fmt.Fprintf(&b, "%s %s/%d\n", title, score, s.MaxGuesses)
	if s.Hard {
		b.WriteString("Hard\n")
	}
	for _, g := range s.Guesses {
		for _, tile := range g.Tiles {
			switch tile {
			case Correct:
				b.WriteString("🟩")
			case Present:
				b.WriteString("🟨")
			default:
				b.WriteString("⬛")
			}
		}
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n")
}

func (s *Session) hardError(word string) error {
	var must [5]byte
	minCount := map[byte]int{}
	for _, g := range s.Guesses {
		counts := map[byte]int{}
		for i := 0; i < 5; i++ {
			if g.Tiles[i] == Correct {
				must[i] = g.Word[i]
				counts[g.Word[i]]++
			}
			if g.Tiles[i] == Present {
				counts[g.Word[i]]++
			}
		}
		for ch, n := range counts {
			if n > minCount[ch] {
				minCount[ch] = n
			}
		}
	}
	for i := 0; i < 5; i++ {
		if must[i] != 0 && word[i] != must[i] {
			return fmt.Errorf("Hard: %c must stay in position %d", must[i], i+1)
		}
	}
	have := map[byte]int{}
	for i := 0; i < 5; i++ {
		have[word[i]]++
	}
	for ch, n := range minCount {
		if have[ch] < n {
			return fmt.Errorf("Hard: guess must include %c", ch)
		}
	}
	return nil
}
