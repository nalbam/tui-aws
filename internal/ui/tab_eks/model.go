package tab_eks

import (
	"context"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	internalaws "tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

// viewState tracks the drill-down level in the EKS hierarchy.
type viewState int

const (
	vsClusterList   viewState = iota // top-level cluster table
	vsClusterSearch                  // search within cluster list
	vsClusterAction                  // action menu on a cluster
	vsClusterDetail                  // cluster detail overlay
	vsNamespaceList                  // namespace list for a cluster
	vsResourceMenu                   // Pods/Deployments/Services selection
	vsPodList                        // pod list for namespace
	vsPodAction                      // action menu on a pod
	vsPodDetail                      // pod detail overlay
	vsPodLogs                        // pod log viewer overlay
	vsDeployList                     // deployment list for namespace
	vsDeployDetail                   // deployment detail overlay
	vsServiceList                    // service list for namespace
	vsServiceDetail                  // service detail overlay
	vsNodeGroupList                  // AWS node groups for cluster
	vsNodeList                       // K8s nodes for cluster
	vsNodeDetail                     // node detail overlay
)

// --- async messages ---

type clustersLoadedMsg struct {
	clusters []internalaws.EKSCluster
	err      error
}

type nodeGroupsLoadedMsg struct {
	nodeGroups []internalaws.EKSNodeGroup
	err        error
}

type namespacesLoadedMsg struct {
	namespaces []internalaws.K8sNamespace
	err        error
}

type podsLoadedMsg struct {
	pods []internalaws.K8sPod
	err  error
}

type deploymentsLoadedMsg struct {
	deployments []internalaws.K8sDeployment
	err         error
}

type servicesLoadedMsg struct {
	services []internalaws.K8sService
	err      error
}

type nodesLoadedMsg struct {
	nodes []internalaws.K8sNode
	err   error
}

type podLogsLoadedMsg struct {
	logs string
	err  error
}

// --- Action menu ---

type Action struct {
	Key   string
	Label string
}

type actionMenu struct {
	title   string
	actions []Action
	cursor  int
}

func (a *actionMenu) MoveUp() {
	if a.cursor > 0 {
		a.cursor--
	}
}

func (a *actionMenu) MoveDown() {
	if a.cursor < len(a.actions)-1 {
		a.cursor++
	}
}

func (a *actionMenu) Selected() string {
	if a.cursor < len(a.actions) {
		return a.actions[a.cursor].Key
	}
	return ""
}

func (a *actionMenu) Render() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  %s\n", a.title))
	b.WriteString("  ─────────────────────────\n")
	for i, act := range a.actions {
		cursor := "  "
		if i == a.cursor {
			cursor = "▸ "
		}
		b.WriteString(fmt.Sprintf("  %s%s\n", cursor, act.Label))
	}
	b.WriteString("\n  Enter: select  Esc: cancel")
	return shared.RenderOverlay(b.String())
}

// --- Search ---

type searchModel struct {
	query  string
	active bool
}

func (s *searchModel) Insert(char rune) { s.query += string(char) }
func (s *searchModel) Backspace() {
	if len(s.query) > 0 {
		s.query = s.query[:len(s.query)-1]
	}
}
func (s *searchModel) Clear() { s.query = ""; s.active = false }
func (s *searchModel) Render(width int) string {
	if !s.active {
		return ""
	}
	prompt := shared.SearchPromptStyle.Render(" /")
	return lipgloss.NewStyle().Width(width).Render(fmt.Sprintf("%s %s█", prompt, s.query))
}

// --- EKSModel ---

// EKSModel implements the shared.TabModel interface for the EKS tab.
type EKSModel struct {
	viewState viewState
	loading   bool
	err       error

	// Cluster level
	clusters      []internalaws.EKSCluster
	filtered      []internalaws.EKSCluster
	clusterCursor int
	search        searchModel

	// Namespace level
	namespaces     []internalaws.K8sNamespace
	namespaceCursor int

	// Pod level
	pods      []internalaws.K8sPod
	podCursor int

	// Deployment level
	deployments     []internalaws.K8sDeployment
	deployCursor    int

	// Service level
	services      []internalaws.K8sService
	serviceCursor int

	// Node Group level
	nodeGroups      []internalaws.EKSNodeGroup
	nodeGroupCursor int

	// Node level
	nodes      []internalaws.K8sNode
	nodeCursor int

	// Logs
	podLogs    string
	podLogsErr error

	// Action menu
	menu actionMenu

	// Selected references (breadcrumb)
	selectedCluster   *internalaws.EKSCluster
	selectedNamespace string
	selectedPod       *internalaws.K8sPod

	// Detail loading
	detailLoading bool
}

