// Package config/config provides basic config constants + methods.
package config

import "github.com/charmbracelet/huh"

const (
	defaultCommitTitleCharLimit = 48
	defaultMessageTemplate      = "{{.Type}}{{if .Scope}}({{.Scope}}){{end}}{{if .IsBreakingChange}}!{{end}}: {{.Message}}"
)

var defaultPrefixes = []huh.Option[string]{
	huh.NewOption("feat - a new feature", "feat"),
	huh.NewOption("fix - a bug fix", "fix"),
	huh.NewOption("build - changes that affect the build system or external dependencies", "build"),
	huh.NewOption("chore - changes to the build process or auxiliary tools and libraries", "chore"),
	huh.NewOption("ci - changes to our CI configuration files and scripts", "ci"),
	huh.NewOption("docs - documentation only changes", "docs"),
	huh.NewOption("perf - a code change that improves performance", "perf"),
	huh.NewOption("refactor - a code change that neither fixes a bug nor adds a feature", "refactor"),
	huh.NewOption("revert - reverts a previous commit", "revert"),
	huh.NewOption("style - changes that do not affect the meaning of the code", "style"),
	huh.NewOption("test - adding missing tests or correcting existing tests", "test"),
}

type LoadConfigReturn struct {
	MessageTemplate           string
	MessageWithTicketTemplate string
	Prefixes                  []huh.Option[string]
	CommitTitleCharLimit      int
	CommitBodyCharLimit       int
	CommitBodyLineLength      int
}

// LoadConfig loads the config file from the current directory or any parent
func LoadConfig() (LoadConfigReturn, error) {
	return LoadConfigReturn{
		MessageTemplate:      defaultMessageTemplate,
		Prefixes:             defaultPrefixes,
		CommitTitleCharLimit: defaultCommitTitleCharLimit,
	}, nil
}
