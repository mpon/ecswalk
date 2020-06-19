package cmd

// GetTaskRow represents output each row
type GetTaskRow struct {
	TaskID               string
	TaskDefinition       string
	Status               string
	ContainerInstanceArn string
	EC2InstanceID        string
	PrivateIP            string
}

// GetTaskRows slice
type GetTaskRows []*GetTaskRow

func (s GetTaskRows) Len() int {
	return len(s)
}

func (s GetTaskRows) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s GetTaskRows) Less(i, j int) bool {
	return s[i].TaskID < s[j].TaskID
}
