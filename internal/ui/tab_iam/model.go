package tab_iam

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

type usersLoadedMsg struct {
	users    []internalaws.IAMUser
	identity *internalaws.CallerIdentity
	err      error
}

type Action struct {
	Key   string
	Label string
}

type ActionMenuModel struct {
	Active  bool
	User    internalaws.IAMUser
	Actions []Action
	Cursor  int
}

func newActionMenu(u internalaws.IAMUser) ActionMenuModel {
	return ActionMenuModel{
		Active: true,
		User:   u,
		Actions: []Action{
			{Key: "detail", Label: "User Details"},
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
	b.WriteString(fmt.Sprintf("  %s\n", a.User.UserName))
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

// IAMModel implements shared.TabModel for the IAM tab.
type IAMModel struct {
	viewState  viewState
	loading    bool
	err        error
	users      []internalaws.IAMUser
	filtered   []internalaws.IAMUser
	identity   *internalaws.CallerIdentity
	cursor     int
	search     SearchModel
	actionMenu ActionMenuModel
	showDetail bool
}

func New() *IAMModel {
	return &IAMModel{viewState: vsTable, loading: true}
}

func (m *IAMModel) Init(s *shared.SharedState) tea.Cmd {
	m.loading = true
	m.err = nil
	return m.loadData(s)
}

func (m *IAMModel) Update(msg tea.Msg, s *shared.SharedState) (shared.TabModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case usersLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.users = msg.users
		m.identity = msg.identity
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

func (m *IAMModel) View(s *shared.SharedState) string {
	var sections []string
	sections = append(sections, renderStatusBar(s.Profile, s.Region, len(m.filtered), m.identity, s.Width))
	if m.search.Active {
		sections = append(sections, m.search.Render(s.Width))
	}
	if m.loading {
		sections = append(sections, lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("Loading IAM Users..."))
	} else if m.err != nil {
		sections = append(sections, lipgloss.NewStyle().Width(s.Width).Padding(1, 2).Render(
			shared.ErrorStyle.Render(fmt.Sprintf("Error: %v\n\nPress R to retry, p to change profile, r to change region", m.err)),
		))
	} else if len(m.filtered) == 0 {
		sections = append(sections, lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("No IAM Users found."))
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
			overlay = RenderUserDetail(m.filtered[m.cursor])
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

func (m *IAMModel) ShortHelp() string {
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

func (m *IAMModel) IsEditing() bool {
	return m.viewState == vsSearch
}

func (m *IAMModel) updateTable(msg tea.Msg, s *shared.SharedState) (shared.TabModel, tea.Cmd) {
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

func (m *IAMModel) updateSearch(msg tea.Msg, s *shared.SharedState) (shared.TabModel, tea.Cmd) {
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

func (m *IAMModel) updateActionMenu(msg tea.Msg, s *shared.SharedState) (shared.TabModel, tea.Cmd) {
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

func (m *IAMModel) updateDetail(msg tea.Msg, _ *shared.SharedState) (shared.TabModel, tea.Cmd) {
	if _, ok := msg.(tea.KeyPressMsg); ok {
		m.showDetail = false
		m.viewState = vsTable
	}
	return m, nil
}

func (m *IAMModel) applyFilters() {
	result := m.users
	if m.search.Query != "" {
		q := strings.ToLower(m.search.Query)
		var filtered []internalaws.IAMUser
		for _, u := range result {
			if strings.Contains(strings.ToLower(u.UserName), q) ||
				strings.Contains(strings.ToLower(u.ARN), q) ||
				strings.Contains(strings.ToLower(u.UserID), q) {
				filtered = append(filtered, u)
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

func (m *IAMModel) loadData(s *shared.SharedState) tea.Cmd {
	profile := s.Profile
	region := s.Region
	return func() tea.Msg {
		ctx := context.Background()
		clients, err := internalaws.NewClients(ctx, profile, region)
		if err != nil {
			return usersLoadedMsg{err: err}
		}
		// Fetch caller identity
		identity, _ := internalaws.FetchCurrentIdentity(ctx, clients.STS)
		users, err := internalaws.FetchIAMUsers(ctx, clients.IAM)
		if err != nil {
			return usersLoadedMsg{err: err}
		}
		return usersLoadedMsg{users: users, identity: identity}
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
