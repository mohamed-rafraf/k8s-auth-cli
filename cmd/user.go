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
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "create user in specific cluster",
	Long: `This command allow you to create user in specific cluster.
Keep in mind that you need either specify group for that user or upload his own permission

Example:
  k8s-auth create user --name mohamed --fullname="Mohamed Rafraf" --mail med.raf@gmail.com --group admins --cluster prod-cluster
  k8s-auth create user --name mohamed --fullname="Mohamed Rafraf" --mail med.raf@gmail.com --file rbac.yaml --cluster prod-cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		cluster, _ := cmd.Flags().GetString("cluster")
		rbac, _ := cmd.Flags().GetString("file")
		name, _ := cmd.Flags().GetString("name")
		fullname, _ := cmd.Flags().GetString("fullname")
		mail, _ := cmd.Flags().GetString("mail")
		group, _ := cmd.Flags().GetString("group")

		if cluster == "" {
			fmt.Println("specify the cluster")
			return
		}

		if name == "" {
			fmt.Println("specify the name")
			return
		}

		if fullname == "" {
			fmt.Println("specify the fullname")
			return
		}

		if mail == "" {
			fmt.Println("specify the mail")
			return
		}

		if group == "" && rbac == "" {
			fmt.Println("You must specify a group or upload the rbac file")
			return
		}

		if group != "" && rbac != "" {
			fmt.Println("You must specify a group or upload the rbac file")
			return
		}
		token, _ := readConfigFile()

		// Create a buffer to store the request body
		body := &bytes.Buffer{}
		// Create a new multipart writer
		writer := multipart.NewWriter(body)

		if group == "" && rbac != "" {
			// Open the file
			file, err := os.Open(rbac)
			if err != nil {
				fmt.Println(err)
				return
			}

			defer file.Close()
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
		}

		urls, err := url.Parse(server + "/users?name=" + name + "&cluster=" + cluster + "&fullname=" + fullname + "&mail=" + mail + "&group=" + group)
		// Create a new request with POST method
		if err != nil {
			fmt.Println(err)
			return
		}

		urls.RawQuery = urls.Query().Encode()
		req, err := http.NewRequest("POST", urls.String(), body)
		if err != nil {
			fmt.Println("Can't make the request")
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
			return
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

var userDeleteCmd = &cobra.Command{
	Use:   "user [name] ",
	Short: "Delete user from cluster",
	Long: `
This command allow you to delete users in specific cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		cluster, _ := cmd.Flags().GetString("cluster")
		email, _ := cmd.Flags().GetString("email")

		if email == "" {
			fmt.Println("specify the mail")
			return
		}

		if cluster == "" {
			fmt.Println("specify the cluster")
			return
		}

		token, _ := readConfigFile()

		url := server + "/users?mail=" + email + "&cluster=" + cluster
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

		if response.Status == "success" {
			fmt.Println(response.Message)
			return
		} else {
			fmt.Println(response.Message)
			return
		}
	},
}

func init() {
	createCmd.AddCommand(userCmd)
	deleteCmd.AddCommand(userDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	userDeleteCmd.Flags().String("email", "", "The mail of user")
	userDeleteCmd.Flags().String("cluster", "", "Server that you want to deal with")
	userCmd.Flags().String("cluster", "", "Server that you want to deal with")
	userCmd.Flags().String("file", "", "The permission of user")
	userCmd.Flags().String("name", "", "The name of user")
	userCmd.Flags().String("fullname", "", "The full name of user")
	userCmd.Flags().String("mail", "", "The mail of user")
	userCmd.Flags().String("group", "", "The group of user")
}
