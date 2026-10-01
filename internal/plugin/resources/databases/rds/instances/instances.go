package instances

import (
	"context"
	"strings"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewRdsInstances(p *plugin.Plugin) *Resource {
	return &Resource{plugin: p}
}

func (r *Resource) Service() string {
	return "rds"
}

func (r *Resource) Kind() string {
	return "rds-instances"
}

func (r *Resource) Title() string {
	return "RDS Instances"
}

func (r *Resource) Aliases() []string {
	return nil
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "type", Title: "TYPE", MinWidth: 10},
		{Key: "flavor", Title: "FLAVOR", MinWidth: 18, Flex: 2},
		{Key: "engine", Title: "ENGINE", MinWidth: 12},
		{Key: "version", Title: "VERSION", MinWidth: 10},
		{Key: "private_ips", Title: "PRIVATE IPS", MinWidth: 15, Flex: 1},
		{Key: "public_ips", Title: "PUBLIC IPS", MinWidth: 15, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 12},
	}
}

func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	client, err := r.plugin.RDSV3(ctx)
	if err != nil {
		return nil, err
	}

	result, err := instances.List(client, instances.ListOpts{})
	if err != nil {
		return nil, err
	}

	rows := make([]pluginsdk.Row, 0, len(result.Instances))

	for _, instance := range result.Instances {
		rows = append(rows, pluginsdk.Row{
			ID: instance.Id,
			Fields: map[string]string{
				"id":                instance.Id,
				"name":              instance.Name,
				"status":            instance.Status,
				"engine":            instance.DataStore.Type,
				"version":           instance.DataStore.Version,
				"type":              strings.ToUpper(instance.Type),
				"flavor":            instance.FlavorRef,
				"region":            instance.Region,
				"vpc_id":            instance.VpcId,
				"subnet_id":         instance.SubnetId,
				"security_group_id": instance.SecurityGroupId,
				"private_ips":       strings.Join(instance.PrivateIps, ", "),
				"public_ips":        strings.Join(instance.PublicIps, ", "),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Commands() []pluginsdk.Command {
	return []pluginsdk.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-n", Description: "Nodes"},
		{Key: "shift-b", Description: "Backups"},
		{Key: "shift-p", Description: "Backup Policy"},
		{Key: "shift-c", Description: "Parameters"},
		{Key: "shift-f", Description: "Flavor"},
		{Key: "shift-g", Description: "Security Group"},
		{Key: "shift-w", Description: "Network"},
	}
}

func (r *Resource) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	switch command.Key {
	case "s":
		return r.show(ctx, row)
	case "shift-g":
		return r.securityGroup(ctx, row)
	case "shift-f":
		return r.flavor(ctx, row)
	case "shift-n":
		return r.nodes(ctx, row)
	case "shift-w":
		return r.network(row)
	case "shift-b":
		return r.backups(ctx, row)
	case "shift-p":
		return r.backupPolicy(ctx, row)
	case "shift-c":
		return r.parameters(ctx, row)
	}

	return pluginsdk.Result{}, nil
}
