package main

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"

	"github.com/ably/ably-go/ably"
	"github.com/labstack/echo/v4"
)


func HandlerIndex(c echo.Context) error {
	c.Logger().Info("index_handler")

	return c.String(http.StatusOK, "alive")
}

func HandlerAblyAuth(c echo.Context) error {
	c.Logger().Info("ably_auth_handler")

	realtime, err := ably.NewREST(ably.WithKey(API_KEY))
	if err != nil {
		return fmt.Errorf("ably client: %w ", err)
	}

	params := &ably.TokenParams{
		ClientID: uuid.NewString(),
	}

	// todo make the ably dependancy interface, seperate interface for
	tokenRequest, err := realtime.Auth.CreateTokenRequest(params)
	if err != nil {
		return fmt.Errorf("ably create token request: %w")
	}

	// maybe set appliucation/json header
	return c.JSON(http.StatusOK, tokenRequest)
}
