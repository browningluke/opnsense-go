//go:build generate
// +build generate

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/browningluke/opnsense-go/internal/generate/generator"
	"github.com/browningluke/opnsense-go/internal/generate/schema"
)

var (
	controller = flag.String("controller", "", "controller to generate; must be set if client is not")
	client     = flag.Bool("client", false, "generate OPNsense client iface; must be set if controller is not")
)

func main() {
	flag.Parse()

	if len(*controller) == 0 && !(*client) {
		flag.Usage()
		os.Exit(2)
	}

	// Generate client
	if *client {
		genClientInterface()
	}

	// Generate controller + resources
	if len(*controller) > 0 {
		// Try load controller schema first
		c := schema.GetController(*controller)
		if c == nil {
			fmt.Printf("Schema %s.yml does not exist in schema folder", *controller)
			os.Exit(2)
		}

		genController(c)

		for _, resource := range c.Resources {
			genResource(c, resource)
		}

		for _, rpc := range c.RPC {
			genRPC(c, rpc)
		}
	}
}

func genRPC(controller *schema.ControllerData, rpc schema.RPCData) {
	filename := fmt.Sprintf("%s.go", rpc.Filename)

	fmt.Printf("Generating internal/%s/%s\n", controller.Name, filename)
	g := generator.NewGenerator(filename)

	err := g.Render(rpcTmpl,
		struct {
			Controller schema.ControllerData
			RPC        schema.RPCData
		}{*controller, rpc},
	)
	if err != nil {
		log.Fatalf("generating file (%s): %s", filename, err)
	}

	if err := g.Write(); err != nil {
		log.Fatalf("generating file (%s): %s", filename, err)
	}
}

func genResource(controller *schema.ControllerData, resource schema.ResourceData) {
	filename := fmt.Sprintf("%s.go", resource.Filename)

	fmt.Printf("Generating internal/%s/%s\n", controller.Name, filename)
	g := generator.NewGenerator(filename)

	err := g.Render(resourceTmpl,
		struct {
			Controller schema.ControllerData
			Resource   schema.ResourceData
		}{*controller, resource},
	)
	if err != nil {
		log.Fatalf("generating file (%s): %s", filename, err)
	}

	if err := g.Write(); err != nil {
		log.Fatalf("generating file (%s): %s", filename, err)
	}
}

func genController(data *schema.ControllerData) {
	const (
		filename = `controller.go`
	)

	fmt.Printf("Generating internal/%s/%s\n", data.Name, filename)
	g := generator.NewGenerator(filename)

	err := g.Render(controllerTmpl, data)
	if err != nil {
		log.Fatalf("generating file (%s): %s", filename, err)
	}

	if err := g.Write(); err != nil {
		log.Fatalf("generating file (%s): %s", filename, err)
	}
}

func genClientInterface() {
	const (
		filename = `client.go`
	)

	fmt.Printf("Generating internal/opnsense/%s\n", filename)
	g := generator.NewGenerator(filename)

	// Get all controller names
	cNames := schema.GetControllerNames()

	err := g.Render(clientTmpl, cNames)
	if err != nil {
		log.Fatalf("generating file (%s): %s", filename, err)
	}

	if err := g.Write(); err != nil {
		log.Fatalf("generating file (%s): %s", filename, err)
	}
}

//go:embed templates/rpc.tmpl
var rpcTmpl string

//go:embed templates/resource.tmpl
var resourceTmpl string

//go:embed templates/controller.tmpl
var controllerTmpl string

//go:embed templates/client.tmpl
var clientTmpl string
