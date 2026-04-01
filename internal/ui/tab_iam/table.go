package tab_iam

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

func DefaultColumns() []shared.Column {
	return []shared.Column{
		{Key: "name", Title: "UserName", Width: 20},
		{Key: "userid", Title: "UserID", Width: 21},
		{Key: "arn", Title: "ARN", Width: 40},
		{Key: "created", Title: "Created", Width: 18},
		{Key: "lastused", Title: "Last Used", Width: 18},
	}
}

func CompactColumns() []shared.Column {
	return []shared.Column{
		{Key: "name", Title: "UserName", Width: 20},
		{Key: "arn", Title: "ARN", Width: 40},
		{Key: "created", Title: "Created", Width: 18},
	}
}

func ColumnsForWidth(width int) []shared.Column {
	if width < 140 {
		return shared.ExpandNameColumn(CompactColumns(), width)
	}
	return shared.ExpandNameColumn(DefaultColumns(), width)
}

func RenderTable(users []aws.IAMUser, columns []shared.Column, cursor, width, height int) string {
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

	for i := offset; i < len(users) && i < offset+maxRows; i++ {
		u := users[i]
		row := shared.RenderRow(columns, func(col shared.Column) string {
			return cellValue(col.Key, u)
		}, func(col shared.Column) lipgloss.Style {
			return lipgloss.Style{}
		})

		if i == cursor {
			row = shared.TableSelectedStyle.Width(width).Render(row)
		}
		b.WriteString(row)
		if i < offset+maxRows-1 && i < len(users)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func cellValue(key string, u aws.IAMUser) string {
	switch key {
	case "name":
		return u.UserName
	case "userid":
		return u.UserID
	case "arn":
		return u.ARN
	case "created":
		if u.CreateDate != "" {
			return u.CreateDate
		}
		return "-"
	case "lastused":
		if u.PasswordLastUsed != "" {
			return u.PasswordLastUsed
		}
		return "Never"
	default:
		return ""
	}
}

func renderStatusBar(profile, region string, count int, identity *aws.CallerIdentity, width int) string {
	profilePart := shared.StatusKeyStyle.Render("Profile: ") + profile
	regionPart := shared.StatusKeyStyle.Render("Region: ") + region
	countPart := fmt.Sprintf("[%d Users]", count)
	content := fmt.Sprintf(" %s  |  %s  |  %s", profilePart, regionPart, countPart)
	if identity != nil {
		content += fmt.Sprintf("  |  Account: %s", identity.Account)
	}
	return shared.StatusBarStyle.Width(width).Render(content)
}
