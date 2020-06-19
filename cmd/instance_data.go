package cmd

// InstanceData represents Container Instance and EC2 instance information
type InstanceData struct {
	ContainerInstanceArn string
	EC2InstanceID        string
	PrivateIP            string
}

// InstanceDatas represents slice of pointerInstanceData
type InstanceDatas []*InstanceData

// UpdateEC2InstanceIDByArn to update EC2 Instance Id
func (instanceDatas InstanceDatas) UpdateEC2InstanceIDByArn(instanceID string, arn string) {
	for _, data := range instanceDatas {
		if data.ContainerInstanceArn == arn {
			data.EC2InstanceID = instanceID
		}
	}
}

// UpdatePrivateIPByInstanceID to update Private Ip address
func (instanceDatas InstanceDatas) UpdatePrivateIPByInstanceID(address string, instanceID string) {
	for _, data := range instanceDatas {
		if data.EC2InstanceID == instanceID {
			data.PrivateIP = address
		}
	}
}

// NewInstanceDatas to create with containerInstanceArns
func NewInstanceDatas(containerInstanceArns []string) InstanceDatas {
	instanceDatas := InstanceDatas{}
	for _, arn := range containerInstanceArns {
		instanceDatas = append(instanceDatas, &InstanceData{
			ContainerInstanceArn: arn,
		})
	}
	return instanceDatas
}
