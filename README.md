# wordle

Wordle in the terminal. Mac and Linux.

## Install

You need [Go 1.25+](https://go.dev/dl/). Then one command:

```bash
go install github.com/parthUltra/wordle@latest
```

Put Go's bin directory on your `PATH` if it isn't already (once):

```bash
# zsh (default on macOS)
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc && source ~/.zshrc

# bash (common on Linux)
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc && source ~/.bashrc
```

Type `wordle` and play.

From a clone instead:

```bash
git clone https://github.com/parthUltra/wordle.git
cd wordle
make install          # copies to ~/.local/bin
```

If `wordle` is not found after `make install`:

```bash
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.zshrc && source ~/.zshrc
```

## Play

```bash
wordle           # menu
wordle daily     # today's puzzle
wordle random    # practice
wordle stats     # streak and guess distribution
wordle help
```

If you already finished today's daily, `wordle daily` opens **that** board — it does not start a new one.

| key | action |
|---|---|
| letters | type a guess |
| enter | submit |
| backspace | delete |
| tab | hint (first free, extras cost a guess) |
| `?` | help |
| esc | menu |

Settings (palette, hard mode) live in the in-app menu. Config and stats: `~/Library/Application Support/wordle` on macOS, `~/.config/wordle` on Linux.

Hard mode is NYT-style: greens stay put, yellows must be reused. Gray letters are still allowed.

```bash
go test ./...
```
