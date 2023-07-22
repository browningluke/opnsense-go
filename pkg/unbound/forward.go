package unbound

import (
	"context"
	"fmt"
	"github.com/browningluke/opnsense-go/pkg/api"
)

const (
	unboundForwardAddEndpoint    = "/unbound/settings/addDot"
	unboundForwardGetEndpoint    = "/unbound/settings/getDot"
	unboundForwardUpdateEndpoint = "/unbound/settings/setDot"
	unboundForwardDeleteEndpoint = "/unbound/settings/delDot"
)

// Data structs

type Forward struct {
	Enabled  string          `json:"enabled"`
	Domain   string          `json:"domain"`
	Type     api.SelectedMap `json:"type"`
	Server   string          `json:"server"`
	Port     string          `json:"port"`
	VerifyCN string          `json:"verify"`
}

// CRUD operations

func (c *Controller) AddForward(ctx context.Context, resource *Forward) (string, error) {
	return api.MakeSetFunc(c, unboundForwardAddEndpoint, unboundReconfigureEndpoint)(ctx,
		map[string]*Forward{
			"dot": resource,
		},
	)
}

func (c *Controller) GetForward(ctx context.Context, id string) (*Forward, error) {
	get, err := api.MakeGetFunc(c, unboundForwardGetEndpoint,
		&struct {
			Dot Forward `json:"dot"`
		}{},
	)(ctx, id)
	if err != nil {
		return nil, err
	}
	return &get.Dot, nil
}

func (c *Controller) UpdateForward(ctx context.Context, id string, resource *Forward) error {
	_, err := api.MakeSetFunc(c, fmt.Sprintf("%s/%s", unboundForwardUpdateEndpoint, id),
		unboundReconfigureEndpoint)(ctx,
		map[string]*Forward{
			"dot": resource,
		},
	)
	return err
}

func (c *Controller) DeleteForward(ctx context.Context, id string) error {
	return api.MakeDeleteFunc(c, unboundForwardDeleteEndpoint, unboundReconfigureEndpoint)(ctx, id)
}
