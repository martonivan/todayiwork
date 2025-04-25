package main

import (
	"context"
	"fmt"
	"os"

	"github.com/1password/onepassword-sdk-go"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"github.com/headzoo/surf/jar"
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

func getNewBrowser() *browser.Browser {
	bow := surf.NewBrowser()
	bow.SetUserAgent(agent.Firefox())
	bow.SetAttributes(browser.AttributeMap{
		browser.SendReferer:         false,
		browser.MetaRefreshHandling: true,
		browser.FollowRedirects:     true,
	})
	bow.SetCookieJar(jar.NewMemoryCookies())
	return bow
}

type Timebutler struct {
	username string
	password string
	browser  *browser.Browser
}

func getNewTimebutler(username, password string, bow *browser.Browser) Timebutler {
	return Timebutler{
		username: username,
		password: password,
		browser:  bow,
	}
}

func (t *Timebutler) login() {
	err := t.browser.Open("https://timebutler.de/login/")
	if err != nil {
		fmt.Println("Opening Login page failed")
		panic(err)
	}

	// Log in to the site.
	fm, err := t.browser.Form("[id='loginform']")
	if err != nil {
		panic(err)
	}
	fm.Input("login", t.username)
	fm.Input("passwort", t.password)
	if fm.Submit() != nil {
		fmt.Println("Submitting login information failed")
		panic(err)
	}

}

func (t *Timebutler) addTimeEntry(timeEntry string) {
	err := t.browser.Open("https://timebutler.de/do?ha=zee&ac=1")
	if err != nil {
		fmt.Println("Opening New time entry page failed")
		panic(err)
	}
	// Submit new Entry
	fm, err := t.browser.Form("[id='formNewEntry']")
	if err != nil {
		panic(err)
	}
	fm.Input("dauer", timeEntry)
	if fm.Submit() != nil {
		fmt.Println("Submitting new time entry failed")
		panic(err)
	}
}

func main() {
	timeEntry := "8h"
	token := os.Getenv("OP_SERVICE_ACCOUNT_TOKEN")

	// Gets your service account token from the OP_SERVICE_ACCOUNT_TOKEN environment variable.
	timebutlerUsername, timebutlerPassword := getTimebutlerCreds(token)

	// Create a new browser
	bow := getNewBrowser()

	// Login to Timebutler
	timeButler := getNewTimebutler(timebutlerUsername, timebutlerPassword, bow)
	timeButler.login()

	// Navigate to the Enter time entry page
	timeButler.addTimeEntry(timeEntry)
	fmt.Printf("%s time entry has been recorded to Timebutler\n", timeEntry)
}
