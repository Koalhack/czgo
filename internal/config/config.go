// Package config/config provides basic config constants + methods.
package config

import (
	"github.com/Koalhack/czgo/internal/styles"
	"github.com/charmbracelet/huh"
)

const (
	defaultCommitTitleCharLimit = 48
	defaultMessageTemplate      = "{{.Type}}{{if .Scope}}({{.Scope}}){{end}}{{if .IsBreakingChange}}!{{end}}: {{.Message}}"
)

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
	prefixStyle := styles.DefaultPrefixStyles()
	return LoadConfigReturn{
		MessageTemplate:      defaultMessageTemplate,
		Prefixes:             defaultPrefixes.Options(prefixStyle),
		CommitTitleCharLimit: defaultCommitTitleCharLimit,
	}, nil
}
