package command

import (
	"fmt"

	"github.com/mpon/ecswalk/internal/pkg/awsecs"
	"github.com/spf13/cobra"
)

var walkTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "describe ECS tasks by selecting cluster and service interactively",
	Run: func(cmd *cobra.Command, args []string) {
		listClustersOutput := awsecs.ListClusters()
		clusterNames := []string{}
		for _, clusterArn := range listClustersOutput.ClusterArns {
			clusterNames = append(clusterNames, awsecs.ShortArn(clusterArn))
		}

		prompt := newPrompt(clusterNames, "Select Cluster")
		_, cluster, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		describeServicesOutputs := awsecs.DescribeAllServices(cluster)
		serviceNames := []string{}
		for _, describeServiceOutput := range describeServicesOutputs {
			for _, service := range describeServiceOutput.Services {
				serviceNames = append(serviceNames, awsecs.ShortArn(*service.ServiceArn))
			}
		}

		prompt2 := newPrompt(serviceNames, "Select Service")
		_, service, err := prompt2.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		getTasksCmdRun(cluster, service)
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
