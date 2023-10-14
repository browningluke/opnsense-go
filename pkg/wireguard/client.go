package wireguard

import (
	"context"
	"github.com/browningluke/opnsense-go/pkg/api"
)

var ClientOpts = api.ReqOpts{
	AddEndpoint:         "/wireguard/client/addClient",
	GetEndpoint:         "/wireguard/client/getClient",
	UpdateEndpoint:      "/wireguard/client/setClient",
	DeleteEndpoint:      "/wireguard/client/delClient",
	ReconfigureEndpoint: wireguardReconfigureEndpoint,
	Monad:               "client",
}

// Data structs

type Client struct {
	Enabled       string              `json:"enabled"`
	Name          string              `json:"name"`
	PublicKey     string              `json:"pubkey"`
	PSK           string              `json:"psk"`
	TunnelAddress api.SelectedMapList `json:"tunneladdress"`
	ServerAddress string              `json:"serveraddress"`
	ServerPort    string              `json:"serverport"`
	KeepAlive     string              `json:"keepalive"`
}

// CRUD operations

func (c *Controller) AddClient(ctx context.Context, resource *Client) (string, error) {
	return api.Add(c.Client(), ctx, ClientOpts, resource)
}

func (c *Controller) GetClient(ctx context.Context, id string) (*Client, error) {
	return api.Get(c.Client(), ctx, ClientOpts, &Client{}, id)
}

func (c *Controller) UpdateClient(ctx context.Context, id string, resource *Client) error {
	return api.Update(c.Client(), ctx, ClientOpts, resource, id)
}

func (c *Controller) DeleteClient(ctx context.Context, id string) error {
	return api.Delete(c.Client(), ctx, ClientOpts, id)
}
