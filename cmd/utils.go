package cmd

import (
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func readConfigFile() (string, error) {
	// Get the current user
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	// Create the path to the config file
	configFile := filepath.Join(currentUser.HomeDir, ".k8s-auth")

	// Open the config file
	file, err := os.Open(configFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read data from the config file
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func GenerateKubeConfiguration(token, cacrt, cluster, apiserver string) error {

	clusters := make(map[string]*api.Cluster)
	clusters[cluster] = &api.Cluster{
		Server:                   apiserver,
		CertificateAuthorityData: []byte(cacrt),
	}

	contexts := make(map[string]*api.Context)
	contexts[cluster] = &api.Context{
		Cluster:   cluster,
		Namespace: "default",
		AuthInfo:  "default",
	}

	authinfos := make(map[string]*api.AuthInfo)
	authinfos["default"] = &api.AuthInfo{
		Token: token,
	}

	clientConfig := api.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       clusters,
		Contexts:       contexts,
		CurrentContext: cluster,
		AuthInfos:      authinfos,
	}
	data, err := clientcmd.Write(clientConfig)
	if err != nil {
		return err
	}
	kubeconfig := string(data)
	kubeconfig = strings.Replace(kubeconfig, "|\n", "", -1)

	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	// Create the path to the config file
	configFile := filepath.Join(currentUser.HomeDir, ".k8s-auth.config")

	// Open the config file
	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read data from the config file
	_, err = file.WriteString(kubeconfig)
	if err != nil {
		return err
	}
	return nil

}
