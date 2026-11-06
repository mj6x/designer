package gif

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Slider contains the values for an input slider.
type Slider struct {
	Min       int
	Max       int
	Value     int
	Step      int
	Width     int
	Dragging  bool
	Editing   bool
	Input     string
	LastClick time.Time
}

var (
	sliderTitle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#d000ff"))
	sliderMuted  = lipgloss.NewStyle().Foreground(lipgloss.Color("#8a648b"))
	sliderActive = lipgloss.NewStyle().Foreground(lipgloss.Color("#f56dff")).Bold(true)
)

// NewSlider creates a new slider.
func NewSlider(min, max, value, step int) Slider {
	if step < 1 {
		step = 1
	}
	if max < min {
		max = min
	}

	s := Slider{
		Min:   min,
		Max:   max,
		Value: value,
		Step:  step,
		Width: 43,
	}

	s.clamp()
	return s
}

func (s *Slider) clamp() {
	if s.Value < s.Min {
		s.Value = s.Min
	}
	if s.Value > s.Max {
		s.Value = s.Max
	}
}

func (s Slider) position() int {
	if s.Max == s.Min {
		return 0
	}

	return (s.Value - s.Min) * s.Width / (s.Max - s.Min)
}

func (s *Slider) setFromPosition(pos int) {
	if pos < 0 {
		pos = 0
	}
	if pos > s.Width {
		pos = s.Width
	}

	if s.Max == s.Min {
		s.Value = s.Min
		return
	}

	raw := s.Min + (pos*(s.Max-s.Min)+s.Width/2)/s.Width

	s.Value = s.Min + ((raw-s.Min)/s.Step)*s.Step
	s.clamp()
}

// Update handles slider input.
func (s *Slider) Update(msg tea.Msg) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if s.Editing {
			switch msg.String() {
			case "enter":
				if value, err := strconv.Atoi(s.Input); err == nil {
					s.Value = value
					s.clamp()
				}
				s.Editing = false

			case "esc":
				s.Editing = false

			case "backspace":
				if len(s.Input) > 0 {
					s.Input = s.Input[:len(s.Input)-1]
				}

			default:
				if len(msg.String()) == 1 {
					c := msg.String()[0]
					if c >= '0' && c <= '9' {
						s.Input += string(c)
					}
				}
			}
			return
		}

		switch msg.String() {
		case "left", "h":
			s.Value -= s.Step
		case "right", "l":
			s.Value += s.Step
		case "home":
			s.Value = s.Min
		case "end":
			s.Value = s.Max
		case "enter":
			s.Editing = true
			s.Input = strconv.Itoa(s.Value)
		}

		s.clamp()

	case tea.MouseMsg:
		switch msg.Action {
		case tea.MouseActionPress:
			if msg.Button == tea.MouseButtonLeft {
				if time.Since(s.LastClick) < 400*time.Millisecond {
					s.Editing = true
					s.Input = strconv.Itoa(s.Value)
					s.Dragging = false
					s.LastClick = time.Time{}
					return
				}

				s.LastClick = time.Now()
				s.Dragging = true
				s.setFromPosition(msg.X)
			}

		case tea.MouseActionMotion:
			if s.Dragging {
				s.setFromPosition(msg.X)
			}

		case tea.MouseActionRelease:
			s.Dragging = false
		}
	}
}

// View renders the slider.
func (s Slider) View() string {
	if s.Editing {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			sliderTitle.Render("Value"),
			sliderActive.Render("["+s.Input+"█]"),
			sliderMuted.Render("enter    Confirm value."),
			sliderMuted.Render("esc      Cancel editing."),
		)
	}

	pos := s.position()

	filled := strings.Repeat("━", pos)
	empty := strings.Repeat("━", s.Width-pos)

	bar := lipgloss.JoinHorizontal(
		lipgloss.Center,
		sliderTitle.Render(filled),
		sliderActive.Render("●"),
		sliderMuted.Render(empty),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		sliderTitle.Render("Value"),
		bar,
		sliderActive.Render(fmt.Sprintf("%d", s.Value)),
		sliderMuted.Render(fmt.Sprintf(
			"%d                                      %d",
			s.Min,
			s.Max,
		)),
		sliderMuted.Render("← → adjust | enter edit value"),
	)
}
