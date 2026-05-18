# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

goksei is an unofficial Go client library for AKSES-KSEI (Indonesian Central Securities Depository). It provides programmatic access to portfolio balances, cash balances, and share holdings through the KSEI web API.

## Development Commands

### Building and Testing
- `go build` - Build the project
- `go test` - Run tests
- `go test ./...` - Run all tests recursively

### Running the Example
The example application demonstrates library usage:
```bash
cd ./cmd/example
go run .
```

### Module Management
- `go mod tidy` - Clean up module dependencies
- `go mod download` - Download dependencies

## Code Architecture

### Core Components

**Client (`client.go`)**
- Main HTTP client with authentication handling
- Automatic token refresh and caching via AuthStore
- Support for both plain and pre-hashed passwords
- Base HTTP methods for API communication

**Authentication (`auth.go`)**
- File-based token storage using gokv library
- Automatic token expiration and refresh

**Data Types (`types.go`)**
- Portfolio types: Equity, MutualFund, Cash, Bond, Other
- Response structures for API endpoints
- Balance and share data models

**Static Data (`data.go`)**
- Embedded CSV data for mutual funds and custodian banks
- Lookup functions for fund and bank information
- Uses Go's embed.FS for static file embedding

**Utilities (`helper.go`)**
- JWT token parsing and expiration checking

### Key Patterns

1. **AuthStore Interface**: Abstraction for token persistence (currently file-based)
2. **Embedded Data**: Static reference data embedded in binary using `//go:embed`
3. **Type Safety**: Strong typing for portfolio types and response structures
4. **Error Handling**: Consistent error propagation throughout API calls

### Authentication Flow

1. Client checks AuthStore for cached valid token
2. If no valid token, performs login with username/password
3. For plain passwords, automatically hashes using KSEI's activation endpoint
4. Stores JWT token in AuthStore for future use
5. Automatically refreshes expired tokens

### API Endpoints Structure

All API calls go through the `/service` base path:
- `/login` - Authentication
- `/myportofolio/summary` - Portfolio overview
- `/myportofolio/summary-detail/{type}` - Detailed balances by type
- `/myaccount/global-identity/` - User account information

### Configuration

The example uses environment variables:
- `GOKSEI_USERNAME` - Login username
- `GOKSEI_PASSWORD` - Password (plain or pre-hashed)
- `GOKSEI_PLAIN_PASSWORD=true` - Enable automatic password hashing

## Writing Tests

- Use `github.com/stretchr/testify` for assertions
- Unit tests should not depend on external services, mock accordingly

