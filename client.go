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

type Client struct {
	baseURL string

	authStore     AuthStore
	username      string
	password      string
	plainPassword bool
}

type ClientOpts struct {
	AuthStore     AuthStore // directory path to store cached authentication data
	Username      string
	Password      string
	PlainPassword bool
}

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
	if c.plainPassword {
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

func (c *Client) SetAuth(username, password string) {
	c.username = username
	c.password = password
}

func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

func (c *Client) SetPlainPassword(plainPassword bool) {
	c.plainPassword = plainPassword
}

func (c *Client) GetPortfolioSummary() (*PortfolioSummaryResponse, error) {
	var response PortfolioSummaryResponse

	if err := c.Get("/myportofolio/summary", &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) GetCashBalances() (*CashBalanceResponse, error) {
	var response CashBalanceResponse

	if err := c.Get("/myportofolio/summary-detail/"+strings.ToLower(string(CashType)), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

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

func (c *Client) GetGlobalIdentity() (*GlobalIdentityResponse, error) {
	var identity GlobalIdentityResponse

	if err := c.Get("/myaccount/global-identity/", &identity); err != nil {
		return nil, err
	}

	return &identity, nil
}
