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
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

func saveToken(token string) error {
	// Get the current user
	currentUser, err := user.Current()
	if err != nil {
		return err
	}

	// Create the path to the config file
	configFile := filepath.Join(currentUser.HomeDir, ".k8s-auth")

	// Open the config file with write-only access and create it if it doesn't exist
	file, err := os.OpenFile(configFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	// Truncate the file to 0 bytes to remove its previous content
	if err := file.Truncate(0); err != nil {
		return err
	}

	// Write the new token value to the file
	if _, err := file.WriteString(token); err != nil {
		return err
	}

	return nil
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to k8s-auth authentication server",
	Long: `This command allow you to authenticate to the authentication server!
Keep in mind that If your mail didn't registred to the server. You can't login

Example:
  k8s-auth login: login as a normal user who want to authenticate to specific cluster
  k8s-auth login --admin: login as an admin who want to manage the server and his resources(users, groups and clusters)`,
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		admin, _ := cmd.Flags().GetBool("admin")
		var url string
		if admin {
			url = server + "/admin"
		} else {
			url = server + "/login"
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Can't make A login request")
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Authentication Server is unreachable")
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		var login Response

		err = json.Unmarshal(data, &login)
		if err != nil {
			fmt.Println("Can't Decode the Response received from Server!")
		}

		decoded, err := base64.StdEncoding.DecodeString(login.Message)
		if err != nil {
			fmt.Println("Can't Decode the URL Login ")
		}
		fmt.Println("This is the Login URL :\n" + string(decoded))
		fmt.Println("\nAfter Logging, you'll recieve a secret token! You need to paste it here!")
		fmt.Print("\nEnter Token here: ")
		var token string
		fmt.Scanln(&token)
		req, err = http.NewRequest("GET", server+"/clusters", nil)
		fmt.Println("\n\nVerifying token ...")
		if err != nil {
			fmt.Println("Can't make A verification request")
			return
		}

		// add the authentication header and bearer token to the request
		req.Header.Add("Authorization", "Bearer "+token)

		// create a new HTTP client and send the request
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Authentication Server is unreachable")
			return
		}
		// read the response body

		defer resp.Body.Close()
		data, _ = io.ReadAll(resp.Body)

		err = json.Unmarshal(data, &login)
		if err != nil {
			fmt.Println("Can't Decode the Response received from Server!")
		}

		if login.Status == "fail" {
			fmt.Println("\n" + "TOKEN IS NOT VERIFIED")
			return
		}
		fmt.Println("\n" + "TOKEN IS VERIFIED")

		saveToken(token)

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//getCmd.PersistentFlags().String("foo", "", "A help for foo")
	//loginCmd.Flags().String("server", Config.Server, "Server that you want to deal with")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	loginCmd.Flags().BoolP("admin", "a", false, "Login as an adminstrator")
}
