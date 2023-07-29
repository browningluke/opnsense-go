package firewall

import (
	"context"
	"github.com/browningluke/opnsense-go/pkg/api"
)

var FilterOpts = api.ReqOpts{
	AddEndpoint:         "/firewall/filter/addRule",
	GetEndpoint:         "/firewall/filter/getRule",
	UpdateEndpoint:      "/firewall/filter/setRule",
	DeleteEndpoint:      "/firewall/filter/delRule",
	ReconfigureEndpoint: "/firewall/filter/apply",
	Monad:               "rule",
}

// Data structs

type Filter struct {
	Enabled           string              `json:"enabled"`
	Sequence          string              `json:"sequence"`
	Action            api.SelectedMap     `json:"action"`
	Quick             string              `json:"quick"`
	Interface         api.SelectedMapList `json:"interface"`
	Direction         api.SelectedMap     `json:"direction"`
	IPProtocol        api.SelectedMap     `json:"ipprotocol"`
	Protocol          api.SelectedMap     `json:"protocol"`
	SourceNet         string              `json:"source_net"`
	SourcePort        string              `json:"source_port"`
	SourceInvert      string              `json:"source_not"`
	DestinationNet    string              `json:"destination_net"`
	DestinationPort   string              `json:"destination_port"`
	DestinationInvert string              `json:"destination_not"`
	Gateway           api.SelectedMap     `json:"gateway"`
	Log               string              `json:"log"`
	Description       string              `json:"description"`
}

// CRUD operations

func (c *Controller) AddFilter(ctx context.Context, resource *Filter) (string, error) {
	return api.Add(c.Client(), ctx, FilterOpts, resource)
}

func (c *Controller) GetFilter(ctx context.Context, id string) (*Filter, error) {
	return api.Get(c.Client(), ctx, FilterOpts, &Filter{}, id)
}

func (c *Controller) UpdateFilter(ctx context.Context, id string, resource *Filter) error {
	return api.Update(c.Client(), ctx, FilterOpts, resource, id)
}

func (c *Controller) DeleteFilter(ctx context.Context, id string) error {
	return api.Delete(c.Client(), ctx, FilterOpts, id)
}
