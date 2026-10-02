package nodepools

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/nodepools"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewCceNodePools(p *plugin.Plugin) *Resource {
	return &Resource{
		plugin: p,
	}
}

func (r *Resource) Service() string {
	return "cce"
}

func (r *Resource) Kind() string {
	return "cce-nodepools"
}

func (r *Resource) Title() string {
	return "CCE Node Pools"
}

func (r *Resource) Aliases() []string {
	return nil
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "flavor", Title: "FLAVOR", MinWidth: 16},
		{Key: "az", Title: "AZ", MinWidth: 8},
		{Key: "os", Title: "OS", MinWidth: 16},
		{Key: "nodes", Title: "NODES", MinWidth: 8},
		{Key: "status", Title: "STATUS", MinWidth: 12, Flex: 1},
	}
}

func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	scope := pluginsdk.Scope(ctx)
	clusterID := scope["cluster_id"]
	if clusterID == "" {
		return nil, fmt.Errorf("cluster_id is required")
	}

	client, err := r.plugin.CCEV3(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting CCE client: %w", err)
	}

	allNodePools, err := nodepools.List(client, clusterID, nodepools.ListOpts{})
	if err != nil {
		return nil, fmt.Errorf("listing CCE node pools: %w", err)
	}

	rows := make([]pluginsdk.Row, 0, len(allNodePools))

	for _, nodePool := range allNodePools {
		rows = append(rows, pluginsdk.Row{
			ID: nodePool.Metadata.Id,
			Fields: map[string]string{
				"id":              nodePool.Metadata.Id,
				"name":            nodePool.Metadata.Name,
				"status":          nodePool.Status.Phase,
				"nodes":           strconv.Itoa(nodePool.Status.CurrentNode),
				"flavor":          nodePool.Spec.NodeTemplate.Flavor,
				"os":              nodePool.Spec.NodeTemplate.Os,
				"cluster_id":      clusterID,
				"security_groups": strings.Join(nodePool.Spec.CustomSecurityGroupIds, ","),
				"az":              nodePool.Spec.NodeTemplate.Az,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Commands() []pluginsdk.Command {
	return []pluginsdk.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-f", Description: "Flavor"},
		{Key: "shift-g", Description: "Security Group"},
		{Key: "ctrl+u", Description: "Scale Up (+1)", StatusLabel: "Scaling up"},
		{Key: "ctrl+d", Description: "Scale Down (-1)", StatusLabel: "Scaling down"},
	}
}

func (r *Resource) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	switch command.Key {
	case "s":
		return r.show(ctx, row)
	case "shift-f":
		return r.flavor(ctx, row)
	case "shift-g":
		return r.securityGroup(ctx, row)
	case "ctrl+u":
		return r.scale(ctx, row, 1)
	case "ctrl+d":
		return r.scale(ctx, row, -1)
	}

	return pluginsdk.Result{}, nil
}
