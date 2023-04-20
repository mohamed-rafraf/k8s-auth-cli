/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/caarlos0/env/v8"
	"github.com/spf13/cobra"
)

type config struct {
	Home   string `env:"HOME"`
	Server string `env:"K8S-AUTH-SERVER" envDefault:"http://localhost:8080"`
}

var Config config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-auth",
	Short: "A brief description of your application",
	Long: `K8s-auth is a command-line interface (CLI) designed to simplify interaction with the k8s-auth authentication server. 
	
With K8s-auth, you can easily register clusters, create users and groups, and manage their permissions. 
You can also authenticate to any registered cluster, and the CLI will automatically generate the kubeconf file for you.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	err := InitConfig()
	if err != nil {
		os.Exit(1)
	}
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.k8s-auth.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func InitConfig() error {
	err := env.Parse(&Config)
	if err != nil {
		return err
	}
	return nil
}
