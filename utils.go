package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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
