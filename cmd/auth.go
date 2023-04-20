/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "authenticate to a specific cluster",
	Long: `This command allow you to authenticate to specific! You need to be authorized to that cluster first!
Once you are authorized, you'll be authenticated and a kubeconfig file will be generated for you!`,
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		cluster, _ := cmd.Flags().GetString("cluster")
		if cluster == "" {
			fmt.Println("specify the cluster")
			return
		}
		if cluster == "" {
			fmt.Println(cmd.Long)
			return
		}
		token, _ := readConfigFile()
		url := server + "/auth?cluster=" + cluster

		// Create a new request with POST method
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		// add the authentication header and bearer token to the request
		req.Header.Add("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Authentication Server not Available")
			return
		}

		defer resp.Body.Close()

		// read the response body
		var response Response
		defer resp.Body.Close()
		dat, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(dat, &response)
		if err != nil {
			fmt.Println(err)
		}
		if response.Status == "success" {
			data := strings.Split(string(response.Message), "---\n")
			token := data[0]
			cacrt, err := base64.StdEncoding.DecodeString(data[1])
			if err != nil {
				panic(err)
			}
			apiserver := data[2]
			err = GenerateKubeConfiguration(token, string(cacrt), cluster, apiserver)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("To use this cluster run this command: export KUBECONFIG=$HOME/.k8s-auth.config")
		} else {
			fmt.Println(response.Message)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	authCmd.Flags().String("server", Config.Server, "Server that you want to deal with")
	authCmd.Flags().String("cluster", "", "Server that you want to deal with")
}
