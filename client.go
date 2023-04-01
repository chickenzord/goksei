package goksei

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/corpix/uarand"
)

var (
	defaultBaseURL   = "https://akses.ksei.co.id/service"
	defaultAuthCache = ".goksei-auth"
)

type Client struct {
	baseURL string

	authStore AuthStore

	username string
	password string
}

type ClientOpts struct {
	AuthStore AuthStore // directory path to store cached authentication data
	Username  string
	Password  string
}

func NewClient(opts ClientOpts) *Client {
	client := &Client{
		baseURL:   defaultBaseURL,
		authStore: opts.AuthStore,
		username:  opts.Username,
		password:  opts.Password,
	}

	return client
}

func (c *Client) login() (string, error) {
	if c.username == "" || c.password == "" {
		return "", fmt.Errorf("username and password are required")
	}

	body, err := json.Marshal(LoginRequest{
		Username: c.username,
		Password: c.password,
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

	req.Header.Set("Referer", "https://akses.ksei.co.id")
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

func (c *Client) get(path string, dst interface{}) error {
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
		return err
	}

	return nil
}

func (c *Client) SetAuth(username, password string) {
	c.username = username
	c.password = password
}

func (c *Client) GetPortfolioSummary() (*PortfolioSummaryResponse, error) {
	var response PortfolioSummaryResponse

	if err := c.get("/myportofolio/summary", &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) GetShareBalances(portfolioType PortfolioType) (*ShareBalanceResponse, error) {
	if portfolioType == CashType {
		return nil, fmt.Errorf("GetShareBalances does not accept cash type")
	}

	var response ShareBalanceResponse

	if err := c.get("/myportofolio/summary-detail/"+strings.ToLower(string(portfolioType)), &response); err != nil {
		return nil, err
	}

	response.RemoveInvalidData()

	return &response, nil
}

func (c *Client) GetGlobalIdentity() (*GlobalIdentityResponse, error) {
	var identity GlobalIdentityResponse

	if err := c.get("/myaccount/global-identity/", &identity); err != nil {
		return nil, err
	}

	return &identity, nil
}
