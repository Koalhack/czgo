// Package config/prefix provides basic config constants + methods for prefix (git type).
package config

import (
	"strings"

	"github.com/Koalhack/czgo/internal/styles"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type Prefix struct {
	name        string
	description string
}

type Prefixes []Prefix

var defaultPrefixes = Prefixes{
	{name: "feat", description: "a new feature"},
	{name: "fix", description: "a bug fix"},
	{name: "build", description: "changes that affect the build system or external dependencies"},
	{name: "chore", description: "changes to the build process or auxiliary tools and libraries"},
	{name: "ci", description: "changes to our CI configuration files and scripts"},
	{name: "docs", description: "documentation only changes"},
	{name: "perf", description: "a code change that improves performance"},
	{name: "refactor", description: "a code change that neither fixes a bug nor adds a feature"},
	{name: "revert", description: "reverts a previous commit"},
	{name: "style", description: "changes that do not affect the meaning of the code"},
	{name: "test", description: "adding missing tests or correcting existing tests"},
}

func (p *Prefixes) Options(s styles.PrefixStyles) []huh.Option[string] {
	prefixes := []Prefix(*p)

	maxWidth := 0
	for _, prefix := range prefixes {
		if w := lipgloss.Width(prefix.name); w > maxWidth {
			maxWidth = w
		}
	}

	var items []huh.Option[string]
	for _, prefix := range prefixes {
		title := s.Title.Render(prefix.name)
		desc := s.Description.Render(prefix.description)

		gap := 2
		padding := maxWidth - lipgloss.Width(prefix.name)
		spaces := strings.Repeat(" ", padding+gap)

		label := title + spaces + desc
		items = append(items, huh.NewOption(label, prefix.name))
	}
	return items
}