func New() *EKSModel {
	return &EKSModel{viewState: vsClusterList, loading: true}
}

func (m *EKSModel) Init(s *shared.SharedState) tea.Cmd {
	m.loading = true
	m.err = nil
	m.viewState = vsClusterList
	return m.loadClusters(s)
}

func (m *EKSModel) Update(msg tea.Msg, s *shared.SharedState) (shared.TabModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil

	case clustersLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.clusters = msg.clusters
		m.applyFilters()
		return m, nil

	case nodeGroupsLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.nodeGroups = msg.nodeGroups
		m.nodeGroupCursor = 0
		return m, nil

	case namespacesLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.namespaces = msg.namespaces
		m.namespaceCursor = 0
		// Default to "default" namespace if present
		for i, ns := range m.namespaces {
			if ns.Name == "default" {
				m.namespaceCursor = i
				break
			}
		}
		return m, nil

	case podsLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.pods = msg.pods
		m.podCursor = 0
		return m, nil

	case deploymentsLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.deployments = msg.deployments
		m.deployCursor = 0
		return m, nil

	case servicesLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.services = msg.services
		m.serviceCursor = 0
		return m, nil

	case nodesLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.nodes = msg.nodes
		m.nodeCursor = 0
		return m, nil

	case podLogsLoadedMsg:
		m.detailLoading = false
		m.podLogsErr = msg.err
		m.podLogs = msg.logs
		return m, nil
	}

	switch m.viewState {
	case vsClusterList:
		return m.updateClusterList(msg, s)
	case vsClusterSearch:
		return m.updateClusterSearch(msg, s)
	case vsClusterAction:
		return m.updateActionMenu(msg, s)
	case vsClusterDetail:
		return m.updateCloseOverlay(msg, s, vsClusterList)
	case vsNamespaceList:
		return m.updateNamespaceList(msg, s)
	case vsResourceMenu:
		return m.updateActionMenu(msg, s)
	case vsPodList:
		return m.updatePodList(msg, s)
	case vsPodAction:
		return m.updateActionMenu(msg, s)
	case vsPodDetail:
		return m.updateCloseOverlay(msg, s, vsPodList)
	case vsPodLogs:
		return m.updateCloseOverlay(msg, s, vsPodList)
	case vsDeployList:
		return m.updateDeployList(msg, s)
	case vsDeployDetail:
		return m.updateCloseOverlay(msg, s, vsDeployList)
	case vsServiceList:
		return m.updateServiceList(msg, s)
	case vsServiceDetail:
		return m.updateCloseOverlay(msg, s, vsServiceList)
	case vsNodeGroupList:
		return m.updateNodeGroupList(msg, s)
	case vsNodeList:
		return m.updateNodeList(msg, s)
	case vsNodeDetail:
		return m.updateCloseOverlay(msg, s, vsNodeList)
	default:
		return m, nil
	}
}

func (m *EKSModel) View(s *shared.SharedState) string {
	var sections []string

	// Status bar with breadcrumb
	sections = append(sections, m.renderStatusBar(s))

	if m.search.active {
		sections = append(sections, m.search.Render(s.Width))
	}

	if m.loading {
		sections = append(sections, lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("Loading EKS Clusters..."))
	} else if m.err != nil && m.viewState == vsClusterList {
		sections = append(sections, lipgloss.NewStyle().Width(s.Width).Padding(1, 2).Render(
			shared.ErrorStyle.Render(fmt.Sprintf("Error: %v\n\nPress R to retry, p to change profile, r to change region", m.err)),
		))
	} else {
		tableHeight := s.Height
		if m.search.active {
			tableHeight--
		}
		sections = append(sections, m.renderContent(s, tableHeight))
	}

	// Overlay
	overlay := m.renderOverlay(s)

	view := strings.Join(sections, "\n")
	if overlay != "" {
		view = shared.PlaceOverlay(s.Width, s.Height, overlay)
	}
	return view
}

