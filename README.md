# linutils-rakesh

A simple TUI-based Linux setup and restore tool for people who distro-hop.
Built with Go, Bubbletea and Huh. Gruvbox Dark Medium theme, ASCII-only
(normal terminal font, no Nerd Font required).

It turns a fresh install — or a messed-up system — back into *your* system:
base packages, desktop environment, dotfiles (via GNU Stow), GNOME
performance tweaks + keybindings, shell/fonts/editors, and your software
catalog. Backups are plain git commits in your dotfiles repo.

## Menu

| Entry | What it does |
|---|---|
| Fresh Install | New system: OS setup, base tools, desktop (GNOME/i3/Hyprland), dotfiles, perf + keybinds, Git + GitHub auth, software picker. Asks every time. |
| Restore My Settings | Re-applies dotfiles, keybinds and perf. Idempotent — safe to re-run when things break. Skips OS setup and GitHub auth. |
| Backup Now to Git | Dumps GNOME `dconf` + package list into the dotfiles repo and pushes. History is git log. |
| Keybindings | Pick terminal/browser/editor/filemanager/launcher, then apply GNOME keybinds. |
| Software Picker | Choose from the catalog. Nothing is auto-installed. |

Extras (Bluetooth, SDDM, icons, scripts, NVIDIA, …) are hidden by default —
nothing is deleted. Show them with `LINUTILS_ADVANCED=1`.

## OS x Desktop support

| | GNOME (primary) | i3 | Hyprland |
|---|---|---|---|
| Arch | yes | yes | yes |
| Fedora | yes | yes | yes |
| Ubuntu / Debian | yes | yes | blocked (message shown) |

Package managers are abstracted (`apt`, `dnf`, `pacman` + AUR via `yay`).

## Quick start

```bash
# one-line installer (downloads the latest release binary)
curl -fsSL https://raw.githubusercontent.com/rak626/linutils-rakesh/main/install.sh | bash

# or build from source (Go 1.26+)
git clone https://github.com/rak626/linutils-rakesh.git
cd linutils-rakesh
go build -o linutils-rakesh .
./linutils-rakesh
```

Keys: `j/k` move, `Space` select, `Enter` run, `/` search, `?` help,
`q` / `Esc` back (on the main menu `q` quits the app).

## Configuration (forks welcome)

| Setting | Default | Override |
|---|---|---|
| Dotfiles repo | `https://github.com/rak626/dotfiles.git` | `DOTFILES_REPO` env var |
| Fonts repo | `https://github.com/rak626/fonts.git` | `FONTS_REPO` env var |
| App choices (`$terminal`, `$browser`, …) | `~/.config/linutils/variables.conf` | asked interactively every run |
| Repo cloner list | none (opt-in) | copy `examples/repos.example.conf` to `~/.config/linutils/repos.conf` |

Conflicting `~/.config` entries are moved to `~/.config.bak` before
stowing — never silently deleted.

## Architecture

- `main.go` — menu loop + Fresh/Restore flows
- `internal/tui/` — Bubbletea list UI, Gruvbox theme (`theme.go`), main menu (`form.go`)
- `internal/system/` — distro/DE/hardware detection
- `internal/pkgmanager/` — `apt` / `dnf` / `pacman` abstraction
- `internal/modules/` — one file per task (setup, dotfiles, backup, software picker, …)
- `internal/config/` — variables + software catalog (`installs.go`)

## Contributing

PRs welcome. Run `go build ./...` and `go vet ./...` before pushing.
CI checks both. Releases are cut by pushing a `v*` tag.

## License

MIT — see [LICENSE](LICENSE).
