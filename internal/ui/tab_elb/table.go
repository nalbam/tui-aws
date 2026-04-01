package tab_elb

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

// Type color styles
var (
	typeALB  = lipgloss.NewStyle().Foreground(lipgloss.Color("#83a598")) // blue
	typeNLB  = lipgloss.NewStyle().Foreground(lipgloss.Color("#b8bb26")) // green
	typeGWLB = lipgloss.NewStyle().Foreground(lipgloss.Color("#fabd2f")) // yellow
	typeCLB  = lipgloss.NewStyle().Foreground(lipgloss.Color("#928374")) // gray
)

// DefaultColumns returns the ELB table columns.
func DefaultColumns() []shared.Column {
	return []shared.Column{
		{Key: "name", Title: "Name", Width: 30},
		{Key: "type", Title: "Type", Width: 8},
		{Key: "scheme", Title: "Scheme", Width: 16},
		{Key: "dns", Title: "DNS Name", Width: 45},
		{Key: "vpc", Title: "VPC", Width: 15},
		{Key: "state", Title: "State", Width: 10},
		{Key: "azs", Title: "AZs", Width: 10},
	}
}

// CompactColumns returns a minimal column set for narrow terminals.
func CompactColumns() []shared.Column {
	return []shared.Column{
		{Key: "name", Title: "Name", Width: 30},
		{Key: "type", Title: "Type", Width: 8},
		{Key: "scheme", Title: "Scheme", Width: 16},
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

// RenderTable renders the ELB table with header, rows, and scrolling.
func RenderTable(lbs []aws.LoadBalancer, columns []shared.Column, cursor, width, height int) string {
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

	for i := offset; i < len(lbs) && i < offset+maxRows; i++ {
		lb := lbs[i]
		row := shared.RenderRow(columns, func(col shared.Column) string {
			return cellValue(col.Key, lb)
		}, func(col shared.Column) lipgloss.Style {
			return cellStyle(col.Key, lb)
		})

		if i == cursor {
			row = shared.TableSelectedStyle.Width(width).Render(row)
		}
		b.WriteString(row)
		if i < offset+maxRows-1 && i < len(lbs)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func cellValue(key string, lb aws.LoadBalancer) string {
	switch key {
	case "name":
		if lb.Name != "" {
			return lb.Name
		}
		return "-"
	case "type":
		return lb.TypeLabel()
	case "scheme":
		return lb.Scheme
	case "dns":
		return lb.DNSName
	case "vpc":
		return lb.VpcID
	case "state":
		return lb.State
	case "azs":
		return fmt.Sprintf("%d", len(lb.AZs))
	default:
		return ""
	}
}

func cellStyle(key string, lb aws.LoadBalancer) lipgloss.Style {
	switch key {
	case "type":
		return typeStyle(lb.Type)
	case "state":
		return stateStyle(lb.State)
	default:
		return lipgloss.Style{}
	}
}

func typeStyle(lbType string) lipgloss.Style {
	switch lbType {
	case "application":
		return typeALB
	case "network":
		return typeNLB
	case "gateway":
		return typeGWLB
	case "classic":
		return typeCLB
	default:
		return lipgloss.Style{}
	}
}

func stateStyle(state string) lipgloss.Style {
	switch state {
	case "active":
		return shared.StateRunning
	case "provisioning":
		return shared.StatePending
	case "active_impaired":
		return shared.StateStopping
	case "failed":
		return shared.StateStopped
	default:
		return lipgloss.Style{}
	}
}

func renderStatusBar(profile, region string, count, width int) string {
	profilePart := shared.StatusKeyStyle.Render("Profile: ") + profile
	regionPart := shared.StatusKeyStyle.Render("Region: ") + region
	countPart := fmt.Sprintf("[%d Load Balancers]", count)

	content := fmt.Sprintf(" %s  |  %s  |  %s", profilePart, regionPart, countPart)
	return shared.StatusBarStyle.Width(width).Render(content)
}
