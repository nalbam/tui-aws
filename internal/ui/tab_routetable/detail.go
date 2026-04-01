package tab_routetable

import (
	"fmt"
	"strings"

	"tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

// RenderRouteEntries renders the route entries detail overlay.
func RenderRouteEntries(rt aws.RouteTable) string {
	var b strings.Builder

	name := rt.Name
	if name == "" {
		name = rt.ID
	}
	b.WriteString(fmt.Sprintf("  %s (%s)\n", name, rt.ID))
	b.WriteString("  ──────────────────────────────────\n")

	if len(rt.Routes) == 0 {
		b.WriteString("  (no routes)\n")
	} else {
		// Header
		b.WriteString(fmt.Sprintf("  %-19s %-17s %s\n", "Destination", "Target", "State"))
		for _, route := range rt.Routes {
			b.WriteString(fmt.Sprintf("  %-19s %-17s %s\n", route.Destination, route.Target, route.State))
		}
	}

	b.WriteString("\n  Esc: close")

	return shared.RenderOverlay(b.String())
}

// RenderSubnets renders the associated subnets detail overlay.
func RenderSubnets(rt aws.RouteTable) string {
	var b strings.Builder

	name := rt.Name
	if name == "" {
		name = rt.ID
	}
	b.WriteString(fmt.Sprintf("  %s (%s)  Associated Subnets\n", name, rt.ID))
	b.WriteString("  ──────────────────────────────────\n")

	if len(rt.Subnets) == 0 {
		b.WriteString("  (no explicit associations)\n")
		if rt.IsMain {
			b.WriteString("  This is the main route table — it applies to all\n")
			b.WriteString("  subnets without explicit associations.\n")
		}
	} else {
		for _, subnet := range rt.Subnets {
			b.WriteString(fmt.Sprintf("  - %s\n", subnet))
		}
	}

	b.WriteString("\n  Esc: close")

	return shared.RenderOverlay(b.String())
}
