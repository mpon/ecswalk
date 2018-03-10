// Copyright © 2018 Masato Oshima
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
