package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/themanselmo/burrow/internal/locale"
	"github.com/themanselmo/burrow/internal/storage"
	"github.com/themanselmo/burrow/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := locale.Load(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to load locale:", err)
		os.Exit(1)
	}

	state, err := storage.LoadState()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load state:", err)
		os.Exit(1)
	}

	app := ui.NewAppModel(state)
	p := tea.NewProgram(app, tea.WithAltScreen())

	// Forward SIGTERM as an abrupt quit so the pet reacts before exit.
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM)
		<-ch
		p.Send(ui.QuitMsg{Abrupt: true})
	}()

	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
