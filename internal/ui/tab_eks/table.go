package tab_eks

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

// ---------------------------------------------------------------------------
// Cluster table
// ---------------------------------------------------------------------------

func clusterColumns(width int) []shared.Column {
	cols := []shared.Column{
		{Key: "name", Title: "Name", Width: 25},
		{Key: "version", Title: "Version", Width: 8},
		{Key: "status", Title: "Status", Width: 10},
		{Key: "endpoint", Title: "Endpoint", Width: 40},
		{Key: "vpc", Title: "VPC", Width: 15},
	}
	if width < 100 {
		cols = []shared.Column{
			{Key: "name", Title: "Name", Width: 25},
			{Key: "version", Title: "Version", Width: 8},
			{Key: "status", Title: "Status", Width: 10},
		}
	}
	return shared.ExpandNameColumn(cols, width)
}

func renderClusterTable(clusters []aws.EKSCluster, cursor, width, height int) string {
	columns := clusterColumns(width)
	return renderGenericTable(len(clusters), cursor, columns, width, height,
		func(i int, col shared.Column) string { return clusterCellValue(col.Key, clusters[i]) },
		func(i int, col shared.Column) lipgloss.Style { return clusterCellStyle(col.Key, clusters[i]) },
	)
}

func clusterCellValue(key string, c aws.EKSCluster) string {
	switch key {
	case "name":
		return c.Name
	case "version":
		return c.Version
	case "status":
		return c.Status
	case "endpoint":
		return c.Endpoint
	case "vpc":
		return c.VpcID
	default:
		return ""
	}
}

func clusterCellStyle(key string, c aws.EKSCluster) lipgloss.Style {
	if key == "status" {
		return eksStatusStyle(c.Status)
	}
	return lipgloss.Style{}
}

// ---------------------------------------------------------------------------
// Namespace table
// ---------------------------------------------------------------------------

func namespaceColumns(width int) []shared.Column {
	cols := []shared.Column{
		{Key: "name", Title: "Namespace", Width: 30},
		{Key: "status", Title: "Status", Width: 10},
	}
	return shared.ExpandNameColumn(cols, width)
}

func renderNamespaceTable(namespaces []aws.K8sNamespace, cursor, width, height int) string {
	columns := namespaceColumns(width)
	return renderGenericTable(len(namespaces), cursor, columns, width, height,
		func(i int, col shared.Column) string {
			ns := namespaces[i]
			switch col.Key {
			case "name":
				return ns.Name
			case "status":
				return ns.Status
			default:
				return ""
			}
		},
		func(i int, col shared.Column) lipgloss.Style {
			if col.Key == "status" {
				return nsStatusStyle(namespaces[i].Status)
			}
			return lipgloss.Style{}
		},
	)
}

// ---------------------------------------------------------------------------
// Pod table
// ---------------------------------------------------------------------------

func podColumns(width int) []shared.Column {
	cols := []shared.Column{
		{Key: "name", Title: "Name", Width: 40},
		{Key: "ready", Title: "Ready", Width: 6},
		{Key: "status", Title: "Status", Width: 18},
		{Key: "restarts", Title: "Restarts", Width: 9},
		{Key: "ip", Title: "IP", Width: 15},
		{Key: "node", Title: "Node", Width: 25},
		{Key: "age", Title: "Age", Width: 8},
	}
	if width < 130 {
		cols = []shared.Column{
			{Key: "name", Title: "Name", Width: 35},
			{Key: "ready", Title: "Ready", Width: 6},
			{Key: "status", Title: "Status", Width: 18},
			{Key: "restarts", Title: "Restarts", Width: 9},
			{Key: "age", Title: "Age", Width: 8},
		}
	}
	if width < 90 {
		cols = []shared.Column{
			{Key: "name", Title: "Name", Width: 30},
			{Key: "ready", Title: "Ready", Width: 6},
			{Key: "status", Title: "Status", Width: 15},
		}
	}
	return cols
}

