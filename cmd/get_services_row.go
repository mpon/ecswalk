package cmd

// GetServiceRow represents output each row
type GetServiceRow struct {
	Name           string
	TaskDefinition string
	Image          string
	Tag            string
	DesiredCount   int64
	RunningCount   int64
}

// GetServiceRows slice
type GetServiceRows []GetServiceRow

func (s GetServiceRows) Len() int {
	return len(s)
}

func (s GetServiceRows) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s GetServiceRows) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}
