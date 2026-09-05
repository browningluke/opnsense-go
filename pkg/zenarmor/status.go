package zenarmor

import (
	"context"
	"fmt"
	"net/http"

	"github.com/browningluke/opnsense-go/pkg/api"
)

type VersionInfo struct {
	Version    string `json:"version"`
	LastUpdate int64  `json:"lastUpdate"`
}

type DatabaseInfo struct {
	Version string `json:"version"`
}

type DatabaseStatus struct {
	Info   DatabaseInfo `json:"info"`
	Status bool         `json:"status"`
}

type ServiceStatus struct {
	Status bool `json:"status"`
}

type InterfaceInfo struct {
	Interface string   `json:"interface"`
	Tags      []string `json:"tags"`
}

// Status is the stable subset currently exposed by the local OPNsense
// /api/zenarmor/status controller. Additional response fields are deliberately
// ignored until they are needed by a typed consumer.
type Status struct {
	AgentInstalled   bool            `json:"agent_installed"`
	AgentEnabled     bool            `json:"agent_enabled"`
	AgentStatus      bool            `json:"agent_status"`
	AgentVersion     VersionInfo     `json:"agent_version"`
	Engine           VersionInfo     `json:"engine"`
	DatabaseVersion  VersionInfo     `json:"database_version"`
	Database         DatabaseStatus  `json:"database"`
	Eastpect         ServiceStatus   `json:"eastpect"`
	License          string          `json:"license"`
	CloudThreatIntel bool            `json:"cloud_threat_intel"`
	DeploymentMode   string          `json:"deployment_mode"`
	InterfacesList   []InterfaceInfo `json:"interfaces_list"`
}

func (c *Controller) Status(ctx context.Context) (*Status, error) {
	result, err := api.Call(c.Client(), ctx, api.RPCOpts{
		BaseEndpoint:    "/zenarmor/status",
		Method:          http.MethodGet,
		QueryParameters: map[string]string{},
		BodyParameters:  map[string]interface{}{},
	}, &Status{})
	if err != nil {
		return nil, fmt.Errorf("get Zenarmor status: %w", err)
	}
	return result, nil
}
