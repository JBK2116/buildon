package ui

import tea "charm.land/bubbletea/v2"

// Screen represents the current display mode of the application.
type Screen int

// MainAction represents the main menu actions.
type MainAction int

// CrudAction represents the available CRUD actions.
type CrudAction int

// SearchAction represents the available search menu actions.
type SearchAction int

const (
	ScreenMain Screen = iota
	ScreenProjectActions
	ScreenProblemActions
	ScreenSearch
)

const (
	ActionProject MainAction = iota
	ActionProblem
	ActionSearch
)

const (
	ActionCreate CrudAction = iota
	ActionEdit
	ActionDelete
)

const (
	SearchProjects SearchAction = iota
	SearchProblems
)

// String represents the string attribute of the underlying type.
func (m MainAction) String() string {
	switch m {
	case ActionProject:
		return "Projects"
	case ActionProblem:
		return "Problems"
	case ActionSearch:
		return "Search"
	default:
		return "Unknown"
	}
}

// String represents the string attribute of the underlying type.
func (c CrudAction) String() string {
	switch c {
	case ActionCreate:
		return "Create"
	case ActionEdit:
		return "Edit"
	case ActionDelete:
		return "Delete"
	default:
		return "Unknown"
	}
}

// String represents the string attribute of the underlying type.
func (s SearchAction) String() string {
	switch s {
	case SearchProjects:
		return "Projects"
	case SearchProblems:
		return "Problems"
	default:
		return "Unknown"
	}
}

// Model represents the structure of the application state.
type Model struct {
	// main represents the main menu actions.
	main []MainAction
	// crd represents the crud actions.
	crud []CrudAction
	// search represents the search menu actions.
	search []SearchAction
	// cursor represents the current selected action.
	cursor int
	// screen represents the mode of the application state.
	screen Screen
	// history tracks the application state.
	history []Screen
	// err handles the current displayed error.
	err error
}

// InitialModel returns the initial internal state of the application.
func InitialModel() Model {
	return Model{
		main:    []MainAction{ActionProject, ActionProblem, ActionSearch},
		crud:    []CrudAction{ActionCreate, ActionEdit, ActionDelete},
		search:  []SearchAction{SearchProjects, SearchProblems},
		cursor:  0,
		screen:  ScreenMain,
		history: make([]Screen, 0),
		err:     nil,
	}
}

// Init represents the initial state of the application.
func (m Model) Init() tea.Cmd {
	return nil
}
