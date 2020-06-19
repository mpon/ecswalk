package awsec2

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/viper"
)

// DescribeInstances to describe instances
func DescribeInstances(instanceIds []string) *ec2.DescribeInstancesOutput {
	svc := newSvc()
	input := &ec2.DescribeInstancesInput{
		InstanceIds: instanceIds,
	}

	req := svc.DescribeInstancesRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
	return result
}

func newSvc() *ec2.EC2 {
	if viper.IsSet("profile") {
		cfg, err := external.LoadDefaultAWSConfig(
			external.WithSharedConfigProfile(viper.GetString("profile")),
		)
		if err != nil {
			panic("failed to load config, " + err.Error())
		}
		return ec2.New(cfg)
	}
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("failed to load config, " + err.Error())
	}
	return ec2.New(cfg)
}
