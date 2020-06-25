package command

import (
	"github.com/mpon/ecswalk/internal/pkg/awsapi"
	"github.com/spf13/cobra"
)

// NewCmdTasks represents tasks cmd
func NewCmdTasks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "describe ECS tasks by selecting cluster and service interactively",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := awsapi.NewClient()
			if err != nil {
				return err
			}
			output, err := client.DescribeECSClusters()
			if err != nil {
				return err
			}
			clusterNames := []string{}
			for _, cluster := range output.Clusters {
				clusterNames = append(clusterNames, *cluster.ClusterName)
			}

			prompt := newPrompt(clusterNames, "Select Cluster")
			_, cluster, err := prompt.Run()
			if err != nil {
				return err
			}

			describeServicesOutputs, err := client.DescribeAllECSServices(cluster)
			if err != nil {
				return err
			}
			serviceNames := []string{}
			for _, describeServiceOutput := range describeServicesOutputs {
				for _, service := range describeServiceOutput.Services {
					serviceNames = append(serviceNames, awsapi.ShortArn(*service.ServiceArn))
				}
			}

			prompt2 := newPrompt(serviceNames, "Select Service")
			_, service, err := prompt2.Run()
			if err != nil {
				return err
			}

			getTasksCmdRun(cluster, service)
			return nil
		},
	}
	return cmd
}
