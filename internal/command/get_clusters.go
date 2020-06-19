package command

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mpon/ecswalk/internal/pkg/awsecs"
	"github.com/spf13/cobra"
)

// clustersCmd represents the clusters command
var getClustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "get ECS clusters",
	Run: func(cmd *cobra.Command, args []string) {
		listClustersOutput := awsecs.ListClusters()
		describeClustersOutput := awsecs.DescribeClusters(listClustersOutput.ClusterArns)

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 1, '\t', 0)
		fmt.Fprintln(w, "Name\tServices\tRunning\tPending\tInstances\t")
		for _, cluster := range describeClustersOutput.Clusters {
			fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t\n",
				*cluster.ClusterName,
				*cluster.ActiveServicesCount,
				*cluster.RunningTasksCount,
				*cluster.PendingTasksCount,
				*cluster.RegisteredContainerInstancesCount)
		}
		w.Flush()
	},
}

func init() {
	getCmd.AddCommand(getClustersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clustersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clustersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