func renderPodTable(pods []aws.K8sPod, cursor, width, height int) string {
	columns := podColumns(width)
	return renderGenericTable(len(pods), cursor, columns, width, height,
		func(i int, col shared.Column) string { return podCellValue(col.Key, pods[i]) },
		func(i int, col shared.Column) lipgloss.Style { return podCellStyle(col.Key, pods[i]) },
	)
}

func podCellValue(key string, p aws.K8sPod) string {
	switch key {
	case "name":
		return p.Name
	case "ready":
		return p.Ready
	case "status":
		return p.Status
	case "restarts":
		return fmt.Sprintf("%d", p.Restarts)
	case "ip":
		if p.IP == "" {
			return "<none>"
		}
		return p.IP
	case "node":
		return p.Node
	case "age":
		return p.Age
	default:
		return ""
	}
}

func podCellStyle(key string, p aws.K8sPod) lipgloss.Style {
	if key == "status" {
		return podStatusStyle(p.Status)
	}
	return lipgloss.Style{}
}

// ---------------------------------------------------------------------------
// Deployment table
// ---------------------------------------------------------------------------

func deployColumns(width int) []shared.Column {
	cols := []shared.Column{
		{Key: "name", Title: "Name", Width: 35},
		{Key: "ready", Title: "Ready", Width: 8},
		{Key: "uptodate", Title: "Up-to-date", Width: 11},
		{Key: "available", Title: "Available", Width: 10},
		{Key: "age", Title: "Age", Width: 8},
	}
	if width < 90 {
		cols = []shared.Column{
			{Key: "name", Title: "Name", Width: 30},
			{Key: "ready", Title: "Ready", Width: 8},
			{Key: "available", Title: "Available", Width: 10},
		}
	}
	return cols
}

func renderDeployTable(deployments []aws.K8sDeployment, cursor, width, height int) string {
	columns := deployColumns(width)
	return renderGenericTable(len(deployments), cursor, columns, width, height,
		func(i int, col shared.Column) string { return deployCellValue(col.Key, deployments[i]) },
		func(i int, col shared.Column) lipgloss.Style { return deployCellStyle(col.Key, deployments[i]) },
	)
}

func deployCellValue(key string, d aws.K8sDeployment) string {
	switch key {
	case "name":
		return d.Name
	case "ready":
		return d.Replicas
	case "uptodate":
		return fmt.Sprintf("%d", d.UpToDate)
	case "available":
		return fmt.Sprintf("%d", d.Available)
	case "age":
		return d.Age
	default:
		return ""
	}
}

func deployCellStyle(key string, d aws.K8sDeployment) lipgloss.Style {
	if key == "ready" {
		if d.Ready == d.Available && d.Available > 0 {
			return shared.StateRunning
		}
		if d.Available == 0 {
			return shared.StateStopped
		}
		return shared.StatePending
	}
	return lipgloss.Style{}
}

// ---------------------------------------------------------------------------
// Service table
// ---------------------------------------------------------------------------

func serviceColumns(width int) []shared.Column {
	cols := []shared.Column{
		{Key: "name", Title: "Name", Width: 30},
		{Key: "type", Title: "Type", Width: 14},
		{Key: "clusterip", Title: "Cluster-IP", Width: 16},
		{Key: "externalip", Title: "External-IP", Width: 30},
		{Key: "ports", Title: "Ports", Width: 20},
		{Key: "age", Title: "Age", Width: 8},
	}
	if width < 120 {
		cols = []shared.Column{
			{Key: "name", Title: "Name", Width: 25},
			{Key: "type", Title: "Type", Width: 14},
			{Key: "clusterip", Title: "Cluster-IP", Width: 16},
			{Key: "ports", Title: "Ports", Width: 18},
			{Key: "age", Title: "Age", Width: 8},
		}
	}
	if width < 90 {
		cols = []shared.Column{
			{Key: "name", Title: "Name", Width: 25},
			{Key: "type", Title: "Type", Width: 14},
			{Key: "ports", Title: "Ports", Width: 15},
		}
	}
	return cols
}

