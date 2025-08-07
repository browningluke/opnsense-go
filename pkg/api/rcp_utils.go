package api

import (
	"context"
	"encoding/json"
)

func Call[R any](c *Client, ctx context.Context, rcpOpts RPCOpts, result *R) (*R, error) {
	// Get generic data
	var reqData json.RawMessage

	var body interface{}
	if len(rcpOpts.BodyParameters)>0{
		body=rcpOpts.BodyParameters
	}

	err := c.doRequest(ctx, rcpOpts.Method, rcpOpts.EndpointURL(), body, &reqData)

	// Handle request errors
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(reqData, result); err != nil {
		return nil, err
	}

	return result, nil
}
