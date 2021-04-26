package main

import (
	"errors"

	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

var vaultNoCredential = errors.New("No credential found")
var vaultNoCredentialValue = errors.New("No value in credential found")

var vaultLogical *api.Logical

func init() {
	//Initialize the Vault client
	vaultClient, err := api.NewClient(&api.Config{
		Address: viper.GetString("vault.address"),
	})

	if err != nil {
		panic(err)
	}

	//Set the authentication token
	vaultClient.SetToken(viper.GetString("vault.token"))

	//Get the logical API manager
	vaultLogical = vaultClient.Logical()
}

//readCredential reads the desired value from a credential in Vault
func readCredential(path string, key string) (map[string]interface{}, error) {
	//Read the credential
	credential, err := vaultLogical.Read(path)

	if err != nil {
		return nil, err
	}

	if credential == nil || credential.Data["data"] == nil {
		return nil, vaultNoCredential
	}

	//Cast credential data
	data := credential.Data["data"].(map[string]interface{})

	if data[key] == nil {
		return data, vaultNoCredentialValue
	}

	return data[key].(map[string]interface{}), nil
}
