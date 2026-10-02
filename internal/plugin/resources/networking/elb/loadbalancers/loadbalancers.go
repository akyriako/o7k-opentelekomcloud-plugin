package loadbalancers

import (
	"context"
	"strconv"
	"strings"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/loadbalancers"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewElbLoadBalancers(p *plugin.Plugin) *Resource {
	return &Resource{plugin: p}
}

func (r *Resource) Service() string {
	return "elb"
}

func (r *Resource) Kind() string {
	return "elb-loadbalancers"
}

func (r *Resource) Title() string {
	return "ELB Load Balancers"
}

func (r *Resource) Aliases() []string {
	return []string{"elb"}
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "provisioning_status", Title: "PROVISIONING", MinWidth: 14},
		{Key: "operating_status", Title: "STATUS", MinWidth: 10},
		{Key: "vip_address", Title: "VIP", MinWidth: 16},
		{Key: "eip_address", Title: "EIP", MinWidth: 16},
		{Key: "az", Title: "AVAILABILITY ZONES", MinWidth: 20},
		{Key: "guaranteed", Title: "DEDICATED", MinWidth: 10},
	}
}

func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return nil, err
	}

	pages, err := loadbalancers.List(client, loadbalancers.ListOpts{}).AllPages()
	if err != nil {
		return nil, err
	}

	result, err := loadbalancers.ExtractLoadbalancers(pages)
	if err != nil {
		return nil, err
	}

	rows := make([]pluginsdk.Row, 0, len(result))

	for _, lb := range result {
		eipAddresses := make([]string, 0, len(lb.Eips))
		eipIDs := make([]string, 0, len(lb.Eips))

		for _, eip := range lb.Eips {
			eipAddresses = append(eipAddresses, eip.EipAddress)
			eipIDs = append(eipIDs, eip.EipID)
		}
		rows = append(rows, pluginsdk.Row{
			ID: lb.ID,
			Fields: map[string]string{
				"id":                  lb.ID,
				"name":                lb.Name,
				"provisioning_status": lb.ProvisioningStatus,
				"operating_status":    lb.OperatingStatus,
				"vip_address":         lb.VipAddress,
				"eip_address":         strings.Join(eipAddresses, ", "),
				"eip_ids":             strings.Join(eipIDs, ","),
				"az":                  strings.Join(lb.AvailabilityZoneList, ", "),
				"guaranteed":          strconv.FormatBool(lb.Guaranteed),

				"vpc_id":             lb.VpcID,
				"subnet_id":          lb.VipSubnetCidrID,
				"vip_port_id":        lb.VipPortID,
				"l4_flavor_id":       lb.L4FlavorID,
				"l7_flavor_id":       lb.L7FlavorID,
				"l4_scale_flavor_id": lb.L4ScaleFlavorID,
				"l7_scale_flavor_id": lb.L7ScaleFlavorID,
				"provider":           lb.Provider,
				"admin_state_up":     strconv.FormatBool(lb.AdminStateUp),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Commands() []pluginsdk.Command {
	return []pluginsdk.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-l", Description: "Listeners"},
		{Key: "shift-p", Description: "Pools"},
		{Key: "shift-f", Description: "Flavors"},
	}
}

func (r *Resource) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	switch command.Key {
	case "s":
		return r.show(ctx, row)
	case "shift-l":
		return r.listeners(ctx, row)
	case "shift-p":
		return r.pools(ctx, row)
	case "shift-f":
		return r.flavors(ctx, row)
	}

	return pluginsdk.Result{}, nil
}