func (m *EKSModel) ShortHelp() string {
	switch m.viewState {
	case vsClusterSearch:
		return helpLine("Esc", "Cancel")
	case vsClusterAction, vsResourceMenu, vsPodAction:
		return helpLine("↑↓", "Navigate", "Enter", "Select", "Esc", "Cancel")
	case vsClusterDetail, vsPodDetail, vsDeployDetail, vsServiceDetail, vsNodeDetail:
		return helpLine("Esc", "Close")
	case vsPodLogs:
		return helpLine("Esc", "Close")
	case vsNamespaceList:
		return helpLine("↑↓", "Navigate", "Enter", "Select", "Esc", "Back")
	case vsPodList:
		return helpLine("↑↓", "Navigate", "Enter", "Actions", "Esc", "Back")
	case vsDeployList:
		return helpLine("↑↓", "Navigate", "Enter", "Details", "Esc", "Back")
	case vsServiceList:
		return helpLine("↑↓", "Navigate", "Enter", "Details", "Esc", "Back")
	case vsNodeGroupList:
		return helpLine("↑↓", "Navigate", "Esc", "Back")
	case vsNodeList:
		return helpLine("↑↓", "Navigate", "Enter", "Details", "Esc", "Back")
	default:
		return helpLine("↑↓", "Navigate", "Enter", "Actions", "/", "Search", "R", "Refresh")
	}
}

func (m *EKSModel) IsEditing() bool {
	return m.viewState == vsClusterSearch
}

// --- Update handlers ---

func (m *EKSModel) updateClusterList(msg tea.Msg, s *shared.SharedState) (shared.TabModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch keyMsg.String() {
	case "up", "k":
		if m.clusterCursor > 0 {
			m.clusterCursor--
		}
	case "down", "j":
		if m.clusterCursor < len(m.filtered)-1 {
			m.clusterCursor++
		}
	case "enter":
		if m.clusterCursor < len(m.filtered) {
			c := m.filtered[m.clusterCursor]
			m.selectedCluster = &c
			m.menu = actionMenu{
				title: fmt.Sprintf("%s (%s)", c.Name, c.Status),
				actions: []Action{
					{Key: "namespaces", Label: "Namespaces (K8s)"},
					{Key: "nodegroups", Label: "Node Groups (AWS)"},
					{Key: "nodes", Label: "Nodes (K8s)"},
					{Key: "detail", Label: "Cluster Details"},
				},
			}
			m.viewState = vsClusterAction
		}
	case "/":
		m.viewState = vsClusterSearch
		m.search.active = true
		m.search.query = ""
	case "R":
		m.loading = true
		m.err = nil
		return m, m.loadClusters(s)
	}
	return m, nil
}

func (m *EKSModel) updateClusterSearch(msg tea.Msg, _ *shared.SharedState) (shared.TabModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch keyMsg.String() {
	case "esc":
		m.search.Clear()
		m.viewState = vsClusterList
		m.applyFilters()
	case "enter":
		m.viewState = vsClusterList
		m.search.active = false
	case "backspace":
		m.search.Backspace()
		m.applyFilters()
	default:
		r := keyMsg.String()
		if len(r) == 1 {
			m.search.Insert(rune(r[0]))
			m.applyFilters()
			m.clusterCursor = 0
		}
	}
	return m, nil
}

func (m *EKSModel) updateActionMenu(msg tea.Msg, s *shared.SharedState) (shared.TabModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}

	parentState := m.parentOfAction()

	switch keyMsg.String() {
	case "esc":
		m.viewState = parentState
		return m, nil
	case "up", "k":
		m.menu.MoveUp()
	case "down", "j":
		m.menu.MoveDown()
	case "enter":
		return m.executeAction(s)
	}
	return m, nil
}

