package ui

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

const (
	NotificationDisplayLength = time.Second * 3
)

// clearNotificationMsg requests that the current notification be dismissed.
type clearNotificationMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyPressMsg:
		// While creating a project, route input to the form. ctrl+c/esc leave
		// the create screen rather than quitting the whole application.
		if m.screen == screenCreateProject {
			return m.handleCreateProjectKey(msg)
		}

		// What was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "backspace" moves the user back to the previous view
		case "backspace", "esc", "h":
			m = m.back()

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < m.currentLen()-1 {
				m.cursor++
			}

		// The "enter" key and the space bar select the item that the
		// cursor is pointing at.
		case "enter", "space", "l":
			switch m.screen {
			case screenMain:
				return m.selectMain()
			case screenProjectActions:
				return m.selectCrud(screenProjectActions)
			case screenProblemActions:
				return m.selectCrud(screenProblemActions)
			case screenSearch:
				return m.selectSearch()
			default:
				panic("unhandled default case")
			}
		}

	// A form was validated and submitted.
	case formSubmittedMsg:
		return m.createProject()

	// The current notification should be dismissed.
	case clearNotificationMsg:
		m.notification = notification{}
		return m, nil
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	return m, nil
}

// back returns the model after navigating to the previous screen in history.
func (m Model) back() Model {
	if len(m.history) == 0 {
		return m // nothing to do
	}
	m.screen = m.history[len(m.history)-1]
	m.history = m.history[:len(m.history)-1]
	m.cursor = 0
	return m
}

// handleCreateProjectKey routes a key press while on the create project screen.
func (m Model) handleCreateProjectKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	// ctrl+c and esc leave the create screen without quitting the app.
	case "ctrl+c", "esc":
		return m.back(), nil
	default:
		cmd := m.form.projectCreate.update(msg)
		return m, cmd
	}
}

// createProject persists the submitted form as a new project.
func (m Model) createProject() (tea.Model, tea.Cmd) {
	if err := m.repository.CreateProject(
		m.form.projectCreate.getTitle(),
		m.form.projectCreate.getDescription(),
	); err != nil {
		m.notification = notification{message: "Failed to create project: " + err.Error(), isError: true}
		return m, clearNotification()
	}
	m.notification = notification{message: "Project created successfully.", isError: false}
	m.form.resetProjectCreate()
	return m.back(), clearNotification()
}

// clearNotification returns a command that dismisses the current notification
// after a short delay.
func clearNotification() tea.Cmd {
	return tea.Tick(NotificationDisplayLength, func(_ time.Time) tea.Msg {
		return clearNotificationMsg{}
	})
}

// currentLen returns the number of selectable options on the current screen.
func (m Model) currentLen() int {
	switch m.screen {
	case screenMain:
		return len(m.main)
	case screenProjectActions, screenProblemActions:
		return len(m.crud)
	case screenSearch:
		return len(m.search)
	default:
		return 0
	}
}

// selectMain handles selecting an option on the main screen.
func (m Model) selectMain() (tea.Model, tea.Cmd) {
	if m.cursor >= len(m.main) {
		return m, nil
	}
	action := m.main[m.cursor]
	m.history = append(m.history, screenMain)
	switch action {
	case actionProject:
		m.screen = screenProjectActions
	case actionProblem:
		m.screen = screenProblemActions
	case actionSearch:
		m.screen = screenSearch
	}
	m.cursor = 0

	return m, nil
}

// selectCrud handles selecting a CRUD option on a project/problem screen.
func (m Model) selectCrud(from screen) (tea.Model, tea.Cmd) {
	if m.cursor >= len(m.crud) {
		return m, nil
	}
	action := m.crud[m.cursor]
	m.history = append(m.history, from)
	if from == screenProjectActions && action == actionCreate {
		m.form.resetProjectCreate()
		m.screen = screenCreateProject
		m.cursor = 0
		return m, nil
	}
	// TODO: handle remaining actions based on screen type
	_ = action

	return m, nil
}

// selectSearch handles selecting a search option on the search screen.
func (m Model) selectSearch() (tea.Model, tea.Cmd) {
	if m.cursor >= len(m.search) {
		return m, nil
	}
	action := m.search[m.cursor]
	m.history = append(m.history, screenSearch)
	// TODO: handle action based on search type
	_ = action

	return m, nil
}
