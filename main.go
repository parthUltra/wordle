package main

import (
	"embed"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/parthUltra/wordle/internal/store"
	"github.com/parthUltra/wordle/internal/ui"
	"github.com/parthUltra/wordle/internal/words"
)

//go:embed words.txt all_words.txt
var wordFS embed.FS

const usage = `wordle — play Wordle in the terminal

Install once, then:

  wordle           menu
  wordle daily     today's puzzle (shows your board if already finished)
  wordle random    a new practice game
  wordle stats     win streak and guess distribution
  wordle help      this text
`

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cmd, err := parseCmd(args)
	if err != nil {
		return err
	}
	if cmd == "help" {
		fmt.Print(usage)
		return nil
	}
	answers, err := wordFS.ReadFile("words.txt")
	if err != nil {
		return err
	}
	allowed, err := wordFS.ReadFile("all_words.txt")
	if err != nil {
		return err
	}
	lists, err := words.Parse(string(answers), string(allowed))
	if err != nil {
		return err
	}
	st, err := store.Open("")
	if err != nil {
		return err
	}
	p := tea.NewProgram(ui.New(lists, st, cmd))
	_, err = p.Run()
	return err
}

func parseCmd(args []string) (string, error) {
	if len(args) == 0 {
		return "", nil
	}
	switch args[0] {
	case "daily", "random", "stats", "help":
		return args[0], nil
	case "-h", "--help":
		return "help", nil
	default:
		return "", fmt.Errorf("unknown command %q\n\n%s", args[0], usage)
	}
}
