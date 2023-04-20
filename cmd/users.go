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

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "A brief description of your command",
	Long:  `This Command Utility is used to list all users authorized by clusters!`,
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		token, err := readConfigFile()
		if err != nil {
			fmt.Println(err)
			return
		}

		req, err := http.NewRequest("GET", server+"/users", nil)
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
		var userResponse UserResponse
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(data, &userResponse)
		if err != nil {
			fmt.Println(err)
		}
		if userResponse.Status == "success" {

			users := userResponse.ClusterUsers

			// Initialize tabwriter
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

			// Print header
			fmt.Fprintln(w, "NAME\t\tFullName\t\tEMAIL\t\tCLUSTERS")
			var cluster_names []string // Print clusters
			for _, c := range users {
				for _, clu := range c.Clusters {
					cluster_names = append(cluster_names, clu.Name)
				}
				fmt.Fprintf(w, "%s\t\t%s\t\t%s\t\t%s\n", c.Name, c.FullName, c.Email, cluster_names)
				cluster_names = []string{}
			}
			// Flush tabwriter buffer
			w.Flush()
		} else {
			fmt.Println(userResponse.Message)
		}

	},
}

func init() {
	getCmd.AddCommand(usersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	usersCmd.Flags().String("server", Config.Server, "Server that you want to deal with")
}
