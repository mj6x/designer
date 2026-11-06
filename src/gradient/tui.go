package gradient

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mj6x/designer/src/colour"
)

// tui contains the bubble tea modules for the gradient interface.
type tui struct {
	focus       int
	colourInput textinput.Model
	fileInput   textinput.Model
	colours     []*colour.RGB
	generated   bool
	message     string
}

// Init initialises the model and starts the text input cursor blink.
func (t tui) Init() tea.Cmd {
	return textinput.Blink
}

var (
	title  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#d000ff"))
	muted  = lipgloss.NewStyle().Foreground(lipgloss.Color("#8a648b"))
	active = lipgloss.NewStyle().Foreground(lipgloss.Color("#f56dff")).Bold(true)
)

// NewTUI creates a new terminal ui for the gradient application.
func NewTUI() tui {
	colourInput := textinput.New()
	colourInput.Placeholder = "e.g #A78BFA or 167,139,250"
	colourInput.Width = 38
	colourInput.CharLimit = 40
	colourInput.Focus()

	fileInput := textinput.New()
	fileInput.Placeholder = "/path/to/file.txt"
	fileInput.Width = 38
	fileInput.CharLimit = 255

	return tui{fileInput: fileInput, colourInput: colourInput}
}

// Update handles keyboard input, manages input focus, and updates
// the model state when colours are added or removed.
func (t tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return t, tea.Quit
		}

		switch msg.String() {
		case "tab":
			t.focus = 1 - t.focus

			if t.focus == 0 {
				t.fileInput.Blur()
				return t, t.colourInput.Focus()
			}

			t.colourInput.Blur()
			return t, t.fileInput.Focus()

		case "ctrl+d":
			if len(t.colours) > 0 {
				c := t.colours[len(t.colours)-1]
				t.colours = t.colours[:len(t.colours)-1]
				t.message = fmt.Sprintf("Successfully removed the colour %s.", fmt.Sprintf("%d,%d,%d", c.R, c.G, c.B))
			}

			return t, nil

		case "ctrl+r":
			bytes, err := os.ReadFile(t.fileInput.Value())
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					t.message = fmt.Sprintf("Invalid render file %s provided!", t.fileInput.Value())
					t.fileInput.SetValue("")
					return t, nil
				}

				t.message = fmt.Sprintf("Failed to open %s! %v", t.fileInput.Value(), err)
				t.fileInput.SetValue("")
				return t, nil
			}

			writer, err := os.Create(strings.Split(filepath.Base(t.fileInput.Value()), filepath.Ext(filepath.Base(t.fileInput.Value())))[0] + ".design.txt")
			if err != nil {
				t.message = fmt.Sprintf("Failed to create output file! %v", err)
				t.fileInput.SetValue("")
				return t, nil
			}
			defer writer.Close()

			data := strings.ReplaceAll(string(bytes), "\r", "")
			render(data, writer, t.colours)
			t.message = fmt.Sprintf("Wrote gradient to the file, %s!", writer.Name())

			return t, tea.Tick(2500*time.Millisecond, func(time.Time) tea.Msg {
				return tea.Quit()
			})

		case "enter":
			if t.focus == 0 {
				colour, err := colour.StringToRGB(t.colourInput.Value())
				if err != nil {
					t.message = fmt.Sprintf("Invalid colour %s!", t.colourInput.Value())
					t.colourInput.SetValue("")
				} else {
					t.colours = append(t.colours, colour)
					t.message = fmt.Sprintf("Colour %s added successfully.", t.colourInput.Value())
					t.colourInput.SetValue("")
				}
			}

			return t, nil
		}

		var cmd tea.Cmd

		if t.focus == 0 {
			t.colourInput, cmd = t.colourInput.Update(msg)
		} else {
			t.fileInput, cmd = t.fileInput.Update(msg)
		}

		return t, cmd
	}

	return t, nil
}

// View renders the gradient interface, including the live preview,
// colour input, file input, and keyboard shortcut hints.
func (t tui) View() string {
	colourLabel := "Colour"
	fileLabel := "Input file"

	if t.focus == 0 {
		colourLabel = active.Render("> Colour")
	} else {
		fileLabel = active.Render("> Input file")
	}

	preview := muted.Render("Please enter a colour to see a preview.")

	if len(t.colours) >= 2 {
		preview = renderGradient(t.colours)
	} else if len(t.colours) > 0 {
		preview = renderGradient([]*colour.RGB{t.colours[0], t.colours[0]})
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title.Render("Gradient generator."),
		"",

		title.Render("Preview"),
		preview,
		"",

		colourLabel,
		t.colourInput.View(),
		"",

		fileLabel,
		t.fileInput.View(),

		"",
		t.message,
		"",

		muted.Render("enter    Concludes an input."),
		muted.Render("tab      Switches input."),
		muted.Render("ctrl+r   Renders the gradient."),
		muted.Render("ctrl+d   Removes the last colour."),
		muted.Render("ctrl+c   Exit out of Design."),
	)
}
