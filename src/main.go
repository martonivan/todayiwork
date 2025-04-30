package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	op_token          string
	op_vault          string
	op_item           string
	op_username_field string
	op_password_field string
	username          string
	password          string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "todayiwork",
		Short: "CLI to let Timebutler know that you are working whenever you have to",
		Run: func(cmd *cobra.Command, args []string) {
			var timebutlerUsername, timebutlerPassword string

			// Validate required flags manually since we're using viper and get timebutler credentials
			timebutlerUsername = viper.GetString("username")
			timebutlerPassword = viper.GetString("password")
			if timebutlerUsername != "" || timebutlerPassword != "" {
				if timebutlerUsername == "" || timebutlerPassword == "" {
					fmt.Println(
						"Error: If --username is defined --password is also required and vice versa",
					)
					os.Exit(1)
				}
			} else {
				token := viper.GetString("op-token")
				vault := viper.GetString("op-vault")
				item := viper.GetString("op-item")
				usernameField := viper.GetString("op-username-field")
				passwordField := viper.GetString("op-password-field")
				if token == "" || vault == "" || item == "" || usernameField == "" || passwordField == "" {
					fmt.Println("Error: If username and password are not defined, 1Password details must be set")
					os.Exit(1)
				}
				// Gets your service account token from the OP_SERVICE_ACCOUNT_TOKEN environment variable.
				timebutlerUsername, timebutlerPassword = getTimebutlerCreds1Password(token, vault, item, usernameField, passwordField)
			}

			todayIWork(timebutlerUsername, timebutlerPassword)
		},
	}

	// Flags
	rootCmd.Flags().
		StringVar(&op_token, "op-token", "", "1Password service token. This flag can be set from OP_SERVICE_ACCOUNT_TOKEN env var.")
	rootCmd.Flags().StringVar(&op_vault, "op-vault", "ONEKEY", "1 Password vault name")
	rootCmd.Flags().StringVar(&op_item, "op-item", "Timebutler", "1Password item name")
	rootCmd.Flags().
		StringVar(&op_username_field, "op-username-field", "username", "1Password username field")
	rootCmd.Flags().
		StringVar(&op_password_field, "op-password-field", "password", "1Password password field")
	rootCmd.Flags().
		StringVar(&username, "username", "", "Timebutler username. If this is set, password must be set too. Overrides the use of 1Password. This flag can be set from TIMEBUTLER_USERNAME env var.")
	rootCmd.Flags().
		StringVar(&password, "password", "", "Timebutler password. If this is set, username must be set too. Overrides the use of 1Password. This flag can be set from TIMEBUTLER_PASSWORD env var.")

	// Bind flags to viper
	viper.BindPFlag("op-token", rootCmd.Flags().Lookup("op-token"))
	viper.BindPFlag("op-vault", rootCmd.Flags().Lookup("op-vault"))
	viper.BindPFlag("op-item", rootCmd.Flags().Lookup("op-item"))
	viper.BindPFlag("op-username-field", rootCmd.Flags().Lookup("op-username-field"))
	viper.BindPFlag("op-password-field", rootCmd.Flags().Lookup("op-password-field"))
	viper.BindPFlag("username", rootCmd.Flags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.Flags().Lookup("password"))

	// Set ENV var binding
	viper.BindEnv("op-token", "OP_SERVICE_ACCOUNT_TOKEN")
	viper.BindEnv("username", "TIMEBUTLER_USERNAME")
	viper.BindEnv("password", "TIMEBUTLER_PASSWORD")

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
