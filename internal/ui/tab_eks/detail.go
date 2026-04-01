package tab_eks

import (
	"fmt"
	"strings"

	internalaws "tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

// renderClusterDetail renders the EKS cluster info overlay.
func renderClusterDetail(c internalaws.EKSCluster) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  Cluster: %s  (v%s)\n", c.Name, c.Version))
	b.WriteString("  ──────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("  Name:             %s\n", c.Name))
	b.WriteString(fmt.Sprintf("  ARN:              %s\n", c.ARN))
	b.WriteString(fmt.Sprintf("  Version:          %s\n", c.Version))
	b.WriteString(fmt.Sprintf("  Platform:         %s\n", c.PlatformVersion))
	b.WriteString(fmt.Sprintf("  Status:           %s\n", c.Status))
	b.WriteString(fmt.Sprintf("  Endpoint:         %s\n", c.Endpoint))
	b.WriteString(fmt.Sprintf("  VPC:              %s\n", c.VpcID))
	if c.CreatedTime != "" {
		b.WriteString(fmt.Sprintf("  Created:          %s\n", c.CreatedTime))
	}

	if len(c.SubnetIDs) > 0 {
		b.WriteString("\n  Subnets:\n")
		for _, sub := range c.SubnetIDs {
			b.WriteString(fmt.Sprintf("    %s\n", sub))
		}
	}

	if len(c.SecurityGroupIDs) > 0 {
		b.WriteString("\n  Security Groups:\n")
		for _, sg := range c.SecurityGroupIDs {
			b.WriteString(fmt.Sprintf("    %s\n", sg))
		}
	}

	b.WriteString("\n  Esc: close")
	return shared.RenderOverlay(b.String())
}

// renderPodDetail renders the pod info overlay.
func renderPodDetail(p internalaws.K8sPod) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  Pod: %s\n", p.Name))
	b.WriteString("  ──────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("  Name:         %s\n", p.Name))
	b.WriteString(fmt.Sprintf("  Namespace:    %s\n", p.Namespace))
	b.WriteString(fmt.Sprintf("  Status:       %s\n", p.Status))
	b.WriteString(fmt.Sprintf("  Ready:        %s\n", p.Ready))
	b.WriteString(fmt.Sprintf("  Restarts:     %d\n", p.Restarts))
	if p.IP != "" {
		b.WriteString(fmt.Sprintf("  Pod IP:       %s\n", p.IP))
	}
	if p.Node != "" {
		b.WriteString(fmt.Sprintf("  Node:         %s\n", p.Node))
	}
	if p.Age != "" {
		b.WriteString(fmt.Sprintf("  Age:          %s\n", p.Age))
	}

	if len(p.Containers) > 0 {
		b.WriteString("\n  Containers:\n")
		for _, c := range p.Containers {
			readyStr := "false"
			if c.Ready {
				readyStr = "true"
			}
			b.WriteString(fmt.Sprintf("    %-25s %-12s ready:%-5s restarts:%d\n",
				c.Name, c.State, readyStr, c.RestartCount))
			b.WriteString(fmt.Sprintf("      Image: %s\n", c.Image))
		}
	}

	b.WriteString("\n  Esc: close")
	return shared.RenderOverlay(b.String())
}

// renderPodLogsOverlay renders the pod log viewer.
func renderPodLogsOverlay(logs string, logsErr error, pod *internalaws.K8sPod) string {
	var b strings.Builder

	title := "Pod Logs"
	if pod != nil {
		title = fmt.Sprintf("Logs: %s", pod.Name)
	}
	b.WriteString(fmt.Sprintf("  %s\n", title))
	b.WriteString("  ──────────────────────────────────────────────────\n")

	if logsErr != nil {
		b.WriteString(fmt.Sprintf("\n  Error: %v\n", logsErr))
	} else if logs == "" {
		b.WriteString("\n  No log output.\n")
	} else {
		for _, line := range strings.Split(strings.TrimRight(logs, "\n"), "\n") {
			b.WriteString(fmt.Sprintf("  %s\n", line))
		}
	}

	b.WriteString("\n  Esc: close")
	return shared.RenderOverlay(b.String())
}

// renderDeployDetail renders the deployment info overlay.
func renderDeployDetail(d internalaws.K8sDeployment) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  Deployment: %s\n", d.Name))
	b.WriteString("  ──────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("  Name:         %s\n", d.Name))
	b.WriteString(fmt.Sprintf("  Namespace:    %s\n", d.Namespace))
	b.WriteString(fmt.Sprintf("  Replicas:     %s\n", d.Replicas))
	b.WriteString(fmt.Sprintf("  Ready:        %d\n", d.Ready))
	b.WriteString(fmt.Sprintf("  Up-to-date:   %d\n", d.UpToDate))
	b.WriteString(fmt.Sprintf("  Available:    %d\n", d.Available))
	if d.Age != "" {
		b.WriteString(fmt.Sprintf("  Age:          %s\n", d.Age))
	}

	b.WriteString("\n  Esc: close")
	return shared.RenderOverlay(b.String())
}

// renderServiceDetail renders the service info overlay.
func renderServiceDetail(s internalaws.K8sService) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  Service: %s\n", s.Name))
	b.WriteString("  ──────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("  Name:         %s\n", s.Name))
	b.WriteString(fmt.Sprintf("  Namespace:    %s\n", s.Namespace))
	b.WriteString(fmt.Sprintf("  Type:         %s\n", s.Type))
	b.WriteString(fmt.Sprintf("  Cluster IP:   %s\n", s.ClusterIP))
	b.WriteString(fmt.Sprintf("  External IP:  %s\n", s.ExternalIP))
	b.WriteString(fmt.Sprintf("  Ports:        %s\n", s.Ports))
	if s.Age != "" {
		b.WriteString(fmt.Sprintf("  Age:          %s\n", s.Age))
	}

	b.WriteString("\n  Esc: close")
	return shared.RenderOverlay(b.String())
}

// renderNodeDetail renders the node info overlay.
func renderNodeDetail(n internalaws.K8sNode) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  Node: %s\n", n.Name))
	b.WriteString("  ──────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("  Name:         %s\n", n.Name))
	b.WriteString(fmt.Sprintf("  Status:       %s\n", n.Status))
	b.WriteString(fmt.Sprintf("  Roles:        %s\n", n.Roles))
	b.WriteString(fmt.Sprintf("  Version:      %s\n", n.Version))
	if n.InternalIP != "" {
		b.WriteString(fmt.Sprintf("  Internal IP:  %s\n", n.InternalIP))
	}
	if n.ExternalIP != "" {
		b.WriteString(fmt.Sprintf("  External IP:  %s\n", n.ExternalIP))
	}
	b.WriteString(fmt.Sprintf("  OS:           %s\n", n.OS))
	b.WriteString(fmt.Sprintf("  Arch:         %s\n", n.Arch))
	if n.CPU != "" {
		b.WriteString(fmt.Sprintf("  CPU:          %s\n", n.CPU))
	}
	if n.Memory != "" {
		b.WriteString(fmt.Sprintf("  Memory:       %s\n", n.Memory))
	}
	if n.Age != "" {
		b.WriteString(fmt.Sprintf("  Age:          %s\n", n.Age))
	}

	b.WriteString("\n  Esc: close")
	return shared.RenderOverlay(b.String())
}
