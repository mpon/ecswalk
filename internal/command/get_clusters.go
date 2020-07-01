package command

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mpon/ecswalk/internal/pkg/awsapi"
	"github.com/spf13/cobra"
)

// NewCmdGetClusters represents the get cluster command
func NewCmdGetClusters() *cobra.Command {
	return &cobra.Command{
		Use:   "clusters",
		Short: "get ECS clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetClustersCmd()
		},
	}
}

func runGetClustersCmd() error {
	client, err := awsapi.NewClient()
	if err != nil {
		return err
	}
	if err := runGetClusters(client); err != nil {
		return err
	}
	return nil
}

func runGetClusters(client *awsapi.Client) error {
	clusters, err := client.GetAllEcsClusters()
	if err != nil {
		return err
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintln(w, "Name\tServices\tRunning\tPending\tInstances\t")
	for _, cluster := range clusters {
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t\n",
			*cluster.ClusterName,
			*cluster.ActiveServicesCount,
			*cluster.RunningTasksCount,
			*cluster.PendingTasksCount,
			*cluster.RegisteredContainerInstancesCount)
	}
	w.Flush()
	return nil
}
