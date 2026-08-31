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
	err := withMutex(ctx, func() error {
		// Get generic data
		var reqData json.RawMessage

		var body interface{}
		if len(rcpOpts.BodyParameters) > 0 {
			body = rcpOpts.BodyParameters
		}

		if err := c.doRequest(ctx, rcpOpts.Method, rcpOpts.EndpointURL(), body, &reqData); err != nil {
			return err
		}

		return json.Unmarshal(reqData, result)
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
