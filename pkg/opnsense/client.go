package opnsense

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/interfaces"
	"github.com/browningluke/opnsense-go/pkg/routes"
	"github.com/browningluke/opnsense-go/pkg/unbound"
)

// Client defines a client interface for the Proxmox Virtual Environment API.
type Client interface {
	Unbound() *unbound.Controller
	Interfaces() *interfaces.Controller
	Routes() *routes.Controller
}

type client struct {
	a *api.Client
}

// NewClient creates a new API client.
func NewClient(a *api.Client) Client {
	return &client{a: a}
}

func (c *client) Unbound() *unbound.Controller {
	return &unbound.Controller{Api: c.a}
}

func (c *client) Interfaces() *interfaces.Controller {
	return &interfaces.Controller{Api: c.a}
}

func (c *client) Routes() *routes.Controller {
	return &routes.Controller{Api: c.a}
}
