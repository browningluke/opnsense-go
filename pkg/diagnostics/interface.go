package diagnostics

import (
	"context"
	"github.com/browningluke/opnsense-go/pkg/api"
)

var InterfaceOpts = api.ReqOpts{
	AddEndpoint:         "",
	GetEndpoint:         "/diagnostics/interface/getInterfaceConfig",
	UpdateEndpoint:      "",
	DeleteEndpoint:      "",
	ReconfigureEndpoint: "",
	Monad:               "",
}

// Data structs

type Ipv4Config struct {
	IpAddr     string `json:"ipaddr"`
	SubnetBits int64  `json:"subnetbits"`
	Tunnel     bool   `json:"tunnel"`
}

type Ipv6Config struct {
	IpAddr     string `json:"ipaddr"`
	SubnetBits int64  `json:"subnetbits"`
	Tunnel     bool   `json:"tunnel"`
	Autoconf   bool   `json:"autoconf"`
	Deprecated bool   `json:"deprecated"`
	LinkLocal  bool   `json:"link-local"`
	Tentative  bool   `json:"tentative"`
}

type Interface struct {
	Device     string `json:"device"`
	Media      string `json:"media"`
	MediaRaw   string `json:"media_raw"`
	MacAddr    string `json:"macaddr"`
	IsPhysical bool   `json:"is_physical"`
	MTU        string `json:"mtu"`
	Status     string `json:"status"`

	Flags          []string `json:"flags"`
	Capabilities   []string `json:"capabilities"`
	Options        []string `json:"options"`
	SupportedMedia []string `json:"supported_media"`
	Groups         []string `json:"groups"`

	Ipv4 []Ipv4Config `json:"ipv4"`
	Ipv6 []Ipv6Config `json:"ipv6"`
}

// CRUD operations

func (c *Controller) GetInterface(ctx context.Context, id string) (*Interface, error) {
	return api.GetFilter(c.Client(), ctx, InterfaceOpts, &Interface{}, id)
}
