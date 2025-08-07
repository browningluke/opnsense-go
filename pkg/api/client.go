package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	clientMaxBackoff = 30
	clientMinBackoff = 1
	clientMaxRetries = 4
)

var clientMutexKey = "OPNSENSE"

type Client struct {
	client *retryablehttp.Client
	opts   Options
}

type Options struct {
	Uri           string
	APIKey        string
	APISecret     string
	AllowInsecure bool

	// Retries
	MaxBackoff int64
	MinBackoff int64
	MaxRetries int64

	Logger *log.Logger
}

func NewClient(options Options) *Client {
	httpClient := retryablehttp.NewClient()
	if options.Logger != nil {
		httpClient.Logger = options.Logger
	}
	client := &Client{
		client: httpClient,
		opts:   options,
	}

	// Configure HTTP client
	client.client.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: options.AllowInsecure},
	}

	//   Set defaults for retries
	client.client.RetryWaitMax = clientMaxBackoff
	client.client.RetryWaitMin = clientMinBackoff
	client.client.RetryMax = clientMaxRetries

	//   Override defaults for retries, if set
	if options.MaxBackoff != 0 {
		client.client.RetryWaitMax = time.Duration(options.MaxBackoff) * time.Second
	}
	if options.MinBackoff != 0 {
		client.client.RetryWaitMin = time.Duration(options.MinBackoff) * time.Second
	}
	if options.MaxRetries != 0 {
		client.client.RetryMax = int(options.MaxRetries)
	}

	return client
}

// Requests

func (c *Client) getAuth() string {
	auth := c.opts.APIKey + ":" + c.opts.APISecret
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *Client) doRequest(ctx context.Context, method, endpoint string, body any, resp any) error {

	// set logger
	var logger *log.Logger
	if c.opts.Logger != nil {
		logger = c.opts.Logger
	} else {
		logger = log.Default()
	}

	// Create IO readers
	var bodyReader io.Reader
	if body != nil {
		// Marshal body into bytes
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}

		// Convert body bytes into reader
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create request
	req, err := retryablehttp.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/api%s", c.opts.Uri, endpoint), bodyReader)
	if err != nil {
		return err
	}

	// Add headers
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", c.getAuth()))
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	// Log request
	dReq, _ := httputil.DumpRequest(req.Request, true)
	logger.Println(fmt.Sprintf("\n%s\n", string(dReq)))

	// Do request
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Log response
	dRes, _ := httputil.DumpResponse(res, true)
	logger.Println(ctx, fmt.Sprintf("\n%s\n", string(dRes)))

	// Check for 200
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code non-200; status code %d", res.StatusCode)
	}

	// Unmarshal resp JSON data to struct
	err = json.NewDecoder(res.Body).Decode(resp)
	if err != nil {
		return err
	}

	return nil
}

// ReconfigureService defined at the endpoint.
func (c *Client) ReconfigureService(ctx context.Context, endpoint string) error {
	// Handle services without a reconfigure endpoint
	if endpoint == "" {
		return nil
	}

	// Send reconfigure request to OPNsense
	respJson := &struct {
		Status string `json:"status,omitempty"`
		Result string `json:"result,omitempty"`
	}{}
	err := c.doRequest(ctx, "POST", endpoint, nil, respJson)
	if err != nil {
		return err
	}

	// Since os-wireguard's reconfigure returns {"result":"ok"}, handle both cases
	status := ""
	if respJson.Status != "" {
		status = respJson.Status
	} else if respJson.Result != "" {
		status = respJson.Result
	} else {
		panic(errors.New("reconfigure returned with unknown status response"))
	}

	// Validate service restarted correctly
	status = cases.Lower(language.English).String(
		strings.TrimSpace(status),
	)
	if status != "ok" {
		return fmt.Errorf("reconfigure failed. status: %s", status)
	}

	return nil
}
