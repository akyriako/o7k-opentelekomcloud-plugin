package parameters

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v2/configs"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewDcsParameters(p *plugin.Plugin) *Resource {
	return &Resource{plugin: p}
}

func (r *Resource) Service() string {
	return "dcs"
}

func (r *Resource) Kind() string {
	return "dcs-parameters"
}

func (r *Resource) Title() string {
	return "DCS Parameters"
}

func (r *Resource) Aliases() []string {
	return nil
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "value", Title: "VALUE", MinWidth: 16},
		{Key: "default", Title: "DEFAULT", MinWidth: 16},
		{Key: "type", Title: "TYPE", MinWidth: 10},
		{Key: "range", Title: "RANGE", MinWidth: 18, Flex: 1},
		{Key: "description", Title: "DESCRIPTION", MinWidth: 30, Flex: 3},
	}
}

func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	scope := pluginsdk.Scope(ctx)
	instanceID := scope["instance_id"]
	if instanceID == "" {
		return nil, fmt.Errorf("DCS instance ID is required")
	}

	client, err := r.plugin.DCSV2(ctx)
	if err != nil {
		return nil, err
	}

	result, err := configs.Get(client, instanceID)
	if err != nil {
		return nil, fmt.Errorf("getting DCS parameters for instance %q: %w", instanceID, err)
	}

	rows := make([]pluginsdk.Row, 0, len(result.RedisConfigs))

	for _, config := range result.RedisConfigs {
		rows = append(rows, pluginsdk.Row{
			ID: config.ParamID,
			Fields: map[string]string{
				"name":        config.ParamName,
				"value":       config.ParamValue,
				"default":     config.DefaultValue,
				"type":        config.ValueType,
				"range":       config.ValueRange,
				"description": config.Description,
				"instance_id": instanceID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Commands() []pluginsdk.Command {
	return nil
}

func (r *Resource) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	return pluginsdk.Result{}, nil
}
