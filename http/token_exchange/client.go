package token_exchange

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Errors returned by NewClient when httpClient is nil, when apiURL, publicToken
// or secretToken is empty, or when apiURL is not a usable endpoint.
var (
	ErrNilHTTPClient    = errors.New("token_exchange: nil httpClient")
	ErrEmptyAPIURL      = errors.New("token_exchange: empty apiURL")
	ErrEmptyPublicToken = errors.New("token_exchange: empty publicToken")
	ErrEmptySecretToken = errors.New("token_exchange: empty secretToken")
	ErrInvalidAPIURL    = errors.New("token_exchange: invalid apiURL")
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

// NewClient creates a new client for the Token Exchange API.
//
// httpClient must be non-nil, and apiURL, publicToken and secretToken must be
// non-empty. apiURL must be an absolute http or https URL with a host.
// Only the scheme and host of the given apiURL will be used.
func NewClient(httpClient *http.Client, apiURL string, publicToken string, secretToken string) (*Client, error) {
	if httpClient == nil {
		return nil, ErrNilHTTPClient
	}
	if apiURL == "" {
		return nil, ErrEmptyAPIURL
	}
	if publicToken == "" {
		return nil, ErrEmptyPublicToken
	}
	if secretToken == "" {
		return nil, ErrEmptySecretToken
	}

	u, err := url.Parse(apiURL)
	if err != nil {
		// url.Error repeats the raw apiURL, which may carry credentials. Keep
		// only the reason, consistent with the redaction below.
		var parseErr *url.Error
		if errors.As(err, &parseErr) {
			err = parseErr.Err
		}
		return nil, fmt.Errorf("%w: %w", ErrInvalidAPIURL, err)
	}

	// net/http.Transport.roundTrip requires an http or https scheme. It would
	// also accept an authority like ":123" and dial localhost, so require a
	// hostname rather than just a non-empty Host.
	if (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" {
		return nil, fmt.Errorf(
			"%w: expected an http or https scheme and a host, got %q",
			ErrInvalidAPIURL, u.Redacted(),
		)
	}

	// Take only the scheme and host from apiURL; any path, query, fragment or
	// userinfo it carries is irrelevant to the token exchange endpoint.
	authUrl := url.URL{Scheme: u.Scheme, Host: u.Host, Path: "/api/auth/token"}

	return &Client{
		httpClient:  httpClient,
		authUrl:     authUrl.String(),
		publicToken: publicToken,
		secretToken: secretToken,
	}, nil
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
