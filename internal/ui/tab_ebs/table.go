package tab_ebs

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
		{Key: "id", Title: "Volume ID", Width: 23},
		{Key: "state", Title: "State", Width: 10},
		{Key: "type", Title: "Type", Width: 8},
		{Key: "size", Title: "Size", Width: 6},
		{Key: "iops", Title: "IOPS", Width: 6},
		{Key: "attached", Title: "Attached To", Width: 20},
		{Key: "az", Title: "AZ", Width: 5},
	}
}

func CompactColumns() []shared.Column {
	return []shared.Column{
		{Key: "name", Title: "Name", Width: 20},
		{Key: "id", Title: "Volume ID", Width: 23},
		{Key: "state", Title: "State", Width: 10},
		{Key: "type", Title: "Type", Width: 8},
		{Key: "size", Title: "Size", Width: 6},
	}
}

func ColumnsForWidth(width int) []shared.Column {
	if width < 120 {
		return shared.ExpandNameColumn(CompactColumns(), width)
	}
	return shared.ExpandNameColumn(DefaultColumns(), width)
}

func RenderTable(volumes []aws.Volume, columns []shared.Column, cursor, width, height int) string {
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

	for i := offset; i < len(volumes) && i < offset+maxRows; i++ {
		v := volumes[i]
		row := shared.RenderRow(columns, func(col shared.Column) string {
			return cellValue(col.Key, v)
		}, func(col shared.Column) lipgloss.Style {
			return cellStyle(col.Key, v)
		})

		if i == cursor {
			row = shared.TableSelectedStyle.Width(width).Render(row)
		}
		b.WriteString(row)
		if i < offset+maxRows-1 && i < len(volumes)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func cellValue(key string, v aws.Volume) string {
	switch key {
	case "name":
		if v.Name != "" {
			return v.Name
		}
		return "-"
	case "id":
		return v.ID
	case "state":
		return v.State
	case "type":
		return v.Type
	case "size":
		return fmt.Sprintf("%dGiB", v.Size)
	case "iops":
		return fmt.Sprintf("%d", v.IOPS)
	case "attached":
		if v.AttachedTo != "" {
			return v.AttachedTo
		}
		return "-"
	case "az":
		az := v.AZ
		if len(az) > 0 {
			return string(az[len(az)-1])
		}
		return az
	default:
		return ""
	}
}

func cellStyle(key string, v aws.Volume) lipgloss.Style {
	switch key {
	case "state":
		return volumeStateStyle(v.State)
	default:
		return lipgloss.Style{}
	}
}

func volumeStateStyle(state string) lipgloss.Style {
	switch state {
	case "in-use":
		return shared.StateRunning
	case "available":
		return shared.StatePending
	case "creating":
		return shared.StatePending
	case "deleting":
		return shared.StateStopping
	case "deleted":
		return shared.StateTerminated
	case "error":
		return shared.StateStopped
	default:
		return lipgloss.Style{}
	}
}

func renderStatusBar(profile, region string, count, width int) string {
	profilePart := shared.StatusKeyStyle.Render("Profile: ") + profile
	regionPart := shared.StatusKeyStyle.Render("Region: ") + region
	countPart := fmt.Sprintf("[%d Volumes]", count)
	content := fmt.Sprintf(" %s  |  %s  |  %s", profilePart, regionPart, countPart)
	return shared.StatusBarStyle.Width(width).Render(content)
}
