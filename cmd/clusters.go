/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// clusterCmd represents the cluster command
var clustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		token, err := readConfigFile()
		if err != nil {
			fmt.Println(err)
			return
		}

		req, err := http.NewRequest("GET", server+"/clusters", nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// add the authentication header and bearer token to the request
		req.Header.Add("Authorization", "Bearer "+token)

		// create a new HTTP client and send the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		// read the response body
		var clusterResponse ClusterResponse
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(data, &clusterResponse)
		if err != nil {
			fmt.Println(err)
		}
		if clusterResponse.Status == "success" {

			clusters := clusterResponse.Clusters
			// Initialize tabwriter
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

			// Print header
			fmt.Fprintln(w, "NAME\t\tSTATUS\t\tTOKEN\t\tAPISERVER")

			// Print clusters
			for _, c := range clusters {
				var stat string
				if c.Status {
					stat = "Active"
				} else {
					stat = "Not Active"
				}
				fmt.Fprintf(w, "%s\t\t%s\t\t%s\t\t%s\n", c.Name, stat, c.Token, c.API)
			}

			// Flush tabwriter buffer
			w.Flush()
		} else {
			fmt.Println(clusterResponse.Message)
		}

	},
}

func init() {
	getCmd.AddCommand(clustersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clustersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clustersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// clustersCmd.PersistentFlags().String("foo", "", "A help for foo")
	// clustersCmd.PersistentFlags().String("foo", "", "A help for foo")

}
