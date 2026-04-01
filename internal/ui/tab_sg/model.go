package tab_sg

import (
	"context"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	internalaws "tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

// viewState tracks the SG tab's internal view mode.
type viewState int

const (
	vsTable viewState = iota
	vsSearch
	vsActionMenu
	vsDetail
)

// sgsLoadedMsg is returned when security groups are fetched.
type sgsLoadedMsg struct {
	sgs []internalaws.SecurityGroup
	err error
}

// naclsLoadedMsg is returned when network ACLs are fetched.
type naclsLoadedMsg struct {
	nacls []internalaws.NetworkACL
	err   error
}

// detailKind distinguishes which detail overlay is active.
type detailKind int

const (
	detailInbound detailKind = iota
	detailOutbound
)

// Action represents a menu action.
type Action struct {
	Key   string
	Label string
}

// ActionMenuModel manages the action menu state.
type ActionMenuModel struct {
	Active  bool
	Title   string
	ID      string
	Actions []Action
	Cursor  int
}

func newSGActionMenu(sg internalaws.SecurityGroup) ActionMenuModel {
	name := sg.Name
	if name == "" {
		name = sg.ID
	}
	return ActionMenuModel{
		Active: true,
		Title:  name,
		ID:     sg.ID,
		Actions: []Action{
			{Key: "inbound", Label: "Inbound Rules"},
			{Key: "outbound", Label: "Outbound Rules"},
		},
		Cursor: 0,
	}
}

func newNACLActionMenu(nacl internalaws.NetworkACL) ActionMenuModel {
	name := nacl.Name
	if name == "" {
		name = nacl.ID
	}
	return ActionMenuModel{
		Active: true,
		Title:  name,
		ID:     nacl.ID,
		Actions: []Action{
			{Key: "inbound", Label: "Inbound Rules"},
			{Key: "outbound", Label: "Outbound Rules"},
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
	b.WriteString(fmt.Sprintf("  %s (%s)\n", a.Title, a.ID))
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

// SearchModel manages the search input state.
type SearchModel struct {
	Query  string
	Active bool
}

func (s *SearchModel) Insert(char rune) {
	s.Query += string(char)
}

func (s *SearchModel) Backspace() {
	if len(s.Query) > 0 {
		s.Query = s.Query[:len(s.Query)-1]
	}
}

func (s *SearchModel) Clear() {
	s.Query = ""
	s.Active = false
}

func (s *SearchModel) Render(width int) string {
	if !s.Active {
		return ""
	}
	prompt := shared.SearchPromptStyle.Render(" /")
	return lipgloss.NewStyle().Width(width).Render(
		fmt.Sprintf("%s %s█", prompt, s.Query),
	)
}

// SGModel implements the shared.TabModel interface for the SG tab.
type SGModel struct {
	viewState viewState
	loading   bool
	err       error

	// Mode: "sg" or "nacl"
	mode string

	// SG data
	sgs         []internalaws.SecurityGroup
	filteredSGs []internalaws.SecurityGroup
	sgLoaded    bool

	// NACL data
	nacls         []internalaws.NetworkACL
	filteredNACLs []internalaws.NetworkACL
	naclLoaded    bool

	cursor int

	search     SearchModel
	actionMenu ActionMenuModel

	// Detail overlay
	showDetail bool
	detailKind detailKind
	// Reference to the selected item for detail rendering
	selectedSG   internalaws.SecurityGroup
	selectedNACL internalaws.NetworkACL
}

// New creates a new SGModel.
func New() *SGModel {
	return &SGModel{
		viewState: vsTable,
		loading:   true,
		mode:      "sg",
	}
}

func (m *SGModel) Init(s *shared.SharedState) tea.Cmd {
	m.loading = true
	m.err = nil
	m.sgLoaded = false
	m.naclLoaded = false
	if m.mode == "sg" {
		return m.loadSGs(s)
	}
	return m.loadNACLs(s)
}

func (m *SGModel) Update(msg tea.Msg, s *shared.SharedState) (shared.TabModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil

	case sgsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.sgs = msg.sgs
		m.sgLoaded = true
		m.applyFilters()
		return m, nil

	case naclsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.nacls = msg.nacls
		m.naclLoaded = true
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

func (m *SGModel) View(s *shared.SharedState) string {
	var sections []string

	// Status bar
	sections = append(sections, renderStatusBar(s.Profile, s.Region, m.mode, m.itemCount(), s.Width))

	// Search bar (if active)
	if m.search.Active {
		sections = append(sections, m.search.Render(s.Width))
	}

	// Main content
	if m.loading {
		label := "Security Groups"
		if m.mode == "nacl" {
			label = "Network ACLs"
		}
		sections = append(sections, lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render(
			fmt.Sprintf("Loading %s...", label),
		))
	} else if m.err != nil {
		sections = append(sections, lipgloss.NewStyle().Width(s.Width).Padding(1, 2).Render(
			shared.ErrorStyle.Render(fmt.Sprintf("Error: %v\n\nPress R to retry, p to change profile, r to change region", m.err)),
		))
	} else if m.itemCount() == 0 {
		label := "Security Groups"
		if m.mode == "nacl" {
			label = "Network ACLs"
		}
		sections = append(sections, lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render(
			fmt.Sprintf("No %s found in this region.", label),
		))
	} else {
		tableHeight := s.Height
		if m.search.Active {
			tableHeight--
		}
		if m.mode == "sg" {
			columns := SGColumnsForWidth(s.Width)
			sections = append(sections, RenderSGTable(m.filteredSGs, columns, m.cursor, s.Width, tableHeight))
		} else {
			columns := NACLColumnsForWidth(s.Width)
			sections = append(sections, RenderNACLTable(m.filteredNACLs, columns, m.cursor, s.Width, tableHeight))
		}
	}

	// Overlay
	overlay := ""
	switch {
	case m.showDetail:
		if m.mode == "sg" {
			overlay = RenderSGRules(m.selectedSG, m.detailKind)
		} else {
			overlay = RenderNACLRules(m.selectedNACL, m.detailKind)
		}
	case m.actionMenu.Active:
		overlay = m.actionMenu.Render(s.Width)
	}

	view := strings.Join(sections, "\n")
	if overlay != "" {
		view += "\n" + shared.PlaceOverlay(s.Width, overlay)
	}

	return view
}

func (m *SGModel) ShortHelp() string {
	switch m.viewState {
	case vsSearch:
		return helpLine("Esc", "Cancel")
	case vsActionMenu:
		return helpLine("↑↓", "Navigate", "Enter", "Select", "Esc", "Cancel")
	case vsDetail:
		return helpLine("Esc", "Close")
	default:
		modeLabel := "NACL"
		if m.mode == "nacl" {
			modeLabel = "SG"
		}
		return helpLine("↑↓", "Navigate", "Enter", "Actions", "/", "Search", "f", modeLabel, "R", "Refresh")
	}
}

// --- Internal update handlers ---

func (m *SGModel) updateTable(msg tea.Msg, s *shared.SharedState) (shared.TabModel, tea.Cmd) {
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
		if m.cursor < m.itemCount()-1 {
			m.cursor++
		}
	case "enter":
		if m.cursor < m.itemCount() {
			if m.mode == "sg" {
				m.actionMenu = newSGActionMenu(m.filteredSGs[m.cursor])
			} else {
				m.actionMenu = newNACLActionMenu(m.filteredNACLs[m.cursor])
			}
			m.viewState = vsActionMenu
		}
	case "/":
		m.viewState = vsSearch
		m.search.Active = true
		m.search.Query = ""
	case "f":
		return m.toggleMode(s)
	case "R":
		m.loading = true
		m.err = nil
		if m.mode == "sg" {
			return m, m.loadSGs(s)
		}
		return m, m.loadNACLs(s)
	}

	return m, nil
}

func (m *SGModel) updateSearch(msg tea.Msg, _ *shared.SharedState) (shared.TabModel, tea.Cmd) {
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

func (m *SGModel) updateActionMenu(msg tea.Msg, _ *shared.SharedState) (shared.TabModel, tea.Cmd) {
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
		action := m.actionMenu.Selected()
		m.actionMenu.Active = false
		m.viewState = vsDetail
		m.showDetail = true
		switch action {
		case "inbound":
			m.detailKind = detailInbound
		case "outbound":
			m.detailKind = detailOutbound
		}
		// Capture the selected item for detail rendering
		if m.mode == "sg" && m.cursor < len(m.filteredSGs) {
			m.selectedSG = m.filteredSGs[m.cursor]
		} else if m.mode == "nacl" && m.cursor < len(m.filteredNACLs) {
			m.selectedNACL = m.filteredNACLs[m.cursor]
		}
	}
	return m, nil
}

func (m *SGModel) updateDetail(msg tea.Msg, _ *shared.SharedState) (shared.TabModel, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		if keyMsg.String() == "esc" {
			m.showDetail = false
			m.viewState = vsTable
		}
	}
	return m, nil
}

// --- Helpers ---

func (m *SGModel) itemCount() int {
	if m.mode == "sg" {
		return len(m.filteredSGs)
	}
	return len(m.filteredNACLs)
}

func (m *SGModel) toggleMode(s *shared.SharedState) (shared.TabModel, tea.Cmd) {
	if m.mode == "sg" {
		m.mode = "nacl"
	} else {
		m.mode = "sg"
	}
	m.cursor = 0
	m.search.Clear()

	// Load data for new mode if not cached
	if m.mode == "sg" && !m.sgLoaded {
		m.loading = true
		m.err = nil
		return m, m.loadSGs(s)
	}
	if m.mode == "nacl" && !m.naclLoaded {
		m.loading = true
		m.err = nil
		return m, m.loadNACLs(s)
	}

	m.applyFilters()
	return m, nil
}

func (m *SGModel) applyFilters() {
	if m.mode == "sg" {
		result := m.sgs
		if m.search.Query != "" {
			q := strings.ToLower(m.search.Query)
			var filtered []internalaws.SecurityGroup
			for _, sg := range result {
				if strings.Contains(strings.ToLower(sg.Name), q) ||
					strings.Contains(strings.ToLower(sg.ID), q) {
					filtered = append(filtered, sg)
				}
			}
			result = filtered
		}
		m.filteredSGs = result
		if m.cursor >= len(m.filteredSGs) {
			m.cursor = len(m.filteredSGs) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}
	} else {
		result := m.nacls
		if m.search.Query != "" {
			q := strings.ToLower(m.search.Query)
			var filtered []internalaws.NetworkACL
			for _, nacl := range result {
				if strings.Contains(strings.ToLower(nacl.Name), q) ||
					strings.Contains(strings.ToLower(nacl.ID), q) {
					filtered = append(filtered, nacl)
				}
			}
			result = filtered
		}
		m.filteredNACLs = result
		if m.cursor >= len(m.filteredNACLs) {
			m.cursor = len(m.filteredNACLs) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}
	}
}

func (m *SGModel) loadSGs(s *shared.SharedState) tea.Cmd {
	profile := s.Profile
	region := s.Region
	return func() tea.Msg {
		ctx := context.Background()
		clients, err := internalaws.NewClients(ctx, profile, region)
		if err != nil {
			return sgsLoadedMsg{err: err}
		}
		sgs, err := internalaws.FetchSecurityGroups(ctx, clients.EC2)
		if err != nil {
			return sgsLoadedMsg{err: err}
		}
		return sgsLoadedMsg{sgs: sgs}
	}
}

func (m *SGModel) loadNACLs(s *shared.SharedState) tea.Cmd {
	profile := s.Profile
	region := s.Region
	return func() tea.Msg {
		ctx := context.Background()
		clients, err := internalaws.NewClients(ctx, profile, region)
		if err != nil {
			return naclsLoadedMsg{err: err}
		}
		nacls, err := internalaws.FetchNetworkACLs(ctx, clients.EC2)
		if err != nil {
			return naclsLoadedMsg{err: err}
		}
		return naclsLoadedMsg{nacls: nacls}
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
