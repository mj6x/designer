package gif

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mj6x/designer/src/animation"
	"github.com/nfnt/resize"
)

type tui struct {
	focus        int
	fileInput    textinput.Model
	formatList   list.Model
	widthSlider  Slider
	heightSlider Slider
	loopsSlider  Slider
	generated    bool
	message      string
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

type formatItem string

// FilterValue returns the format's filter value.
func (f formatItem) FilterValue() string { return string(f) }

// Title returns the format's title.
func (f formatItem) Title() string { return string(f) }

// Description returns the format's description.
func (f formatItem) Description() string { return "" }

// NewTUI creates a new terminal ui for the GIF application.
func NewTUI() tui {
	fileInput := textinput.New()
	fileInput.Placeholder = "/path/to/file.gif or https://link.to/gif.gif"
	fileInput.Width = 38
	fileInput.CharLimit = 255
	fileInput.Focus()

	items := []list.Item{
		formatItem("ANSI"),
		formatItem("ASCII"),
	}

	formatList := list.New(items, list.NewDefaultDelegate(), 38, 4)
	formatList.SetShowTitle(false)
	formatList.SetShowStatusBar(false)
	formatList.SetShowHelp(false)
	formatList.SetFilteringEnabled(false)
	formatList.SetShowPagination(false)

	return tui{
		fileInput:    fileInput,
		formatList:   formatList,
		widthSlider:  NewSlider(10, 1000, 80, 1),
		heightSlider: NewSlider(10, 1000, 24, 1),
		loopsSlider:  NewSlider(0, 10000, 0, 1),
	}
}

// Update handles mouse and keyboard inputs and updates the tui state.
func (t tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return t, tea.Quit
		}

		if msg.String() == "ctrl+r" {
			input := t.fileInput.Value()
			t.fileInput.SetValue("")

			if input == "" {
				t.message = "Please provide a GIF file or URL."
				return t, nil
			}

			var data []byte
			var err error

			u, parseErr := url.ParseRequestURI(input)

			if parseErr == nil && (u.Scheme == "http" || u.Scheme == "https") {
				resp, err := http.Get(input)
				if err != nil {
					t.message = "Failed to download GIF. " + muted.Render(err.Error())
					return t, nil
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					t.message = "Failed to download GIF. " + muted.Render(resp.Status)
					return t, nil
				}

				data, err = io.ReadAll(resp.Body)
			} else {
				data, err = os.ReadFile(input)
			}

			if err != nil {
				t.message = "Failed to read GIF. " + muted.Render(err.Error())
				return t, nil
			}

			decoded, err := gif.DecodeAll(bytes.NewReader(data))
			if err != nil {
				t.message = "Invalid GIF file. " + muted.Render(err.Error())
				return t, nil
			}

			if len(decoded.Image) == 0 {
				t.message = "GIF contains no frames."
				return t, nil
			}

			bounds := decoded.Image[0].Bounds()
			base := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))

			filename := filepath.Base(input)
			if u != nil && (u.Scheme == "http" || u.Scheme == "https") {
				filename = filepath.Base(u.Path)
			}

			filename = strings.TrimSuffix(filename, filepath.Ext(filename))
			if filename == "" || filename == "." {
				filename = "output"
			}

			writer, err := os.Create(filename + ".design.tfx")
			if err != nil {
				t.message = "Failed to create output file. " + muted.Render(err.Error())
				return t, nil
			}

			defer writer.Close()

			if _, err := writer.WriteString("\x1b[?25l"); err != nil {
				t.message = "Failed to write output. " + muted.Render(err.Error())
				return t, nil
			}

			format := t.formatList.SelectedItem().(formatItem)
			width := uint(t.widthSlider.Value)
			height := uint(t.heightSlider.Value)

			loops := t.loopsSlider.Value + 1

			for x := 0; x < loops; x++ {
				for i, frame := range decoded.Image {
					draw.Draw(base, base.Bounds(), frame, frame.Bounds().Min, draw.Over)

					delay := decoded.Delay[i]
					if delay == 0 {
						delay = 1
					}

					resizeHeight := height
					if format == "ANSI" {
						resizeHeight *= 2
					}

					resized := resize.Resize(width, resizeHeight, base, resize.Lanczos2)

					var output string
					switch format {
					case "ANSI":
						output = animation.ANSI(resized)
					case "ASCII":
						output = animation.ASCII(resized)
					}

					output = strings.TrimSuffix(output, "\r\n")
					if _, err := fmt.Fprintf(writer, "\x1b[H%s<<sleep(%d)>>", output, delay*10); err != nil {
						t.message = "Failed to write output. " + muted.Render(err.Error())
						return t, nil
					}
				}
			}

			if _, err := writer.WriteString("\x1b[?25h\x1b[0m"); err != nil {
				t.message = "Failed to finish output. " + muted.Render(err.Error())
				return t, nil
			}

			if err := writer.Close(); err != nil {
				t.message = "Failed to close output file. " + muted.Render(err.Error())
				return t, nil
			}

			t.message = fmt.Sprintf("Wrote GIF to the file, %s!", writer.Name())

			return t, tea.Tick(2500*time.Millisecond, func(time.Time) tea.Msg {
				return tea.Quit()
			})
		}

		switch msg.String() {
		case "tab", "down":
			t.focus++
			if t.focus > 4 {
				t.focus = 0
				t.fileInput.Focus()
			} else {
				t.fileInput.Blur()
			}
			return t, nil

		case "up":
			t.focus--
			if t.focus < 0 {
				t.focus = 4
			}
			if t.focus == 0 {
				t.fileInput.Focus()
			} else {
				t.fileInput.Blur()
			}
			return t, nil
		}

		switch t.focus {
		case 0:
			var cmd tea.Cmd
			t.fileInput, cmd = t.fileInput.Update(msg)
			return t, cmd
		case 1:
			var cmd tea.Cmd
			t.formatList, cmd = t.formatList.Update(msg)
			return t, cmd
		case 2:
			t.widthSlider.Update(msg)
		case 3:
			t.heightSlider.Update(msg)
		case 4:
			t.loopsSlider.Update(msg)
		}

	case tea.MouseMsg:
		switch t.focus {
		case 1:
			var cmd tea.Cmd
			t.formatList, cmd = t.formatList.Update(msg)
			return t, cmd
		case 2:
			t.widthSlider.Update(msg)
		case 3:
			t.heightSlider.Update(msg)
		case 4:
			t.loopsSlider.Update(msg)
		}
	}

	if t.focus == 0 {
		var cmd tea.Cmd
		t.fileInput, cmd = t.fileInput.Update(msg)
		return t, cmd
	}

	return t, nil
}

