package aws

import (
	"context"
	"fmt"

	sdkaws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

// Alarm represents a CloudWatch alarm.
type Alarm struct {
	Name               string
	State              string
	MetricName         string
	Namespace          string
	Threshold          float64
	ComparisonOperator string
	Period             string
	Dimensions         map[string]string
	UpdatedTime        string
	AlarmARN           string
	AlarmActions       []string
	OKActions          []string
	InsufficientActions []string
}

// FetchAlarms returns all CloudWatch alarms via DescribeAlarms (paginated).
func FetchAlarms(ctx context.Context, cwClient *cloudwatch.Client) ([]Alarm, error) {
	paginator := cloudwatch.NewDescribeAlarmsPaginator(cwClient, &cloudwatch.DescribeAlarmsInput{})
	var alarms []Alarm
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, a := range page.MetricAlarms {
			dims := make(map[string]string)
			for _, d := range a.Dimensions {
				dims[sdkaws.ToString(d.Name)] = sdkaws.ToString(d.Value)
			}

			updatedTime := ""
			if a.StateUpdatedTimestamp != nil {
				updatedTime = a.StateUpdatedTimestamp.Format("2006-01-02 15:04")
			}

			threshold := 0.0
			if a.Threshold != nil {
				threshold = sdkaws.ToFloat64(a.Threshold)
			}

			period := ""
			if a.Period != nil {
				period = fmt.Sprintf("%ds", sdkaws.ToInt32(a.Period))
			}

			alarms = append(alarms, Alarm{
				Name:               sdkaws.ToString(a.AlarmName),
				State:              string(a.StateValue),
				MetricName:         sdkaws.ToString(a.MetricName),
				Namespace:          sdkaws.ToString(a.Namespace),
				Threshold:          threshold,
				ComparisonOperator: string(a.ComparisonOperator),
				Period:             period,
				Dimensions:         dims,
				UpdatedTime:        updatedTime,
				AlarmARN:           sdkaws.ToString(a.AlarmArn),
				AlarmActions:       a.AlarmActions,
				OKActions:          a.OKActions,
				InsufficientActions: a.InsufficientDataActions,
			})
		}
	}
	return alarms, nil
}
