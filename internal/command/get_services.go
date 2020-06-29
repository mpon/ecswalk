package command

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mpon/ecswalk/internal/pkg/awsapi"
	"github.com/spf13/cobra"
)

// NewCmdGetServices represents the get services command
func NewCmdGetServices() *cobra.Command {
	var getServicesCmdFlagCluster string
	cmd := &cobra.Command{
		Use:   "services",
		Short: "get all ECS services specified cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetServicesCmd(getServicesCmdFlagCluster)
		},
	}
	cmd.Flags().StringVarP(&getServicesCmdFlagCluster, "cluster", "c", "", "AWS ECS cluster")
	cmd.MarkFlagRequired("cluster")

	return cmd
}

func runGetServicesCmd(clusterName string) error {
	client, err := awsapi.NewClient()
	if err != nil {
		return err
	}

	cluster, err := client.GetECSCluster(clusterName)
	if err != nil {
		return err
	}

	return runGetServices(client, cluster)
}

func runGetServices(client *awsapi.Client, cluster *ecs.Cluster) error {
	services, err := client.GetAllECSServices(cluster)
	if err != nil {
		return nil
	}

	if len(services) == 0 {
		fmt.Printf("%s has no services\n", *cluster.ClusterName)
		return nil
	}

	taskDefinitions, err := client.GetECSTaskDefinitions(cluster, services)
	if err != nil {
		return err
	}

	rows := GetServiceRows{}
	for _, service := range services {
		td := ecs.TaskDefinition{}
		for _, t := range taskDefinitions {
			if *service.TaskDefinition == *t.TaskDefinitionArn {
				td = t
			}
		}
		for _, containerDefinition := range td.ContainerDefinitions {
			image, tag := awsapi.ShortDockerImage(*containerDefinition.Image)
			rows = append(rows, GetServiceRow{
				Name:           *service.ServiceName,
				TaskDefinition: awsapi.ShortArn(*td.TaskDefinitionArn),
				Image:          image,
				Tag:            tag,
				DesiredCount:   *service.DesiredCount,
				RunningCount:   *service.RunningCount,
			})
		}
	}
	sort.Sort(rows)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintln(w, "Name\tTaskDefinition\tImage\tTag\tDesired\tRunning\t")
	for _, row := range rows {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%d\t\n",
			row.Name,
			row.TaskDefinition,
			row.Image,
			row.Tag,
			row.DesiredCount,
			row.RunningCount)
	}
	w.Flush()
	return nil
}
