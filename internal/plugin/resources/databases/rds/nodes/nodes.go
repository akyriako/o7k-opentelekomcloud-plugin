package nodes

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	rdsinstances "github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewRdsNodes(p *plugin.Plugin) *Resource {
	return &Resource{plugin: p}
}

func (r *Resource) Service() string {
	return "rds"
}

func (r *Resource) Kind() string {
	return "rds-nodes"
}

func (r *Resource) Title() string {
	return "RDS Nodes"
}

func (r *Resource) Aliases() []string {
	return nil
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "role", Title: "ROLE", MinWidth: 12},
		{Key: "status", Title: "STATUS", MinWidth: 12},
		{Key: "availability_zone", Title: "AVAILABILITY ZONE", MinWidth: 18},
	}
}

func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	scope := pluginsdk.Scope(ctx)
	instanceID := scope["instance_id"]

	if instanceID == "" {
		return nil, fmt.Errorf("instance_id is required")
	}

	client, err := r.plugin.RDSV3(ctx)
	if err != nil {
		return nil, err
	}

	result, err := rdsinstances.List(client, rdsinstances.ListOpts{
		Id: instanceID,
	})
	if err != nil {
		return nil, fmt.Errorf("getting RDS instance %q: %w", instanceID, err)
	}

	if len(result.Instances) == 0 {
		return nil, fmt.Errorf("RDS instance %q not found", instanceID)
	}

	instance := result.Instances[0]
	rows := make([]pluginsdk.Row, 0, len(instance.Nodes))

	for _, node := range instance.Nodes {
		rows = append(rows, pluginsdk.Row{
			ID: node.Id,
			Fields: map[string]string{
				"id":                node.Id,
				"name":              node.Name,
				"role":              node.Role,
				"status":            node.Status,
				"availability_zone": node.AvailabilityZone,
				"instance_id":       instanceID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Commands() []pluginsdk.Command {
	return nil
}

func (r *Resource) Execute(context.Context, pluginsdk.Command, pluginsdk.Row) (pluginsdk.Result, error) {
	return pluginsdk.Result{}, nil
}
