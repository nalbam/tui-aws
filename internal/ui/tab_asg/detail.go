package tab_asg

import (
	"fmt"
	"strings"

	internalaws "tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

func RenderASGDetail(g internalaws.AutoScalingGroup) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("  %s\n", g.Name))
	b.WriteString("  ──────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("  Name:             %s\n", g.Name))
	b.WriteString(fmt.Sprintf("  ARN:              %s\n", g.ARN))
	if g.LaunchConfig != "" {
		b.WriteString(fmt.Sprintf("  Launch Config:    %s\n", g.LaunchConfig))
	}
	if g.LaunchTemplate != "" {
		b.WriteString(fmt.Sprintf("  Launch Template:  %s\n", g.LaunchTemplate))
	}
	b.WriteString(fmt.Sprintf("  Min / Max / Des:  %d / %d / %d\n", g.MinSize, g.MaxSize, g.DesiredCapacity))
	b.WriteString(fmt.Sprintf("  Health Check:     %s\n", g.HealthCheckType))
	b.WriteString(fmt.Sprintf("  Status:           %s\n", g.Status))
	b.WriteString(fmt.Sprintf("  Created:          %s\n", g.CreatedTime))
	b.WriteString(fmt.Sprintf("  AZs:              %s\n", strings.Join(g.AZs, ", ")))

	if len(g.Instances) > 0 {
		b.WriteString("\n  Instances:\n")
		for _, id := range g.Instances {
			b.WriteString(fmt.Sprintf("    %s\n", id))
		}
	}

	if len(g.TargetGroups) > 0 {
		b.WriteString("\n  Target Groups:\n")
		for _, tg := range g.TargetGroups {
			b.WriteString(fmt.Sprintf("    %s\n", tg))
		}
	}

	b.WriteString("\n  Press any key to close")
	return shared.RenderOverlay(b.String())
}
