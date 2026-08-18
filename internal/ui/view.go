package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Styles used to render the interface.
var (
	appStyle = lipgloss.NewStyle().
			Padding(1, 2) //nolint:mnd // This is not worth converting to named variable

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("63")).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Faint(true).
			MarginBottom(1)

	optionStyle = lipgloss.NewStyle().
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			Padding(0, 1)

	footerStyle = lipgloss.NewStyle().
			Faint(true).
			MarginTop(1)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			MarginTop(1)
)

// View displays the application state to the user.
func (m Model) View() tea.View {
	if m.screen == screenCreateProject {
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			headerStyle.Render("Create Project (Title required)"),
			subtitleStyle.Render("Description is optional. Fill in the details below."),
			m.renderFormError(),
			m.form.projectCreate.view(),
			m.renderNotification(),
			footerStyle.Render("tab / shift+tab switch field  •  ctrl+s submit  •  esc back"),
		)
		v := tea.NewView(appStyle.Render(content))
		v.AltScreen = true
		return v
	}

	var header, subtitle string
	var options []string

	switch m.screen {
	case screenMain:
		header = "Main Menu"
		subtitle = "What would you like to work on?"
		for i, action := range m.main {
			options = append(options, m.renderOption(i, action.string()))
		}
	case screenProjectActions, screenProblemActions:
		if m.screen == screenProjectActions {
			header = "Projects"
		} else {
			header = "Problems"
		}
		subtitle = "Choose an action to perform."
		for i, action := range m.crud {
			options = append(options, m.renderOption(i, action.string()))
		}
	case screenSearch:
		header = "Search"
		subtitle = "What would you like to search?"
		for i, action := range m.search {
			options = append(options, m.renderOption(i, action.string()))
		}
	default:
		v := tea.NewView("")
		v.AltScreen = true
		return v
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		headerStyle.Render(header),
		subtitleStyle.Render(subtitle),
		lipgloss.JoinVertical(lipgloss.Left, options...),
		m.renderNotification(),
		footerStyle.Render("↑/↓ navigate  •  enter confirm  •  backspace back  •  q quit"),
		footerStyle.Render("vim: h/j/k/l navigate"),
	)

	v := tea.NewView(appStyle.Render(content))
	v.AltScreen = true
	return v
}

// renderFormError renders the current form validation error, if any.
func (m Model) renderFormError() string {
	if err := m.form.projectCreate.getError(); err != nil {
		return errorStyle.Render(err.Error())
	}
	return ""
}

// renderNotification renders the current notification, or an empty string when
// none is set.
func (m Model) renderNotification() string {
	if m.notification.message == "" {
		return ""
	}
	style := successStyle
	if m.notification.isError {
		style = errorStyle
	}
	return style.Render(m.notification.message)
}

// renderOption formats a single option line, highlighting the one under cursor.
func (m Model) renderOption(i int, label string) string {
	if i == m.cursor {
		return selectedStyle.Render("▸ " + label)
	}

	return optionStyle.Render("  " + label)
}
