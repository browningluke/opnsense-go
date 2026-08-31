package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/browningluke/opnsense-go/pkg/errs"
)

func withMutex(ctx context.Context, fn func() error) error {
	if err := GlobalMutexKV.Lock(clientMutexKey, ctx); err != nil {
		return err
	}
	defer GlobalMutexKV.Unlock(clientMutexKey, ctx)
	return fn()
}

func set[K any](c *Client, ctx context.Context, opts ReqOpts, resource *K, endpoint string) (string, error) {
	// Since the OPNsense controller has to be reconfigured after every change, locking the mutex prevents
	// the API from being written to while it's reconfiguring, which results in data loss.
	var uuid string
	err := withMutex(ctx, func() error {
		// Wrap resource
		wrapped := map[string](*K){opts.Monad: resource}

		// Make request to OPNsense
		respJson := &addResp{}
		if err := c.doRequest(ctx, "POST", endpoint, wrapped, respJson); err != nil {
			return err
		}

		// Validate result
		if respJson.Result != "saved" {
			return fmt.Errorf("resource not changed. result: %s. errors: %s", respJson.Result, respJson.Validations)
		}

		uuid = respJson.UUID

		// Reconfigure (i.e. restart) the OPNsense service
		return c.ReconfigureService(ctx, opts.ReconfigureEndpoint)
	})
	return uuid, err
}

func get(c *Client, ctx context.Context, endpoint string) (map[string]json.RawMessage, error) {
	// Get generic data
	var reqData map[string]json.RawMessage

	// Make request to OPNsense
	err := c.doRequest(ctx, "GET", endpoint, nil, &reqData)

	// Handle request errors
	if err != nil {
		switch err.(type) {
		case *json.UnmarshalTypeError:
			// Handle unmarshal error (means ID is invalid, or was deleted upstream)
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return reqData, err
}

func Add[K any](c *Client, ctx context.Context, opts ReqOpts, resource *K) (string, error) {
	return set(c, ctx, opts, resource, opts.AddEndpoint)
}

func Update[K any](c *Client, ctx context.Context, opts ReqOpts, resource *K, id string) error {
	_, err := set(c, ctx, opts, resource, fmt.Sprintf("%s/%s", opts.UpdateEndpoint, id))
	return err
}

func Get[K any](c *Client, ctx context.Context, opts ReqOpts, resource *K, id string) (*K, error) {
	// Get resource data
	reqData, err := get(c, ctx, fmt.Sprintf("%s/%s", opts.GetEndpoint, id))
	if err != nil {
		return nil, err
	}

	// Unwrap json
	wrapped, ok := reqData[opts.Monad]
	if !ok || len(wrapped) == 0 {
		// Upstream returned 200 but the response did not include the
		// configured monad — treat as not-found.
		return nil, errs.ErrNotFound
	}
	if err := json.Unmarshal(wrapped, resource); err != nil {
		return nil, err
	}

	return resource, nil
}

func GetFilter[K any](c *Client, ctx context.Context, opts ReqOpts, resource *K, key string) (*K, error) {
	// Get resource data
	reqData, err := get(c, ctx, opts.GetEndpoint)
	if err != nil {
		return nil, err
	}

	// Find key in returned list
	for i, _ := range reqData {
		if i == key {
			if err := json.Unmarshal(reqData[i], resource); err != nil {
				return nil, err
			}
			return resource, nil
		}
	}

	// If loop exits without match, key doesn't exist in list
	return nil, errs.ErrNotFound
}

func GetAll[K any](c *Client, ctx context.Context, opts ReqOpts, resources []K) ([]K, error) {
	// Get resource data
	reqData, err := get(c, ctx, opts.GetEndpoint)
	if err != nil {
		return nil, err
	}

	// Find key in returned list
	for key := range reqData {
		r := new(K)
		if err := json.Unmarshal(reqData[key], r); err != nil {
			return nil, err
		}

		resources = append(resources, *r)
	}

	if len(resources) == 0 {
		// If no resources returned, exit with NotFound error
		return nil, errs.ErrNotFound
	}

	return resources, nil
}

func Delete(c *Client, ctx context.Context, opts ReqOpts, id string) error {
	// Since the OPNsense controller has to be reconfigured after every change, locking the mutex prevents
	// the API from being written to while it's reconfiguring, which results in data loss.
	return withMutex(ctx, func() error {
		respJson := &deleteResp{}
		if err := c.doRequest(ctx, "POST", fmt.Sprintf("%s/%s", opts.DeleteEndpoint, id), nil, respJson); err != nil {
			return err
		}

		// Validate that override was deleted
		if respJson.Result != "deleted" {
			return fmt.Errorf("resource not deleted. result: %s", respJson.Result)
		}

		// Reconfigure (i.e. restart) the OPNsense service
		return c.ReconfigureService(ctx, opts.ReconfigureEndpoint)
	})
}
