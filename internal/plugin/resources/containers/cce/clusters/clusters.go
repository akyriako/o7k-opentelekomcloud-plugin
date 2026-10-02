package clusters

import (
	"context"
	"fmt"
	"strings"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/clusters"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewCceClusters(p *plugin.Plugin) *Resource {
	return &Resource{
		plugin: p,
	}
}

func (r *Resource) Service() string {
	return "cce"
}

func (r *Resource) Kind() string {
	return "cce-clusters"
}

func (r *Resource) Title() string {
	return "CCE Clusters"
}

func (r *Resource) Aliases() []string {
	return []string{"cce"}
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 12},
		{Key: "version", Title: "VERSION", MinWidth: 12},
		{Key: "flavor", Title: "CCE FLAVOR", MinWidth: 16},
		{Key: "availability_zone", Title: "AVAILABILITY ZONES", MinWidth: 18, Flex: 1},
		{Key: "type", Title: "TYPE", MinWidth: 16},
	}
}

func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	client, err := r.plugin.CCEV3(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting CCE client: %w", err)
	}

	allClusters, err := clusters.List(client, clusters.ListOpts{})
	if err != nil {
		return nil, fmt.Errorf("listing CCE clusters: %w", err)
	}

	rows := make([]pluginsdk.Row, 0, len(allClusters))

	for _, cluster := range allClusters {
		availabilityZones := make([]string, 0, len(cluster.Spec.Masters))

		for _, master := range cluster.Spec.Masters {
			if master.AvailabilityZone != "" {
				availabilityZones = append(
					availabilityZones,
					master.AvailabilityZone,
				)
			}
		}

		rows = append(rows, pluginsdk.Row{
			ID: cluster.Metadata.Id,
			Fields: map[string]string{
				"id":                cluster.Metadata.Id,
				"name":              cluster.Metadata.Name,
				"version":           cluster.Spec.Version,
				"flavor":            cluster.Spec.Flavor,
				"availability_zone": strings.Join(availabilityZones, ", "),
				"type":              cluster.Spec.Type,
				"status":            cluster.Status.Phase,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Commands() []pluginsdk.Command {
	return []pluginsdk.Command{
		{Key: "s", Description: "Show"},
		{Key: "shift-n", Description: "Nodes", Default: true},
		{Key: "shift-p", Description: "Node Pools"},
		{Key: "shift-k", Description: "Show Kubeconfig"},
		{Key: "ctrl+k", Description: "Get Kubeconfig", StatusLabel: "Downloading"},
	}
}

func (r *Resource) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	switch command.Key {
	case "s":
		return r.show(ctx, row.ID)
	case "shift-n":
		return r.nodes(ctx, row)
	case "shift-p":
		return r.nodePools(ctx, row)
	case "shift-k":
		return r.kubeconfig(ctx, row.ID)
	case "ctrl+k":
		return r.getKubeconfig(ctx, row)
	}

	return pluginsdk.Result{}, nil
}
