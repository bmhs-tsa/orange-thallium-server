package main

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

//credentialRoutes registers all credential routes
func credentialRoutes() {
	//Create a new route group
	group := app.Group("/credentials")

	//Register routes
	group.GET("/:platform/:accountID", getCredential, auth(viper.GetString("openid.roles.user")))
}

//getCredentialResponse represents a credential fetch request body
type getCredentialResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
	//MfaCode  string `json:"mfaCode,omitempty"`
}

//getCredential is used to get a credential
func getCredential(ctx echo.Context) error {
	//Get the platform
	platform := platforms[ctx.Param("platform")]

	//Make sure the platform is valid and enabled
	if platform == nil || !platform.Enabled {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid or disabled platform")
	}

	//Decode the account ID
	accountID, err := url.QueryUnescape(ctx.Param("accountID"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	//Read the credential from Vault
	credential, err := readCredential(platform.Key, accountID)

	if err == vaultNoCredential {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	//Generate the response
	res := getCredentialResponse{
		Username: credential["username"].(string),
		Password: credential["password"].(string),
	}

	//MFA code
	/*if credential["mfaSecret"] != nil && len(credential["mfaSecret"].(string)) > 0 {
		//Compute the desired time
		advance := time.Duration(viper.GetInt64("mfa.time_shift"))
		desiredTime := time.Now().Add(time.Second * time.Duration(advance))

		//Generate the current MFA code
		mfaCode, err := totp.GenerateCode(credential["mfaSecret"].(string), desiredTime)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		res.MfaCode = mfaCode
	}*/

	return ctx.JSON(http.StatusOK, res)
}
