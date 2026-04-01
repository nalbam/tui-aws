package aws

import (
	"context"

	sdkaws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// Volume represents an EBS volume.
type Volume struct {
	ID           string
	Name         string
	State        string
	Type         string
	Size         int // GiB
	IOPS         int
	Throughput   int // MB/s
	Encrypted    bool
	AZ           string
	AttachedTo   string // instance ID
	AttachDevice string // /dev/xvda
	CreateTime   string
	SnapshotID   string
}

// FetchVolumes returns all EBS volumes via DescribeVolumes (paginated).
func FetchVolumes(ctx context.Context, ec2Client *ec2.Client) ([]Volume, error) {
	paginator := ec2.NewDescribeVolumesPaginator(ec2Client, &ec2.DescribeVolumesInput{})
	var volumes []Volume
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range page.Volumes {
			name := ""
			for _, tag := range v.Tags {
				if sdkaws.ToString(tag.Key) == "Name" {
					name = sdkaws.ToString(tag.Value)
					break
				}
			}

			attachedTo := ""
			attachDevice := ""
			if len(v.Attachments) > 0 {
				attachedTo = sdkaws.ToString(v.Attachments[0].InstanceId)
				attachDevice = sdkaws.ToString(v.Attachments[0].Device)
			}

			createTime := ""
			if v.CreateTime != nil {
				createTime = v.CreateTime.Format("2006-01-02 15:04")
			}

			iops := 0
			if v.Iops != nil {
				iops = int(sdkaws.ToInt32(v.Iops))
			}

			throughput := 0
			if v.Throughput != nil {
				throughput = int(sdkaws.ToInt32(v.Throughput))
			}

			volumes = append(volumes, Volume{
				ID:           sdkaws.ToString(v.VolumeId),
				Name:         name,
				State:        string(v.State),
				Type:         string(v.VolumeType),
				Size:         int(sdkaws.ToInt32(v.Size)),
				IOPS:         iops,
				Throughput:   throughput,
				Encrypted:    sdkaws.ToBool(v.Encrypted),
				AZ:           sdkaws.ToString(v.AvailabilityZone),
				AttachedTo:   attachedTo,
				AttachDevice: attachDevice,
				CreateTime:   createTime,
				SnapshotID:   sdkaws.ToString(v.SnapshotId),
			})
		}
	}
	return volumes, nil
}
