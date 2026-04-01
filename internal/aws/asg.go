package aws

import (
	"context"
	"strings"

	sdkaws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
)

// AutoScalingGroup represents an ASG.
type AutoScalingGroup struct {
	Name             string
	ARN              string
	LaunchConfig     string
	LaunchTemplate   string
	MinSize          int
	MaxSize          int
	DesiredCapacity  int
	Instances        []string // instance IDs
	AZs             []string
	HealthCheckType  string // EC2 or ELB
	Status           string
	TargetGroups     []string // TG ARNs
	CreatedTime      string
}

// FetchAutoScalingGroups returns all ASGs via DescribeAutoScalingGroups.
func FetchAutoScalingGroups(ctx context.Context, asgClient *autoscaling.Client) ([]AutoScalingGroup, error) {
	paginator := autoscaling.NewDescribeAutoScalingGroupsPaginator(asgClient, &autoscaling.DescribeAutoScalingGroupsInput{})
	var groups []AutoScalingGroup
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, g := range page.AutoScalingGroups {
			name := sdkaws.ToString(g.AutoScalingGroupName)

			launchConfig := ""
			if g.LaunchConfigurationName != nil {
				launchConfig = sdkaws.ToString(g.LaunchConfigurationName)
			}

			launchTemplate := ""
			if g.LaunchTemplate != nil {
				ltName := sdkaws.ToString(g.LaunchTemplate.LaunchTemplateName)
				ltVer := sdkaws.ToString(g.LaunchTemplate.Version)
				launchTemplate = ltName
				if ltVer != "" {
					launchTemplate += " (" + ltVer + ")"
				}
			}

			instances := make([]string, 0, len(g.Instances))
			for _, inst := range g.Instances {
				instances = append(instances, sdkaws.ToString(inst.InstanceId))
			}

			azs := make([]string, len(g.AvailabilityZones))
			copy(azs, g.AvailabilityZones)

			tgs := make([]string, len(g.TargetGroupARNs))
			copy(tgs, g.TargetGroupARNs)

			createdTime := ""
			if g.CreatedTime != nil {
				createdTime = g.CreatedTime.Format("2006-01-02 15:04")
			}

			status := sdkaws.ToString(g.Status)
			if status == "" {
				status = "InService"
			}

			groups = append(groups, AutoScalingGroup{
				Name:            name,
				ARN:             sdkaws.ToString(g.AutoScalingGroupARN),
				LaunchConfig:    launchConfig,
				LaunchTemplate:  launchTemplate,
				MinSize:         int(sdkaws.ToInt32(g.MinSize)),
				MaxSize:         int(sdkaws.ToInt32(g.MaxSize)),
				DesiredCapacity: int(sdkaws.ToInt32(g.DesiredCapacity)),
				Instances:       instances,
				AZs:             azs,
				HealthCheckType: sdkaws.ToString(g.HealthCheckType),
				Status:          status,
				TargetGroups:    tgs,
				CreatedTime:     createdTime,
			})
		}
	}
	return groups, nil
}

// AZsShort returns a comma-separated short AZ string (just the suffix letter).
func (a *AutoScalingGroup) AZsShort() string {
	var shorts []string
	for _, az := range a.AZs {
		if len(az) > 0 {
			shorts = append(shorts, string(az[len(az)-1]))
		}
	}
	return strings.Join(shorts, ",")
}
