package tab_asg

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

func DefaultColumns() []shared.Column {
	return []shared.Column{
		{Key: "name", Title: "Name", Width: 30},
		{Key: "min", Title: "Min", Width: 4},
		{Key: "max", Title: "Max", Width: 4},
		{Key: "desired", Title: "Desired", Width: 4},
		{Key: "instances", Title: "Inst", Width: 5},
		{Key: "health", Title: "Health", Width: 6},
		{Key: "azs", Title: "AZs", Width: 15},
	}
}

func CompactColumns() []shared.Column {
	return []shared.Column{
		{Key: "name", Title: "Name", Width: 30},
		{Key: "min", Title: "Min", Width: 4},
		{Key: "max", Title: "Max", Width: 4},
		{Key: "desired", Title: "Des", Width: 4},
		{Key: "instances", Title: "Inst", Width: 5},
	}
}

func ColumnsForWidth(width int) []shared.Column {
	if width < 100 {
		return shared.ExpandNameColumn(CompactColumns(), width)
	}
	return shared.ExpandNameColumn(DefaultColumns(), width)
}

func RenderTable(groups []aws.AutoScalingGroup, columns []shared.Column, cursor, width, height int) string {
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

	for i := offset; i < len(groups) && i < offset+maxRows; i++ {
		g := groups[i]
		row := shared.RenderRow(columns, func(col shared.Column) string {
			return cellValue(col.Key, g)
		}, func(col shared.Column) lipgloss.Style {
			return cellStyle(col.Key, g)
		})

		if i == cursor {
			row = shared.TableSelectedStyle.Width(width).Render(row)
		}
		b.WriteString(row)
		if i < offset+maxRows-1 && i < len(groups)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func cellValue(key string, g aws.AutoScalingGroup) string {
	switch key {
	case "name":
		if g.Name != "" {
			return g.Name
		}
		return "-"
	case "min":
		return fmt.Sprintf("%d", g.MinSize)
	case "max":
		return fmt.Sprintf("%d", g.MaxSize)
	case "desired":
		return fmt.Sprintf("%d", g.DesiredCapacity)
	case "instances":
		return fmt.Sprintf("%d", len(g.Instances))
	case "health":
		return g.HealthCheckType
	case "azs":
		return g.AZsShort()
	default:
		return ""
	}
}

func cellStyle(key string, g aws.AutoScalingGroup) lipgloss.Style {
	switch key {
	case "health":
		if g.HealthCheckType == "ELB" {
			return shared.StateRunning
		}
		return lipgloss.Style{}
	default:
		return lipgloss.Style{}
	}
}

func renderStatusBar(profile, region string, count, width int) string {
	profilePart := shared.StatusKeyStyle.Render("Profile: ") + profile
	regionPart := shared.StatusKeyStyle.Render("Region: ") + region
	countPart := fmt.Sprintf("[%d ASGs]", count)
	content := fmt.Sprintf(" %s  |  %s  |  %s", profilePart, regionPart, countPart)
	return shared.StatusBarStyle.Width(width).Render(content)
}
