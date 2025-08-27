// Package goksei provides an unofficial Go client library for AKSES-KSEI,
// the Indonesian Central Securities Depository API.
//
// This library allows programmatic access to portfolio information including:
//   - Cash balances across different custodian banks and currencies
//   - Share/security holdings for equities, mutual funds, bonds, and other assets
//   - Portfolio summaries with total values and asset breakdowns
//   - Account identity information
//
// The client handles authentication automatically, including token caching and renewal.
// It supports both plain text passwords (automatically hashed) and pre-hashed passwords.
//
// Basic usage:
//
//	authStore, err := goksei.NewFileAuthStore(".goksei-auth")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	client := goksei.NewClient(goksei.ClientOpts{
//		Username:      "your.email@domain.com",
//		Password:      "your-password",
//		PlainPassword: true,
//		AuthStore:     authStore,
//	})
//
//	summary, err := client.GetPortfolioSummary()
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Security note: This library uses SHA1 for password hashing as required by the KSEI API.
// This is a requirement from KSEI and cannot be changed.
package goksei
