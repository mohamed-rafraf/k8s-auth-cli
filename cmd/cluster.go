/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	//"os"
	"encoding/json"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "create/register cluster ",
	Long: `This command allow you to create/regsiter cluster.
Once the registration process work successfully, you'll recieve a secret Token that must be used by the controller`,
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		token, _ := readConfigFile()
		url := server + "/clusters?name=" + args[0]
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
			fmt.Println("Error creating request:", err)
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
			fmt.Println("The cluster `" + args[0] + "` is created and this is the Token : " + response.Message)
			return
		} else {
			fmt.Println(response.Message)
			return
		}

	},
}

var clusterDeleteCmd = &cobra.Command{
	Use:   "cluster",
	Short: "create/register cluster ",
	Long: `This command allow you to create/regsiter cluster.
Once the registration process work successfully, you'll recieve a secret Token that must be used by the controller`,
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		token, _ := readConfigFile()
		url := server + "/clusters?name=" + args[0]
		// Create a new request with POST method
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		// add the authentication header and bearer token to the request
		req.Header.Add("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error creating request:", err)
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

		fmt.Println(response.Message)
	},
}

func init() {
	createCmd.AddCommand(clusterCmd)
	deleteCmd.AddCommand(clusterDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clusterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clusterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	clusterDeleteCmd.Flags().String("server", Config.Server, "Server that you want to deal with")
	clusterCmd.Flags().String("server", Config.Server, "Server that you want to deal with")

}