func renderServiceTable(services []aws.K8sService, cursor, width, height int) string {
	columns := serviceColumns(width)
	return renderGenericTable(len(services), cursor, columns, width, height,
		func(i int, col shared.Column) string { return serviceCellValue(col.Key, services[i]) },
		func(i int, col shared.Column) lipgloss.Style { return lipgloss.Style{} },
	)
}

func serviceCellValue(key string, s aws.K8sService) string {
	switch key {
	case "name":
		return s.Name
	case "type":
		return s.Type
	case "clusterip":
		return s.ClusterIP
	case "externalip":
		return s.ExternalIP
	case "ports":
		return s.Ports
	case "age":
		return s.Age
	default:
		return ""
	}
}

// ---------------------------------------------------------------------------
// Node Group table (AWS)
// ---------------------------------------------------------------------------

func nodeGroupColumns(width int) []shared.Column {
	cols := []shared.Column{
		{Key: "name", Title: "Name", Width: 25},
		{Key: "status", Title: "Status", Width: 10},
		{Key: "instances", Title: "Instances", Width: 15},
		{Key: "scaling", Title: "Min/Des/Max", Width: 14},
		{Key: "ami", Title: "AMI Type", Width: 15},
	}
	if width < 90 {
		cols = []shared.Column{
			{Key: "name", Title: "Name", Width: 25},
			{Key: "status", Title: "Status", Width: 10},
			{Key: "scaling", Title: "Min/Des/Max", Width: 14},
		}
	}
	return shared.ExpandNameColumn(cols, width)
}

func renderNodeGroupTable(nodeGroups []aws.EKSNodeGroup, cursor, width, height int) string {
	columns := nodeGroupColumns(width)
	return renderGenericTable(len(nodeGroups), cursor, columns, width, height,
		func(i int, col shared.Column) string { return nodeGroupCellValue(col.Key, nodeGroups[i]) },
		func(i int, col shared.Column) lipgloss.Style { return nodeGroupCellStyle(col.Key, nodeGroups[i]) },
	)
}

func nodeGroupCellValue(key string, ng aws.EKSNodeGroup) string {
	switch key {
	case "name":
		return ng.Name
	case "status":
		return ng.Status
	case "instances":
		return ng.InstanceTypes
	case "scaling":
		return fmt.Sprintf("%d/%d/%d", ng.MinSize, ng.DesiredSize, ng.MaxSize)
	case "ami":
		return ng.AmiType
	default:
		return ""
	}
}

func nodeGroupCellStyle(key string, ng aws.EKSNodeGroup) lipgloss.Style {
	if key == "status" {
		return eksStatusStyle(ng.Status)
	}
	return lipgloss.Style{}
}

// ---------------------------------------------------------------------------
// Node table (K8s)
// ---------------------------------------------------------------------------

func nodeColumns(width int) []shared.Column {
	cols := []shared.Column{
		{Key: "name", Title: "Name", Width: 35},
		{Key: "status", Title: "Status", Width: 10},
		{Key: "roles", Title: "Roles", Width: 15},
		{Key: "version", Title: "Version", Width: 15},
		{Key: "ip", Title: "Internal-IP", Width: 16},
		{Key: "cpu", Title: "CPU", Width: 5},
		{Key: "memory", Title: "Memory", Width: 8},
		{Key: "age", Title: "Age", Width: 8},
	}
	if width < 130 {
		cols = []shared.Column{
			{Key: "name", Title: "Name", Width: 30},
			{Key: "status", Title: "Status", Width: 10},
			{Key: "version", Title: "Version", Width: 15},
			{Key: "ip", Title: "Internal-IP", Width: 16},
			{Key: "age", Title: "Age", Width: 8},
		}
	}
	if width < 90 {
		cols = []shared.Column{
			{Key: "name", Title: "Name", Width: 30},
			{Key: "status", Title: "Status", Width: 10},
			{Key: "age", Title: "Age", Width: 8},
		}
	}
	return cols
}

