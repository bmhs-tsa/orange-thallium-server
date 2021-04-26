package main

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var oidcVerifier *oidc.IDTokenVerifier

func init() {
	//Initialize the OpenID client
	oidcProvider, err := oidc.NewProvider(context.Background(), viper.GetString("openid.address"))

	if err != nil {
		panic(err)
	}

	//Initialize the OpenID verifier
	oidcVerifier = oidcProvider.Verifier(&oidc.Config{
		ClientID: viper.GetString("openid.client_id"),
	})

}

//authByIdToken authenticates the user by an Open ID Connect ID token and authorized the user by a role
func authByIdToken(idToken string, requiredRole string) error {
	//Verify the ID token
	token, err := oidcVerifier.Verify(context.Background(), idToken)

	if err != nil {
		return err
	}

	//Unmarshal the raw claims
	var rawClaims interface{}
	err = token.Claims(&rawClaims)

	if err != nil {
		return err
	}

	//Re-marshal the claims (In order to get the bytes)
	claims, err := json.Marshal(rawClaims)

	if err != nil {
		return err
	}

	//Query the roles using the specified path
	roles := gjson.GetBytes(claims, viper.GetString("openid.role_path"))

	//Ensure roles are set
	if !roles.Exists() {
		return errors.New("No OpenID Connect roles are set!")
	}

	//Only allow users with the read role
	for _, role := range roles.Array() {
		if role.String() == requiredRole {
			return nil
		}
	}

	//Catch users without the read role
	return errors.New("You're missing the OpenID connect credential-read role! (Contact an administrator)")
}
