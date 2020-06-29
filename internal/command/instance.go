package command

import "github.com/aws/aws-sdk-go-v2/service/ecs"

// Instance represents Container Instance and EC2 instance information
type Instance struct {
	ContainerInstanceArn string
	EC2InstanceID        string
	PrivateIP            string
}

// Instances represents slice of pointerInstance
type Instances []*Instance

// UpdateEC2InstanceIDByArn to update EC2 Instance Id
func (Instances Instances) UpdateEC2InstanceIDByArn(instanceID string, arn string) {
	for _, data := range Instances {
		if data.ContainerInstanceArn == arn {
			data.EC2InstanceID = instanceID
		}
	}
}

// UpdatePrivateIPByInstanceID to update Private Ip address
func (Instances Instances) UpdatePrivateIPByInstanceID(address string, instanceID string) {
	for _, data := range Instances {
		if data.EC2InstanceID == instanceID {
			data.PrivateIP = address
		}
	}
}

// NewInstances to create with containerInstanceArns
func NewInstances(containerInstanceArns []string) Instances {
	Instances := Instances{}
	for _, arn := range containerInstanceArns {
		Instances = append(Instances, &Instance{
			ContainerInstanceArn: arn,
		})
	}
	return Instances
}

// CreateInstances to create with containerInstances
func CreateInstances(containerInstances []ecs.ContainerInstance) Instances {
	Instances := Instances{}
	for _, c := range containerInstances {
		Instances = append(Instances, &Instance{
			ContainerInstanceArn: *c.ContainerInstanceArn,
			EC2InstanceID:        *c.Ec2InstanceId,
		})
	}
	return Instances
}

// EC2InstanceIDs return ec2InstanceID list
func EC2InstanceIDs(instances Instances) []string {
	var ids []string
	for _, i := range instances {
		ids = append(ids, i.EC2InstanceID)
	}
	return ids
}

// FindPrivateIP find private IP Addres
func FindPrivateIP(instances Instances, instanceID string) string {
	for _, i := range instances {
		if i.EC2InstanceID == instanceID {
			return i.PrivateIP
		}
	}
	return ""
}
