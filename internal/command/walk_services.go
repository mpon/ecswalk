package command

import (
	"fmt"

	"github.com/mpon/ecswalk/internal/pkg/awsapi"
	"github.com/spf13/cobra"
)

var walkServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "describe ECS services by selecting cluster interactively",
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
			fmt.Printf("Prompt failed %v\n", err)
			return err
		}
		getServicesCmdRun(cluster)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(walkServicesCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// servicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// servicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
