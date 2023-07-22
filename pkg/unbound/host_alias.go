package unbound

import (
	"context"
	"fmt"
	"github.com/browningluke/opnsense-go/pkg/api"
)

const (
	unboundHostAliasAddEndpoint    = "/unbound/settings/addHostAlias"
	unboundHostAliasGetEndpoint    = "/unbound/settings/getHostAlias"
	unboundHostAliasUpdateEndpoint = "/unbound/settings/setHostAlias"
	unboundHostAliasDeleteEndpoint = "/unbound/settings/delHostAlias"
)

// Data structs

type HostAlias struct {
	Enabled     string          `json:"enabled"`
	Host        api.SelectedMap `json:"host"`
	Hostname    string          `json:"hostname"`
	Domain      string          `json:"domain"`
	Description string          `json:"description"`
}

// CRUD operations

func (c *Controller) AddHostAlias(ctx context.Context, resource *HostAlias) (string, error) {
	return api.MakeSetFunc(c, unboundHostAliasAddEndpoint, unboundReconfigureEndpoint)(ctx,
		map[string]*HostAlias{
			"alias": resource,
		},
	)
}

func (c *Controller) GetHostAlias(ctx context.Context, id string) (*HostAlias, error) {
	get, err := api.MakeGetFunc(c, unboundHostAliasGetEndpoint,
		&struct {
			Alias HostAlias `json:"alias"`
		}{},
	)(ctx, id)
	if err != nil {
		return nil, err
	}
	return &get.Alias, nil
}

func (c *Controller) UpdateHostAlias(ctx context.Context, id string, resource *HostAlias) error {
	_, err := api.MakeSetFunc(c, fmt.Sprintf("%s/%s", unboundHostAliasUpdateEndpoint, id),
		unboundReconfigureEndpoint)(ctx,
		map[string]*HostAlias{
			"alias": resource,
		},
	)
	return err
}

func (c *Controller) DeleteHostAlias(ctx context.Context, id string) error {
	return api.MakeDeleteFunc(c, unboundHostAliasDeleteEndpoint, unboundReconfigureEndpoint)(ctx, id)
}
