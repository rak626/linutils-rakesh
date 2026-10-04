package modules

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/rak626/linutils-rakesh/internal/config"
	"github.com/rak626/linutils-rakesh/internal/pkgmanager"
	"github.com/rak626/linutils-rakesh/internal/tui"
)

const defaultFontSize = 11

// FontSizes holds the global base size plus per-app overrides.
// App IDs match config.FontSizeApps (e.g. "alacritty", "gnome", "zed").
type FontSizes struct {
	Base   int
	PerApp map[string]int
}

// Get returns the size for an app, falling back to Base then defaultFontSize.
func (f FontSizes) Get(app string) int {
	if f.PerApp != nil {
		if v, ok := f.PerApp[app]; ok && v > 0 {
			return v
		}
	}
	if f.Base > 0 {
		return f.Base
	}
	return defaultFontSize
}

// UniformFontSizes returns sizes where every app inherits base.
func UniformFontSizes(base int) FontSizes {
	if base <= 0 {
		base = defaultFontSize
	}
	return FontSizes{Base: base, PerApp: map[string]int{}}
}

// FontSizesFromConfig rebuilds sizes from persisted variables.conf,
// with per-app keys falling back to $font_size.
func FontSizesFromConfig() FontSizes {
	base := config.FontSizeFor("") // empty app => base fallback chain
	per := map[string]int{}
	for _, app := range config.FontSizeApps {
		per[app] = config.FontSizeFor(app)
	}
	return FontSizes{Base: base, PerApp: per}
}

// persistFontChoice writes family + base + per-app sizes to variables.conf.
func persistFontChoice(family string, sizes FontSizes) {
	config.UserVars["$font_mono"] = family
	if sizes.Base <= 0 {
		sizes.Base = defaultFontSize
	}
	config.UserVars["$font_size"] = strconv.Itoa(sizes.Base)
	for _, app := range config.FontSizeApps {
		config.UserVars[config.FontSizeKey(app)] = strconv.Itoa(sizes.Get(app))
	}
	config.SaveVariables()
}

// fontAppLabel is the human-facing name + hint shown during per-app prompts.
type fontAppLabel struct {
	title string
	hint  string
}

var fontAppLabels = map[string]fontAppLabel{
	"gnome":          {"GNOME system", "gsettings + gtk-3/4 settings.ini"},
	"gnome_terminal": {"GNOME Terminal", "all profiles via gsettings"},
	"alacritty":      {"Alacritty", "~/.config/alacritty/alacritty.toml"},
	"kitty":          {"Kitty", "~/.config/kitty/kitty.conf"},
	"foot":           {"Foot", "~/.config/foot/foot.ini"},
	"rofi":           {"Rofi", "config.rasi font"},
	"wofi":           {"Wofi", "~/.config/wofi/style.css"},
	"waybar":         {"Waybar", "~/.config/waybar/style.css"},
	"mako":           {"Mako", "~/.config/mako/config"},
	"qt":             {"Qt (qt5ct/qt6ct)", "fixed_font point size"},
	"i3":             {"i3", "font pango lines in i3 config"},
	"i3status":       {"i3status", "bar config font"},
	"polybar":        {"Polybar", "font-0"},
	"zed":            {"Zed", "buffer + terminal font size"},
	"vscode":         {"VS Code / VSCodium", "editor + terminal font size"},
}

func fontAppTitle(app string) string {
	if l, ok := fontAppLabels[app]; ok {
		return l.title
	}
	return app
}

// SetupFontSwitcher is the interactive entry point wired to the TUI menu.
// It lists only installed fonts (fc-list), lets the user pick one + size,
// persists the choice to variables.conf, then applies everywhere with
// per-app installed guards (skip if not installed).
func SetupFontSwitcher(manager pkgmanager.PackageManager) error {
	fmt.Println("\n--- Font Switcher (all apps, guard-first) ---")

	families, err := ListInstalledFontFamilies()
	if err != nil {
		return fmt.Errorf("failed to list installed fonts (is fontconfig installed?): %v", err)
	}
	if len(families) == 0 {
		return fmt.Errorf("no installed fonts found via fc-list; run Fonts Setup first")
	}

	// Prefer JetBrains Mono variants at the top, keep the rest alphabetical.
	jetbrains, rest := splitJetBrainsFirst(families)
	ordered := append(jetbrains, rest...)

	var items []tui.ListItem
	for _, f := range ordered {
		desc := "Apply " + f + " to GNOME + Hyprland + i3 apps"
		if isJetBrainsMono(f) {
			desc = "Recommended mono — " + desc
		}
		items = append(items, tui.ListItem{Key: f, Name: f, Description: desc})
	}

	action, results, err := tui.RunListUIWithDesc("Font Switcher", "Pick an INSTALLED font. Missing fonts are never installed here.", items)
	if err != nil || action == "" || action == "back" || action == "quit" {
		fmt.Println("Font switch cancelled.")
		return err
	}
	chosen := ""
	for _, it := range results {
		if it.Selected {
			chosen = it.Key
			break
		}
	}
	if chosen == "" {
		fmt.Println("No font selected.")
		return nil
	}

	if !IsFontInstalled(chosen) {
		return fmt.Errorf("font %q not found via fc-list; run Fonts Setup first, nothing was changed", chosen)
	}

	size := pickFontSize()
	base := size
	sizes := pickPerAppSizes(base)
	fmt.Printf("\nSwitching every installed app to %q (base %d, per-app overrides)...\n", chosen, base)
	printFontSizes(sizes)

	// Persist choice so Restore/Re-run can reuse it.
	persistFontChoice(chosen, sizes)

	return SwitchFontEverywhere(chosen, sizes, manager)
}

