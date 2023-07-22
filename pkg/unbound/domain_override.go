package unbound

import (
	"context"
	"fmt"
	"github.com/browningluke/opnsense-go/pkg/api"
)

const (
	unboundDomainOverrideAddEndpoint    = "/unbound/settings/addDomainOverride"
	unboundDomainOverrideGetEndpoint    = "/unbound/settings/getDomainOverride"
	unboundDomainOverrideUpdateEndpoint = "/unbound/settings/setDomainOverride"
	unboundDomainOverrideDeleteEndpoint = "/unbound/settings/delDomainOverride"
)

// Data structs

type DomainOverride struct {
	Enabled     string `json:"enabled"`
	Domain      string `json:"domain"`
	Server      string `json:"server"`
	Description string `json:"description"`
}

// CRUD operations

func (c *Controller) AddDomainOverride(ctx context.Context, resource *DomainOverride) (string, error) {
	return api.MakeSetFunc(c, unboundDomainOverrideAddEndpoint, unboundReconfigureEndpoint)(ctx,
		map[string]*DomainOverride{
			"domain": resource,
		},
	)
}

func (c *Controller) GetDomainOverride(ctx context.Context, id string) (*DomainOverride, error) {
	get, err := api.MakeGetFunc(c, unboundDomainOverrideGetEndpoint,
		&struct {
			Domain DomainOverride `json:"domain"`
		}{},
	)(ctx, id)
	if err != nil {
		return nil, err
	}
	return &get.Domain, nil
}

func (c *Controller) UpdateDomainOverride(ctx context.Context, id string, resource *DomainOverride) error {
	_, err := api.MakeSetFunc(c, fmt.Sprintf("%s/%s", unboundDomainOverrideUpdateEndpoint, id),
		unboundReconfigureEndpoint)(ctx,
		map[string]*DomainOverride{
			"domain": resource,
		},
	)
	return err
}

func (c *Controller) DeleteDomainOverride(ctx context.Context, id string) error {
	return api.MakeDeleteFunc(c, unboundDomainOverrideDeleteEndpoint, unboundReconfigureEndpoint)(ctx, id)
}
