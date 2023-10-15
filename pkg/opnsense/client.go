package opnsense

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/firewall"
	"github.com/browningluke/opnsense-go/pkg/interfaces"
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/browningluke/opnsense-go/pkg/routes"
	"github.com/browningluke/opnsense-go/pkg/unbound"
	"github.com/browningluke/opnsense-go/pkg/wireguard"
)

// Client defines a client interface for the Proxmox Virtual Environment API.
type Client interface {
	Unbound() *unbound.Controller
	Wireguard() *wireguard.Controller
	Quagga() *quagga.Controller
	Interfaces() *interfaces.Controller
	Routes() *routes.Controller
	Firewall() *firewall.Controller
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

func (c *client) Wireguard() *wireguard.Controller {
	return &wireguard.Controller{Api: c.a}
}

func (c *client) Quagga() *quagga.Controller {
	return &quagga.Controller{Api: c.a}
}

func (c *client) Interfaces() *interfaces.Controller {
	return &interfaces.Controller{Api: c.a}
}

func (c *client) Routes() *routes.Controller {
	return &routes.Controller{Api: c.a}
}

func (c *client) Firewall() *firewall.Controller {
	return &firewall.Controller{Api: c.a}
}