func (m *EKSModel) updateNamespaceList(msg tea.Msg, s *shared.SharedState) (shared.TabModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch keyMsg.String() {
	case "esc":
		m.viewState = vsClusterList
		m.namespaces = nil
	case "up", "k":
		if m.namespaceCursor > 0 {
			m.namespaceCursor--
		}
	case "down", "j":
		if m.namespaceCursor < len(m.namespaces)-1 {
			m.namespaceCursor++
		}
	case "enter":
		if m.namespaceCursor < len(m.namespaces) {
			ns := m.namespaces[m.namespaceCursor]
			m.selectedNamespace = ns.Name
			m.menu = actionMenu{
				title: fmt.Sprintf("Namespace: %s", ns.Name),
				actions: []Action{
					{Key: "pods", Label: "Pods"},
					{Key: "deployments", Label: "Deployments"},
					{Key: "services", Label: "Services"},
				},
			}
			m.viewState = vsResourceMenu
		}
	}
	return m, nil
}

func (m *EKSModel) updatePodList(msg tea.Msg, s *shared.SharedState) (shared.TabModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch keyMsg.String() {
	case "esc":
		m.viewState = vsNamespaceList
		m.pods = nil
		m.selectedPod = nil
	case "up", "k":
		if m.podCursor > 0 {
			m.podCursor--
		}
	case "down", "j":
		if m.podCursor < len(m.pods)-1 {
			m.podCursor++
		}
	case "enter":
		if m.podCursor < len(m.pods) {
			p := m.pods[m.podCursor]
			m.selectedPod = &p
			actions := []Action{
				{Key: "pod_logs", Label: "Pod Logs (last 50 lines)"},
				{Key: "pod_detail", Label: "Pod Details"},
			}
			// Add per-container log options if multiple containers
			if len(p.Containers) > 1 {
				for _, c := range p.Containers {
					actions = append(actions, Action{
						Key:   "container_logs:" + c.Name,
						Label: fmt.Sprintf("Logs: %s", c.Name),
					})
				}
			}
			m.menu = actionMenu{
				title:   fmt.Sprintf("%s (%s)", p.Name, p.Status),
				actions: actions,
			}
			m.viewState = vsPodAction
		}
	}
	return m, nil
}

func (m *EKSModel) updateDeployList(msg tea.Msg, _ *shared.SharedState) (shared.TabModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch keyMsg.String() {
	case "esc":
		m.viewState = vsNamespaceList
		m.deployments = nil
	case "up", "k":
		if m.deployCursor > 0 {
			m.deployCursor--
		}
	case "down", "j":
		if m.deployCursor < len(m.deployments)-1 {
			m.deployCursor++
		}
	case "enter":
		if m.deployCursor < len(m.deployments) {
			m.viewState = vsDeployDetail
		}
	}
	return m, nil
}

func (m *EKSModel) updateServiceList(msg tea.Msg, _ *shared.SharedState) (shared.TabModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch keyMsg.String() {
	case "esc":
		m.viewState = vsNamespaceList
		m.services = nil
	case "up", "k":
		if m.serviceCursor > 0 {
			m.serviceCursor--
		}
	case "down", "j":
		if m.serviceCursor < len(m.services)-1 {
			m.serviceCursor++
		}
	case "enter":
		if m.serviceCursor < len(m.services) {
			m.viewState = vsServiceDetail
		}
	}
	return m, nil
}

func (m *EKSModel) updateNodeGroupList(msg tea.Msg, _ *shared.SharedState) (shared.TabModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch keyMsg.String() {
	case "esc":
		m.viewState = vsClusterList
		m.nodeGroups = nil
	case "up", "k":
		if m.nodeGroupCursor > 0 {
			m.nodeGroupCursor--
		}
	case "down", "j":
		if m.nodeGroupCursor < len(m.nodeGroups)-1 {
			m.nodeGroupCursor++
		}
	}
	return m, nil
}

