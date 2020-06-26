package command

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/ktr0731/go-fuzzyfinder"
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

			if len(output.Clusters) == 0 {
				fmt.Println("cluster not found")
				return nil
			}

			idx, _ := fuzzyfinder.Find(output.Clusters,
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

			cluster := output.Clusters[idx]
			describeServicesOutputs, err := client.DescribeAllECSServices(*cluster.ClusterName)
			if err != nil {
				return err
			}

			services := []ecs.Service{}
			for _, describeServiceOutput := range describeServicesOutputs {
				for _, service := range describeServiceOutput.Services {
					services = append(services, service)
				}
			}

			if len(services) == 0 {
				fmt.Printf("%s has no services\n", *cluster.ClusterName)
				return nil
			}

			idx2, _ := fuzzyfinder.Find(services,
				func(i int) string {
					s := services[i]
					return fmt.Sprintf("%s", *s.ServiceName)
				},
				fuzzyfinder.WithPromptString("Select Service:"),
				fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
					s := services[i]
					return fmt.Sprintf("%s\n\nTask Definition: %s\nDesired tasks: %d\nRunning tasks: %d",
						*s.ServiceName,
						awsapi.ShortArn(*s.TaskDefinition),
						*s.DesiredCount,
						*s.RunningCount,
					)
				}),
			)

			getTasksCmdRun(*cluster.ClusterName, *services[idx2].ServiceName)
			return nil
		},
	}
	return cmd
}
