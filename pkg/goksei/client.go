package goksei

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/corpix/uarand"
)

var (
	defaultBaseURL   = "https://akses.ksei.co.id/service"
	defaultAuthCache = ".goksei-auth"
)

type Client struct {
	baseURL  string
	username string
	password string

	authCache   string
	token       string
	tokenExpire time.Time
}

func NewClient(username, password string) *Client {
	return &Client{
		baseURL:   defaultBaseURL,
		authCache: defaultAuthCache,

		username: username,
		password: password,
	}
}

func (c *Client) loadToken() error {
	bytes, err := ioutil.ReadFile(c.authCache)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	c.token = string(bytes)

	expire, err := getExpireTime(c.token)
	if err != nil {
		return err
	}

	c.tokenExpire = *expire

	return nil
}

func (c *Client) saveToken(token string, expire time.Time) error {
	c.token = token
	c.tokenExpire = expire

	if err := ioutil.WriteFile(c.authCache, []byte(token), 0700); err != nil {
		return err
	}

	return nil
}

func (c *Client) login() error {
	if err := c.loadToken(); err != nil {
		return err
	}

	if c.token != "" && c.tokenExpire.After(time.Now()) {
		return nil
	}

	if c.username == "" || c.password == "" {
		return fmt.Errorf("username and password are required")
	}

	body, err := json.Marshal(LoginRequest{
		Username: c.username,
		Password: c.password,
		ID:       "1",
		AppType:  "web",
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/login?lang=id", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Referer", "https://akses.ksei.co.id")
	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	var loginResponse LoginResponse

	if err := json.NewDecoder(res.Body).Decode(&loginResponse); err != nil {
		return err
	}

	expire, err := getExpireTime(loginResponse.Validation)
	if err != nil {
		return err
	}

	if c.saveToken(loginResponse.Validation, *expire); err != nil {
		return nil
	}

	return nil
}

func (c *Client) get(path string, dst interface{}) error {
	if err := c.login(); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Referer", "https://akses.ksei.co.id")
	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Header.Set("Authorization", "Bearer "+c.token)

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
		return nil, fmt.Errorf("this method does not accept cash type")
	}

	var response ShareBalanceResponse

	if err := c.get("/myportofolio/summary-detail/"+strings.ToLower(string(portfolioType)), &response); err != nil {
		return nil, err
	}

	return &response, nil
}