func (m *EKSModel) updateNodeList(msg tea.Msg, _ *shared.SharedState) (shared.TabModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch keyMsg.String() {
	case "esc":
		m.viewState = vsClusterList
		m.nodes = nil
	case "up", "k":
		if m.nodeCursor > 0 {
			m.nodeCursor--
		}
	case "down", "j":
		if m.nodeCursor < len(m.nodes)-1 {
			m.nodeCursor++
		}
	case "enter":
		if m.nodeCursor < len(m.nodes) {
			m.viewState = vsNodeDetail
		}
	}
	return m, nil
}

func (m *EKSModel) updateCloseOverlay(msg tea.Msg, _ *shared.SharedState, backState viewState) (shared.TabModel, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		if keyMsg.String() == "esc" {
			m.viewState = backState
		}
	}
	return m, nil
}

// parentOfAction returns the view state to return to when Esc is pressed on a menu.
func (m *EKSModel) parentOfAction() viewState {
	switch m.viewState {
	case vsClusterAction:
		return vsClusterList
	case vsResourceMenu:
		return vsNamespaceList
	case vsPodAction:
		return vsPodList
	default:
		return vsClusterList
	}
}

// executeAction performs the selected action from the current menu.
func (m *EKSModel) executeAction(s *shared.SharedState) (shared.TabModel, tea.Cmd) {
	key := m.menu.Selected()

	switch key {
	// --- Cluster actions ---
	case "namespaces":
		m.viewState = vsNamespaceList
		m.detailLoading = true
		m.namespaces = nil
		m.err = nil
		return m, m.loadNamespaces(s)
	case "nodegroups":
		m.viewState = vsNodeGroupList
		m.detailLoading = true
		m.nodeGroups = nil
		m.err = nil
		return m, m.loadNodeGroups(s, m.selectedCluster.Name)
	case "nodes":
		m.viewState = vsNodeList
		m.detailLoading = true
		m.nodes = nil
		m.err = nil
		return m, m.loadNodes(s)
	case "detail":
		m.viewState = vsClusterDetail
		return m, nil

	// --- Resource menu actions ---
	case "pods":
		m.viewState = vsPodList
		m.detailLoading = true
		m.pods = nil
		m.err = nil
		return m, m.loadPods(s, m.selectedNamespace)
	case "deployments":
		m.viewState = vsDeployList
		m.detailLoading = true
		m.deployments = nil
		m.err = nil
		return m, m.loadDeployments(s, m.selectedNamespace)
	case "services":
		m.viewState = vsServiceList
		m.detailLoading = true
		m.services = nil
		m.err = nil
		return m, m.loadServices(s, m.selectedNamespace)

	// --- Pod actions ---
	case "pod_logs":
		m.viewState = vsPodLogs
		m.detailLoading = true
		m.podLogs = ""
		m.podLogsErr = nil
		// Use first container if available, otherwise empty (all containers)
		containerName := ""
		if len(m.selectedPod.Containers) > 0 {
			containerName = m.selectedPod.Containers[0].Name
		}
		return m, m.loadPodLogs(s, m.selectedPod.Namespace, m.selectedPod.Name, containerName)
	case "pod_detail":
		m.viewState = vsPodDetail
		return m, nil
	}

	// Container-specific logs: "container_logs:<name>"
	if strings.HasPrefix(key, "container_logs:") {
		containerName := strings.TrimPrefix(key, "container_logs:")
		m.viewState = vsPodLogs
		m.detailLoading = true
		m.podLogs = ""
		m.podLogsErr = nil
		return m, m.loadPodLogs(s, m.selectedPod.Namespace, m.selectedPod.Name, containerName)
	}

	return m, nil
}

// --- Rendering ---

func (m *EKSModel) renderStatusBar(s *shared.SharedState) string {
	profilePart := shared.StatusKeyStyle.Render("Profile: ") + s.Profile
	regionPart := shared.StatusKeyStyle.Render("Region: ") + s.Region

	breadcrumb := m.breadcrumb()
	content := fmt.Sprintf(" %s  |  %s  |  %s", profilePart, regionPart, breadcrumb)
	return shared.StatusBarStyle.Width(s.Width).Render(content)
}

