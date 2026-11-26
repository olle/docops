package tui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
)

type keyMap struct {
	Quit       key.Binding
	Reload     key.Binding
	Help       key.Binding
	FilterKind key.Binding
	TogglePane key.Binding
	Up         key.Binding
	Down       key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Quit:       key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc", "quit")),
		Reload:     key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reload")),
		Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		FilterKind: key.NewBinding(key.WithKeys("f", "tab"), key.WithHelp("f/tab", "filter kind")),
		TogglePane: key.NewBinding(key.WithKeys("right", "left", "enter"), key.WithHelp("←/→", "focus pane")),
		Up:         key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:       key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.FilterKind, k.TogglePane, k.Reload, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.TogglePane},
		{k.FilterKind, k.Reload, k.Help, k.Quit},
	}
}

var _ help.KeyMap = keyMap{}
