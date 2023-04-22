/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// groupCmd represents the group command
var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "create a group inside a specific cluster",
	Long: `This command allow you to create a group in specific cluster!
Keep in mind that you need to upload the roles of this cluster!

Example:
  k8s-auth create group dev --cluster prod-cluster --file rbac.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		cluster, _ := cmd.Flags().GetString("cluster")
		rbac, _ := cmd.Flags().GetString("file")
		if cluster == "" {
			fmt.Println("specify the cluster")
			return
		}

		if rbac == "" {
			fmt.Println("You Must Specify the RBAC file!")
			return
		}
		token, _ := readConfigFile()

		// Open the file
		file, err := os.Open(rbac)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		// Create a buffer to store the request body
		body := &bytes.Buffer{}

		// Create a new multipart writer
		writer := multipart.NewWriter(body)

		// Create a new form file field and add it to the writer
		part, err := writer.CreateFormFile("file", filepath.Base(rbac))
		if err != nil {
			fmt.Println(err)
			return
		}

		// Copy the file to the form field
		_, err = io.Copy(part, file)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Close the multipart writer
		err = writer.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

		url := server + "/groups?name=" + args[0] + "&cluster=" + cluster

		// Create a new request with POST method
		req, err := http.NewRequest("POST", url, body)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Set the Content-Type header to the value returned by the FormDataContentType method
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Add the authentication header and bearer token to the request
		req.Header.Add("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		defer resp.Body.Close()

		// Read the response body
		var response Response
		dat, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(dat, &response)
		if err != nil {
			fmt.Println(err)
		}

		if response.Status == "success" {
			fmt.Println(response.Message)
			return
		} else {
			fmt.Println(response.Message)
			return
		}
	},
}

var groupDeleteCmd = &cobra.Command{
	Use:   "group [name] ",
	Short: "Delete a group from specific cluster",
	Long: `
This command allow you to delete groups in specific cluster.
You must to keep in mind that you can't delete group that have users on it!
You need to move these users or delete them!`,
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		cluster, _ := cmd.Flags().GetString("cluster")

		token, _ := readConfigFile()

		url := server + "/groups?name=" + args[0] + "&cluster=" + cluster
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
	createCmd.AddCommand(groupCmd)
	deleteCmd.AddCommand(groupDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// groupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// groupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	groupDeleteCmd.Flags().String("cluster", "", "Server that you want to deal with")
	groupCmd.Flags().String("cluster", "", "Server that you want to deal with")
	groupCmd.Flags().String("file", "", "Server that you want to deal with")

}
