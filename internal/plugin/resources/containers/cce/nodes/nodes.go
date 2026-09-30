package nodes

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	ccenodes "github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/nodes"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewCceNodes(p *plugin.Plugin) *Resource {
	return &Resource{
		plugin: p,
	}
}

func (r *Resource) Service() string {
	return "cce"
}

func (r *Resource) Kind() string {
	return "cce-nodes"
}

func (r *Resource) Title() string {
	return "CCE Nodes"
}

func (r *Resource) Aliases() []string {
	return nil
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 40},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 12},
		{Key: "az", Title: "AVAILABILITY ZONE", MinWidth: 8, Flex: 1},
		{Key: "internal_ip", Title: "INTERNAL IP", MinWidth: 16},
		{Key: "flavor", Title: "FLAVOR", MinWidth: 14},
		{Key: "os", Title: "OS", MinWidth: 20, Flex: 1},
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

	nodes, err := ccenodes.List(client, clusterID, ccenodes.ListOpts{})
	if err != nil {
		return nil, fmt.Errorf(
			"listing CCE nodes for cluster %q: %w",
			clusterID,
			err,
		)
	}

	rows := make([]pluginsdk.Row, 0, len(nodes))

	for _, node := range nodes {
		rows = append(rows, pluginsdk.Row{
			ID: node.Metadata.Id,
			Fields: map[string]string{
				"id":          node.Metadata.Id,
				"name":        node.Metadata.Name,
				"az":          node.Spec.Az,
				"status":      node.Status.Phase,
				"internal_ip": node.Status.PrivateIP,
				"flavor":      node.Spec.Flavor,
				"os":          node.Spec.Os,
				"cluster_id":  clusterID,
				"server_id":   node.Status.ServerID,
				"subnet_id":   node.Spec.NodeNicSpec.PrimaryNic.SubnetId,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Commands() []pluginsdk.Command {
	return []pluginsdk.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-s", Description: "Server"},
		{Key: "shift-f", Description: "Flavor"},
		{Key: "shift-n", Description: "Network"},
	}
}

func (r *Resource) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	clusterID := row.Fields["cluster_id"]

	if clusterID == "" {
		scope := pluginsdk.Scope(ctx)
		clusterID = scope["cluster_id"]
	}

	switch command.Key {
	case "s":
		return r.show(ctx, row.ID, clusterID)
	case "shift-f":
		return r.flavor(row)
	case "shift-n":
		return r.network(row)
	case "shift-s":
		return r.server(row)
	}

	return pluginsdk.Result{}, nil
}