// pickFontSize offers the base size via the same list UI.
func pickFontSize() int {
	sizes := []tui.ListItem{
		{Key: "9", Name: "9 — compact", Description: "Small UI / laptop"},
		{Key: "10", Name: "10 — compact", Description: "Small UI / laptop"},
		{Key: "11", Name: "11 — default (Recommended)", Description: "Balanced for GNOME + tiled WMs"},
		{Key: "12", Name: "12 — large", Description: "Terminals + HiDPI"},
		{Key: "13", Name: "13 — XL", Description: "Large text"},
		{Key: "14", Name: "14 — XXL", Description: "Accessibility"},
		{Key: "16", Name: "16 — huge", Description: "Low vision / demo"},
	}
	action, results, err := tui.RunListUIWithDesc("Font Size (base)", "Base size. Next: per-app overrides, Enter keeps base.", sizes)
	if err != nil || action == "" || action == "back" || action == "quit" {
		return defaultFontSize
	}
	for _, it := range results {
		if it.Selected {
			if n, err := strconv.Atoi(it.Key); err == nil {
				return n
			}
		}
	}
	return defaultFontSize
}

// pickPerAppSizes prompts one size per app, pre-filled with base.
// Enter on "Keep base" keeps it — only change where needed.
func pickPerAppSizes(base int) FontSizes {
	if base <= 0 {
		base = defaultFontSize
	}
	sizes := FontSizes{Base: base, PerApp: map[string]int{}}
	for i, app := range config.FontSizeApps {
		label, ok := fontAppLabels[app]
		if !ok {
			label = fontAppLabel{title: app}
		}
		title := fmt.Sprintf("Font Size %d/%d — %s", i+1, len(config.FontSizeApps), label.title)
		desc := fmt.Sprintf("Base %d — Enter keeps it. %s", base, label.hint)
		sizes.PerApp[app] = pickSingleAppSize(title, desc, base)
	}
	return sizes
}

func pickSingleAppSize(title, desc string, base int) int {
	items := []tui.ListItem{
		{Key: strconv.Itoa(base), Name: fmt.Sprintf("Keep base (%d) (Recommended)", base), Description: "Use the base size for this app"},
	}
	for _, n := range []int{8, 9, 10, 11, 12, 13, 14, 16, 18} {
		if n == base {
			continue
		}
		items = append(items, tui.ListItem{
			Key:         strconv.Itoa(n),
			Name:        fmt.Sprintf("%d", n),
			Description: fmt.Sprintf("Use %dpt for this app", n),
		})
	}
	action, results, err := tui.RunListUIWithDesc(title, desc, items)
	if err != nil || action == "" || action == "back" || action == "quit" {
		return base
	}
	for _, it := range results {
		if it.Selected {
			if n, err := strconv.Atoi(it.Key); err == nil && n > 0 {
				return n
			}
		}
	}
	return base
}

func printFontSizes(sizes FontSizes) {
	fmt.Printf("  base: %d\n", sizes.Base)
	for _, app := range config.FontSizeApps {
		fmt.Printf("  %-14s %d\n", fontAppTitle(app)+":", sizes.Get(app))
	}
}

// SwitchFontEverywhere applies family + per-app sizes to GNOME, Hyprland and i3 blocks.
// Every single target is guarded: missing binary AND missing package AND
// missing config file => SKIP with a log line, never create new configs.
func SwitchFontEverywhere(family string, sizes FontSizes, manager pkgmanager.PackageManager) error {
	home, _ := os.UserHomeDir()
	ok, failed := 0, 0
	note := func(applied bool, msg string) {
		if applied {
			ok++
			fmt.Printf("  [OK]   %s\n", msg)
		} else {
			fmt.Printf("  [SKIP] %s\n", msg)
		}
	}

	fmt.Println("\n[1/3] GNOME / system block")
	for _, r := range applyGnomeBlock(home, family, sizes.Get("gnome"), sizes.Get("gnome_terminal")) {
		note(r.applied, r.msg)
		if !r.applied && r.fatal {
			failed++
		}
	}

	fmt.Println("\n[2/3] Hyprland block")
	for _, r := range applyHyprlandBlock(home, family, sizes, manager) {
		note(r.applied, r.msg)
	}

	fmt.Println("\n[3/3] i3 block")
	for _, r := range applyI3Block(home, family, sizes, manager) {
		note(r.applied, r.msg)
	}

	fmt.Println("\n[4/4] Editors block (shared)")
	for _, r := range applyEditorsBlock(home, family, sizes, manager) {
		note(r.applied, r.msg)
	}

	reloadAfterSwitch(family, sizes.Base)

	// Persist even for programmatic callers.
	persistFontChoice(family, sizes)

	fmt.Printf("\nFont switch complete: %d applied, %d skipped. (%q base %dpt)\n", ok, failed, family, sizes.Base)
	if failed > 0 {
		return fmt.Errorf("%d target(s) failed", failed)
	}
	return nil
}

