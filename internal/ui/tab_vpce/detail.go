package tab_vpce

import (
	"fmt"
	"strings"

	internalaws "tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

// RenderEndpointDetail renders the VPC endpoint detail overlay.
func RenderEndpointDetail(ep internalaws.VPCEndpoint) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("  %s  %s\n", ep.ID, ep.ServiceName))
	b.WriteString("  ──────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("  Name:          %s\n", displayName(ep.Name)))
	b.WriteString(fmt.Sprintf("  ID:            %s\n", ep.ID))
	b.WriteString(fmt.Sprintf("  Service:       %s\n", ep.ServiceName))
	b.WriteString(fmt.Sprintf("  Type:          %s\n", ep.Type))
	b.WriteString(fmt.Sprintf("  State:         %s\n", ep.State))
	b.WriteString(fmt.Sprintf("  VPC:           %s\n", ep.VpcID))
	b.WriteString(fmt.Sprintf("  Private DNS:   %s\n", boolLabel(ep.PrivateDNS)))
	if ep.CreationTime != "" {
		b.WriteString(fmt.Sprintf("  Created:       %s\n", ep.CreationTime))
	}

	if len(ep.RouteTableIDs) > 0 {
		b.WriteString("\n  Route Tables:\n")
		for _, rtb := range ep.RouteTableIDs {
			b.WriteString(fmt.Sprintf("    %s\n", rtb))
		}
	}

	if len(ep.SubnetIDs) > 0 {
		b.WriteString("\n  Subnets:\n")
		for _, sub := range ep.SubnetIDs {
			b.WriteString(fmt.Sprintf("    %s\n", sub))
		}
	}

	if len(ep.SecurityGroupIDs) > 0 {
		b.WriteString("\n  Security Groups:\n")
		for _, sg := range ep.SecurityGroupIDs {
			b.WriteString(fmt.Sprintf("    %s\n", sg))
		}
	}

	if len(ep.NetworkInterfaceIDs) > 0 {
		b.WriteString("\n  Network Interfaces:\n")
		for _, eni := range ep.NetworkInterfaceIDs {
			b.WriteString(fmt.Sprintf("    %s\n", eni))
		}
	}

	b.WriteString("\n  Press any key to close")
	return shared.RenderOverlay(b.String())
}

func displayName(name string) string {
	if name == "" {
		return "-"
	}
	return name
}

func boolLabel(v bool) string {
	if v {
		return "enabled"
	}
	return "disabled"
}
