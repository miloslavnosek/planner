package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
	"path/filepath"
	"planner/database"
	"planner/ui"
	"runtime"
)

func getDatabasePath() string {
	var dbDir string
	var dbPath string

	if os.Getenv("PLANNER_ENV") == "dev" {
		dbDir = "dev_db"

		dbPath = filepath.Join(dbDir, "planner.dev.db")
	} else {
		switch runtime.GOOS {
		case "linux":
			dbDir = filepath.Join(os.Getenv("HOME"), ".local", "share", "io.github.miloslavnosek.planner")
			os.MkdirAll(dbDir, 0755)

		case "darwin":
			dbDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "io.github.miloslavnosek.planner")
			os.MkdirAll(dbDir, 0755)

		case "windows":
			dbDir = filepath.Join(os.Getenv("LOCALAPPDATA"), "io.github.miloslavnosek.planner")
			os.MkdirAll(dbDir, 0755)

		default:
			fmt.Println("Unsupported operating system")
			os.Exit(1)
		}

		dbPath = filepath.Join(dbDir, "planner.db")
	}

	return dbPath
}

func main() {
	dbPath := getDatabasePath()

	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	p := tea.NewProgram(ui.InitialModel(db))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