func (m *EKSModel) breadcrumb() string {
	parts := []string{fmt.Sprintf("[%d Clusters]", len(m.filtered))}

	if m.selectedCluster != nil && m.viewState > vsClusterList && m.viewState != vsClusterSearch {
		parts = append(parts, m.selectedCluster.Name)
	}
	if m.selectedNamespace != "" && m.viewState >= vsPodList && m.viewState <= vsServiceDetail {
		parts = append(parts, "ns:"+m.selectedNamespace)
	}
	if m.selectedPod != nil && (m.viewState == vsPodAction || m.viewState == vsPodDetail || m.viewState == vsPodLogs) {
		parts = append(parts, m.selectedPod.Name)
	}

	return strings.Join(parts, " > ")
}

func (m *EKSModel) renderContent(s *shared.SharedState, tableHeight int) string {
	switch m.viewState {
	case vsClusterList, vsClusterSearch, vsClusterAction, vsClusterDetail:
		if len(m.filtered) == 0 {
			return lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("No EKS clusters found in this region.")
		}
		return renderClusterTable(m.filtered, m.clusterCursor, s.Width, tableHeight)

	case vsNamespaceList, vsResourceMenu:
		if m.detailLoading {
			return lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("Loading namespaces...")
		}
		if m.err != nil {
			return lipgloss.NewStyle().Width(s.Width).Padding(1, 2).Render(
				shared.ErrorStyle.Render(fmt.Sprintf("Error: %v\n\nPress Esc to go back", m.err)))
		}
		if len(m.namespaces) == 0 {
			return lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("No namespaces found. Press Esc to go back.")
		}
		return renderNamespaceTable(m.namespaces, m.namespaceCursor, s.Width, tableHeight)

	case vsPodList, vsPodAction, vsPodDetail, vsPodLogs:
		if m.detailLoading && len(m.pods) == 0 {
			return lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("Loading pods...")
		}
		if m.err != nil && len(m.pods) == 0 {
			return lipgloss.NewStyle().Width(s.Width).Padding(1, 2).Render(
				shared.ErrorStyle.Render(fmt.Sprintf("Error: %v\n\nPress Esc to go back", m.err)))
		}
		if len(m.pods) == 0 {
			return lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("No pods found. Press Esc to go back.")
		}
		return renderPodTable(m.pods, m.podCursor, s.Width, tableHeight)

	case vsDeployList, vsDeployDetail:
		if m.detailLoading {
			return lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("Loading deployments...")
		}
		if m.err != nil {
			return lipgloss.NewStyle().Width(s.Width).Padding(1, 2).Render(
				shared.ErrorStyle.Render(fmt.Sprintf("Error: %v\n\nPress Esc to go back", m.err)))
		}
		if len(m.deployments) == 0 {
			return lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("No deployments found. Press Esc to go back.")
		}
		return renderDeployTable(m.deployments, m.deployCursor, s.Width, tableHeight)

	case vsServiceList, vsServiceDetail:
		if m.detailLoading {
			return lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("Loading services...")
		}
		if m.err != nil {
			return lipgloss.NewStyle().Width(s.Width).Padding(1, 2).Render(
				shared.ErrorStyle.Render(fmt.Sprintf("Error: %v\n\nPress Esc to go back", m.err)))
		}
		if len(m.services) == 0 {
			return lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("No services found. Press Esc to go back.")
		}
		return renderServiceTable(m.services, m.serviceCursor, s.Width, tableHeight)

	case vsNodeGroupList:
		if m.detailLoading {
			return lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("Loading node groups...")
		}
		if m.err != nil {
			return lipgloss.NewStyle().Width(s.Width).Padding(1, 2).Render(
				shared.ErrorStyle.Render(fmt.Sprintf("Error: %v\n\nPress Esc to go back", m.err)))
		}
		if len(m.nodeGroups) == 0 {
			return lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("No node groups found. Press Esc to go back.")
		}
		return renderNodeGroupTable(m.nodeGroups, m.nodeGroupCursor, s.Width, tableHeight)

	case vsNodeList, vsNodeDetail:
		if m.detailLoading {
			return lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("Loading nodes...")
		}
		if m.err != nil {
			return lipgloss.NewStyle().Width(s.Width).Padding(1, 2).Render(
				shared.ErrorStyle.Render(fmt.Sprintf("Error: %v\n\nPress Esc to go back", m.err)))
		}
		if len(m.nodes) == 0 {
			return lipgloss.NewStyle().Width(s.Width).Padding(2, 2).Render("No nodes found. Press Esc to go back.")
		}
		return renderNodeTable(m.nodes, m.nodeCursor, s.Width, tableHeight)

	default:
		return ""
	}
}

