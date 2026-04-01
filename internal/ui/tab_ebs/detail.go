package tab_ebs

import (
	"fmt"
	"strings"

	internalaws "tui-aws/internal/aws"
	"tui-aws/internal/ui/shared"
)

func RenderVolumeDetail(v internalaws.Volume) string {
	var b strings.Builder

	name := v.Name
	if name == "" {
		name = v.ID
	}
	b.WriteString(fmt.Sprintf("  %s\n", name))
	b.WriteString("  ──────────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("  Volume ID:     %s\n", v.ID))
	b.WriteString(fmt.Sprintf("  Name:          %s\n", displayStr(v.Name)))
	b.WriteString(fmt.Sprintf("  State:         %s\n", v.State))
	b.WriteString(fmt.Sprintf("  Type:          %s\n", v.Type))
	b.WriteString(fmt.Sprintf("  Size:          %d GiB\n", v.Size))
	b.WriteString(fmt.Sprintf("  IOPS:          %d\n", v.IOPS))
	if v.Throughput > 0 {
		b.WriteString(fmt.Sprintf("  Throughput:    %d MB/s\n", v.Throughput))
	}
	b.WriteString(fmt.Sprintf("  Encrypted:     %s\n", boolLabel(v.Encrypted)))
	b.WriteString(fmt.Sprintf("  AZ:            %s\n", v.AZ))
	b.WriteString(fmt.Sprintf("  Created:       %s\n", v.CreateTime))
	if v.SnapshotID != "" {
		b.WriteString(fmt.Sprintf("  Snapshot:      %s\n", v.SnapshotID))
	}

	if v.AttachedTo != "" {
		b.WriteString("\n  Attachment:\n")
		b.WriteString(fmt.Sprintf("    Instance:    %s\n", v.AttachedTo))
		b.WriteString(fmt.Sprintf("    Device:      %s\n", v.AttachDevice))
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

func boolLabel(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}
