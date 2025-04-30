package main

import (
	"os"
)

func main() {
	token := os.Getenv("OP_SERVICE_ACCOUNT_TOKEN")

	// Gets your service account token from the OP_SERVICE_ACCOUNT_TOKEN environment variable.
	timebutlerUsername, timebutlerPassword := getTimebutlerCreds(token)

	todayIWork(timebutlerUsername, timebutlerPassword)

}