// SwitchFontEverywhereSingle keeps the old one-size behavior for programmatic
// callers: every app inherits size.
func SwitchFontEverywhereSingle(family string, size int, manager pkgmanager.PackageManager) error {
	return SwitchFontEverywhere(family, UniformFontSizes(size), manager)
}

// ListInstalledFontFamilies returns sorted unique families from fc-list.
// Only installed fonts are ever offered — nothing is downloaded here.
func ListInstalledFontFamilies() ([]string, error) {
	// Prefer strict mono spacing list; fall back to all families.
	if out, err := exec.Command("fc-list", ":spacing=mono", "family").Output(); err == nil {
		if fams := parseFcListFamilies(string(out)); len(fams) > 0 {
			return fams, nil
		}
	}
	out, err := exec.Command("fc-list", ":", "family").Output()
	if err != nil {
		return nil, err
	}
	return parseFcListFamilies(string(out)), nil
}

func parseFcListFamilies(out string) []string {
	seen := map[string]bool{}
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(strings.TrimSuffix(line, "\n"))
		if line == "" {
			continue
		}
		for _, part := range strings.Split(line, ",") {
			f := strings.TrimSpace(part)
			if f != "" && !seen[f] {
				seen[f] = true
			}
		}
	}
	fams := make([]string, 0, len(seen))
	for f := range seen {
		fams = append(fams, f)
	}
	sort.Strings(fams)
	return fams
}

// IsFontInstalled reports whether family resolves via fontconfig.
func IsFontInstalled(family string) bool {
	if family == "" {
		return false
	}
	out, err := exec.Command("fc-list", family, "family").Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) != ""
}

func splitJetBrainsFirst(families []string) ([]string, []string) {
	var jb, rest []string
	for _, f := range families {
		if isJetBrainsMono(f) {
			jb = append(jb, f)
		} else {
			rest = append(rest, f)
		}
	}
	// Strict Mono variant first (terminal-safe), then other JetBrains.
	sort.Slice(jb, func(a, b int) bool {
		ma := strings.Contains(jb[a], "Mono")
		mb := strings.Contains(jb[b], "Mono")
		if ma != mb {
			return ma
		}
		return jb[a] < jb[b]
	})
	return jb, rest
}

func isJetBrainsMono(f string) bool {
	l := strings.ToLower(f)
	return strings.Contains(l, "jetbrains")
}

// ---------- guard helpers ----------

