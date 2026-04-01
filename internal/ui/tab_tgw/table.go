package tab_tgw

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

func DefaultColumns() []shared.Column {
	return []shared.Column{
		{Key: "name", Title: "Name", Width: 20},
		{Key: "id", Title: "TGW ID", Width: 23},
		{Key: "state", Title: "State", Width: 10},
		{Key: "asn", Title: "ASN", Width: 10},
		{Key: "attachments", Title: "Attach", Width: 5},
	}
}

func CompactColumns() []shared.Column {
	return []shared.Column{
		{Key: "name", Title: "Name", Width: 20},
		{Key: "id", Title: "TGW ID", Width: 23},
		{Key: "state", Title: "State", Width: 10},
	}
}

func ColumnsForWidth(width int) []shared.Column {
	if width < 100 {
		return shared.ExpandNameColumn(CompactColumns(), width)
	}
	return shared.ExpandNameColumn(DefaultColumns(), width)
}

func RenderTable(gateways []aws.TransitGateway, columns []shared.Column, cursor, width, height int) string {
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

	for i := offset; i < len(gateways) && i < offset+maxRows; i++ {
		g := gateways[i]
		row := shared.RenderRow(columns, func(col shared.Column) string {
			return cellValue(col.Key, g)
		}, func(col shared.Column) lipgloss.Style {
			return cellStyle(col.Key, g)
		})

		if i == cursor {
			row = shared.TableSelectedStyle.Width(width).Render(row)
		}
		b.WriteString(row)
		if i < offset+maxRows-1 && i < len(gateways)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func cellValue(key string, g aws.TransitGateway) string {
	switch key {
	case "name":
		if g.Name != "" {
			return g.Name
		}
		return "-"
	case "id":
		return g.ID
	case "state":
		return g.State
	case "asn":
		if g.ASN > 0 {
			return fmt.Sprintf("%d", g.ASN)
		}
		return "-"
	case "attachments":
		return fmt.Sprintf("%d", len(g.Attachments))
	default:
		return ""
	}
}

func cellStyle(key string, g aws.TransitGateway) lipgloss.Style {
	switch key {
	case "state":
		return tgwStateStyle(g.State)
	default:
		return lipgloss.Style{}
	}
}

func tgwStateStyle(state string) lipgloss.Style {
	switch state {
	case "available":
		return shared.StateRunning
	case "pending":
		return shared.StatePending
	case "modifying":
		return shared.StatePending
	case "deleting":
		return shared.StateStopping
	case "deleted":
		return shared.StateTerminated
	default:
		return lipgloss.Style{}
	}
}

func renderStatusBar(profile, region string, count, width int) string {
	profilePart := shared.StatusKeyStyle.Render("Profile: ") + profile
	regionPart := shared.StatusKeyStyle.Render("Region: ") + region
	countPart := fmt.Sprintf("[%d TGWs]", count)
	content := fmt.Sprintf(" %s  |  %s  |  %s", profilePart, regionPart, countPart)
	return shared.StatusBarStyle.Width(width).Render(content)
}
