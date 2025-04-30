# Today I Work!

> Just leave me alone! I want to work and not clicking around Timebutler. OK, I see, Timebutler wants to know about it. Have this and just leave me alone!

## Description

This CLI tool is client for [Timebutler](https://timebutler.de). Uses [surf](https://github.com/headzoo/surf) to impersonate you, fetch your missing hours and submit them as completed work.

Authentication supports:

- Username and password
- 1Password

Username can also be defined in `TIMEBUTLER_USERNAME`, while password in `TIMEBUTLER_PASSWORD` environment variables.
Supports the use of 1Password service token only, which can be set in `OP_SERVICE_ACCOUNT_TOKEN` environment variable.

## Usage

```
Usage:
  todayiwork [flags]

Flags:
  -h, --help                       help for todayiwork
      --op-item string             1Password item name (default "Timebutler")
      --op-password-field string   1Password password field (default "password")
      --op-token string            1Password service token. This flag can be set from OP_SERVICE_ACCOUNT_TOKEN env var.
      --op-username-field string   1Password username field (default "username")
      --op-vault string            1 Password vault name (default "ONEKEY")
      --password string            Timebutler password. If this is set, username must be set too. Overrides the use of 1Password. This flag can be set from TIMEBUTLER_PASSWORD env var.
      --username string            Timebutler username. If this is set, password must be set too. Overrides the use of 1Password. This flag can be set from TIMEBUTLER_USERNAME env var.
```
