package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	gommonLog "github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

//Initialize echo
var app = echo.New()

func init() {
	//Viper config
	viper.AddConfigPath("config")
	viper.SetConfigType("toml")

	// Use the local config if present otherwise use the default config
	if exists("config/local.toml") {
		viper.SetConfigName("local")
	} else {
		viper.SetConfigName("default")
	}

	//Read the config
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	//Enable environment variables
	viper.AutomaticEnv()
}

func main() {
	//Get environment
	env := viper.GetString("env")

	//Configure echo
	app.Server.IdleTimeout = time.Minute
	app.Server.ReadTimeout = time.Second
	app.Server.WriteTimeout = 2 * time.Second
	app.HideBanner = true
	app.Validator = &CustomValidator{
		Validator: validator.New(),
	}

	//Request ID
	app.Pre(middleware.RequestID())

	//Log level
	level := gommonLog.WARN

	if env == "development" {
		level = gommonLog.DEBUG
	}

	//Logger middleware
	app.Logger.SetLevel(level)
	app.Logger.SetHeader(`{"time":"${time_rfc3339}"}`)
	app.Pre(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339}","id":"${id}","remote_ip":"${remote_ip}","method":"${method}","path":"${path}","status":"${status}"}` + "\n",
	}))

	//Security Middleware
	app.Pre(middleware.SecureWithConfig(middleware.SecureConfig{
		//CSP
		ContentSecurityPolicy: "default-src 'none';",

		//HSTS for 30 days
		HSTSMaxAge: 60 * 60 * 24 * 30,
	}))

	//Remove trailing slashes
	app.Pre(middleware.RemoveTrailingSlash())

	//Error recoverer
	if env != "development" {
		app.Use(middleware.Recover())
	}

	//Error handler
	app.HTTPErrorHandler = func(err error, ctx echo.Context) {
		//Get time
		time := time.Now().Format(time.RFC3339)

		//Get request ID
		id := ctx.Response().Header().Get(echo.HeaderXRequestID)

		//Generate message
		msg := fmt.Sprintf(`{"time":"%s","id":"%s","error":"%s"}`+"\n", time, id, err.Error())

		//Log
		_, internalError := ctx.Logger().Output().Write([]byte(msg))

		if internalError != nil {
			panic(internalError)
		}

		//Send
		app.DefaultHTTPErrorHandler(err, ctx)
	}

	//Add routes
	credentialRoutes()
	platformRoutes()

	//Start echo
	address := viper.GetString("http.address")
	if viper.GetBool("http.tls") {
		//Get the cert and key
		cert := viper.GetString("http.cert")
		key := viper.GetString("http.key")

		//Ensure cert and key exist
		if !exists(cert) {
			fmt.Printf("TLS certificate %s does not exist!", cert)
			os.Exit(1)
		}

		if !exists(key) {
			fmt.Printf("TLS key %s does not exist!", key)
			os.Exit(1)
		}

		//Enforce minimum TLS version
		app.Server.TLSConfig = &tls.Config{
			//Enforce minimum TLS version
			MinVersion: tls.VersionTLS13,
		}

		//Start the server
		app.Logger.Fatal(app.StartTLS(address, cert, key))
	} else {
		//Start the server
		app.Logger.Fatal(app.Start(address))
	}
}

//Check wether a file exists
func exists(filename string) bool {
	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}
