package zenarmor

import "github.com/browningluke/opnsense-go/pkg/api"

type Controller struct{ Api *api.Client }

func (c *Controller) Client() *api.Client { return c.Api }
