# goksei

Unofficial client library for AKSES-KSEI

[![Go Reference](https://pkg.go.dev/badge/github.com/chickenzord/goksei.svg)](https://pkg.go.dev/github.com/chickenzord/goksei)
[![Go Report Card](https://goreportcard.com/badge/github.com/chickenzord/goksei)](https://goreportcard.com/report/github.com/chickenzord/goksei)

## Project status

Unstable proof of concept

## Features

- [x] Login with username and (salted) password
- [x] Cache token on disk with auto relogin when expired
- [x] Get balance overview
- [x] Get balance for Equities, Mutual Funds, Bonds, and "Others"
- [x] Get cash balance
- [ ] Command-line interface

## Using as library

Get it as dependency

```sh
go get -u github.com/chickenzord/goksei
```

Example usages:

```go

import "github.com/chickenzord/goksei"

func main() {
	username := "myemail@domain.com"
	password := "myinsecurepassword"
	plainPassword := true

	authStore, err := goksei.NewFileAuthStore(".goksei-auth")
	if err != nil {
		panic(err)
	}

	client := goksei.NewClient(goksei.ClientOpts{
		Username:      username,
		Password:      password,
		PlainPassword: plainPassword,
		AuthStore:     authStore,
	})

	equityBalance, err := client.GetShareBalances(goksei.EquityType)
	if err != nil {
		panic(err)
	}

	// ...
}

```

## Trying out the example

Create `.env` file with following content:

```sh
GOKSEI_USERNAME=youremail@domain.com
GOKSEI_PASSWORD=yoursaltedpassword
```

The salted password can be obtained by logging in with your account on https://akses.ksei.co.id/login and inspect the request payload sent by JS code.

(New feature) Alternatively, you can also supply your plaintext password in `GOKSEI_PASSWORD` then set `GOKSEI_PLAIN_PASSWORD` to `true`. Goksei will automate the hashing process for every login attempts.

```sh
GOKSEI_USERNAME=youremail@domain.com
GOKSEI_PASSWORD=yourplainpassword
GOKSEI_PLAIN_PASSWORD=true
```

Then you can run the example using this command:

```sh
cd ./cmd/example
go run .
```

## Disclaimer

This project is only for personal and educational purpose.
Use on your own risk, there is no guarantee this project will always work when KSEI changed their API or policies.
