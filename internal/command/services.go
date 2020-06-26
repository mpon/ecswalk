package command

import (
	"fmt"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mpon/ecswalk/internal/pkg/awsapi"
	"github.com/spf13/cobra"
)

// NewCmdServices represents services cmd
func NewCmdServices() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "services",
		Short: "describe ECS services by selecting cluster interactively",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := awsapi.NewClient()
			if err != nil {
				return err
			}
			output, err := client.DescribeECSClusters()
			if err != nil {
				return err
			}

			if len(output.Clusters) == 0 {
				fmt.Println("cluster not found")
				return nil
			}

			idx, err := fuzzyfinder.Find(output.Clusters,
				func(i int) string {
					return fmt.Sprintf("%s", *output.Clusters[i].ClusterName)
				},
				fuzzyfinder.WithPromptString("Select Cluster:"),
				fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
					cluster := output.Clusters[i]
					return fmt.Sprintf("%s\n\nServices: %d\nRunning Tasks: %d\nPending Tasks: %d",
						*cluster.ClusterName,
						*cluster.ActiveServicesCount,
						*cluster.RunningTasksCount,
						*cluster.PendingTasksCount)
				}),
			)

			if err != nil {
				return nil
			}

			getServicesCmdRun(*output.Clusters[idx].ClusterName)
			return nil
		},
	}
	return cmd
}
