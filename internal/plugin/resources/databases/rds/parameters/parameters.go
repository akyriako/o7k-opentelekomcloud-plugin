package parameters

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	rdsconfigurations "github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/configurations"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewRdsParameters(p *plugin.Plugin) *Resource {
	return &Resource{plugin: p}
}

func (r *Resource) Service() string {
	return "rds"
}

func (r *Resource) Kind() string {
	return "rds-parameters"
}

func (r *Resource) Title() string {
	return "RDS Parameters"
}

func (r *Resource) Aliases() []string {
	return nil
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "name", Title: "NAME", MinWidth: 20, Flex: 1},
		{Key: "value", Title: "VALUE", MinWidth: 16, Flex: 1},
		{Key: "type", Title: "TYPE", MinWidth: 10},
		{Key: "restart_required", Title: "RESTART REQUIRED", MinWidth: 16},
		{Key: "read_only", Title: "READ ONLY", MinWidth: 10},
		{Key: "value_range", Title: "VALUE RANGE", MinWidth: 20, Flex: 2},
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

	configuration, err := rdsconfigurations.GetForInstance(client, instanceID)
	if err != nil {
		return nil, fmt.Errorf("getting configuration for RDS instance %q: %w", instanceID, err)
	}

	rows := make([]pluginsdk.Row, 0, len(configuration.Parameters))

	for _, parameter := range configuration.Parameters {
		rows = append(rows, pluginsdk.Row{
			ID: parameter.Name,
			Fields: map[string]string{
				"name":             parameter.Name,
				"value":            parameter.Value,
				"type":             parameter.Type,
				"restart_required": strconv.FormatBool(parameter.RestartRequired),
				"read_only":        strconv.FormatBool(parameter.ReadOnly),
				"value_range":      parameter.ValueRange,
				"description":      parameter.Description,
				"instance_id":      instanceID,
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
