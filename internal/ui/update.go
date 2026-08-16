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
			switch m.screen {
			case ScreenMain:
				if m.cursor < len(m.main)-1 {
					m.cursor++
				}
			case ScreenProjectActions, ScreenProblemActions:
				if m.cursor < len(m.crud)-1 {
					m.cursor++
				}
			case ScreenSearch:
				if m.cursor < len(m.search)-1 {
					m.cursor++
				}
			}

		// The "enter" key and the space bar select the item that the
		// cursor is pointing at.
		case "enter", "space", "l":
			switch m.screen {
			case ScreenMain:
				if m.cursor >= len(m.main) {
					return m, nil
				}
				action := m.main[m.cursor]
				m.history = append(m.history, ScreenMain)
				switch action {
				case ActionProject:
					m.screen = ScreenProjectActions
				case ActionProblem:
					m.screen = ScreenProblemActions
				case ActionSearch:
					m.screen = ScreenSearch
				}
				m.cursor = 0
			case ScreenProjectActions:
				if m.cursor >= len(m.crud) {
					return m, nil
				}
				action := m.crud[m.cursor]
				m.history = append(m.history, ScreenProjectActions)
				// TODO: Handle the actions here via switch
				_ = action
			case ScreenProblemActions:
				if m.cursor >= len(m.crud) {
					return m, nil
				}
				action := m.crud[m.cursor]
				m.history = append(m.history, ScreenProblemActions)
				// TODO: Handle the actions here via switch
				_ = action
			case ScreenSearch:
				if m.cursor >= len(m.search) {
					return m, nil
				}
				action := m.search[m.cursor]
				m.history = append(m.history, ScreenSearch)
				// TODO: Handle the actions here via switch
				_ = action
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	return m, nil
}
