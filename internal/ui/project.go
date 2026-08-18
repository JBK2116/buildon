package ui

import (
	"errors"
	"strings"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("63")).
			MarginBottom(1)

	formFieldStyle = lipgloss.NewStyle().
			MarginBottom(1)
)

const (
	TitleLimit        int = 100
	DescriptionWidth  int = 50
	DescriptionHeight int = 3
)

// projectForm encapsulates the bubble components used to collect data when
// creating a project, keeping them out of the top-level Model.
type projectForm struct {
	// Title captures the project title.
	title textinput.Model
	// description captures the project description.
	description textarea.Model
	// focusIndex tracks the current focused on component
	focusIndex int
	// err holds the last validation error, if any.
	err error
}

// formSubmittedMsg is sent when a form has been validated and submitted.
type formSubmittedMsg struct{}

// newProjectForm returns a fresh projectForm with its components initialized.
func newProjectForm() projectForm {
	input := textinput.New()
	input.Placeholder = ""
	input.Prompt = "> "
	input.CharLimit = TitleLimit

	description := textarea.New()
	description.Placeholder = ""
	description.SetWidth(DescriptionWidth)
	description.SetHeight(DescriptionHeight)

	form := projectForm{title: input, description: description}
	form.title.Focus()
	return form
}

// getTitle returns the value of the title input.
func (p *projectForm) getTitle() string {
	val := p.title.Value()
	val = strings.TrimSpace(val)
	return val
}

// getDescription returns the value of the description textarea.
func (p *projectForm) getDescription() string {
	val := p.description.Value()
	val = strings.TrimSpace(val)
	if val == "" {
		return "No Description Provided"
	}
	return val
}

// view renders the createForm fields stacked vertically. The title placeholder
// signals it is required while the description placeholder signals it is
// optional.
func (p *projectForm) view() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		formFieldStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				labelStyle.Render("Title"),
				p.title.View(),
			),
		),
		formFieldStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				labelStyle.Render("Description"),
				p.description.View(),
			),
		),
	)
}

// getError returns the last validation error, if any.
func (p *projectForm) getError() error {
	return p.err
}

// validate ensures that the createForm is valid.
func (p *projectForm) validate() error {
	if len(p.getTitle()) == 0 {
		return errors.New("title is required")
	}
	return nil
}

func (p *projectForm) update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		// ctrl+s submits the form once it validates. enter is intentionally left
		// to the focused field so newlines can be entered in the textarea.
		case "ctrl+s":
			if err := p.validate(); err != nil {
				p.err = err
				return nil
			}
			p.err = nil
			return p.submit()
		// tab moves focus to the next field.
		case "tab":
			return p.moveFocus(1)
		// shift+tab moves focus to the previous field.
		case "shift+tab":
			return p.moveFocus(-1)
		default:
			// Typing clears any previous validation error.
			p.err = nil
			// Route everything else to the focused field.
			return p.updateFocused(msg)
		}
	}
	return nil
}

// submit signals the model that the form has been validated and submitted.
func (p *projectForm) submit() tea.Cmd {
	return func() tea.Msg {
		return formSubmittedMsg{}
	}
}

// moveFocus shifts focus by delta fields, wrapping around the available fields.
func (p *projectForm) moveFocus(delta int) tea.Cmd {
	count := p.fieldCount()
	next := (p.focusIndex + delta) % count
	if next < 0 {
		next += count
	}
	if next == p.focusIndex {
		return nil
	}
	p.blurFocused()
	p.focusIndex = next
	return p.focusIndexed()
}

// updateFocused routes a key press to the currently focused field.
func (p *projectForm) updateFocused(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd
	switch p.focusIndex {
	case 0:
		p.title, cmd = p.title.Update(msg)
	case 1:
		p.description, cmd = p.description.Update(msg)
	}
	return cmd
}

// fieldCount returns the number of fields in the form.
func (p *projectForm) fieldCount() int {
	return 2 //nolint:mnd // no need for a global var to replace this
}

// blurFocused blurs the currently focused field.
func (p *projectForm) blurFocused() {
	switch p.focusIndex {
	case 0:
		p.title.Blur()
	case 1:
		p.description.Blur()
	}
}

// focusIndexed focuses the field at the current focusIndex.
func (p *projectForm) focusIndexed() tea.Cmd {
	switch p.focusIndex {
	case 0:
		return p.title.Focus()
	case 1:
		return p.description.Focus()
	}
	return nil
}
