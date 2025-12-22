package token_exchange

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type tokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	UserID       string `json:"user_id"`
}

type Client struct {
	httpClient             *http.Client
	authUrl                string
	publicToken            string
	secretToken            string
	accessToken            string
	accessTokenExpiration  time.Time
	refreshToken           string
	refreshTokenExpiration time.Time
	lock                   sync.Mutex
}

// NewClient creates a new client for the Token Exchange API
func NewClient(httpClient *http.Client, apiURL string, publicToken string, secretToken string) *Client {
	client := &Client{
		httpClient:  httpClient,
		authUrl:     "https://app.underline.com/api/auth/token",
		publicToken: publicToken,
		secretToken: secretToken,
	}

	if apiURL != "" {
		u, err := url.Parse(apiURL)
		if err != nil {
			log.Printf("Error parsing url for Client Token Exchange: %s using default %s", err, client.authUrl)
			return client
		}

		u.Path = "/api/auth/token"

		client.authUrl = u.String()
	}

	return client
}

// GetAccessToken gets the access token for the client
// if access token isn't expired use it
// if access token is expired and refresh token isn't expired use it to get a new access token
// if access token is expired and refresh token is expired use api token to get a new refresh token and access token
func (c *Client) GetAccessToken() (token string, err error) {
	return c.GetAccessTokenCtx(context.Background())
}

func (c *Client) GetAccessTokenCtx(ctx context.Context) (token string, err error) {
	defer func() {
		if err != nil {
			slog.Error(
				"Error getting access token",
				slog.String("package", "toolkit-go/http/token_exchange"),
				slog.String("refresh_token_expiry", c.refreshTokenExpiration.String()),
				slog.String("access_token_expiry", c.accessTokenExpiration.String()),
				slog.Any("error", err),
			)
		}
	}()

	// Since multiple requests can be made at the same time we are going to take a lock while in here in case another request could be trying to refresh the token
	c.lock.Lock()
	defer c.lock.Unlock()

	// If access token exists and won't be expired in 30 seconds return it
	if c.accessToken != "" && c.accessTokenExpiration.After(time.Now().Add(30*time.Second)) {
		return c.accessToken, err
	}

	// Values to use if we need to get a new refesh token
	requestValues := url.Values{
		"grant_type":   {"api_token"},
		"public_token": {c.publicToken},
		"secret_token": {c.secretToken},
	}

	// If refresh token exits and won't be expired in 30 seconds use it to get a new access token
	if c.refreshToken != "" && c.refreshTokenExpiration.After(time.Now().Add(30*time.Second)) {
		requestValues = url.Values{
			"grant_type":    {"refresh_token"},
			"refresh_token": {c.refreshToken},
		}
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.authUrl,
		strings.NewReader(requestValues.Encode()),
	)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := c.httpClient.Do(req)
	if err != nil {
		return
	}

	if err = c.handleAccessTokenResponse(response); err != nil {
		return
	}

	return c.accessToken, nil
}

func (c *Client) handleAccessTokenResponse(response *http.Response) error {
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New(string(body))
	}

	var token tokenResponse

	if err := json.Unmarshal(body, &token); err != nil {
		return err
	}

	c.refreshToken = token.RefreshToken
	c.refreshTokenExpiration = time.Now().Add(16 * time.Hour)
	c.accessToken = token.AccessToken
	c.accessTokenExpiration = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	return nil
}
