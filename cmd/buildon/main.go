package main

import (
	"database/sql"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	_ "modernc.org/sqlite"

	"github.com/JBK2116/buildon/internal/db"
	"github.com/JBK2116/buildon/internal/logging"

	"github.com/JBK2116/buildon/internal/ui"
)

func main() {
	logger := logging.GetLogger()

	conn, err := db.OpenConn(logger)
	if err != nil {
		panic(err)
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			logger.Error("potential resource leak with closing db connection", "error", err)
		}
	}(conn)

	repository := db.NewRepository(conn, logger)

	model := ui.InitialModel(repository)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there has been a problem in running the application %v", err)
		os.Exit(1)
	}
}
