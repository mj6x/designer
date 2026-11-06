package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mj6x/designer/src/gif"
	"github.com/mj6x/designer/src/gradient"
)

var (
	title        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#d000ff"))
	active       = lipgloss.NewStyle().Foreground(lipgloss.Color("#f56dff")).Bold(true)
	card         = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#49304f")).Padding(1, 3).Width(46)
	selectedCard = card.Copy().BorderForeground(lipgloss.Color("#f56dff"))
	description  = lipgloss.NewStyle().Foreground(lipgloss.Color("#a98ba9")).Width(40)
	footer       = lipgloss.NewStyle().Foreground(lipgloss.Color("#8a648b"))
)

type tui struct {
	cursor int
	width  int
	height int
}

// NewTUI creates a new terminal UI.
func NewTUI() tui {
	return tui{}
}

// Init initialises the TUI.
func (t tui) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input and manages menu navigation.
func (t tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.width = msg.Width
		t.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return t, tea.Quit

		case "up", "k":
			if t.cursor > 0 {
				t.cursor--
			}

		case "down", "j":
			if t.cursor < 1 {
				t.cursor++
			}

		case "enter":
			switch t.cursor {
			case 0:
				return gradient.NewTUI(), nil
			case 1:
				return gif.NewTUI(), nil
			}
		}
	}

	return t, nil
}

// View renders the main menu.
func (t tui) View() string {
	gradientCard := lipgloss.JoinVertical(
		lipgloss.Left,
		title.Render("01  GRADIENT GENERATOR"),
		"",
		description.Render("Create beautiful colour gradients and export them to terminal designs."),
		"",
		active.Render("● ● ● ● ●"),
	)

	gifCard := lipgloss.JoinVertical(
		lipgloss.Left,
		title.Render("02  GIF CONVERTER"),
		"",
		description.Render("Convert GIF animations into ANSI or ASCII TermFX files."),
	)

	first := card.Render(gradientCard)
	second := card.Render(gifCard)

	if t.cursor == 0 {
		first = selectedCard.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				active.Render("01  GRADIENT GENERATOR"),
				"",
				description.Render("Create beautiful colour gradients and export them to terminal designs."),
				"",
				active.Render("● ● ● ● ●"),
			),
		)
	} else {
		second = selectedCard.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				active.Render("02  GIF CONVERTER"),
				"",
				description.Render("Convert GIF animations into ANSI or ASCII TermFX files."),
			),
		)
	}

	header := lipgloss.JoinVertical(
		lipgloss.Center,
		title.Render("Designer! by mj6x."),
		"",
		footer.Render("Choose a tool to get started."),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		first,
		"",
		second,
	)

	hints := lipgloss.JoinHorizontal(
		lipgloss.Center,
		footer.Render("↑/↓ Navigate"),
		footer.Render("    •    "),
		active.Render("Enter Select"),
		footer.Render("    •    "),
		footer.Render("Q Quit"),
	)

	view := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		"",
		content,
		"",
		hints,
	)

	return lipgloss.Place(
		t.width,
		t.height,
		lipgloss.Center,
		lipgloss.Center,
		view,
	)
}

func main() {
	p := tea.NewProgram(
		NewTUI(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