func (m *EKSModel) renderOverlay(s *shared.SharedState) string {
	switch m.viewState {
	case vsClusterAction, vsResourceMenu, vsPodAction:
		return m.menu.Render()
	case vsClusterDetail:
		if m.selectedCluster != nil {
			return renderClusterDetail(*m.selectedCluster)
		}
	case vsPodDetail:
		if m.selectedPod != nil {
			return renderPodDetail(*m.selectedPod)
		}
	case vsPodLogs:
		if m.detailLoading {
			return shared.RenderOverlay("  Loading pod logs...")
		}
		return renderPodLogsOverlay(m.podLogs, m.podLogsErr, m.selectedPod)
	case vsDeployDetail:
		if m.deployCursor < len(m.deployments) {
			return renderDeployDetail(m.deployments[m.deployCursor])
		}
	case vsServiceDetail:
		if m.serviceCursor < len(m.services) {
			return renderServiceDetail(m.services[m.serviceCursor])
		}
	case vsNodeDetail:
		if m.nodeCursor < len(m.nodes) {
			return renderNodeDetail(m.nodes[m.nodeCursor])
		}
	}
	return ""
}

// --- Filters ---

func (m *EKSModel) applyFilters() {
	result := m.clusters
	if m.search.query != "" {
		q := strings.ToLower(m.search.query)
		var filtered []internalaws.EKSCluster
		for _, c := range result {
			if strings.Contains(internalaws.EKSSearchFields(c), q) {
				filtered = append(filtered, c)
			}
		}
		result = filtered
	}
	m.filtered = result
	if m.clusterCursor >= len(m.filtered) {
		m.clusterCursor = len(m.filtered) - 1
	}
	if m.clusterCursor < 0 {
		m.clusterCursor = 0
	}
}

// --- Loaders ---

func (m *EKSModel) loadClusters(s *shared.SharedState) tea.Cmd {
	profile := s.Profile
	region := s.Region
	return func() tea.Msg {
		ctx := context.Background()
		clients, err := internalaws.NewClients(ctx, profile, region)
		if err != nil {
			return clustersLoadedMsg{err: err}
		}
		clusters, err := internalaws.FetchEKSClusters(ctx, clients.EKS)
		return clustersLoadedMsg{clusters: clusters, err: err}
	}
}

func (m *EKSModel) loadNodeGroups(s *shared.SharedState, clusterName string) tea.Cmd {
	profile := s.Profile
	region := s.Region
	return func() tea.Msg {
		ctx := context.Background()
		clients, err := internalaws.NewClients(ctx, profile, region)
		if err != nil {
			return nodeGroupsLoadedMsg{err: err}
		}
		nodeGroups, err := internalaws.FetchEKSNodeGroups(ctx, clients.EKS, clusterName)
		return nodeGroupsLoadedMsg{nodeGroups: nodeGroups, err: err}
	}
}

func (m *EKSModel) newK8sClient(s *shared.SharedState) (*internalaws.K8sClient, error) {
	if m.selectedCluster == nil {
		return nil, fmt.Errorf("no cluster selected")
	}
	ctx := context.Background()
	token, err := internalaws.GetEKSToken(ctx, m.selectedCluster.Name, s.Profile, s.Region)
	if err != nil {
		return nil, fmt.Errorf("getting token: %w", err)
	}
	client, err := internalaws.NewK8sClient(m.selectedCluster.Endpoint, token, m.selectedCluster.CACertData)
	if err != nil {
		return nil, fmt.Errorf("creating K8s client: %w", err)
	}
	return client, nil
}

