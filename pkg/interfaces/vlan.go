package interfaces

import (
	"context"
	"fmt"
	"github.com/browningluke/opnsense-go/pkg/api"
)

const (
	interfacesVlanReconfigureEndpoint = "/interfaces/vlan_settings/reconfigure"
	interfacesVlanAddEndpoint         = "/interfaces/vlan_settings/addItem"
	interfacesVlanGetEndpoint         = "/interfaces/vlan_settings/getItem"
	interfacesVlanUpdateEndpoint      = "/interfaces/vlan_settings/setItem"
	interfacesVlanDeleteEndpoint      = "/interfaces/vlan_settings/delItem"
)

// Data structs

type Vlan struct {
	Description string          `json:"descr"`
	Tag         string          `json:"tag"`
	Priority    api.SelectedMap `json:"pcp"`
	Parent      api.SelectedMap `json:"if"`
	Device      string          `json:"vlanif"`
}

// CRUD operations

func (c *Controller) AddVlan(ctx context.Context, vlan *Vlan) (string, error) {
	return api.MakeSetFunc(c, interfacesVlanAddEndpoint, interfacesVlanReconfigureEndpoint)(ctx,
		map[string]*Vlan{
			"vlan": vlan,
		},
	)
}

func (c *Controller) GetVlan(ctx context.Context, id string) (*Vlan, error) {
	get, err := api.MakeGetFunc(c, interfacesVlanGetEndpoint,
		&struct {
			Vlan Vlan `json:"vlan"`
		}{},
	)(ctx, id)
	if err != nil {
		return nil, err
	}
	return &get.Vlan, nil
}

func (c *Controller) UpdateVlan(ctx context.Context, id string, vlan *Vlan) error {
	_, err := api.MakeSetFunc(c, fmt.Sprintf("%s/%s", interfacesVlanUpdateEndpoint, id),
		interfacesVlanReconfigureEndpoint)(ctx,
		map[string]*Vlan{
			"vlan": vlan,
		},
	)
	return err
}

func (c *Controller) DeleteVlan(ctx context.Context, id string) error {
	return api.MakeDeleteFunc(c, interfacesVlanDeleteEndpoint, interfacesVlanReconfigureEndpoint)(ctx, id)
}
