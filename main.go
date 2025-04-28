package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/1password/onepassword-sdk-go"
	"github.com/PuerkitoBio/goquery"
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

func (t *Timebutler) enterWorkingTime(date, timeEntry string) {
	err := t.browser.Open(fmt.Sprintf("https://timebutler.de/do?ha=zee&ac=1&setstart=%s-0", date))
	if err != nil {
		fmt.Println("Opening Enter working time page failed")
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

func (t *Timebutler) getMissingEntries() map[string]string {
	if err := t.browser.Open("https://timebutler.de/do?ha=zee&ac=41"); err != nil {
		fmt.Println("Opening Calendar view failed")
		panic(err)
	}
	var missing map[string]string = make(map[string]string)

	t.browser.Find("td.day").Each(func(i int, s *goquery.Selection) {
		missingHours := s.Find("span.ov").Text()
		classStr, exists := s.Attr("class")
		if exists == false {
			fmt.Println("Class attribute is missing for a day")
			return
		}
		if strings.HasPrefix(missingHours, "-") &&
			!strings.Contains(classStr, " dea ") {
			for _, class := range strings.Split(classStr, " ") {
				if strings.HasPrefix(class, "sch-") {
					missing[strings.TrimPrefix(class, "sch-")] = missingHours
					return
				}
			}
		}
	})

	return missing
}

func getDateFromDOY(doy string) string {
	parts := strings.Split(doy, "-")

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Println("Failed to cast date class string")
		panic(err)
	}

	nthDay, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println("Failed to cast date class string")
		panic(err)

	}
	// Start from the first day of the year
	startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)

	// Add (nthDay - 1) days to the start of the year
	targetDate := startOfYear.AddDate(0, 0, nthDay-1)

	return targetDate.Format("2006-1-2")
}

func tidyMissingHour(missingHour string) string {
	return strings.ReplaceAll(strings.ReplaceAll(missingHour, " ", ""), "-", "")
}

func main() {
	token := os.Getenv("OP_SERVICE_ACCOUNT_TOKEN")

	// Gets your service account token from the OP_SERVICE_ACCOUNT_TOKEN environment variable.
	timebutlerUsername, timebutlerPassword := getTimebutlerCreds(token)

	// Create a new browser
	bow := getNewBrowser()

	// Login to Timebutler
	timeButler := getNewTimebutler(timebutlerUsername, timebutlerPassword, bow)
	timeButler.login()

	// Fetch the missing entries from the Calendar view
	missing := timeButler.getMissingEntries()

	if len(missing) == 0 {
		fmt.Println("There is nothing to be added to Timebutler!")
	} else {
		for dayOfYear, missingHour := range missing {
			dateStr := getDateFromDOY(dayOfYear)
			missingHourStr := tidyMissingHour(missingHour)
			fmt.Printf("Adding %s to %s\n", missingHourStr, dateStr)
			// Add the missing time entries per day
			timeButler.enterWorkingTime(dateStr, missingHourStr)
		}
	}
}
