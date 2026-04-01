package tab_vpce

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

// DefaultColumns returns the VPCE table columns.
func DefaultColumns() []shared.Column {
	return []shared.Column{
		{Key: "name", Title: "Name", Width: 25},
		{Key: "id", Title: "Endpoint ID", Width: 25},
		{Key: "service", Title: "Service Name", Width: 40},
		{Key: "type", Title: "Type", Width: 10},
		{Key: "vpc", Title: "VPC", Width: 15},
		{Key: "state", Title: "State", Width: 10},
	}
}

// CompactColumns returns a minimal column set for narrow terminals.
func CompactColumns() []shared.Column {
	return []shared.Column{
		{Key: "name", Title: "Name", Width: 25},
		{Key: "id", Title: "Endpoint ID", Width: 25},
		{Key: "service", Title: "Service Name", Width: 40},
		{Key: "state", Title: "State", Width: 10},
	}
}

// ColumnsForWidth returns the appropriate column set for the given terminal width.
func ColumnsForWidth(width int) []shared.Column {
	if width < 100 {
		return shared.ExpandNameColumn(CompactColumns(), width)
	}
	return shared.ExpandNameColumn(DefaultColumns(), width)
}

// RenderTable renders the VPCE table with header, rows, and scrolling.
func RenderTable(endpoints []aws.VPCEndpoint, columns []shared.Column, cursor, width, height int) string {
	var b strings.Builder

	// Header
	header := shared.RenderRow(columns, func(col shared.Column) string {
		return col.Title
	}, nil)
	b.WriteString(shared.TableHeaderStyle.Width(width).Render(header))
	b.WriteString("\n")

	// Available rows: total height minus statusbar(1) + helpbar(1) + header(1) + tabbar(1)
	maxRows := height - 4
	if maxRows < 1 {
		maxRows = 1
	}

	// Calculate scroll offset
	offset := 0
	if cursor >= maxRows {
		offset = cursor - maxRows + 1
	}

	for i := offset; i < len(endpoints) && i < offset+maxRows; i++ {
		ep := endpoints[i]
		row := shared.RenderRow(columns, func(col shared.Column) string {
			return cellValue(col.Key, ep)
		}, func(col shared.Column) lipgloss.Style {
			return cellStyle(col.Key, ep)
		})

		if i == cursor {
			row = shared.TableSelectedStyle.Width(width).Render(row)
		}
		b.WriteString(row)
		if i < offset+maxRows-1 && i < len(endpoints)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func cellValue(key string, ep aws.VPCEndpoint) string {
	switch key {
	case "name":
		if ep.Name != "" {
			return ep.Name
		}
		return "-"
	case "id":
		return ep.ID
	case "service":
		return ep.ServiceName
	case "type":
		return ep.Type
	case "vpc":
		return ep.VpcID
	case "state":
		return ep.State
	default:
		return ""
	}
}

func cellStyle(key string, ep aws.VPCEndpoint) lipgloss.Style {
	switch key {
	case "state":
		return endpointStateStyle(ep.State)
	default:
		return lipgloss.Style{}
	}
}

func endpointStateStyle(state string) lipgloss.Style {
	switch state {
	case "available":
		return shared.StateRunning
	case "pending", "pendingAcceptance":
		return shared.StatePending
	case "deleting":
		return shared.StateStopping
	case "deleted":
		return shared.StateTerminated
	case "rejected", "failed":
		return shared.StateStopped
	default:
		return lipgloss.Style{}
	}
}

func renderStatusBar(profile, region string, count, width int) string {
	profilePart := shared.StatusKeyStyle.Render("Profile: ") + profile
	regionPart := shared.StatusKeyStyle.Render("Region: ") + region
	countPart := fmt.Sprintf("[%d Endpoints]", count)

	content := fmt.Sprintf(" %s  |  %s  |  %s", profilePart, regionPart, countPart)
	return shared.StatusBarStyle.Width(width).Render(content)
}
