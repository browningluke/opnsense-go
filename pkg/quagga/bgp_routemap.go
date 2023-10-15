package quagga

import (
	"context"
	"github.com/browningluke/opnsense-go/pkg/api"
)

var BGPRouteMapOpts = api.ReqOpts{
	AddEndpoint:         "/quagga/bgp/addRoutemap",
	GetEndpoint:         "/quagga/bgp/getRoutemap",
	UpdateEndpoint:      "/quagga/bgp/setRoutemap",
	DeleteEndpoint:      "/quagga/bgp/delRoutemap",
	ReconfigureEndpoint: quaggaReconfigureEndpoint,
	Monad:               "routemap",
}

// Data structs

type BGPRouteMap struct {
	Enabled       string              `json:"enabled"`
	Description   string              `json:"description"`
	Name          string              `json:"name"`
	Action        api.SelectedMap     `json:"action"`
	RouteMapID    string              `json:"id"`
	ASPathList    api.SelectedMapList `json:"match"`
	PrefixList    api.SelectedMapList `json:"match2"`
	CommunityList api.SelectedMapList `json:"match3"`
	Set           string              `json:"set"`
}

// CRUD operations

func (c *Controller) AddBGPRouteMap(ctx context.Context, resource *BGPRouteMap) (string, error) {
	return api.Add(c.Client(), ctx, BGPRouteMapOpts, resource)
}

func (c *Controller) GetBGPRouteMap(ctx context.Context, id string) (*BGPRouteMap, error) {
	return api.Get(c.Client(), ctx, BGPRouteMapOpts, &BGPRouteMap{}, id)
}

func (c *Controller) UpdateBGPRouteMap(ctx context.Context, id string, resource *BGPRouteMap) error {
	return api.Update(c.Client(), ctx, BGPRouteMapOpts, resource, id)
}

func (c *Controller) DeleteBGPRouteMap(ctx context.Context, id string) error {
	return api.Delete(c.Client(), ctx, BGPRouteMapOpts, id)
}
