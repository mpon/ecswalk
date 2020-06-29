package command

// GetInstanceRow represents output each row
type GetInstanceRow struct {
	ContainerInstanceArn string
	EC2InstanceID        string
	AgentConnected       bool
	Status               string
	RunningTasksCount    int64
	CPUAvailable         int64
	MemoryAvailable      int64
	AgentVersion         string
	DockerVersion        string
	PrivateIP            string
}

// GetInstanceRows slice
type GetInstanceRows []*GetInstanceRow

func (s GetInstanceRows) Len() int {
	return len(s)
}

func (s GetInstanceRows) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s GetInstanceRows) Less(i, j int) bool {
	return s[i].ContainerInstanceArn < s[j].ContainerInstanceArn
}
