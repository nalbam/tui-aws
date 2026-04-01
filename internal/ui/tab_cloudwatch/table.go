package tab_cloudwatch

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
		{Key: "state", Title: "State", Width: 10},
		{Key: "metric", Title: "Metric", Width: 20},
		{Key: "namespace", Title: "Namespace", Width: 20},
		{Key: "threshold", Title: "Threshold", Width: 10},
	}
}

func CompactColumns() []shared.Column {
	return []shared.Column{
		{Key: "name", Title: "Name", Width: 30},
		{Key: "state", Title: "State", Width: 10},
		{Key: "metric", Title: "Metric", Width: 20},
	}
}

func ColumnsForWidth(width int) []shared.Column {
	if width < 110 {
		return shared.ExpandNameColumn(CompactColumns(), width)
	}
	return shared.ExpandNameColumn(DefaultColumns(), width)
}

func RenderTable(alarms []aws.Alarm, columns []shared.Column, cursor, width, height int) string {
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

	for i := offset; i < len(alarms) && i < offset+maxRows; i++ {
		a := alarms[i]
		row := shared.RenderRow(columns, func(col shared.Column) string {
			return cellValue(col.Key, a)
		}, func(col shared.Column) lipgloss.Style {
			return cellStyle(col.Key, a)
		})

		if i == cursor {
			row = shared.TableSelectedStyle.Width(width).Render(row)
		}
		b.WriteString(row)
		if i < offset+maxRows-1 && i < len(alarms)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func cellValue(key string, a aws.Alarm) string {
	switch key {
	case "name":
		return a.Name
	case "state":
		return a.State
	case "metric":
		return a.MetricName
	case "namespace":
		return a.Namespace
	case "threshold":
		return fmt.Sprintf("%.2f", a.Threshold)
	default:
		return ""
	}
}

func cellStyle(key string, a aws.Alarm) lipgloss.Style {
	switch key {
	case "state":
		return alarmStateStyle(a.State)
	default:
		return lipgloss.Style{}
	}
}

func alarmStateStyle(state string) lipgloss.Style {
	switch state {
	case "OK":
		return shared.StateRunning // green
	case "ALARM":
		return shared.StateStopped // red
	case "INSUFFICIENT_DATA":
		return shared.StatePending // yellow
	default:
		return lipgloss.Style{}
	}
}

func renderStatusBar(profile, region string, count, width int) string {
	profilePart := shared.StatusKeyStyle.Render("Profile: ") + profile
	regionPart := shared.StatusKeyStyle.Render("Region: ") + region
	countPart := fmt.Sprintf("[%d Alarms]", count)
	content := fmt.Sprintf(" %s  |  %s  |  %s", profilePart, regionPart, countPart)
	return shared.StatusBarStyle.Width(width).Render(content)
}
