package ui

// forms encapsulates all form used in the application.
type forms struct {
	projectCreate projectForm
}

// getForms returns all forms to be used in the application.
func getForms() forms {
	return forms{
		projectCreate: newProjectForm(),
	}
}

// resetProjectCreate resets the state of the projectCreate form.
func (f *forms) resetProjectCreate() {
	f.projectCreate = newProjectForm()
}
