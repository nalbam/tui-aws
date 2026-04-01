package tab_tgw

import (
	"fmt"
	"strings"

	internalaws "tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

func RenderTGWDetail(g internalaws.TransitGateway) string {
	var b strings.Builder

	name := g.Name
	if name == "" {
		name = g.ID
	}
	b.WriteString(fmt.Sprintf("  %s\n", name))
	b.WriteString("  ──────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("  TGW ID:        %s\n", g.ID))
	b.WriteString(fmt.Sprintf("  Name:          %s\n", displayStr(g.Name)))
	b.WriteString(fmt.Sprintf("  State:         %s\n", g.State))
	b.WriteString(fmt.Sprintf("  Owner:         %s\n", g.OwnerID))
	if g.ASN > 0 {
		b.WriteString(fmt.Sprintf("  ASN:           %d\n", g.ASN))
	}
	if len(g.CIDR) > 0 {
		b.WriteString(fmt.Sprintf("  CIDRs:         %s\n", strings.Join(g.CIDR, ", ")))
	}

	if len(g.Attachments) > 0 {
		b.WriteString("\n  Attachments:\n")
		for _, att := range g.Attachments {
			b.WriteString(fmt.Sprintf("    %s  %s  %s  [%s]\n", att.ID, att.ResourceType, att.ResourceID, att.State))
		}
	}

	if len(g.RouteTables) > 0 {
		b.WriteString("\n  Route Tables:\n")
		for _, rt := range g.RouteTables {
			label := rt.ID
			if rt.Name != "" {
				label = rt.Name + " (" + rt.ID + ")"
			}
			b.WriteString(fmt.Sprintf("    %s\n", label))
			for _, route := range rt.Routes {
				b.WriteString(fmt.Sprintf("      %s -> %s (%s) [%s]\n", route.DestCIDR, route.AttachmentID, route.ResourceType, route.State))
			}
		}
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
