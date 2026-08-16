package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/JBK2116/buildon/internal/ui"
)

func main() {
	p := tea.NewProgram(ui.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there has been a problem in running the application %v", err)
		os.Exit(1)
	}
}
