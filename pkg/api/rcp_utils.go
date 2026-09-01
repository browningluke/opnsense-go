package api

import (
	"context"
	"encoding/json"
)

func Call[R any](c *Client, ctx context.Context, rcpOpts RPCOpts, result *R) (*R, error) {
	// Serialise every RPC against the same global mutex the CRUD helpers use.
	// Some RPC endpoints mutate config or apply a reconfigure; running them in
	// parallel with a CRUD set/Delete (which holds this lock for its own
	// reconfigure window) is the data-loss scenario the global mutex exists to
	// prevent.
	if err := GlobalMutexKV.Lock(clientMutexKey, ctx); err != nil {
		return nil, err
	}
	defer GlobalMutexKV.Unlock(clientMutexKey, ctx)

	// Get generic data
	var reqData json.RawMessage

	var body interface{}
	if len(rcpOpts.BodyParameters) > 0 {
		body = rcpOpts.BodyParameters
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
