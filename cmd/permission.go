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

// permissionCmd represents the permission command
var getPermissionCmd = &cobra.Command{
	Use:   "permission",
	Short: "check permissions/roles for groups/users in cluster",
	Long: `This command allow you to check the permission for group or users in specific cluster
Example:
  k8s-auth get permission --cluster ctf-cluster --group dev
  k8s-auth get permission --cluster ctf-cluster --user mohamedrafraf@gmail.com`,
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		cluster, _ := cmd.Flags().GetString("cluster")
		group, _ := cmd.Flags().GetString("group")
		user, _ := cmd.Flags().GetString("user")

		if cluster == "" {
			fmt.Println("specify the cluster")
			return
		}

		if group == "" && user == "" {
			fmt.Println("You must specify a group or/and user")
			return
		}

		if group != "" && user != "" {
			fmt.Println("You must specify a group or/and user")
			return
		}

		token, err := readConfigFile()
		if err != nil {
			fmt.Println(err)
			return
		}

		var Type string
		var Name string
		if group == "" {
			Type = "user"
			Name = user
		} else {
			Name = group
			Type = "group"
		}

		req, err := http.NewRequest("GET", server+"/permissions?type="+Type+"&cluster="+cluster+"&name="+Name, nil)
		if err != nil {
			fmt.Println("Can't make the request")
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
		var response Response
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(data, &response)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(response.Message)

	},
}

var updatePermissionCmd = &cobra.Command{
	Use:   "permission",
	Short: "update permissions",
	Long: `This command allow you to update permission (will be used for groups/users/clusters in the future)
Example:
  k8s-auth update permission --cluster ctf-cluster --group dev --file RBAC.yaml
  k8s-auth update permission --cluster ctf-cluster --user mohamedrafraf@gmail.com --file RBAC.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		cluster, _ := cmd.Flags().GetString("cluster")
		group, _ := cmd.Flags().GetString("group")
		user, _ := cmd.Flags().GetString("user")
		rbac, _ := cmd.Flags().GetString("file")

		if cluster == "" {
			fmt.Println("specify the cluster")
			return
		}

		if group == "" && user == "" {
			fmt.Println("You must specify a group or/and user")
			return
		}

		if group != "" && user != "" {
			fmt.Println("You must specify a group or/and user")
			return
		}

		token, err := readConfigFile()
		if err != nil {
			fmt.Println(err)
			return
		}

		var Type string
		var Name string
		if group == "" {
			Type = "user"
			Name = user
		} else {
			Name = group
			Type = "group"
		}

		if rbac == "" {
			fmt.Println("You Must Specify the RBAC file!")
			return
		}

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

		url := server + "/permissions?type=" + Type + "&cluster=" + cluster + "&name=" + Name
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
		fmt.Println(response.Message)
	},
}

func init() {
	getCmd.AddCommand(getPermissionCmd)
	updateCmd.AddCommand(updatePermissionCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// permissionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// permissionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getPermissionCmd.Flags().String("cluster", "", "Select the cluster that you want to check it ")
	getPermissionCmd.Flags().String("user", "", "The mail of user")
	getPermissionCmd.Flags().String("group", "", "The group of user")
	updatePermissionCmd.Flags().String("cluster", "", "Select the cluster that you want to check it ")
	updatePermissionCmd.Flags().String("user", "", "The mail of user")
	updatePermissionCmd.Flags().String("group", "", "The group of user")
	updatePermissionCmd.Flags().String("file", "", "The file that have roles for your user/group")
}
