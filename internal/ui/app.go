package ui

import tea "charm.land/bubbletea/v2"

// screen represents the current display mode of the application.
type screen int

// mainAction represents the main menu actions.
type mainAction int

// crudAction represents the available CRUD actions.
type crudAction int

// searchAction represents the available search menu actions.
type searchAction int

const (
	screenMain screen = iota
	screenProjectActions
	screenProblemActions
	screenSearch
)

const (
	actionProject mainAction = iota
	actionProblem
	actionSearch
)

const (
	actionCreate crudAction = iota
	actionEdit
	actionDelete
)

const (
	searchProjects searchAction = iota
	searchProblems
)

// string represents the string attribute of the underlying type.
func (m mainAction) string() string {
	switch m {
	case actionProject:
		return "Projects"
	case actionProblem:
		return "Problems"
	case actionSearch:
		return "Search"
	default:
		return "Unknown"
	}
}

// string represents the string attribute of the underlying type.
func (c crudAction) string() string {
	switch c {
	case actionCreate:
		return "Create"
	case actionEdit:
		return "Edit"
	case actionDelete:
		return "Delete"
	default:
		return "Unknown"
	}
}

// string represents the string attribute of the underlying type.
func (s searchAction) string() string {
	switch s {
	case searchProjects:
		return "Projects"
	case searchProblems:
		return "Problems"
	default:
		return "Unknown"
	}
}

// Model represents the structure of the application state.
type Model struct {
	// main represents the main menu actions.
	main []mainAction
	// crd represents the crud actions.
	crud []crudAction
	// search represents the search menu actions.
	search []searchAction
	// cursor represents the current selected action.
	cursor int
	// screen represents the mode of the application state.
	screen screen
	// history tracks the application state.
	history []screen
	// err handles the current displayed error.
	err error
}

// InitialModel returns the initial internal state of the application.
func InitialModel() Model {
	return Model{
		main:    []mainAction{actionProject, actionProblem, actionSearch},
		crud:    []crudAction{actionCreate, actionEdit, actionDelete},
		search:  []searchAction{searchProjects, searchProblems},
		cursor:  0,
		screen:  screenMain,
		history: make([]screen, 0),
		err:     nil,
	}
}

// Init represents the initial state of the application.
func (m Model) Init() tea.Cmd {
	return nil
}
