package tui

import (
	"context"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/olle/docops/internal/app"
	"github.com/olle/docops/internal/view"
)

type loadedMsg struct {
	snapshot view.Snapshot
	err      error
}

type model struct {
	ctx           context.Context
	source        app.SnapshotSource
	keys          keyMap
	list          list.Model
	detail        viewport.Model
	help          help.Model
	snapshot      view.Snapshot
	err           error
	width         int
	height        int
	showHelp      bool
	detailFocused bool
	kindFilter    string
	status        string
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	dimStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	warnStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	paneStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(0, 1)
)

var kindFilters = []struct {
	token string
	label string
}{
	{token: "", label: "all"},
	{token: "asyncapi", label: "AsyncAPI"},
	{token: "avro", label: "Avro"},
	{token: "jsonschema", label: "JSON Schema"},
	{token: "openapi", label: "OpenAPI"},
}

func newModel(ctx context.Context, source app.SnapshotSource, snapshot view.Snapshot, err error) *model {
	keys := newKeyMap()
	items := list.New(nil, list.NewDefaultDelegate(), 80, 20)
	items.Title = "Specifications"
	items.SetShowTitle(true)
	items.SetShowStatusBar(true)
	items.SetStatusBarItemName("specification", "specifications")
	items.SetShowHelp(false)
	items.DisableQuitKeybindings()

	m := &model{
		ctx:      ctx,
		source:   source,
		keys:     keys,
		list:     items,
		detail:   viewport.New(),
		help:     help.New(),
		snapshot: snapshot,
		err:      err,
		width:    100,
		height:   40,
	}
	m.applyItems()
	m.layout()
	m.status = m.statusText()
	m.refreshDetail()
	return m
}

func load(ctx context.Context, source app.SnapshotSource) tea.Cmd {
	return func() tea.Msg {
		snapshot, err := source.Snapshot(ctx)
		return loadedMsg{snapshot: snapshot, err: err}
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typed := msg.(type) {
	case loadedMsg:
		m.snapshot = typed.snapshot
		m.err = typed.err
		m.applyItems()
		m.status = m.statusText()
		m.refreshDetail()
		return m, nil

	case tea.WindowSizeMsg:
		m.width = typed.Width
		m.height = typed.Height
		m.layout()
		m.refreshDetail()
		return m, nil

	case tea.KeyPressMsg:
		if !m.list.SettingFilter() {
			switch {
			case key.Matches(typed, m.keys.Quit):
				return m, tea.Quit
			case key.Matches(typed, m.keys.Reload):
				m.status = "reloading…"
				return m, load(m.ctx, m.source)
			case key.Matches(typed, m.keys.Help):
				m.showHelp = !m.showHelp
				m.layout()
				return m, nil
			case key.Matches(typed, m.keys.FilterKind):
				m.cycleKindFilter()
				m.status = m.statusText()
				return m, nil
			case key.Matches(typed, m.keys.TogglePane):
				m.detailFocused = !m.detailFocused
				m.status = m.statusText()
				return m, nil
			}
		}
	}

	cursor := m.list.Cursor()
	var cmds []tea.Cmd
	if m.detailFocused {
		if _, isKey := msg.(tea.KeyPressMsg); isKey {
			var cmd tea.Cmd
			m.detail, cmd = m.detail.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}
	} else {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}
	if _, isKey := msg.(tea.KeyPressMsg); !isKey {
		var cmd tea.Cmd
		m.detail, cmd = m.detail.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.list.Cursor() != cursor {
		m.refreshDetail()
	}
	return m, tea.Batch(cmds...)
}

func (m *model) View() tea.View {
	header := titleStyle.Render(m.snapshot.Greeting)
	if m.status != "" {
		header += "  " + dimStyle.Render(m.status)
	}
	if len(m.snapshot.Warnings) > 0 {
		header += "  " + warnStyle.Render("! "+strings.Join(m.snapshot.Warnings, " | "))
	}

	left := paneStyle.Width(m.list.Width()).Height(m.list.Height()).Render(m.list.View())
	right := paneStyle.Width(m.detail.Width()).Height(m.detail.Height()).Render(m.detail.View())

	content := lipgloss.JoinVertical(lipgloss.Left, header, lipgloss.JoinHorizontal(lipgloss.Top, left, right))
	if m.showHelp {
		content += "\n" + m.help.View(m.keys)
	}

	view := tea.NewView(content)
	view.AltScreen = true
	view.WindowTitle = "docops"
	return view
}

func (m *model) applyItems() {
	query := m.list.FilterValue()
	items := make([]list.Item, 0, len(m.snapshot.Documents))
	for _, doc := range view.Summaries(m.snapshot) {
		if m.kindFilter != "" && doc.Kind != m.kindFilter {
			continue
		}
		items = append(items, docItem{doc: doc})
	}
	_ = m.list.SetItems(items)
	if query != "" {
		m.list.SetFilterText(query)
	}
}

func (m *model) cycleKindFilter() {
	current := 0
	for i, filter := range kindFilters {
		if filter.token == m.kindFilter {
			current = i
			break
		}
	}
	m.kindFilter = kindFilters[(current+1)%len(kindFilters)].token
	m.applyItems()
	m.refreshDetail()
}

func (m *model) refreshDetail() {
	item, ok := m.list.SelectedItem().(docItem)
	if !ok {
		m.detail.SetContent("No specifications found.\n\nPress ? for help.")
	} else {
		m.detail.SetContent(detailContent(item.doc))
	}
	m.detail.GotoTop()
}

func (m *model) statusText() string {
	focus := "list"
	if m.detailFocused {
		focus = "detail"
	}
	label := "all"
	for _, filter := range kindFilters {
		if filter.token == m.kindFilter {
			label = filter.label
		}
	}
	return "kind: " + label + " · focus: " + focus
}

func (m *model) layout() {
	width := m.width
	if width < 1 {
		width = 80
	}
	height := m.height
	if height < 1 {
		height = 24
	}
	m.help.SetWidth(width)

	if width < 80 {
		innerHeight := (height - 8) / 2
		if innerHeight < 3 {
			innerHeight = 3
		}
		m.list.SetSize(width-2, innerHeight)
		m.detail.SetWidth(width - 4)
		m.detail.SetHeight(innerHeight)
		return
	}

	listWidth := width/2 - 2
	detailWidth := width - listWidth - 4
	innerHeight := height - 5
	if innerHeight < 3 {
		innerHeight = 3
	}
	m.list.SetSize(listWidth, innerHeight)
	m.detail.SetWidth(detailWidth)
	m.detail.SetHeight(innerHeight)
}
