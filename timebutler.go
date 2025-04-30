package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf/browser"
)

type Timebutler struct {
	username string
	password string
	browser  *browser.Browser
}

func GetNewTimebutler(username, password string, bow *browser.Browser) Timebutler {
	return Timebutler{
		username: username,
		password: password,
		browser:  bow,
	}
}

func (t *Timebutler) Login() {
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

func (t *Timebutler) EnterWorkingTime(date, timeEntry string) {
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

func (t *Timebutler) GetMissingEntries() map[string]string {
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
