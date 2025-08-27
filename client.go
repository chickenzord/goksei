package goksei

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/corpix/uarand"
)

var (
	defaultBaseReferer = "https://akses.ksei.co.id"
	defaultBaseURL     = "https://akses.ksei.co.id/service"
)

// Client provides access to the KSEI (Indonesian Central Securities Depository) API.
// It handles authentication, token management, and provides methods to retrieve
// portfolio information including cash balances, share holdings, and account details.
type Client struct {
	baseURL string

	authStore     AuthStore
	username      string
	password      string
	plainPassword bool
}

// ClientOpts contains configuration options for creating a new Client.
type ClientOpts struct {
	AuthStore     AuthStore // directory path to store cached authentication data
	Username      string
	Password      string
	PlainPassword bool
}

// NewClient creates a new KSEI API client with the provided options.
// The client will use the provided AuthStore for token caching and automatic re-authentication.
func NewClient(opts ClientOpts) *Client {
	client := &Client{
		baseURL:       defaultBaseURL,
		authStore:     opts.AuthStore,
		username:      opts.Username,
		password:      opts.Password,
		plainPassword: opts.PlainPassword,
	}

	return client
}

func (c *Client) hashPassword() (string, error) {
	if !c.plainPassword {
		return c.password, nil
	}

	passwordSHA1 := fmt.Sprintf("%x", sha1.Sum([]byte(c.password)))
	timestamp := time.Now().Unix()
	param := fmt.Sprintf("%s@@!!@@%d", passwordSHA1, timestamp)
	encodedParam := base64.StdEncoding.EncodeToString([]byte(param))

	url := fmt.Sprintf("%s/activation/generated?param=%s", c.baseURL, url.QueryEscape(encodedParam))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating hashed password request: %w", err)
	}

	req.Header.Set("Referer", defaultBaseReferer)
	req.Header.Set("User-Agent", uarand.GetRandom())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error getting hashed password: %w", err)
	}

	var activationResponse struct {
		Code   string `json:"code"`   // e.g. "200"
		Status string `json:"status"` // e.g. "success"
		Data   []struct {
			Pass string `json:"pass"`
		} `json:"data"`
	}

	if err := json.NewDecoder(res.Body).Decode(&activationResponse); err != nil {
		return "", fmt.Errorf("error decoding activation response body: %w", err)
	}

	if len(activationResponse.Data) == 0 {
		return "", fmt.Errorf("no data found in activation response: %v", activationResponse)
	}

	return activationResponse.Data[0].Pass, nil
}

func (c *Client) login() (string, error) {
	if c.username == "" || c.password == "" {
		return "", fmt.Errorf("username and password are required")
	}

	hashedPassword, err := c.hashPassword()
	if err != nil {
		return "", err
	}

	body, err := json.Marshal(LoginRequest{
		Username: c.username,
		Password: hashedPassword,
		ID:       "1",
		AppType:  "web",
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/login?lang=id", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Referer", defaultBaseReferer)
	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	var loginResponse LoginResponse

	if err := json.NewDecoder(res.Body).Decode(&loginResponse); err != nil {
		return "", err
	}

	token := loginResponse.Validation

	if c.authStore != nil {
		if err := c.authStore.Set(c.username, token); err != nil {
			return "", err
		}
	}

	return token, nil
}

func (c *Client) getToken() (string, error) {
	if c.authStore == nil {
		return c.login()
	}

	var token string

	found, err := c.authStore.Get(c.username, &token)
	if err != nil {
		return "", err
	}

	if !found || token == "" {
		return c.login()
	}

	expire, err := getExpireTime(token)
	if err != nil {
		return "", err
	}

	if expire.Before(time.Now()) {
		return c.login()
	}

	return token, nil
}

// Get performs an authenticated GET request to the specified API path and decodes
// the JSON response into dst. It automatically handles authentication and token refresh.
func (c *Client) Get(path string, dst interface{}) error {
	token, err := c.getToken()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Referer", "https://akses.ksei.co.id")
	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(res.Body).Decode(dst); err != nil {
		return fmt.Errorf("error decoding body: %w", err)
	}

	return nil
}

// SetAuth updates the client's authentication credentials.
// This will invalidate any cached tokens and require re-authentication on the next API call.
func (c *Client) SetAuth(username, password string) {
	c.username = username
	c.password = password
}

// SetBaseURL updates the base URL for API requests.
// This is primarily useful for testing or if KSEI changes their API endpoint.
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// SetPlainPassword configures whether the password should be automatically hashed.
// When true, the client will hash plain text passwords using KSEI's hashing service.
// When false, the password is expected to be pre-hashed.
func (c *Client) SetPlainPassword(plainPassword bool) {
	c.plainPassword = plainPassword
}

// GetPortfolioSummary retrieves a summary of all portfolio holdings including
// total values and breakdown by asset type (equity, mutual funds, bonds, etc.).
func (c *Client) GetPortfolioSummary() (*PortfolioSummaryResponse, error) {
	var response PortfolioSummaryResponse

	if err := c.Get("/myportofolio/summary", &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetCashBalances retrieves detailed cash balance information across all accounts,
// including different currencies and custodian banks.
func (c *Client) GetCashBalances() (*CashBalanceResponse, error) {
	var response CashBalanceResponse

	if err := c.Get("/myportofolio/summary-detail/"+strings.ToLower(string(CashType)), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetShareBalances retrieves detailed share/security holdings for the specified portfolio type.
// Valid portfolio types are EquityType, MutualFundType, BondType, and OtherType.
// Use GetCashBalances() for cash holdings instead.
func (c *Client) GetShareBalances(portfolioType PortfolioType) (*ShareBalanceResponse, error) {
	if portfolioType == CashType {
		return nil, fmt.Errorf("GetShareBalances does not accept cash type")
	}

	var response ShareBalanceResponse

	if err := c.Get("/myportofolio/summary-detail/"+strings.ToLower(string(portfolioType)), &response); err != nil {
		return nil, err
	}

	response.RemoveInvalidData()

	return &response, nil
}

// GetGlobalIdentity retrieves detailed account and identity information
// including investor ID, tax numbers, and other personal details.
func (c *Client) GetGlobalIdentity() (*GlobalIdentityResponse, error) {
	var identity GlobalIdentityResponse

	if err := c.Get("/myaccount/global-identity/", &identity); err != nil {
		return nil, err
	}

	return &identity, nil
}