type applyResult struct {
	applied bool
	fatal   bool // hard failure (vs expected skip)
	msg     string
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func installed(manager pkgmanager.PackageManager, pkgs []string, bins []string, files []string) bool {
	for _, b := range bins {
		if pkgmanager.IsCommandAvailable(b) {
			return true
		}
	}
	if manager != nil {
		for _, p := range pkgs {
			if manager.IsInstalled(p) {
				return true
			}
		}
	}
	for _, f := range files {
		if fileExists(f) {
			return true
		}
	}
	return false
}

func backupOnce(path string) {
	bak := path + ".bak"
	if fileExists(bak) {
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	_ = os.WriteFile(bak, data, 0644)
}

// ---------- per-DE blocks ----------

func applyGnomeBlock(home, family string, gnomeSize, termSize int) []applyResult {
	var out []applyResult
	fontStr := fmt.Sprintf("%s %d", family, gnomeSize)
	boldStr := fmt.Sprintf("%s Bold %d", family, gnomeSize)

	// gsettings — only if the binary exists (GNOME or GTK apps on any DE).
	if pkgmanager.IsCommandAvailable("gsettings") {
		sets := [][]string{
			{"org.gnome.desktop.interface", "monospace-font-name", fontStr},
			{"org.gnome.desktop.interface", "font-name", fontStr},
			{"org.gnome.desktop.interface", "document-font-name", fontStr},
			{"org.gnome.desktop.wm.preferences", "titlebar-font", boldStr},
		}
		for _, s := range sets {
			if err := pkgmanager.RunCommand("gsettings", "set", s[0], s[1], fmt.Sprintf("'%s'", s[2])); err != nil {
				out = append(out, applyResult{msg: fmt.Sprintf("gsettings %s %s (%v)", s[0], s[1], err)})
			} else {
				out = append(out, applyResult{applied: true, msg: fmt.Sprintf("gsettings %s %s = '%s'", s[0], s[1], s[2])})
			}
		}
		// GNOME Terminal profiles — best effort, skip silently if schema absent.
		out = append(out, applyGnomeTerminalProfiles(family, termSize)...)
	} else {
		out = append(out, applyResult{msg: "gsettings not found — GNOME system fonts skipped"})
	}

	// GTK ini files — only if the file already exists.
	for _, p := range []string{
		filepath.Join(home, ".config", "gtk-3.0", "settings.ini"),
		filepath.Join(home, ".config", "gtk-4.0", "settings.ini"),
	} {
		if !fileExists(p) {
			out = append(out, applyResult{msg: fmt.Sprintf("%s not present — skipped", p)})
			continue
		}
		backupOnce(p)
		if patchIniKey(p, "gtk-font-name", fontStr) {
			out = append(out, applyResult{applied: true, msg: fmt.Sprintf("%s gtk-font-name=%s", p, fontStr)})
		} else {
			out = append(out, applyResult{msg: fmt.Sprintf("%s patch failed", p)})
		}
	}

	// fontconfig fallback — only if user already has one.
	fcPath := filepath.Join(home, ".config", "fontconfig", "fonts.conf")
	if !fileExists(fcPath) {
		out = append(out, applyResult{msg: fcPath + " not present — skipped (no new file created)"})
	} else {
		backupOnce(fcPath)
		if patchFontconfigMonospace(fcPath, family) {
			out = append(out, applyResult{applied: true, msg: fcPath + " monospace prefer -> " + family})
		} else {
			out = append(out, applyResult{msg: fcPath + " patch failed"})
		}
	}
	return out
}

func applyGnomeTerminalProfiles(family string, size int) []applyResult {
	// Enumerate: gsettings get org.gnome.Terminal.ProfilesList list
	cmd := exec.Command("gsettings", "get", "org.gnome.Terminal.ProfilesList", "list")
	data, err := cmd.Output()
	if err != nil {
		return []applyResult{{msg: "gnome-terminal profiles absent — skipped"}}
	}
	// Output like ['<uuid>', ...] — extract quoted ids.
	re := regexp.MustCompile(`'([^']+)'`)
	ids := re.FindAllStringSubmatch(string(data), -1)
	if len(ids) == 0 {
		return []applyResult{{msg: "gnome-terminal has no profiles — skipped"}}
	}
	var out []applyResult
	fontStr := fmt.Sprintf("%s %d", family, size)
	for _, m := range ids {
		id := m[1]
		base := fmt.Sprintf("org.gnome.Terminal.Legacy.Profile:/org/gnome/Terminal/Legacy/Profiles:/:%s/", id)
		if err := pkgmanager.RunCommand("gsettings", "set", base, "font", fmt.Sprintf("'%s'", fontStr)); err != nil {
			continue
		}
		_ = pkgmanager.RunCommand("gsettings", "set", base, "use-system-font", "false")
		out = append(out, applyResult{applied: true, msg: fmt.Sprintf("gnome-terminal profile %s font = '%s'", shortID(id), fontStr)})
	}
	if len(out) == 0 {
		return []applyResult{{msg: "gnome-terminal profiles found but patch failed — skipped"}}
	}
	return out
}

func shortID(id string) string {
	if len(id) > 8 {
		return id[:8]
	}
	return id
}

func applyHyprlandBlock(home, family string, sizes FontSizes, manager pkgmanager.PackageManager) []applyResult {
	var out []applyResult

	// Shared terminals (also used on GNOME/i3, but grouped here once).
	out = append(out, applyAlacritty(home, family, sizes.Get("alacritty"), manager))
	out = append(out, applyKitty(home, family, sizes.Get("kitty"), manager))
	out = append(out, applyFoot(home, family, sizes.Get("foot"), manager))

	// Launchers / bar.
	wofiCSS := filepath.Join(home, ".config", "wofi", "style.css")
	if installed(manager, []string{"wofi"}, []string{"wofi"}, []string{wofiCSS}) && fileExists(wofiCSS) {
		backupOnce(wofiCSS)
		famOK := patchCssFontFamily(wofiCSS, family)
		sizeOK := patchCssFontSize(wofiCSS, fontPtToPx(sizes.Get("wofi")))
		if famOK && sizeOK {
			out = append(out, applyResult{applied: true, msg: fmt.Sprintf("wofi style.css font-family -> %s size %dpx", family, fontPtToPx(sizes.Get("wofi")))})
		} else if famOK || sizeOK {
			out = append(out, applyResult{applied: true, msg: "wofi style.css partially patched (family or size)"})
		} else {
			out = append(out, applyResult{msg: "wofi style.css patch failed"})
		}
	} else {
		out = append(out, applyResult{msg: "wofi not installed — skipped"})
	}

	for _, rp := range []string{
		filepath.Join(home, ".config", "rofi", "config.rasi"),
		filepath.Join(home, ".config", "rofi", "config", "config.rasi"),
	} {
		_ = rp
	}
	rofiPath := firstExisting(
		filepath.Join(home, ".config", "rofi", "config.rasi"),
		filepath.Join(home, ".config", "rofi", "config", "config.rasi"),
	)
	if installed(manager, []string{"rofi", "rofi-wayland"}, []string{"rofi"}, []string{rofiPath}) && rofiPath != "" {
		backupOnce(rofiPath)
		rofiSize := sizes.Get("rofi")
		if patchRasiFont(rofiPath, family, rofiSize) {
			out = append(out, applyResult{applied: true, msg: fmt.Sprintf("rofi %s font -> %s %d", rofiPath, family, rofiSize)})
		} else {
			out = append(out, applyResult{msg: "rofi config.rasi patch failed"})
		}
	} else {
		out = append(out, applyResult{msg: "rofi not installed — skipped"})
	}

	waybarCSS := filepath.Join(home, ".config", "waybar", "style.css")
	if installed(manager, []string{"waybar"}, []string{"waybar"}, []string{waybarCSS}) && fileExists(waybarCSS) {
		backupOnce(waybarCSS)
		famOK := patchCssFontFamily(waybarCSS, family)
		sizeOK := patchCssFontSize(waybarCSS, fontPtToPx(sizes.Get("waybar")))
		if famOK && sizeOK {
			out = append(out, applyResult{applied: true, msg: fmt.Sprintf("waybar style.css font-family -> %s size %dpx", family, fontPtToPx(sizes.Get("waybar")))})
		} else if famOK || sizeOK {
			out = append(out, applyResult{applied: true, msg: "waybar style.css partially patched (family or size)"})
		} else {
			out = append(out, applyResult{msg: "waybar style.css patch failed"})
		}
	} else {
		out = append(out, applyResult{msg: "waybar not installed — skipped"})
	}

	// Mako notifications.
	makoCfg := firstExisting(
		filepath.Join(home, ".config", "mako", "config"),
	)
	makoSize := sizes.Get("mako")
	if installed(manager, []string{"mako"}, []string{"mako"}, []string{makoCfg}) && makoCfg != "" {
		backupOnce(makoCfg)
		if patchIniKeyGeneric(makoCfg, "font", fmt.Sprintf("%s %d", family, makoSize)) {
			out = append(out, applyResult{applied: true, msg: fmt.Sprintf("mako font -> %s %d", family, makoSize)})
		} else {
			out = append(out, applyResult{msg: "mako config patch failed"})
		}
	} else {
		out = append(out, applyResult{msg: "mako not installed — skipped"})
	}

	// Qt fixed fonts (Hyprland needs qt5ct/qt6ct ini).
	qtSize := sizes.Get("qt")
	for _, qp := range []string{
		filepath.Join(home, ".config", "qt5ct", "qt5ct.conf"),
		filepath.Join(home, ".config", "qt6ct", "qt6ct.conf"),
	} {
		if !fileExists(qp) {
			continue // silent: Qt may simply be unconfigured
		}
		backupOnce(qp)
		if patchQtFixedFont(qp, family, qtSize) {
			out = append(out, applyResult{applied: true, msg: fmt.Sprintf("%s fixed_font -> %s %d", qp, family, qtSize)})
		}
	}
	return out
}

func applyI3Block(home, family string, sizes FontSizes, manager pkgmanager.PackageManager) []applyResult {
	var out []applyResult

	i3Size := sizes.Get("i3")
	i3cfg := firstExisting(
		filepath.Join(home, ".config", "i3", "config"),
		filepath.Join(home, ".i3", "config"),
	)
	if installed(manager, []string{"i3", "i3-wm"}, []string{"i3"}, []string{i3cfg}) && i3cfg != "" {
		backupOnce(i3cfg)
		if n := patchI3Fonts(i3cfg, family, i3Size); n > 0 {
			out = append(out, applyResult{applied: true, msg: fmt.Sprintf("i3 %s: %d font line(s) -> %s %d", i3cfg, n, family, i3Size)})
		} else {
			out = append(out, applyResult{msg: "i3 config: no pango font line patched"})
		}
	} else {
		out = append(out, applyResult{msg: "i3 not installed — skipped"})
	}

	// i3status / polybar bars.
	i3statusSize := sizes.Get("i3status")
	i3status := firstExisting(
		filepath.Join(home, ".config", "i3status", "config"),
		filepath.Join(home, ".i3status.conf"),
	)
	if fileExists(i3status) {
		backupOnce(i3status)
		if patchIniKeyGeneric(i3status, "font", fmt.Sprintf("pango:%s %d", family, i3statusSize)) {
			out = append(out, applyResult{applied: true, msg: fmt.Sprintf("i3status font -> %s %d", family, i3statusSize)})
		}
	} else {
		out = append(out, applyResult{msg: "i3status config not present — skipped"})
	}

	polySize := sizes.Get("polybar")
	poly := firstExisting(
		filepath.Join(home, ".config", "polybar", "config.ini"),
		filepath.Join(home, ".config", "polybar", "config"),
	)
	if installed(manager, []string{"polybar"}, []string{"polybar"}, []string{poly}) && poly != "" {
		backupOnce(poly)
		if patchIniKeyGeneric(poly, "font-0", fmt.Sprintf("%s:size=%d;1", family, polySize)) {
			out = append(out, applyResult{applied: true, msg: fmt.Sprintf("polybar font-0 -> %s %d", family, polySize)})
		} else {
			out = append(out, applyResult{msg: "polybar patch failed"})
		}
	} else {
		out = append(out, applyResult{msg: "polybar not installed — skipped"})
	}

	// dmenu has no config file — runtime -fn only.
	if pkgmanager.IsCommandAvailable("dmenu") {
		out = append(out, applyResult{msg: "dmenu uses -fn at runtime — set DMENU_FN='" + family + "' if your scripts support it"})
	}
	return out
}

func applyEditorsBlock(home, family string, sizes FontSizes, manager pkgmanager.PackageManager) []applyResult {
	var out []applyResult

	zedSize := sizes.Get("zed")
	zedPath := filepath.Join(home, ".config", "zed", "settings.json")
	if installed(manager, []string{"zed"}, []string{"zed", "zeditor"}, []string{zedPath}) && fileExists(zedPath) {
		backupOnce(zedPath)
		famOK := patchJsonFontKeys(zedPath, map[string]string{
			"buffer_font_family": family,
			"ui_font_family":     family,
		}, []string{"terminal", "font_family"}, family)
		sizeOK := patchJsonSizeKeys(zedPath,
			map[string]int{"buffer_font_size": zedSize},
			[]string{"terminal", "font_size"}, zedSize)
		if famOK && sizeOK {
			out = append(out, applyResult{applied: true, msg: fmt.Sprintf("zed settings.json font -> %s %d", family, zedSize)})
		} else if famOK || sizeOK {
			out = append(out, applyResult{applied: true, msg: "zed settings.json partially patched (family or size)"})
		} else {
			out = append(out, applyResult{msg: "zed settings.json patch failed"})
		}
	} else {
		out = append(out, applyResult{msg: "zed not installed — skipped"})
	}

	// VS Code / VSCodium.
	vscodeSize := sizes.Get("vscode")
	for _, cp := range []string{
		filepath.Join(home, ".config", "Code", "User", "settings.json"),
		filepath.Join(home, ".config", "VSCodium", "User", "settings.json"),
	} {
		bin := "code"
		if strings.Contains(cp, "VSCodium") {
			bin = "codium"
		}
		if !installed(manager, []string{"code", "visual-studio-code-bin", "vscodium"}, []string{bin, "code", "codium"}, []string{cp}) || !fileExists(cp) {
			continue
		}
		backupOnce(cp)
		famOK := patchJsonFontKeys(cp, map[string]string{
			"editor.fontFamily":              "'" + family + "', 'Symbols Nerd Font', monospace",
			"terminal.integrated.fontFamily": family,
		}, nil, "")
		sizeOK := patchJsonSizeKeys(cp,
			map[string]int{
				"editor.fontSize":              vscodeSize,
				"terminal.integrated.fontSize": vscodeSize,
			}, nil, 0)
		if famOK && sizeOK {
			out = append(out, applyResult{applied: true, msg: fmt.Sprintf("%s fontFamily -> %s size %d", cp, family, vscodeSize)})
		} else if famOK || sizeOK {
			out = append(out, applyResult{applied: true, msg: cp + " partially patched (family or size)"})
		}
	}

	// Neovim terminal buffers inherit the terminal — nothing to patch unless
	// a GUI guifont is explicitly configured.
	if pkgmanager.IsCommandAvailable("nvim") {
		out = append(out, applyResult{msg: "nvim terminal inherits terminal font — no file change (GUI guifont untouched)"})
	} else {
		out = append(out, applyResult{msg: "nvim not installed — skipped"})
	}
	return out
}

// ---------- per-app patchers ----------

func applyAlacritty(home, family string, size int, manager pkgmanager.PackageManager) applyResult {
	p := filepath.Join(home, ".config", "alacritty", "alacritty.toml")
	if !installed(manager, []string{"alacritty"}, []string{"alacritty"}, []string{p}) || !fileExists(p) {
		return applyResult{msg: "alacritty not installed — skipped"}
	}
	backupOnce(p)
	data, err := os.ReadFile(p)
	if err != nil {
		return applyResult{msg: "alacritty read failed: " + err.Error()}
	}
	s := string(data)
	famRe := regexp.MustCompile(`(?m)^(\s*family\s*=\s*")[^"]*(")`)
	if famRe.MatchString(s) {
		s = famRe.ReplaceAllString(s, `${1}`+family+`${2}`)
	} else {
		s += "\n[font.normal]\nfamily = \"" + family + "\"\n"
	}
	sizeRe := regexp.MustCompile(`(?m)^(\s*size\s*=\s*)[0-9]+(\.[0-9]+)?`)
	if sizeRe.MatchString(s) {
		s = sizeRe.ReplaceAllString(s, `${1}`+strconv.Itoa(size)+`.0`)
	}
	if err := os.WriteFile(p, []byte(s), 0644); err != nil {
		return applyResult{msg: "alacritty write failed: " + err.Error()}
	}
	return applyResult{applied: true, msg: fmt.Sprintf("alacritty.toml family -> %s size %d", family, size)}
}

func applyKitty(home, family string, size int, manager pkgmanager.PackageManager) applyResult {
	p := filepath.Join(home, ".config", "kitty", "kitty.conf")
	if !installed(manager, []string{"kitty"}, []string{"kitty"}, []string{p}) || !fileExists(p) {
		return applyResult{msg: "kitty not installed — skipped"}
	}
	backupOnce(p)
	data, _ := os.ReadFile(p)
	s := string(data)
	s = regexp.MustCompile(`(?m)^\s*font_family\s+.*$`).ReplaceAllString(s, "font_family "+family)
	if regexp.MustCompile(`(?m)^\s*font_size\s+`).MatchString(s) {
		s = regexp.MustCompile(`(?m)^\s*font_size\s+.*$`).ReplaceAllString(s, fmt.Sprintf("font_size %d.0", size))
	}
	_ = os.WriteFile(p, []byte(s), 0644)
	return applyResult{applied: true, msg: "kitty.conf font_family -> " + family}
}

func applyFoot(home, family string, size int, manager pkgmanager.PackageManager) applyResult {
	p := filepath.Join(home, ".config", "foot", "foot.ini")
	if !installed(manager, []string{"foot"}, []string{"foot"}, []string{p}) || !fileExists(p) {
		return applyResult{msg: "foot not installed — skipped"}
	}
	backupOnce(p)
	data, _ := os.ReadFile(p)
	s := string(data)
	want := fmt.Sprintf("font=%s:size=%d", family, size)
	if regexp.MustCompile(`(?m)^\s*font\s*=.*$`).MatchString(s) {
		s = regexp.MustCompile(`(?m)^\s*font\s*=.*$`).ReplaceAllString(s, want)
	} else {
		s += "\n[main]\n" + want + "\n"
	}
	_ = os.WriteFile(p, []byte(s), 0644)
	return applyResult{applied: true, msg: "foot.ini " + want}
}

func patchCssFontFamily(path, family string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	s := string(data)
	re := regexp.MustCompile(`(?i)font-family\s*:\s*[^;]+;`)
	want := `font-family: "` + family + `", "Symbols Nerd Font", monospace;`
	if re.MatchString(s) {
		s = re.ReplaceAllString(s, want)
	} else {
		s = "* { " + want + " }\n" + s
	}
	return os.WriteFile(path, []byte(s), 0644) == nil
}

// fontPtToPx maps terminal pt sizes to CSS px (96/72 ratio).
// 10pt->13px, 11pt->15px, 12pt->16px, 14pt->19px.
func fontPtToPx(pt int) int {
	if pt <= 0 {
		pt = defaultFontSize
	}
	return (pt*4 + 1) / 3
}

func patchCssFontSize(path string, px int) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	s := string(data)
	re := regexp.MustCompile(`(?i)font-size\s*:\s*[0-9.]+(px|pt)\s*;`)
	want := fmt.Sprintf("font-size: %dpx;", px)
	if re.MatchString(s) {
		s = re.ReplaceAllString(s, want)
	} else if strings.Contains(s, "{") {
		// Insert into the first rule block so existing selectors keep styling.
		s = strings.Replace(s, "{", "{\n  "+want, 1)
	} else {
		s = "* { " + want + " }\n" + s
	}
	return os.WriteFile(path, []byte(s), 0644) == nil
}

// patchQtFixedFont rewrites the qt5ct/qt6ct fixed_font value preserving the
// trailing style fields, only swapping family + point size.
// Format: "Family,Pt,-1,5,400,0,0,0,0,0,0,0,0,0,0,1"
func patchQtFixedFont(path, family string, size int) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	s := string(data)
	lines := strings.Split(s, "\n")
	changed := false
	re := regexp.MustCompile(`^(\s*fixed_font\s*=).*?$`)
	for i, ln := range lines {
		if !re.MatchString(ln) {
			continue
		}
		eq := strings.Index(ln, "=")
		if eq < 0 {
			continue
		}
		val := strings.TrimSpace(ln[eq+1:])
		parts := strings.Split(val, ",")
		if len(parts) >= 2 {
			parts[0] = family
			parts[1] = strconv.Itoa(size)
			lines[i] = ln[:eq+1] + " " + strings.Join(parts, ",")
		} else {
			lines[i] = ln[:eq+1] + " " + family + "," + strconv.Itoa(size) + ",-1,5,400,0,0,0,0,0,0,0,0,0,0,1"
		}
		changed = true
	}
	if !changed {
		return false
	}
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644) == nil
}

