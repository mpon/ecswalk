package command

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// Version is dynamically set by the toolchain or overridden by the Makefile.
var Version = "DEV"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "ecswalk",
	Short:        fmt.Sprintf("ecswalk version %s", Version),
	SilenceUsage: true,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .ecswalk.yaml, next $HOME/.ecswalk.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cmdVersion := NewCmdVersion()

	cmdGet := NewCmdGet()
	cmdGet.AddCommand(NewCmdGetClusters())
	cmdGet.AddCommand(NewCmdGetServices())
	cmdGet.AddCommand(NewCmdGetTasks())

	cmdServices := NewCmdServices()

	cmdTasks := NewCmdTasks()

	rootCmd.AddCommand(cmdVersion)
	rootCmd.AddCommand(cmdGet)
	rootCmd.AddCommand(cmdServices)
	rootCmd.AddCommand(cmdTasks)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ecswalk" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigName(".ecswalk")
	}

	// viper.AutomaticEnv() // read in environment variables that match
	viper.ReadInConfig() // If a config file is found, read it in.
}

func newPrompt(elements []string, label string) promptui.Select {
	searcher := func(input string, index int) bool {
		cluster := strings.ToLower(elements[index])
		return strings.Contains(cluster, input)
	}

	return promptui.Select{
		Label:    label,
		Items:    elements,
		Size:     20,
		Searcher: searcher,
	}
}
