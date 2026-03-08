package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Koalhack/czgo/internal/config"
	"github.com/Koalhack/czgo/internal/template"
	"github.com/Koalhack/czgo/internal/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const maxWidth = 80

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

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
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

type state int

const (
	statusMessage state = iota
	statusNormal
	stateDone
)

type Model struct {
	width  int
	state  state
	lg     *lipgloss.Renderer
	styles *Styles
	form   *huh.Form
	commit *types.Commit
	config config.LoadConfigReturn
}

func msgForm(m *Model) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("type").
				Options(m.config.Prefixes...).
				Value(&m.commit.Type).
				Title("Type"),
		),
		huh.NewGroup(
			huh.NewInput().
				Key("scope").
				Value(&m.commit.Scope).
				Title("Scope"),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Key("break").
				Value(&m.commit.IsBreakingChange).
				Title("Breaking Change").
				Description("Is this a breaking change?").
				Affirmative("Yes!").
				Negative("Nope."),
		),
	).
		WithWidth(45).
		WithShowHelp(false).
		WithShowErrors(false)
}

func mainForm(m *Model) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("message").
				Title("Message").
				Value(&m.commit.Message),
			huh.NewText().
				Key("body").
				Title("Body").
				ShowLineNumbers(true).
				Lines(8),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Key("done").
				Title("Ready to commit?").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("Welp, finish up then")
					}
					return nil
				}).
				Affirmative("Yep").
				Negative("Wait, no"),
		),
	).
		WithWidth(45).
		WithShowHelp(false).
		WithShowErrors(false)
}

func NewModel(cfg config.LoadConfigReturn) Model {
	m := Model{width: maxWidth}
	m.commit = &types.Commit{}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.config = cfg

	m.form = msgForm(&m)
	return m
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		case "esc":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	switch m.state {
	case statusMessage:
		if m.form.State == huh.StateCompleted {
			commitMsg, err := template.RenderCommitMsg(m.config.MessageTemplate, m.commit)
			if err != nil {
				cmds = append(cmds, tea.Quit)
				break
			}

			m.state = statusNormal
			m.commit.Message = commitMsg
			m.form = mainForm(&m)

			cmds = append(cmds, m.form.Init())
		}

	case statusNormal:
		if m.form.State == huh.StateCompleted {
			m.state = stateDone

			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := m.styles

	switch m.form.State {
	case huh.StateCompleted:
		var b strings.Builder
		return s.Status.Margin(0, 1).Padding(1, 2).Width(48).Render(b.String()) + "\n\n"
	default:

		// Form (left side)
		v := strings.TrimSuffix(m.form.View(), "\n\n")
		form := m.lg.NewStyle().Margin(1, 0).Render(v)

		// Status (right side)
		var status string

		errors := m.form.Errors()
		header := m.appBoundaryView("Charm Employment Application")
		if len(errors) > 0 {
			header = m.appErrorBoundaryView(m.errorView())
		}
		body := lipgloss.JoinHorizontal(lipgloss.Left, form, status)

		footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))
		if len(errors) > 0 {
			footer = m.appErrorBoundaryView("")
		}

		return s.Base.Render(header + "\n" + body + "\n\n" + footer)
	}
}

func (m Model) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += err.Error()
	}
	return s
}

func (m Model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}

func (m Model) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(red),
	)
}

func main() {
	if err := run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	p := tea.NewProgram(NewModel(cfg))

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
