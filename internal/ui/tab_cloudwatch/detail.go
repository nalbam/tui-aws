package tab_cloudwatch

import (
	"fmt"
	"strings"

	internalaws "tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

func RenderAlarmDetail(a internalaws.Alarm) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("  %s\n", a.Name))
	b.WriteString("  ──────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("  Name:           %s\n", a.Name))
	b.WriteString(fmt.Sprintf("  State:          %s\n", a.State))
	b.WriteString(fmt.Sprintf("  Metric:         %s\n", a.MetricName))
	b.WriteString(fmt.Sprintf("  Namespace:      %s\n", a.Namespace))
	b.WriteString(fmt.Sprintf("  Threshold:      %.2f\n", a.Threshold))
	b.WriteString(fmt.Sprintf("  Comparison:     %s\n", a.ComparisonOperator))
	b.WriteString(fmt.Sprintf("  Period:         %s\n", a.Period))
	if a.UpdatedTime != "" {
		b.WriteString(fmt.Sprintf("  Last Updated:   %s\n", a.UpdatedTime))
	}
	if a.AlarmARN != "" {
		b.WriteString(fmt.Sprintf("  ARN:            %s\n", a.AlarmARN))
	}

	if len(a.Dimensions) > 0 {
		b.WriteString("\n  Dimensions:\n")
		for k, v := range a.Dimensions {
			b.WriteString(fmt.Sprintf("    %s: %s\n", k, v))
		}
	}

	if len(a.AlarmActions) > 0 {
		b.WriteString("\n  Alarm Actions:\n")
		for _, action := range a.AlarmActions {
			b.WriteString(fmt.Sprintf("    %s\n", action))
		}
	}

	if len(a.OKActions) > 0 {
		b.WriteString("\n  OK Actions:\n")
		for _, action := range a.OKActions {
			b.WriteString(fmt.Sprintf("    %s\n", action))
		}
	}

	if len(a.InsufficientActions) > 0 {
		b.WriteString("\n  Insufficient Data Actions:\n")
		for _, action := range a.InsufficientActions {
			b.WriteString(fmt.Sprintf("    %s\n", action))
		}
	}

	b.WriteString("\n  Press any key to close")
	return shared.RenderOverlay(b.String())
}
