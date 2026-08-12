# wordle

Wordle in the terminal. Mac and Linux.

## Install

### One command

**macOS**

```bash
brew install go && go install github.com/parthUltra/wordle@latest && echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc && source ~/.zshrc && wordle help
```

**Linux (apt)**

```bash
sudo apt install -y golang-go && go install github.com/parthUltra/wordle@latest && echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc && source ~/.bashrc && wordle help
```

### Step by step

**1. Install Go** (skip if you already have it)

macOS:

```bash
brew install go
```

Linux:

```bash
sudo apt install golang-go
```

Other distros / latest Go: https://go.dev/dl/

**2. Install wordle**

```bash
go install github.com/parthUltra/wordle@latest
```

**3. Put it on your PATH** (once)

macOS (zsh):

```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc && source ~/.zshrc
```

Linux (bash):

```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc && source ~/.bashrc
```

Then type `wordle` and play.

From a clone instead:

```bash
git clone https://github.com/parthUltra/wordle.git
cd wordle
make install
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
