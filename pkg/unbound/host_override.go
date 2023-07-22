package unbound

import (
	"context"
	"fmt"
	"github.com/browningluke/opnsense-go/pkg/api"
)

const (
	unboundHostOverrideAddEndpoint    = "/unbound/settings/addHostOverride"
	unboundHostOverrideGetEndpoint    = "/unbound/settings/getHostOverride"
	unboundHostOverrideUpdateEndpoint = "/unbound/settings/setHostOverride"
	unboundHostOverrideDeleteEndpoint = "/unbound/settings/delHostOverride"
)

// Data structs

type HostOverride struct {
	Enabled     string          `json:"enabled"`
	Hostname    string          `json:"hostname"`
	Domain      string          `json:"domain"`
	Type        api.SelectedMap `json:"rr"`
	Server      string          `json:"server"`
	MXPriority  string          `json:"mxprio"`
	MXDomain    string          `json:"mx"`
	Description string          `json:"description"`
}

// CRUD operations

func (c *Controller) AddHostOverride(ctx context.Context, resource *HostOverride) (string, error) {
	return api.MakeSetFunc(c, unboundHostOverrideAddEndpoint, unboundReconfigureEndpoint)(ctx,
		map[string]*HostOverride{
			"host": resource,
		},
	)
}

func (c *Controller) GetHostOverride(ctx context.Context, id string) (*HostOverride, error) {
	get, err := api.MakeGetFunc(c, unboundHostOverrideGetEndpoint,
		&struct {
			Host HostOverride `json:"host"`
		}{},
	)(ctx, id)
	if err != nil {
		return nil, err
	}
	return &get.Host, nil
}

func (c *Controller) UpdateHostOverride(ctx context.Context, id string, resource *HostOverride) error {
	_, err := api.MakeSetFunc(c, fmt.Sprintf("%s/%s", unboundHostOverrideUpdateEndpoint, id),
		unboundReconfigureEndpoint)(ctx,
		map[string]*HostOverride{
			"host": resource,
		},
	)
	return err
}

func (c *Controller) DeleteHostOverride(ctx context.Context, id string) error {
	return api.MakeDeleteFunc(c, unboundHostOverrideDeleteEndpoint, unboundReconfigureEndpoint)(ctx, id)
}
