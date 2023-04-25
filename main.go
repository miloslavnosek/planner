package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"planner/database"
	"planner/task"
	"planner/ui"
)

func main() {
	// todo - this should be in os-specific config directory
	dbPath := "./planner.db"

	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	if err != nil {
		log.Fatalf("Error adding task: %v", err)
	}

	p := tea.NewProgram(ui.InitialModel(db))
	if err := p.Start(); err != nil {
		fmt.Printf("Error running program: %v", err)
		panic(err)
	}
}
