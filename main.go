package main

import (
	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/containers/cce/clusters"
	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/containers/cce/nodepools"
	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/containers/cce/nodes"
	rdsbackups "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/databases/rds/backups"
	rdsflavors "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/databases/rds/flavors"
	rdsinstances "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/databases/rds/instances"
	rdsnodes "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/databases/rds/nodes"
	rdsparameters "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/databases/rds/parameters"
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
		rdsinstances.NewRdsInstances(p),
		rdsflavors.NewRdsFlavors(p),
		rdsnodes.NewRdsNodes(p),
		rdsbackups.NewRdsBackups(p),
		rdsparameters.NewRdsParameters(p),
	)

	pluginsdk.Serve(p)
}
