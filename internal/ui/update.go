package ui

import tea "charm.land/bubbletea/v2"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyPressMsg:
		// What was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "backspace" moves the user back to the previous view
		case "backspace", "esc", "h":
			if len(m.history) == 0 {
				break // nothing to do
			}
			m.screen = m.history[len(m.history)-1]
			m.history = m.history[:len(m.history)-1]
			m.cursor = 0

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
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	return m, nil
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
	// TODO: handle action based on screen type
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
