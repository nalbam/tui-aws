package tab_sg

import (
	"fmt"
	"strings"

	"tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

// RenderSGRules renders the SG rules detail overlay (inbound or outbound).
func RenderSGRules(sg aws.SecurityGroup, kind detailKind) string {
	var b strings.Builder

	name := sg.Name
	if name == "" {
		name = sg.ID
	}

	direction := "Inbound Rules"
	rules := sg.InboundRules
	if kind == detailOutbound {
		direction = "Outbound Rules"
		rules = sg.OutboundRules
	}

	b.WriteString(fmt.Sprintf("  %s (%s)  %s\n", name, sg.ID, direction))
	b.WriteString("  ──────────────────────────────────────\n")

	if len(rules) == 0 {
		b.WriteString("  (no rules)\n")
	} else {
		// Header
		b.WriteString(fmt.Sprintf("  %-6s %-9s %-15s %s\n", "Proto", "Ports", "Source", "Description"))
		for _, rule := range rules {
			desc := rule.Description
			if len(desc) > 30 {
				desc = desc[:27] + "..."
			}
			b.WriteString(fmt.Sprintf("  %-6s %-9s %-15s %s\n", rule.Protocol, rule.PortRange, rule.Source, desc))
		}
	}

	b.WriteString("\n  Esc: close")

	return shared.RenderOverlay(b.String())
}

// RenderNACLRules renders the NACL rules detail overlay (inbound or outbound).
func RenderNACLRules(nacl aws.NetworkACL, kind detailKind) string {
	var b strings.Builder

	name := nacl.Name
	if name == "" {
		name = nacl.ID
	}

	direction := "Inbound Rules"
	rules := nacl.InboundRules
	if kind == detailOutbound {
		direction = "Outbound Rules"
		rules = nacl.OutboundRules
	}

	b.WriteString(fmt.Sprintf("  %s  %s\n", name, direction))
	b.WriteString("  ──────────────────────────────────────\n")

	if len(rules) == 0 {
		b.WriteString("  (no rules)\n")
	} else {
		// Header
		b.WriteString(fmt.Sprintf("  %-6s %-6s %-7s %-15s %s\n", "Rule#", "Proto", "Ports", "CIDR", "Action"))
		for _, rule := range rules {
			ruleNum := fmt.Sprintf("%d", rule.RuleNumber)
			if rule.RuleNumber == 32767 {
				ruleNum = "*"
			}
			action := strings.ToUpper(rule.Action)
			b.WriteString(fmt.Sprintf("  %-6s %-6s %-7s %-15s %s\n", ruleNum, rule.Protocol, rule.PortRange, rule.CIDRBlock, action))
		}
	}

	b.WriteString("\n  Esc: close")

	return shared.RenderOverlay(b.String())
}
