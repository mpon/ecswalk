package command

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/mpon/ecswalk/internal/pkg/awsapi"
	"github.com/mpon/ecswalk/internal/pkg/fuzzyfinder"
	"github.com/mpon/ecswalk/internal/pkg/sliceutil"
	"github.com/spf13/cobra"
)

// NewCmdInstances represents instances cmd
func NewCmdInstances() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instances",
		Short: "describe ECS container instances by selecting cluster interactively",
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

			cluster, err := fuzzyfinder.FindCluster(output.Clusters)
			if err != nil {
				// Abort fuzzyfinder
				return nil
			}

			res, err := client.DescribeAllContainerInstances(*cluster.ClusterName)
			if err != nil {
				return err
			}

			instances := CreateInstances(res.ContainerInstances)
			ec2List := sliceutil.DistinctSlice(EC2InstanceIDs(instances))
			output2, err := client.DescribeEC2Instances(ec2List)
			if err != nil {
				return err
			}

			for _, reserv := range output2.Reservations {
				for _, i := range reserv.Instances {
					instances.UpdatePrivateIPByInstanceID(*i.PrivateIpAddress, *i.InstanceId)
				}
			}

			rows := GetInstanceRows{}
			for _, c := range res.ContainerInstances {
				var cpuAvailable int64
				var memoryAvailable int64
				for _, r := range c.RemainingResources {
					if *r.Name == "CPU" {
						cpuAvailable = *r.IntegerValue
					}
					if *r.Name == "MEMORY" {
						memoryAvailable = *r.IntegerValue
					}
				}

				rows = append(rows, &GetInstanceRow{
					ContainerInstanceArn: awsapi.ShortArn(*c.ContainerInstanceArn),
					EC2InstanceID:        *c.Ec2InstanceId,
					AgentConnected:       *c.AgentConnected,
					Status:               *c.Status,
					RunningTasksCount:    *c.RunningTasksCount,
					CPUAvailable:         cpuAvailable,
					MemoryAvailable:      memoryAvailable,
					AgentVersion:         *c.VersionInfo.AgentVersion,
					DockerVersion:        strings.Replace(*c.VersionInfo.DockerVersion, "DockerVersion: ", "", -1),
					PrivateIP:            FindPrivateIP(instances, *c.Ec2InstanceId),
				})
			}
			sort.Sort(rows)

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 1, '\t', 0)
			fmt.Fprintln(w, "ContainerInstance\tEC2Instance\tPrivateIP\tConnected\tStatus\tRunning\tCPU\tMemory\tAgent\tDocker")
			for _, row := range rows {
				fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\t%d\t%d\t%d\t%s\t%s\n",
					row.ContainerInstanceArn,
					row.EC2InstanceID,
					row.PrivateIP,
					row.AgentConnected,
					row.Status,
					row.RunningTasksCount,
					row.CPUAvailable,
					row.MemoryAvailable,
					row.AgentVersion,
					row.DockerVersion,
				)
			}
			w.Flush()

			return nil
		},
	}
	return cmd
}
