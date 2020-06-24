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

var getServicesCmdFlagCluster string

// servicesCmd represents the services command
var getServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "get all ECS services specified cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return getServicesCmdRun(getServicesCmdFlagCluster)
	},
}

func getServicesCmdRun(cluster string) error {
	client, err := awsapi.NewClient()
	if err != nil {
		return err
	}
	describeServicesOutputs, err := client.DescribeAllECSServices(cluster)
	if err != nil {
		return err
	}
	services := []ecs.Service{}
	serviceArns := []string{}
	for _, describeServiceOutput := range describeServicesOutputs {
		for _, service := range describeServiceOutput.Services {
			services = append(services, service)
			serviceArns = append(serviceArns, *service.ServiceArn)
		}
	}
	describeTaskDefinitionOutputs, err := client.DescribeTaskDefinitions(cluster, serviceArns)
	if err != nil {
		return err
	}

	rows := GetServiceRows{}
	for _, service := range services {
		td := ecs.TaskDefinition{}
		for _, describeTaskDefinitionOutput := range describeTaskDefinitionOutputs {
			if *service.TaskDefinition == *describeTaskDefinitionOutput.TaskDefinition.TaskDefinitionArn {
				td = *describeTaskDefinitionOutput.TaskDefinition
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

func init() {
	getCmd.AddCommand(getServicesCmd)
	getServicesCmd.Flags().StringVarP(&getServicesCmdFlagCluster, "cluster", "c", "", "AWS ECS cluster")
	getServicesCmd.MarkFlagRequired("cluster")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// servicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// servicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