func renderNodeTable(nodes []aws.K8sNode, cursor, width, height int) string {
	columns := nodeColumns(width)
	return renderGenericTable(len(nodes), cursor, columns, width, height,
		func(i int, col shared.Column) string { return nodeCellValue(col.Key, nodes[i]) },
		func(i int, col shared.Column) lipgloss.Style { return nodeCellStyle(col.Key, nodes[i]) },
	)
}

func nodeCellValue(key string, n aws.K8sNode) string {
	switch key {
	case "name":
		return n.Name
	case "status":
		return n.Status
	case "roles":
		return n.Roles
	case "version":
		return n.Version
	case "ip":
		return n.InternalIP
	case "cpu":
		return n.CPU
	case "memory":
		return n.Memory
	case "age":
		return n.Age
	default:
		return ""
	}
}

func nodeCellStyle(key string, n aws.K8sNode) lipgloss.Style {
	if key == "status" {
		return nodeStatusStyle(n.Status)
	}
	return lipgloss.Style{}
}

// ---------------------------------------------------------------------------
// Generic table renderer
// ---------------------------------------------------------------------------

func renderGenericTable(
	count, cursor int,
	columns []shared.Column,
	width, height int,
	getText func(rowIdx int, col shared.Column) string,
	getStyle func(rowIdx int, col shared.Column) lipgloss.Style,
) string {
	var b strings.Builder

	header := shared.RenderRow(columns, func(col shared.Column) string {
		return col.Title
	}, nil)
	b.WriteString(shared.TableHeaderStyle.Width(width).Render(header))
	b.WriteString("\n")

	maxRows := height - 4
	if maxRows < 1 {
		maxRows = 1
	}

	offset := 0
	if cursor >= maxRows {
		offset = cursor - maxRows + 1
	}

	for i := offset; i < count && i < offset+maxRows; i++ {
		idx := i
		row := shared.RenderRow(columns,
			func(col shared.Column) string { return getText(idx, col) },
			func(col shared.Column) lipgloss.Style { return getStyle(idx, col) },
		)
		if i == cursor {
			row = shared.TableSelectedStyle.Width(width).Render(row)
		}
		b.WriteString(row)
		if i < offset+maxRows-1 && i < count-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

// ---------------------------------------------------------------------------
// Status styles
// ---------------------------------------------------------------------------

func eksStatusStyle(status string) lipgloss.Style {
	switch status {
	case "ACTIVE":
		return shared.StateRunning
	case "CREATING", "UPDATING":
		return shared.StatePending
	case "DELETING", "FAILED":
		return shared.StateStopped
	default:
		return lipgloss.Style{}
	}
}

func nsStatusStyle(status string) lipgloss.Style {
	switch status {
	case "Active":
		return shared.StateRunning
	case "Terminating":
		return shared.StateStopping
	default:
		return lipgloss.Style{}
	}
}

func podStatusStyle(status string) lipgloss.Style {
	switch status {
	case "Running":
		return shared.StateRunning
	case "Succeeded":
		return shared.StateRunning
	case "Pending", "ContainerCreating", "PodInitializing":
		return shared.StatePending
	case "Failed", "CrashLoopBackOff", "Error", "OOMKilled", "ImagePullBackOff", "ErrImagePull":
		return shared.StateStopped
	case "Terminating":
		return shared.StateStopping
	case "Completed":
		return shared.StateTerminated
	default:
		return lipgloss.Style{}
	}
}

func nodeStatusStyle(status string) lipgloss.Style {
	switch status {
	case "Ready":
		return shared.StateRunning
	case "NotReady":
		return shared.StateStopped
	default:
		return lipgloss.Style{}
	}
}
