package zenarmor

import (
	"context"
	"fmt"
	"net/http"

	"github.com/browningluke/opnsense-go/pkg/api"
)

type PolicyNode struct {
	ID string `json:"id"`
}

type Policy struct {
	ID            int64        `json:"id"`
	LocalID       int64        `json:"local_id"`
	CloudPolicyID string       `json:"cloud_policyid"`
	Name          string       `json:"name"`
	IsCentralized bool         `json:"isCentralized"`
	IsActive      bool         `json:"isActive"`
	IsDefault     bool         `json:"isDefault"`
	User          string       `json:"user"`
	Nodes         []PolicyNode `json:"nodes"`
	Tags          []string     `json:"tags"`
	Projects      []string     `json:"projects"`
	Checksum      int64        `json:"checksum"`
}

type PoliciesResponse struct {
	Policies []Policy `json:"policies"`
}

func (c *Controller) Policies(ctx context.Context) ([]Policy, error) {
	result, err := api.Call(c.Client(), ctx, api.RPCOpts{
		BaseEndpoint:    "/zenarmor/policy",
		Method:          http.MethodGet,
		QueryParameters: map[string]string{},
		BodyParameters:  map[string]interface{}{},
	}, &PoliciesResponse{})
	if err != nil {
		return nil, fmt.Errorf("list Zenarmor policies: %w", err)
	}
	if result.Policies == nil {
		return []Policy{}, nil
	}
	return result.Policies, nil
}
