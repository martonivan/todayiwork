package main

import "fmt"

func todayIWork(timebutlerUsername, timebutlerPassword string) {
	// Create a new browser
	bow := getNewBrowser()

	// Login to Timebutler
	timeButler := GetNewTimebutler(timebutlerUsername, timebutlerPassword, bow)
	timeButler.Login()

	// Fetch the missing entries from the Calendar view
	missing := timeButler.GetMissingEntries()

	if len(missing) == 0 {
		fmt.Println("There is nothing to be added to Timebutler!")
	} else {
		for dayOfYear, missingHour := range missing {
			dateStr := getDateFromDOY(dayOfYear)
			missingHourStr := tidyMissingHour(missingHour)
			fmt.Printf("Adding %s to %s\n", missingHourStr, dateStr)
			// Add the missing time entries per day
			timeButler.EnterWorkingTime(dateStr, missingHourStr)
		}
	}
}
