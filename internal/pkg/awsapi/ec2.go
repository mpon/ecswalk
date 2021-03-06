package awsapi

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// GetEc2Instances to describe instances
func (client Client) GetEc2Instances(instanceIds []string) ([]ec2.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: instanceIds,
	}

	req := client.Ec2Client.DescribeInstancesRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}

	var instances []ec2.Instance

	for _, r := range result.DescribeInstancesOutput.Reservations {
		instances = append(instances, r.Instances...)
	}
	return instances, nil
}