func patchRasiFont(path, family string, size int) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	s := string(data)
	want := fmt.Sprintf(`font: "%s %d";`, family, size)
	re := regexp.MustCompile(`(?m)font\s*:\s*"[^"]*"\s*;`)
	if re.MatchString(s) {
		s = re.ReplaceAllString(s, want)
	} else if strings.Contains(s, "configuration") {
		s = strings.Replace(s, "configuration {", "configuration {\n  "+want, 1)
	} else {
		s += "\nconfiguration {\n  " + want + "\n}\n"
	}
	return os.WriteFile(path, []byte(s), 0644) == nil
}

// patchIniKey replaces "key = ..." under any section; appends under
// [Settings] (or end of file) if the key is absent.
func patchIniKey(path, key, value string) bool {
	return patchIniKeyGeneric(path, key, value)
}

func patchIniKeyGeneric(path, key, value string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	s := string(data)
	// Quote value if it contains spaces and caller did not quote already.
	re := regexp.MustCompile(`(?m)^(\s*` + regexp.QuoteMeta(key) + `\s*=).*?$`)
	if re.MatchString(s) {
		// Preserve "=" spacing style by reusing prefix before "=".
		lines := strings.Split(s, "\n")
		for i, ln := range lines {
			trim := strings.TrimSpace(ln)
			if strings.HasPrefix(trim, key) {
				eq := strings.Index(ln, "=")
				if eq >= 0 {
					after := ln[eq+1:]
					// Preserve original spacing after "=": if there was
					// at least one space/tab, keep a single space.
					pad := ""
					if len(after) > 0 && (after[0] == ' ' || after[0] == '\t') {
						pad = " "
					}
					lines[i] = ln[:eq+1] + pad + strings.TrimSpace(value)
				}
			}
		}
		s = strings.Join(lines, "\n")
	} else {
		if strings.Contains(s, "[Settings]") {
			s = strings.Replace(s, "[Settings]", "[Settings]\n"+key+"="+value, 1)
		} else {
			s += "\n" + key + "=" + value + "\n"
		}
	}
	return os.WriteFile(path, []byte(s), 0644) == nil
}

