package awsapi

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/spf13/viper"
)

// Client provides each AWS client
type Client struct {
	EC2Client *ec2.Client
	ECSClient *ecs.Client
}

// NewClient is constructor
func NewClient() (*Client, error) {
	var cfg aws.Config
	var err error
	if viper.IsSet("profile") {
		cfg, err = external.LoadDefaultAWSConfig(
			external.WithSharedConfigProfile(viper.GetString("profile")),
		)
	} else {
		cfg, err = external.LoadDefaultAWSConfig()
	}
	if err != nil {
		return nil, err
	}

	c := &Client{
		EC2Client: ec2.New(cfg),
		ECSClient: ecs.New(cfg),
	}

	return c, nil
}