// View renders the GIF user interface to the terminal.
func (t tui) View() string {
	fileLabel := muted.Render("Input file")
	formatLabel := muted.Render("Output format")
	widthLabel := muted.Render("Width")
	heightLabel := muted.Render("Height")
	loopsLabel := muted.Render("Loops")

	switch t.focus {
	case 0:
		fileLabel = active.Render("> Input file")
	case 1:
		formatLabel = active.Render("> Output format")
	case 2:
		widthLabel = active.Render("> Width")
	case 3:
		heightLabel = active.Render("> Height")
	case 4:
		loopsLabel = active.Render("> Loops")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title.Render("GIF Converter."),
		"",
		fileLabel,
		t.fileInput.View(),
		"",
		formatLabel,
		t.formatList.View(),
		"",
		widthLabel,
		t.widthSlider.View(),
		"",
		heightLabel,
		t.heightSlider.View(),
		"",
		loopsLabel,
		t.loopsSlider.View(),
		"",
		t.message,
		"",
		muted.Render("tab       Switch input / slider."),
		muted.Render("up/down   Navigate options."),
		muted.Render("left/right Adjust values."),
		muted.Render("enter     Select / edit value."),
		muted.Render("ctrl+r    Generate."),
		muted.Render("ctrl+c    Exit."),
	)
}
