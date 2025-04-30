package main

import (
	"context"
	"fmt"

	"github.com/1password/onepassword-sdk-go"
)

func getTimebutlerCreds(token string) (string, string) {
	// Authenticates with your service account token and connects to 1Password.
	client, err := onepassword.NewClient(context.Background(),
		onepassword.WithServiceAccountToken(token),
		onepassword.WithIntegrationInfo("Today I Work", "v1.0.0"),
	)
	if err != nil {
		panic(err)
	}
	// Retrieves a secret from 1Password.
	// Takes a secret reference as input and returns the secret to which it points.
	vaultID := "ONEKEY"
	itemID := "Timebutler"
	usernameField := "username"
	passwordField := "password"
	username, err := client.Secrets().
		Resolve(context.Background(), fmt.Sprintf("op://%s/%s/%s", vaultID, itemID, usernameField))
	if err != nil {
		panic(err)
	}
	password, err := client.Secrets().
		Resolve(context.Background(), fmt.Sprintf("op://%s/%s/%s", vaultID, itemID, passwordField))
	if err != nil {
		panic(err)

	}
	return username, password
}
