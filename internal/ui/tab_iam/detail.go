package tab_iam

import (
	"fmt"
	"strings"

	internalaws "tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

func RenderUserDetail(u internalaws.IAMUser) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("  %s\n", u.UserName))
	b.WriteString("  ──────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("  UserName:         %s\n", u.UserName))
	b.WriteString(fmt.Sprintf("  UserID:           %s\n", u.UserID))
	b.WriteString(fmt.Sprintf("  ARN:              %s\n", u.ARN))
	b.WriteString(fmt.Sprintf("  Created:          %s\n", displayStr(u.CreateDate)))
	lastUsed := u.PasswordLastUsed
	if lastUsed == "" {
		lastUsed = "Never"
	}
	b.WriteString(fmt.Sprintf("  Password Used:    %s\n", lastUsed))

	if len(u.Groups) > 0 {
		b.WriteString("\n  Groups:\n")
		for _, g := range u.Groups {
			b.WriteString(fmt.Sprintf("    %s\n", g))
		}
	} else {
		b.WriteString("\n  Groups: (none)\n")
	}

	if len(u.Policies) > 0 {
		b.WriteString("\n  Attached Policies:\n")
		for _, p := range u.Policies {
			b.WriteString(fmt.Sprintf("    %s\n", p))
		}
	} else {
		b.WriteString("\n  Attached Policies: (none)\n")
	}

	b.WriteString("\n  Press any key to close")
	return shared.RenderOverlay(b.String())
}

func displayStr(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
