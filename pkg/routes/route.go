package routes

import (
	"context"
	"fmt"
	"github.com/browningluke/opnsense-go/pkg/api"
)

const (
	routesRouteAddEndpoint    = "/routes/routes/addroute"
	routesRouteGetEndpoint    = "/routes/routes/getroute"
	routesRouteUpdateEndpoint = "/routes/routes/setroute"
	routesRouteDeleteEndpoint = "/routes/routes/delroute"
)

// Data structs

type Route struct {
	Disabled    string          `json:"disabled"`
	Description string          `json:"descr"`
	Gateway     api.SelectedMap `json:"gateway"`
	Network     string          `json:"network"`
}

// CRUD operations

func (c *Controller) AddRoute(ctx context.Context, route *Route) (string, error) {
	return api.MakeSetFunc(c, routesRouteAddEndpoint, routesRouteReconfigureEndpoint)(ctx,
		map[string]*Route{
			"route": route,
		},
	)
}

func (c *Controller) GetRoute(ctx context.Context, id string) (*Route, error) {
	get, err := api.MakeGetFunc(c, routesRouteGetEndpoint,
		&struct {
			Route Route `json:"route"`
		}{},
	)(ctx, id)
	if err != nil {
		return nil, err
	}
	return &get.Route, nil
}

func (c *Controller) UpdateRoute(ctx context.Context, id string, route *Route) error {
	_, err := api.MakeSetFunc(c, fmt.Sprintf("%s/%s", routesRouteUpdateEndpoint, id),
		routesRouteReconfigureEndpoint)(ctx,
		map[string]*Route{
			"route": route,
		},
	)
	return err
}

func (c *Controller) DeleteRoute(ctx context.Context, id string) error {
	return api.MakeDeleteFunc(c, routesRouteDeleteEndpoint, routesRouteReconfigureEndpoint)(ctx, id)
}
