package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// WatchlistModel represents the model for the watchlist TUI
type WatchlistModel struct {
	choices  []string         // items in the watchlist
	cursor   int              // which item the cursor is pointing at
	selected map[int]struct{} // which items are selected
	done     bool             // whether the user is done editing
	saved    bool             // whether the changes have been saved
}

// NewWatchlistModel creates a new watchlist model
func NewWatchlistModel(watchlist []string) WatchlistModel {
	return WatchlistModel{
		choices:  watchlist,
		selected: make(map[int]struct{}),
	}
}

// Init initializes the model
func (m WatchlistModel) Init() tea.Cmd {
	return nil
}

// Update updates the model based on messages
func (m WatchlistModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			// Save changes
			m.saved = true
			m.done = true
			return m, tea.Quit

		case "ctrl+c", "q":
			// Quit without saving
			m.done = true
			return m, tea.Quit

		case "up", "k":
			// Move cursor up
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			// Move cursor down
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			// Toggle selection
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

// View renders the model
func (m WatchlistModel) View() string {
	// The header
	s := "Watchlist\n\n"

	// Iterate over the choices
	for i, choice := range m.choices {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := "x" // not selected
		if _, ok := m.selected[i]; ok {
			checked = " " // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress s to save.\nPress q to quit without saving.\n"

	return s
}

// GetRemainingChoices returns the choices that were not selected
func (m WatchlistModel) GetRemainingChoices() []string {
	var remaining []string

	for i, choice := range m.choices {
		if _, ok := m.selected[i]; !ok {
			remaining = append(remaining, choice)
		}
	}

	return remaining
}

// IsSaved returns whether the changes have been saved
func (m WatchlistModel) IsSaved() bool {
	return m.saved
}

// IsDone returns whether the user is done editing
func (m WatchlistModel) IsDone() bool {
	return m.done
}

// RunWatchlistEditor runs the watchlist editor and returns the updated watchlist
func RunWatchlistEditor(watchlist []string) ([]string, bool) {
	p := tea.NewProgram(NewWatchlistModel(watchlist))
	model, err := p.Run()
	if err != nil {
		fmt.Printf("Error running watchlist editor: %v\n", err)
		return watchlist, false
	}

	watchlistModel, ok := model.(WatchlistModel)
	if !ok {
		fmt.Println("Error: could not convert model to WatchlistModel")
		return watchlist, false
	}

	if !watchlistModel.IsSaved() {
		return watchlist, false
	}

	return watchlistModel.GetRemainingChoices(), true
}
