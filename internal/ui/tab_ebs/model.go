package tab_ebs

import (
	"context"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	internalaws "tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

type viewState int

const (
	vsTable viewState = iota
	vsSearch
	vsActionMenu
	vsDetail
)

type volumesLoadedMsg struct {
	volumes []internalaws.Volume
	err     error
}

type Action struct {
	Key   string
	Label string
}

type ActionMenuModel struct {
	Active  bool
	Volume  internalaws.Volume
	Actions []Action
	Cursor  int
}

func newActionMenu(v internalaws.Volume) ActionMenuModel {
	return ActionMenuModel{
		Active: true,
		Volume: v,
		Actions: []Action{
			{Key: "detail", Label: "Volume Details"},
		},
		Cursor: 0,
	}
}

func (a *ActionMenuModel) MoveUp() {
	if a.Cursor > 0 {
		a.Cursor--
	}
}

func (a *ActionMenuModel) MoveDown() {
	if a.Cursor < len(a.Actions)-1 {
		a.Cursor++
	}
}

func (a *ActionMenuModel) Selected() string {
	if a.Cursor < len(a.Actions) {
		return a.Actions[a.Cursor].Key
	}
	return ""
}

func (a *ActionMenuModel) Render(width int) string {
	if !a.Active {
		return ""
	}
	var b strings.Builder
	name := a.Volume.Name
	if name == "" {
		name = a.Volume.ID
	}
	b.WriteString(fmt.Sprintf("  %s (%s)\n", name, a.Volume.ID))
	b.WriteString("  ─────────────────────────\n")
	for i, action := range a.Actions {
		cursor := "  "
		if i == a.Cursor {
			cursor = "▸ "
		}
		b.WriteString(fmt.Sprintf("  %s%s\n", cursor, action.Label))
	}
	b.WriteString("\n  Enter: select  Esc: cancel")
	return shared.RenderOverlay(b.String())
}

type SearchModel struct {
	Query  string
	Active bool
}

func (s *SearchModel) Insert(char rune) { s.Query += string(char) }
func (s *SearchModel) Backspace() {
	if len(s.Query) > 0 {
		s.Query = s.Query[:len(s.Query)-1]
	}
}
func (s *SearchModel) Clear() { s.Query = ""; s.Active = false }
func (s *SearchModel) Render(width int) string {
	if !s.Active {
		return ""
	}
	prompt := shared.SearchPromptStyle.Render(" /")
	return lipgloss.NewStyle().Width(width).Render(fmt.Sprintf("%s %s█", prompt, s.Query))
}

// EBSModel implements shared.TabModel for the EBS tab.
type EBSModel struct {
	viewState  viewState
	loading    bool
	err        error
	volumes    []internalaws.Volume
	filtered   []internalaws.Volume
	cursor     int
	search     SearchModel
	actionMenu ActionMenuModel
	showDetail bool
}

func New() *EBSModel {
	return &EBSModel{viewState: vsTable, loading: true}
}

func (m *EBSModel) Init(s *shared.SharedState) tea.Cmd {
	m.loading = true
	m.err = nil
	return m.loadData(s)
}

func (m *EBSModel) Update(msg tea.Msg, s *shared.SharedState) (shared.TabModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case volumesLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.volumes = msg.volumes
		m.applyFilters()
		return m, nil
	}

	switch m.viewState {
	case vsSearch:
		return m.updateSearch(msg, s)
	case vsActionMenu:
		return m.updateActionMenu(msg, s)
	case vsDetail:
		return m.updateDetail(msg, s)
	default:
		return m.updateTable(msg, s)
	}
}

func (m *EBSModel) View(s *shared.SharedState) string {
	var sections []string
	sections = append(sections, renderStatusBar(s.Profile, s.Region, len(m.filtered), s.Width))
	if m.search.Active {
		sections = append(sections, m.search.Render(s.Width))
	}
	if m.loading {
		sections = append(sections, lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("Loading EBS Volumes..."))
	} else if m.err != nil {
		sections = append(sections, lipgloss.NewStyle().Width(s.Width).Padding(1, 2).Render(
			shared.ErrorStyle.Render(fmt.Sprintf("Error: %v\n\nPress R to retry, p to change profile, r to change region", m.err)),
		))
	} else if len(m.filtered) == 0 {
		sections = append(sections, lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("No EBS Volumes found in this region."))
	} else {
		columns := ColumnsForWidth(s.Width)
		tableHeight := s.Height
		if m.search.Active {
			tableHeight--
		}
		sections = append(sections, RenderTable(m.filtered, columns, m.cursor, s.Width, tableHeight))
	}

	overlay := ""
	switch {
	case m.showDetail:
		if m.cursor < len(m.filtered) {
			overlay = RenderVolumeDetail(m.filtered[m.cursor])
		}
	case m.actionMenu.Active:
		overlay = m.actionMenu.Render(s.Width)
	}

	view := strings.Join(sections, "\n")
	if overlay != "" {
		view = shared.PlaceOverlay(s.Width, s.Height, overlay)
	}
	return view
}

