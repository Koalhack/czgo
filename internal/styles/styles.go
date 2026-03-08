// Package styles/styles provides styling rules.
package styles

import "github.com/charmbracelet/lipgloss"

// ─── Palette ─────────────────────────────────────────────────────────────────────────

var Colors = struct {
	Red    lipgloss.AdaptiveColor
	Indigo lipgloss.AdaptiveColor
	Green  lipgloss.AdaptiveColor
	Gray   lipgloss.AdaptiveColor
}{
	Red: lipgloss.AdaptiveColor{
		Light: "#FE5F86",
		Dark:  "#FE5F86",
	},
	Indigo: lipgloss.AdaptiveColor{
		Light: "#5A56E0",
		Dark:  "#7571F9",
	},
	Green: lipgloss.AdaptiveColor{
		Light: "#02BA84",
		Dark:  "#02BF87",
	},
	Gray: lipgloss.AdaptiveColor{
		Light: "#25262C",
		Dark:  "#878080",
	},
}

// ─── General ─────────────────────────────────────────────────────────────────────────

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(Colors.Indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(Colors.Indigo).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(Colors.Red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

// ─── Prefixes ────────────────────────────────────────────────────────────────────────

type PrefixStyles struct {
	Title       lipgloss.Style
	Description lipgloss.Style
	Separator   lipgloss.Style
}

func DefaultPrefixStyles() PrefixStyles {
	return PrefixStyles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(Colors.Green),

		Description: lipgloss.NewStyle().
			Faint(true).
			Foreground(Colors.Gray),
	}
}
