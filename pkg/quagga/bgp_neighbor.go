package quagga

import (
	"context"
	"github.com/browningluke/opnsense-go/pkg/api"
)

var BGPNeighborOpts = api.ReqOpts{
	AddEndpoint:         "/quagga/bgp/addNeighbor",
	GetEndpoint:         "/quagga/bgp/getNeighbor",
	UpdateEndpoint:      "/quagga/bgp/setNeighbor",
	DeleteEndpoint:      "/quagga/bgp/delNeighbor",
	ReconfigureEndpoint: quaggaReconfigureEndpoint,
	Monad:               "neighbor",
}

// Data structs

type BGPNeighbor struct {
	Enabled               string          `json:"enabled"`
	Description           string          `json:"description"`
	PeerIP                string          `json:"address"`
	RemoteAS              string          `json:"remoteas"`
	Password              string          `json:"password"`
	Weight                string          `json:"weight"`
	LocalIP               string          `json:"localip"`
	UpdateSource          api.SelectedMap `json:"updatesource"`
	LinkLocalInterface    api.SelectedMap `json:"linklocalinterface"`
	NextHopSelf           string          `json:"nexthopself"`
	NextHopSelfAll        string          `json:"nexthopselfall"`
	MultiHop              string          `json:"multihop"`
	MultiProtocol         string          `json:"multiprotocol"`
	RRClient              string          `json:"rrclient"`
	BFD                   string          `json:"bfd"`
	KeepAlive             string          `json:"keepalive"`
	HoldDown              string          `json:"holddown"`
	ConnectTimer          string          `json:"connecttimer"`
	DefaultRoute          string          `json:"defaultoriginate"`
	ASOverride            string          `json:"asoverride"`
	DisableConnectedCheck string          `json:"disable_connected_check"`
	AttributeUnchanged    api.SelectedMap `json:"attributeunchanged"`
	PrefixListIn          api.SelectedMap `json:"linkedPrefixlistIn"`
	PrefixListOut         api.SelectedMap `json:"linkedPrefixlistOut"`
	RouteMapIn            api.SelectedMap `json:"linkedRoutemapIn"`
	RouteMapOut           api.SelectedMap `json:"linkedRoutemapOut"`
}

// CRUD operations

func (c *Controller) AddBGPNeighbor(ctx context.Context, resource *BGPNeighbor) (string, error) {
	return api.Add(c.Client(), ctx, BGPNeighborOpts, resource)
}

func (c *Controller) GetBGPNeighbor(ctx context.Context, id string) (*BGPNeighbor, error) {
	return api.Get(c.Client(), ctx, BGPNeighborOpts, &BGPNeighbor{}, id)
}

func (c *Controller) UpdateBGPNeighbor(ctx context.Context, id string, resource *BGPNeighbor) error {
	return api.Update(c.Client(), ctx, BGPNeighborOpts, resource, id)
}

func (c *Controller) DeleteBGPNeighbor(ctx context.Context, id string) error {
	return api.Delete(c.Client(), ctx, BGPNeighborOpts, id)
}
