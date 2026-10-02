package main

import (
	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	cceclusters "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/containers/cce/clusters"
	ccenodepools "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/containers/cce/nodepools"
	ccenodes "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/containers/cce/nodes"
	dcsinstances "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/databases/dcs/instances"
	dcsparameters "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/databases/dcs/parameters"
	dcswhitelists "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/databases/dcs/whitelists"
	rdsbackups "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/databases/rds/backups"
	rdsflavors "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/databases/rds/flavors"
	rdsinstances "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/databases/rds/instances"
	rdsnodes "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/databases/rds/nodes"
	rdsparameters "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/databases/rds/parameters"
	elbflavors "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/networking/elb/flavors"
	elbhealthmonitors "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/networking/elb/healthmonitors"
	elbipgroups "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/networking/elb/ipgroups"
	elblisteners "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/networking/elb/listeners"
	elbloadbalancers "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/networking/elb/loadbalancers"
	elbmembers "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/networking/elb/members"
	elbpolicies "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/networking/elb/policies"
	elbpools "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/networking/elb/pools"
	elbrules "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/networking/elb/rules"
	obsbuckets "github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin/resources/objectstorage/obs/buckets"
	"github.com/akyriako/o7k/pluginsdk"
)

func main() {
	p := plugin.New()

	p.Register(
		cceclusters.NewCceClusters(p),
		ccenodes.NewCceNodes(p),
		ccenodepools.NewCceNodePools(p),
		obsbuckets.NewObsBuckets(p),
		rdsinstances.NewRdsInstances(p),
		rdsflavors.NewRdsFlavors(p),
		rdsnodes.NewRdsNodes(p),
		rdsbackups.NewRdsBackups(p),
		rdsparameters.NewRdsParameters(p),
		dcsinstances.NewDcsInstances(p),
		dcsparameters.NewDcsParameters(p),
		dcswhitelists.NewDcsWhitelists(p),
		elbloadbalancers.NewElbLoadBalancers(p),
		elblisteners.NewElbListeners(p),
		elbpools.NewElbPools(p),
		elbmembers.NewElbMembers(p),
		elbhealthmonitors.NewElbHealthMonitors(p),
		elbpolicies.NewElbPolicies(p),
		elbrules.NewElbRules(p),
		elbflavors.NewElbFlavors(p),
		elbipgroups.NewElbIpGroups(p),
	)

	pluginsdk.Serve(p)
}
