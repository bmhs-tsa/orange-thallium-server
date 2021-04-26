package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

//platformConfig represents an esport-platform configuration
type platformConfig struct {
	Enabled bool
	Key     string
}

var platforms map[string]*platformConfig

//platformRoutes registers all platform routes
func platformRoutes() {
	//Parse platforms
	platforms = map[string]*platformConfig{}
	rawPlatforms := viper.GetStringMap("platforms")

	for id, platform := range rawPlatforms {
		//Cast
		config := platform.(map[string]interface{})

		//Update
		platforms[id] = &platformConfig{
			Enabled: config["enabled"].(bool),
			Key:     config["key"].(string),
		}
	}

	//Create a new route group
	group := app.Group("/platforms")

	//Register routes
	group.GET("/all", getPlatforms, auth(viper.GetString("openid.roles.user")))
}

//getPlatforms is used to retrieve a list of all platforms
func getPlatforms(ctx echo.Context) error {
	//Generate a list of all enabled platforms
	platformList := []string{}

	for platformName, platform := range platforms {
		if platform.Enabled {
			platformList = append(platformList, platformName)
		}
	}

	return ctx.JSON(http.StatusOK, platformList)
}
