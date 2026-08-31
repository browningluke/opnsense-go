package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type ReqOpts struct {
	AddEndpoint         string
	GetEndpoint         string
	UpdateEndpoint      string
	DeleteEndpoint      string
	ReconfigureEndpoint string

	Monad string
}

// Response structs
type addResp struct {
	Result      string                 `json:"result"`
	UUID        string                 `json:"uuid"`
	Validations map[string]interface{} `json:"validations,omitempty"`
}

type deleteResp struct {
	Result string `json:"result"`
}

// RCP Options
type RPCOpts struct {
	BaseEndpoint    string
	Method          string
	PathParameters  []string
	QueryParameters map[string]string
	BodyParameters  map[string]interface{}
}

func (p *RPCOpts) EndpointURL() string {
	currentPath := p.BaseEndpoint
	for _, param := range p.PathParameters {
		escapedParam := url.PathEscape(param)

		if currentPath == "" {
			currentPath = escapedParam
		} else if strings.HasSuffix(currentPath, "/") {
			currentPath += escapedParam
		} else {
			currentPath += "/" + escapedParam
		}
	}

	if len(p.QueryParameters) > 0 {
		keys := make([]string, 0, len(p.QueryParameters))
		for k := range p.QueryParameters {
			keys = append(keys, k)
		}
		// Sort so the URL is deterministic — important for any test that
		// asserts on the request URL and for caches/loggers that key on it.
		sort.Strings(keys)

		values := url.Values{}
		for _, k := range keys {
			values.Set(k, p.QueryParameters[k])
		}

		separator := "?"
		if strings.Contains(currentPath, "?") {
			separator = "&"
		}
		currentPath += separator + values.Encode()
	}
	return currentPath
}

func (p *RPCOpts) Body() (string, error) {
	if len(p.BodyParameters) == 0 {
		return "{}", nil
	}
	jsonBytes, err := json.Marshal(p.BodyParameters)
	if err != nil {
		return "", fmt.Errorf("failed to marshal BodyParameters to JSON: %w", err)
	}
	return string(jsonBytes), nil
}

type ActionResult struct {
	Result string `json:"result"`
}

type ReconfigureStatusResult struct {
	Status string `json:"status"`
}
