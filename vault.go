package main

import (
	"encoding/json"

	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

var vaultLogical *api.Logical

func init() {
	//Vault config
	config := &api.Config{
		Address:   viper.GetString("vault.address"),
		SRVLookup: viper.GetBool("vault.srv_lookup"),
	}

	err := config.ConfigureTLS(&api.TLSConfig{
		CACert:        viper.GetString("vault.ca_cert"),
		CAPath:        viper.GetString("vault.ca_path"),
		ClientCert:    viper.GetString("vault.client_cert"),
		ClientKey:     viper.GetString("vault.client_key"),
		Insecure:      viper.GetBool("vault.insecure"),
		TLSServerName: viper.GetString("vault.sni_host"),
	})

	if err != nil {
		panic(err)
	}

	//Initialize the Vault client
	vaultClient, err := api.NewClient(config)

	if err != nil {
		panic(err)
	}

	//Set the authentication token
	vaultClient.SetToken(viper.GetString("vault.token"))

	//Get the logical API manager
	vaultLogical = vaultClient.Logical()
}

type vaultError struct {
	Message  string   `json:"message"`
	Warnings []string `json:"warnings"`
}

func (err vaultError) Error() string {
	//Marshal the error
	raw, _ := json.Marshal(err)

	//Convert to string
	errorString := string(raw)

	return errorString
}

//readCredential reads the desired value from a credential in Vault
func readCredential(path string, key string) (map[string]interface{}, error) {
	//Read the credential
	credential, err := vaultLogical.Read(path)

	if err != nil {
		return nil, err
	}

	if credential == nil || credential.Data["data"] == nil {
		return nil, vaultError{
			Message:  "No credential found! (Verify the Vault secret path)",
			Warnings: credential.Warnings,
		}
	}

	//Cast credential data
	data := credential.Data["data"].(map[string]interface{})

	if data[key] == nil {
		return data, vaultError{
			Message:  "Invalid account! (Verify the account ID)",
			Warnings: credential.Warnings,
		}
	}

	return data[key].(map[string]interface{}), nil
}

//writeCredential writes the desired value to a credential in Vault
func writeCredential(path string, key string, data map[string]interface{}) error {
	//Read the credential
	credential, err := vaultLogical.Read(path)

	if err != nil {
		return err
	}

	if credential == nil || credential.Data["data"] == nil {
		return vaultError{
			Message:  "No credential found! (Verify the Vault secret path)",
			Warnings: credential.Warnings,
		}
	}

	//Update the credential
	credential.Data["data"].(map[string]interface{})[key] = data

	//Write the credential
	_, err = vaultLogical.Write(path, credential.Data)

	return err
}
