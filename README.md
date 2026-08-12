# wordle

Wordle in the terminal. Mac and Linux.

## Install

Copy one block. It installs Go if needed, installs `wordle`, and puts it on your PATH.

**macOS**

```bash
brew install go && go install github.com/parthUltra/wordle@latest && echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc && source ~/.zshrc && wordle help
```

**Linux (apt)**

```bash
sudo apt install -y golang-go && go install github.com/parthUltra/wordle@latest && echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc && source ~/.bashrc && wordle help
```

If Go is already installed, you only need:

```bash
go install github.com/parthUltra/wordle@latest
```

…and make sure `$(go env GOPATH)/bin` is on your PATH (the one-liners above do that for you).

Other Linux distros / latest Go: https://go.dev/dl/

From a clone instead:

```bash
git clone https://github.com/parthUltra/wordle.git && cd wordle && make install && echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.zshrc && source ~/.zshrc
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