func (m *EBSModel) ShortHelp() string {
	switch m.viewState {
	case vsSearch:
		return helpLine("Esc", "Cancel")
	case vsActionMenu:
		return helpLine("↑↓", "Navigate", "Enter", "Select", "Esc", "Cancel")
	case vsDetail:
		return helpLine("any key", "Close")
	default:
		return helpLine("↑↓", "Navigate", "Enter", "Actions", "/", "Search", "R", "Refresh")
	}
}

func (m *EBSModel) updateTable(msg tea.Msg, s *shared.SharedState) (shared.TabModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch keyMsg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}
	case "enter":
		if m.cursor < len(m.filtered) {
			m.actionMenu = newActionMenu(m.filtered[m.cursor])
			m.viewState = vsActionMenu
		}
	case "/":
		m.viewState = vsSearch
		m.search.Active = true
		m.search.Query = ""
	case "R":
		m.loading = true
		m.err = nil
		return m, m.loadData(s)
	}
	return m, nil
}

func (m *EBSModel) updateSearch(msg tea.Msg, s *shared.SharedState) (shared.TabModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch keyMsg.String() {
	case "esc":
		m.search.Clear()
		m.viewState = vsTable
		m.applyFilters()
	case "enter":
		m.viewState = vsTable
		m.search.Active = false
	case "backspace":
		m.search.Backspace()
		m.applyFilters()
	default:
		r := keyMsg.String()
		if len(r) == 1 {
			m.search.Insert(rune(r[0]))
			m.applyFilters()
			m.cursor = 0
		}
	}
	return m, nil
}

func (m *EBSModel) updateActionMenu(msg tea.Msg, s *shared.SharedState) (shared.TabModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch keyMsg.String() {
	case "esc":
		m.actionMenu.Active = false
		m.viewState = vsTable
	case "up", "k":
		m.actionMenu.MoveUp()
	case "down", "j":
		m.actionMenu.MoveDown()
	case "enter":
		if m.actionMenu.Selected() == "detail" {
			m.actionMenu.Active = false
			m.viewState = vsDetail
			m.showDetail = true
		}
	}
	return m, nil
}

func (m *EBSModel) updateDetail(msg tea.Msg, _ *shared.SharedState) (shared.TabModel, tea.Cmd) {
	if _, ok := msg.(tea.KeyPressMsg); ok {
		m.showDetail = false
		m.viewState = vsTable
	}
	return m, nil
}

func (m *EBSModel) applyFilters() {
	result := m.volumes
	if m.search.Query != "" {
		q := strings.ToLower(m.search.Query)
		var filtered []internalaws.Volume
		for _, v := range result {
			if strings.Contains(strings.ToLower(v.Name), q) ||
				strings.Contains(strings.ToLower(v.ID), q) ||
				strings.Contains(strings.ToLower(v.State), q) ||
				strings.Contains(strings.ToLower(v.Type), q) ||
				strings.Contains(strings.ToLower(v.AttachedTo), q) {
				filtered = append(filtered, v)
			}
		}
		result = filtered
	}
	m.filtered = result
	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m *EBSModel) loadData(s *shared.SharedState) tea.Cmd {
	profile := s.Profile
	region := s.Region
	return func() tea.Msg {
		ctx := context.Background()
		clients, err := internalaws.NewClients(ctx, profile, region)
		if err != nil {
			return volumesLoadedMsg{err: err}
		}
		volumes, err := internalaws.FetchVolumes(ctx, clients.EC2)
		if err != nil {
			return volumesLoadedMsg{err: err}
		}
		return volumesLoadedMsg{volumes: volumes}
	}
}

func helpLine(keyvals ...string) string {
	var s string
	for i := 0; i < len(keyvals)-1; i += 2 {
		if s != "" {
			s += "  "
		}
		s += fmt.Sprintf("%s: %s", shared.HelpKeyStyle.Render(keyvals[i]), keyvals[i+1])
	}
	return " " + s
}
