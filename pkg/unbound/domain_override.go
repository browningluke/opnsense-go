package unbound

import (
	"context"
	"github.com/browningluke/opnsense-go/pkg/api"
)

var DomainOverrideOpts = api.ReqOpts{
	AddEndpoint:         "/unbound/settings/addDomainOverride",
	GetEndpoint:         "/unbound/settings/getDomainOverride",
	UpdateEndpoint:      "/unbound/settings/setDomainOverride",
	DeleteEndpoint:      "/unbound/settings/delDomainOverride",
	ReconfigureEndpoint: unboundReconfigureEndpoint,
	Monad:               "domain",
}

// Data structs

type DomainOverride struct {
	Enabled     string `json:"enabled"`
	Domain      string `json:"domain"`
	Server      string `json:"server"`
	Description string `json:"description"`
}

// CRUD operations

func (c *Controller) AddDomainOverride(ctx context.Context, resource *DomainOverride) (string, error) {
	return api.Add(c.Client(), ctx, DomainOverrideOpts, resource)
}

func (c *Controller) GetDomainOverride(ctx context.Context, id string) (*DomainOverride, error) {
	return api.Get(c.Client(), ctx, DomainOverrideOpts, &DomainOverride{}, id)
}

func (c *Controller) UpdateDomainOverride(ctx context.Context, id string, resource *DomainOverride) error {
	return api.Update(c.Client(), ctx, DomainOverrideOpts, resource, id)
}

func (c *Controller) DeleteDomainOverride(ctx context.Context, id string) error {
	return api.Delete(c.Client(), ctx, DomainOverrideOpts, id)
}
