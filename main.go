package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
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

type loggedTransport struct {
	rt http.RoundTripper
}

func (lt *loggedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("Request:\n%s\n", string(dump))
	return lt.rt.RoundTrip(req)
}

func main() {
	timeEntry := "8h"
	// Gets your service account token from the OP_SERVICE_ACCOUNT_TOKEN environment variable.
	token := os.Getenv("OP_SERVICE_ACCOUNT_TOKEN")

	timebutlerUsername, timebutlerPassword := getTimebutlerCreds(token)

	// Create a new browser and open timebutler login page
	bow := surf.NewBrowser()
	bow.SetUserAgent(agent.Firefox())
	bow.SetAttributes(browser.AttributeMap{
		browser.SendReferer:         false,
		browser.MetaRefreshHandling: true,
		browser.FollowRedirects:     true,
	})
	bow.SetCookieJar(jar.NewMemoryCookies())
	//bow.SetTransport(&loggedTransport{http.DefaultTransport})
	err := bow.Open("https://timebutler.de/login/")
	if err != nil {
		fmt.Println("Opening Login page failed")
		panic(err)
	}

	// Log in to the site.
	fm, err := bow.Form("[id='loginform']")
	if err != nil {
		panic(err)
	}
	fm.Input("login", timebutlerUsername)
	fm.Input("passwort", timebutlerPassword)
	if fm.Submit() != nil {
		fmt.Println("Submitting login information failed")
		panic(err)
	}

	// Navigate to the Enter time entry page
	err = bow.Open("https://timebutler.de/do?ha=zee&ac=1")
	if err != nil {
		fmt.Println("Opening New time entry page failed")
		panic(err)
	}
	// Submit new Entry
	fm, err = bow.Form("[id='formNewEntry']")
	if err != nil {
		panic(err)
	}
	fm.Input("dauer", timeEntry)
	if fm.Submit() != nil {
		fmt.Println("Submitting new time entry failed")
		panic(err)
	}
	fmt.Printf("%s time entry has been recorded to Timebutler\n", timeEntry)
}
