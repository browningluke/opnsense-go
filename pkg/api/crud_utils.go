package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/browningluke/opnsense-go/pkg/errs"
)

func resourceWrap[K any](monad string, resource K) map[string]K {
	return map[string]K{
		monad: resource,
	}
}

func resourceUnwrap[K any](monad string, resource K, reqData map[string]json.RawMessage) error {
	wrapped := reqData[monad]

	if err := json.Unmarshal(wrapped, resource); err != nil {
		return err
	}

	return nil
}

func set[K any](c *Client, ctx context.Context, opts ReqOpts, resource *K, endpoint string) (string, error) {
	// Since the OPNsense controller has to be reconfigured after every change, locking the mutex prevents
	// the API from being written to while it's reconfiguring, which results in data loss.
	GlobalMutexKV.Lock(clientMutexKey, ctx)
	defer GlobalMutexKV.Unlock(clientMutexKey, ctx)

	// Wrap resource
	wrapped := resourceWrap(opts.Monad, resource)

	// Make request to OPNsense
	respJson := &addResp{}
	err := c.doRequest(ctx, "POST", endpoint, wrapped, respJson)
	if err != nil {
		return "", err
	}

	// Validate result
	if respJson.Result != "saved" {
		return "", fmt.Errorf("resource not changed. result: %s. errors: %s", respJson.Result, respJson.Validations)
	}

	// Reconfigure (i.e. restart) the OPNsense service
	err = c.ReconfigureService(ctx, opts.ReconfigureEndpoint)
	if err != nil {
		return respJson.UUID, err
	}

	return respJson.UUID, nil
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
			return nil, errs.NewNotFoundError()
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
	err = resourceUnwrap(opts.Monad, resource, reqData)
	if err != nil {
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
	return nil, errs.NewNotFoundError()
}

func Delete(c *Client, ctx context.Context, opts ReqOpts, id string) error {
	// Since the OPNsense controller has to be reconfigured after every change, locking the mutex prevents
	// the API from being written to while it's reconfiguring, which results in data loss.
	GlobalMutexKV.Lock(clientMutexKey, ctx)
	defer GlobalMutexKV.Unlock(clientMutexKey, ctx)

	respJson := &deleteResp{}
	err := c.doRequest(ctx, "POST", fmt.Sprintf("%s/%s", opts.DeleteEndpoint, id), nil, respJson)
	if err != nil {
		return err
	}

	// Validate that override was deleted
	if respJson.Result != "deleted" {
		return fmt.Errorf("resource not deleted. result: %s", respJson.Result)
	}

	// Reconfigure (i.e. restart) the OPNsense service
	err = c.ReconfigureService(ctx, opts.ReconfigureEndpoint)
	if err != nil {
		return err
	}

	return nil
}