func (m *EKSModel) loadNamespaces(s *shared.SharedState) tea.Cmd {
	cluster := m.selectedCluster
	profile := s.Profile
	region := s.Region
	return func() tea.Msg {
		ctx := context.Background()
		token, err := internalaws.GetEKSToken(ctx, cluster.Name, profile, region)
		if err != nil {
			return namespacesLoadedMsg{err: err}
		}
		client, err := internalaws.NewK8sClient(cluster.Endpoint, token, cluster.CACertData)
		if err != nil {
			return namespacesLoadedMsg{err: err}
		}
		namespaces, err := client.ListNamespaces(ctx)
		return namespacesLoadedMsg{namespaces: namespaces, err: err}
	}
}

func (m *EKSModel) loadPods(s *shared.SharedState, namespace string) tea.Cmd {
	cluster := m.selectedCluster
	profile := s.Profile
	region := s.Region
	return func() tea.Msg {
		ctx := context.Background()
		token, err := internalaws.GetEKSToken(ctx, cluster.Name, profile, region)
		if err != nil {
			return podsLoadedMsg{err: err}
		}
		client, err := internalaws.NewK8sClient(cluster.Endpoint, token, cluster.CACertData)
		if err != nil {
			return podsLoadedMsg{err: err}
		}
		pods, err := client.ListPods(ctx, namespace)
		return podsLoadedMsg{pods: pods, err: err}
	}
}

func (m *EKSModel) loadDeployments(s *shared.SharedState, namespace string) tea.Cmd {
	cluster := m.selectedCluster
	profile := s.Profile
	region := s.Region
	return func() tea.Msg {
		ctx := context.Background()
		token, err := internalaws.GetEKSToken(ctx, cluster.Name, profile, region)
		if err != nil {
			return deploymentsLoadedMsg{err: err}
		}
		client, err := internalaws.NewK8sClient(cluster.Endpoint, token, cluster.CACertData)
		if err != nil {
			return deploymentsLoadedMsg{err: err}
		}
		deployments, err := client.ListDeployments(ctx, namespace)
		return deploymentsLoadedMsg{deployments: deployments, err: err}
	}
}

func (m *EKSModel) loadServices(s *shared.SharedState, namespace string) tea.Cmd {
	cluster := m.selectedCluster
	profile := s.Profile
	region := s.Region
	return func() tea.Msg {
		ctx := context.Background()
		token, err := internalaws.GetEKSToken(ctx, cluster.Name, profile, region)
		if err != nil {
			return servicesLoadedMsg{err: err}
		}
		client, err := internalaws.NewK8sClient(cluster.Endpoint, token, cluster.CACertData)
		if err != nil {
			return servicesLoadedMsg{err: err}
		}
		services, err := client.ListServices(ctx, namespace)
		return servicesLoadedMsg{services: services, err: err}
	}
}

func (m *EKSModel) loadNodes(s *shared.SharedState) tea.Cmd {
	cluster := m.selectedCluster
	profile := s.Profile
	region := s.Region
	return func() tea.Msg {
		ctx := context.Background()
		token, err := internalaws.GetEKSToken(ctx, cluster.Name, profile, region)
		if err != nil {
			return nodesLoadedMsg{err: err}
		}
		client, err := internalaws.NewK8sClient(cluster.Endpoint, token, cluster.CACertData)
		if err != nil {
			return nodesLoadedMsg{err: err}
		}
		nodes, err := client.ListNodes(ctx)
		return nodesLoadedMsg{nodes: nodes, err: err}
	}
}

func (m *EKSModel) loadPodLogs(s *shared.SharedState, namespace, podName, containerName string) tea.Cmd {
	cluster := m.selectedCluster
	profile := s.Profile
	region := s.Region
	return func() tea.Msg {
		ctx := context.Background()
		token, err := internalaws.GetEKSToken(ctx, cluster.Name, profile, region)
		if err != nil {
			return podLogsLoadedMsg{err: err}
		}
		client, err := internalaws.NewK8sClient(cluster.Endpoint, token, cluster.CACertData)
		if err != nil {
			return podLogsLoadedMsg{err: err}
		}
		logs, err := client.GetPodLogs(ctx, namespace, podName, containerName, 50)
		return podLogsLoadedMsg{logs: logs, err: err}
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
