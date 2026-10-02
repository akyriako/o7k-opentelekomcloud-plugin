package whitelists

import (
	"context"
	"fmt"
	"strings"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	dcswhitelists "github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v2/whitelists"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewDcsWhitelists(p *plugin.Plugin) *Resource {
	return &Resource{plugin: p}
}

func (r *Resource) Service() string {
	return "dcs"
}

func (r *Resource) Kind() string {
	return "dcs-whitelists"
}

func (r *Resource) Title() string {
	return "DCS Whitelists"
}

func (r *Resource) Aliases() []string {
	return nil
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "group", Title: "GROUP", MinWidth: 24},
		{Key: "enabled", Title: "ENABLED", MinWidth: 10},
		{Key: "ips", Title: "IP ADDRESSES", MinWidth: 32, Flex: 1},
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

	result, err := dcswhitelists.Get(client, instanceID)
	if err != nil {
		return nil, fmt.Errorf("getting DCS whitelist for instance %q: %w", instanceID, err)
	}

	rows := make([]pluginsdk.Row, 0, len(result.Groups))

	for _, group := range result.Groups {
		rows = append(rows, pluginsdk.Row{
			ID: group.GroupName,
			Fields: map[string]string{
				"group":       group.GroupName,
				"enabled":     fmt.Sprintf("%t", result.Enable),
				"ips":         strings.Join(group.IPList, ", "),
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
