package routes

import (
	"github.com/browningluke/opnsense-go/pkg/api"
)

const routesRouteReconfigureEndpoint = "/routes/routes/reconfigure"

// Controller for routes
type Controller struct {
	Api *api.Client
}

func (c *Controller) Client() *api.Client {
	return c.Api
}

func (c *Controller) Name() string {
	return "routes"
}