func patchI3Fonts(path, family string, size int) int {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	s := string(data)
	re := regexp.MustCompile(`(?m)^(\s*font\s+pango:).*$`)
	matches := re.FindAllString(s, -1)
	if len(matches) == 0 {
		return 0
	}
	s = re.ReplaceAllString(s, `${1}`+family+` `+strconv.Itoa(size))
	if os.WriteFile(path, []byte(s), 0644) != nil {
		return 0
	}
	return len(matches)
}

func patchFontconfigMonospace(path, family string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	s := string(data)
	if strings.Contains(s, "<family>"+family+"</family>") {
		return true // already correct
	}
	// Replace the first <family> inside the monospace <prefer> alias.
	re := regexp.MustCompile(`(?s)(<alias>\s*<family>monospace</family>\s*<prefer>\s*<family>)[^<]*(</family>)`)
	if re.MatchString(s) {
		s = re.ReplaceAllString(s, `${1}`+family+`${2}`)
		return os.WriteFile(path, []byte(s), 0644) == nil
	}
	return false
}

// patchJsonFontKeys sets top-level string keys and one nested path (e.g.
// terminal.font_family).
func patchJsonFontKeys(path string, top map[string]string, nested []string, nestedVal string) bool {
	data, err := os.ReadFile(path)
	if err != nil || len(bytesTrimSpace(data)) == 0 {
		return false
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return false
	}
	changed := 0
	for k, v := range top {
		m[k] = v
		changed++
	}
	if len(nested) == 2 && nestedVal != "" {
		sub, ok := m[nested[0]].(map[string]interface{})
		if !ok {
			sub = map[string]interface{}{}
		}
		sub[nested[1]] = nestedVal
		m[nested[0]] = sub
		changed++
	}
	out, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return false
	}
	out = append(out, '\n')
	if os.WriteFile(path, out, 0644) != nil {
		return false
	}
	return changed > 0
}

