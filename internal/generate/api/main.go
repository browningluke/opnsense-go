//go:build generate
// +build generate

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/browningluke/opnsense-go/internal/generate/generator"
	"github.com/browningluke/opnsense-go/internal/generate/schema"
	"log"
	"os"
)

var (
	controller = flag.String("controller", "", "controller to generate; must be set if client is not")
)

func main() {
	flag.Parse()

	if len(*controller) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// Generate controller + resources
	// Try load controller schema first
	c := schema.GetController(*controller)
	if c == nil {
		fmt.Printf("Schema %s.yml does not exist in schema folder", *controller)
		os.Exit(2)
	}

	genController(c)
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

//go:embed templates/controller.tmpl
var controllerTmpl string
