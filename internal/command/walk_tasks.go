package command

import (
	"github.com/mpon/ecswalk/internal/pkg/awsapi"
	"github.com/spf13/cobra"
)

var walkTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "describe ECS tasks by selecting cluster and service interactively",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := awsapi.NewClient()
		if err != nil {
			return err
		}
		listClustersOutput, err := client.ListECSClusters()
		if err != nil {
			return err
		}
		clusterNames := []string{}
		for _, clusterArn := range listClustersOutput.ClusterArns {
			clusterNames = append(clusterNames, awsapi.ShortArn(clusterArn))
		}

		prompt := newPrompt(clusterNames, "Select Cluster")
		_, cluster, err := prompt.Run()
		if err != nil {
			return err
		}

		describeServicesOutputs := client.DescribeAllServices(cluster)
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

func init() {
	rootCmd.AddCommand(walkTasksCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// servicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// servicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
