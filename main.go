package main

import (
	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/containers/cce/clusters"
	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/containers/cce/nodepools"
	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/containers/cce/nodes"
	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/objectstorage/obs/buckets"
	"github.com/akyriako/o7k/pluginsdk"
)

func main() {
	p := plugin.New()

	p.Register(
		clusters.NewCceClusters(p),
		nodes.NewCceNodes(p),
		nodepools.NewCceNodePools(p),
		buckets.NewObsBuckets(p),
	)

	pluginsdk.Serve(p)
}
