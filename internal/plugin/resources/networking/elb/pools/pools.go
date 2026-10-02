package pools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	elbpools "github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/pools"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewElbPools(p *plugin.Plugin) *Resource {
	return &Resource{plugin: p}
}

func (r *Resource) Service() string {
	return "elb"
}

func (r *Resource) Kind() string {
	return "elb-pools"
}

func (r *Resource) Title() string {
	return "ELB Pools"
}

func (r *Resource) Aliases() []string {
	return nil
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "protocol", Title: "PROTOCOL", MinWidth: 10},
		{Key: "algorithm", Title: "ALGORITHM", MinWidth: 18},
		{Key: "admin_state_up", Title: "ADMIN UP", MinWidth: 10},
		{Key: "type", Title: "TYPE", MinWidth: 12},
		{Key: "ip_version", Title: "IP VERSION", MinWidth: 10},
	}
}

func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	scope := pluginsdk.Scope(ctx)
	loadbalancerID := scope["loadbalancer_id"]
	if loadbalancerID == "" {
		return nil, fmt.Errorf("ELB load balancer ID is required")
	}

	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return nil, err
	}

	pages, err := elbpools.List(client, elbpools.ListOpts{
		LoadbalancerID: []string{loadbalancerID},
	}).AllPages()
	if err != nil {
		return nil, err
	}

	result, err := elbpools.ExtractPools(pages)
	if err != nil {
		return nil, err
	}

	rows := make([]pluginsdk.Row, 0, len(result))

	for _, pool := range result {
		rows = append(rows, pluginsdk.Row{
			ID: pool.ID,
			Fields: map[string]string{
				"id":              pool.ID,
				"name":            pool.Name,
				"protocol":        pool.Protocol,
				"algorithm":       pool.LBMethod,
				"admin_state_up":  strconv.FormatBool(pool.AdminStateUp),
				"type":            pool.Type,
				"ip_version":      pool.IpVersion,
				"monitor_id":      pool.MonitorID,
				"vpc_id":          pool.VpcId,
				"loadbalancer_id": loadbalancerID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Commands() []pluginsdk.Command {
	return []pluginsdk.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-m", Description: "Members"},
		{Key: "shift-h", Description: "Health Monitor"},
	}
}

func (r *Resource) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	switch command.Key {
	case "s":
		return r.show(ctx, row)
	case "shift-m":
		return r.members(ctx, row)
	case "shift-h":
		return r.healthMonitor(ctx, row)
	}

	return pluginsdk.Result{}, nil
}
