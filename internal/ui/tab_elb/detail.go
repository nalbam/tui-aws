package tab_elb

import (
	"fmt"
	"strings"

	internalaws "tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

// RenderELBDetail renders the load balancer detail overlay.
func RenderELBDetail(lb internalaws.LoadBalancer, loading bool) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("  %s  (%s)\n", lb.Name, lb.TypeLabel()))
	b.WriteString("  ──────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("  Name:          %s\n", lb.Name))
	if lb.ARN != "" {
		b.WriteString(fmt.Sprintf("  ARN:           %s\n", lb.ARN))
	}
	b.WriteString(fmt.Sprintf("  DNS:           %s\n", lb.DNSName))
	b.WriteString(fmt.Sprintf("  Type:          %s (%s)\n", lb.Type, lb.TypeLabel()))
	b.WriteString(fmt.Sprintf("  Scheme:        %s\n", lb.Scheme))
	b.WriteString(fmt.Sprintf("  State:         %s\n", lb.State))
	b.WriteString(fmt.Sprintf("  VPC:           %s\n", lb.VpcID))
	if lb.CreatedTime != "" {
		b.WriteString(fmt.Sprintf("  Created:       %s\n", lb.CreatedTime))
	}

	if len(lb.AZs) > 0 {
		b.WriteString("\n  Availability Zones:\n")
		for _, az := range lb.AZs {
			b.WriteString(fmt.Sprintf("    %s\n", az))
		}
	}

	if len(lb.SecurityGroups) > 0 {
		b.WriteString("\n  Security Groups:\n")
		for _, sg := range lb.SecurityGroups {
			b.WriteString(fmt.Sprintf("    %s\n", sg))
		}
	}

	if loading {
		b.WriteString("\n  Loading listeners and target groups...")
	} else {
		if len(lb.Listeners) > 0 {
			b.WriteString("\n  Listeners:\n")
			for _, l := range lb.Listeners {
				rulesInfo := ""
				if l.Rules > 0 {
					rulesInfo = fmt.Sprintf("  (%d rules)", l.Rules)
				}
				b.WriteString(fmt.Sprintf("    %s :%d%s\n", l.Protocol, l.Port, rulesInfo))
			}
		} else if lb.Type != "classic" {
			b.WriteString("\n  Listeners:     (none)\n")
		}

		if len(lb.TargetGroups) > 0 {
			b.WriteString("\n  Target Groups:\n")
			for _, tg := range lb.TargetGroups {
				health := tg.HealthCheck
				if health == "" {
					health = "-"
				}
				b.WriteString(fmt.Sprintf("    %s  %s:%d  %s  %s\n", tg.Name, tg.Protocol, tg.Port, tg.TargetType, health))
			}
		} else if lb.Type != "classic" {
			b.WriteString("\n  Target Groups: (none)\n")
		}
	}

	b.WriteString("\n  Press any key to close")
	return shared.RenderOverlay(b.String())
}
