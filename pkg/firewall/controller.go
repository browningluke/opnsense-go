// Code generated by internal/generate/api/main.go; DO NOT EDIT.

package firewall

import "github.com/browningluke/opnsense-go/pkg/api"

// Controller for firewall
type Controller struct {
	Api *api.Client
}

func (c *Controller) Client() *api.Client {
	return c.Api
}
