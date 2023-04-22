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
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// groupsCmd represents the groups command
var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List groups and number of users",
	Long:  `This command allow you to list groups inside the cluster and check the number of users `,
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		cluster, _ := cmd.Flags().GetString("cluster")

		token, err := readConfigFile()
		if err != nil {
			fmt.Println(err)
			return
		}

		req, err := http.NewRequest("GET", server+"/groups?cluster="+cluster, nil)
		if err != nil {
			fmt.Println("Can't make A login request")
			return
		}

		// add the authentication header and bearer token to the request
		req.Header.Add("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Authentication Server is unreachable")
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		var response Response

		err = json.Unmarshal(data, &response)
		if err != nil {
			fmt.Println("Can't Decode the Response received from Server!")
		}
		if response.Status == "success" {
			groups_names := strings.Split(response.Message, ",")
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print header
			fmt.Fprintln(w, "NAME\t\tUSER NUMBER\t\t")
			for _, group_name := range groups_names {

				name := strings.Split(group_name, "-")
				fmt.Fprintf(w, "%s\t\t%s\t\t\n", name[0], name[1])

			}
			w.Flush()
		}
	},
}

func init() {
	getCmd.AddCommand(groupsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// groupsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// groupsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	groupsCmd.Flags().String("cluster", "", "Server that you want to deal with")
}