// patchJsonSizeKeys sets top-level numeric keys and one nested numeric path
// (e.g. terminal.font_size). A nil nested path (or size <= 0 there) skips the
// nested write; top entries with size <= 0 are skipped.
func patchJsonSizeKeys(path string, top map[string]int, nested []string, nestedVal int) bool {
	data, err := os.ReadFile(path)
	if err != nil || len(bytesTrimSpace(data)) == 0 {
		return false
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return false
	}
	changed := 0
	for k, v := range top {
		if v <= 0 {
			continue
		}
		m[k] = v
		changed++
	}
	if len(nested) == 2 && nestedVal > 0 {
		sub, ok := m[nested[0]].(map[string]interface{})
		if !ok {
			sub = map[string]interface{}{}
		}
		sub[nested[1]] = nestedVal
		m[nested[0]] = sub
		changed++
	}
	if changed == 0 {
		return false
	}
	out, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return false
	}
	out = append(out, '\n')
	if os.WriteFile(path, out, 0644) != nil {
		return false
	}
	return true
}

func bytesTrimSpace(b []byte) []byte {
	return []byte(strings.TrimSpace(string(b)))
}

func firstExisting(paths ...string) string {
	for _, p := range paths {
		if p != "" && fileExists(p) {
			return p
		}
	}
	return ""
}

func reloadAfterSwitch(family string, size int) {
	_ = family
	_ = size
	_ = pkgmanager.RunCommand("fc-cache", "-f")
	if pkgmanager.IsCommandAvailable("hyprctl") {
		_ = pkgmanager.RunCommand("hyprctl", "reload")
	}
	if pkgmanager.IsCommandAvailable("i3-msg") {
		// Only restarts when i3 is actually running; harmless otherwise.
		_ = exec.Command("i3-msg", "restart").Run()
	}
	if _, err := exec.Command("pgrep", "waybar").Output(); err == nil {
		_ = exec.Command("killall", "-SIGUSR2", "waybar").Run()
	}
}
