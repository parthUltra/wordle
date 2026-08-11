package words

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"
)

type Lists struct {
	Answers []string
	Allowed map[string]struct{}
}

func Parse(answersText, allowedText string) (*Lists, error) {
	answers := parseLines(answersText)
	if len(answers) == 0 {
		return nil, fmt.Errorf("no answers in word list")
	}
	allowed := make(map[string]struct{}, 12000)
	for _, w := range parseLines(allowedText) {
		allowed[w] = struct{}{}
	}
	for _, w := range answers {
		allowed[w] = struct{}{}
	}
	return &Lists{Answers: answers, Allowed: allowed}, nil
}

func parseLines(s string) []string {
	raw := strings.Split(s, "\n")
	out := make([]string, 0, len(raw))
	for _, line := range raw {
		w := strings.ToUpper(strings.TrimSpace(line))
		if len(w) == 5 {
			out = append(out, w)
		}
	}
	return out
}

func (l *Lists) Random(seed uint64) string {
	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
	}
	r := rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15))
	return l.Answers[r.IntN(len(l.Answers))]
}

func (l *Lists) Daily(day time.Time) string {
	key := day.In(time.Local).Format("2006-01-02")
	sum := sha256.Sum256([]byte("wordle-daily:" + key))
	idx := binary.BigEndian.Uint64(sum[:8]) % uint64(len(l.Answers))
	return l.Answers[idx]
}

func (l *Lists) BySeed(seed uint64) string {
	if seed == 0 {
		return l.Random(0)
	}
	return l.Answers[seed%uint64(len(l.Answers))]
}
