package main

import (
	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/containers/cce/clusters"
	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/containers/cce/nodes"
	"github.com/akyriako/o7k/pluginsdk"
)

func main() {
	p := plugin.New()

	p.Register(
		clusters.NewCceClusters(p),
		nodes.NewCceNodes(p),
	)

	pluginsdk.Serve(p)
}
