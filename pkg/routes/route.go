package routes

import (
	"context"
	"github.com/browningluke/opnsense-go/pkg/api"
)

var RouteOpts = api.ReqOpts{
	AddEndpoint:         "/routes/routes/addroute",
	GetEndpoint:         "/routes/routes/getroute",
	UpdateEndpoint:      "/routes/routes/setroute",
	DeleteEndpoint:      "/routes/routes/delroute",
	ReconfigureEndpoint: routesRouteReconfigureEndpoint,
	Monad:               "route",
}

// Data structs

type Route struct {
	Disabled    string          `json:"disabled"`
	Description string          `json:"descr"`
	Gateway     api.SelectedMap `json:"gateway"`
	Network     string          `json:"network"`
}

// CRUD operations

func (c *Controller) AddRoute(ctx context.Context, resource *Route) (string, error) {
	return api.Add(c.Client(), ctx, RouteOpts, resource)
}

func (c *Controller) GetRoute(ctx context.Context, id string) (*Route, error) {
	return api.Get(c.Client(), ctx, RouteOpts, &Route{}, id)
}

func (c *Controller) UpdateRoute(ctx context.Context, id string, resource *Route) error {
	return api.Update(c.Client(), ctx, RouteOpts, resource, id)
}

func (c *Controller) DeleteRoute(ctx context.Context, id string) error {
	return api.Delete(c.Client(), ctx, RouteOpts, id)
}
