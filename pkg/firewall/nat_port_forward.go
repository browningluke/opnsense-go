package firewall

import (
	"context"

	"github.com/browningluke/opnsense-go/pkg/api"
)

var NatPortForwardOpts = api.ReqOpts{
	AddEndpoint:         "/firewall/d_nat/addRule",
	GetEndpoint:         "/firewall/d_nat/getRule",
	UpdateEndpoint:      "/firewall/d_nat/setRule",
	DeleteEndpoint:      "/firewall/d_nat/delRule",
	ReconfigureEndpoint: "/firewall/filter/apply",
	Monad:               "rule",
}

// Data structs

// NatPortForwardLocation represents a source or destination block
// in a destination NAT rule. The OPNsense d_nat API returns these
// as nested objects on GET (e.g. {"network": "wanip", "port": "80"}).
type NatPortForwardLocation struct {
	Network string `json:"network"`
	Address string `json:"address"`
	Port    string `json:"port"`
	Invert  string `json:"not"`
}

type NatPortForward struct {
	// Disabled is "0" for enabled, "1" for disabled (matches OPNsense API).
	Disabled      string                 `json:"disabled"`
	NoRDR         string                 `json:"nordr"`
	Sequence      string                 `json:"sequence"`
	Interface     api.SelectedMap        `json:"interface"`
	IPProtocol    api.SelectedMap        `json:"ipprotocol"`
	Protocol      api.SelectedMap        `json:"protocol"`
	Source        NatPortForwardLocation `json:"source"`
	Destination   NatPortForwardLocation `json:"destination"`
	Target        string                 `json:"target"`
	TargetPort    string                 `json:"local-port"`
	PoolOpts      api.SelectedMap        `json:"poolopts"`
	Log           string                 `json:"log"`
	Category      api.SelectedMapList    `json:"category"`
	Description   string                 `json:"descr"`
	Tag           string                 `json:"tag"`
	Tagged        string                 `json:"tagged"`
	NoSync        string                 `json:"nosync"`
	NatReflection api.SelectedMap        `json:"natreflection"`
}

// CRUD operations

func (c *Controller) AddNatPortForward(ctx context.Context, resource *NatPortForward) (string, error) {
	return api.Add(c.Client(), ctx, NatPortForwardOpts, resource)
}

func (c *Controller) GetNatPortForward(ctx context.Context, id string) (*NatPortForward, error) {
	return api.Get(c.Client(), ctx, NatPortForwardOpts, &NatPortForward{}, id)
}

func (c *Controller) UpdateNatPortForward(ctx context.Context, id string, resource *NatPortForward) error {
	return api.Update(c.Client(), ctx, NatPortForwardOpts, resource, id)
}

func (c *Controller) DeleteNatPortForward(ctx context.Context, id string) error {
	return api.Delete(c.Client(), ctx, NatPortForwardOpts, id)
}
