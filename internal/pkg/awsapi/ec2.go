package awsapi

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// DescribeEC2Instances to describe instances
func (client Client) DescribeEC2Instances(instanceIds []string) (*ec2.DescribeInstancesOutput, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: instanceIds,
	}

	req := client.EC2Client.DescribeInstancesRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return result.DescribeInstancesOutput, nil
}
