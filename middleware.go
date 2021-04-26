package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

//auth is the authentication and authorization middleware
func auth(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			//Get the authorization header
			authorization := ctx.Request().Header.Get("Authorization")

			//Parse the access token
			accessToken, err := parseAuthorization(authorization)

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			//Authenticate and authorize the user
			err = authByIdToken(accessToken, role)

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			return next(ctx)
		}
	}
}
