package instances

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	dcsinstances "github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v2/instance"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewDcsInstances(p *plugin.Plugin) *Resource {
	return &Resource{plugin: p}
}

func (r *Resource) Service() string {
	return "dcs"
}

func (r *Resource) Kind() string {
	return "dcs-instances"
}

func (r *Resource) Title() string {
	return "DCS Instances"
}

func (r *Resource) Aliases() []string {
	return []string{"dcs"}
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 12},
		{Key: "engine", Title: "ENGINE", MinWidth: 10},
		{Key: "version", Title: "VERSION", MinWidth: 10},
		{Key: "capacity", Title: "CAPACITY", MinWidth: 10},
		{Key: "cache_size", Title: "CACHE SIZE", MinWidth: 12},
		{Key: "memory", Title: "MEMORY", MinWidth: 16},
		{Key: "ip", Title: "IP", MinWidth: 16},
		{Key: "port", Title: "PORT", MinWidth: 8},
		{Key: "ssl", Title: "SSL", MinWidth: 6},
	}
}

func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	client, err := r.plugin.DCSV2(ctx)
	if err != nil {
		return nil, err
	}

	result, err := dcsinstances.List(client, dcsinstances.ListDcsInstanceOpts{
		Limit: 100,
	})
	if err != nil {
		return nil, fmt.Errorf("listing DCS instances: %w", err)
	}

	rows := make([]pluginsdk.Row, 0, len(result.Instances))

	for _, instance := range result.Instances {
		cacheSize := instance.CapacityMinor
		if strings.HasPrefix(cacheSize, ".") {
			cacheSize = "0" + cacheSize
		}
		rows = append(rows, pluginsdk.Row{
			ID: instance.InstanceId,
			Fields: map[string]string{
				"id":                instance.InstanceId,
				"name":              instance.Name,
				"status":            instance.Status,
				"engine":            instance.Engine,
				"version":           instance.EngineVersion,
				"capacity":          strconv.FormatFloat(instance.Capacity, 'f', -1, 64) + " GB",
				"cache_size":        cacheSize + " GB",
				"memory":            fmt.Sprintf("%d/%d MB", instance.UsedMemory, instance.MaxMemory),
				"ip":                instance.Ip,
				"port":              strconv.Itoa(instance.Port),
				"vpc_id":            instance.VpcId,
				"subnet_id":         instance.SubnetId,
				"security_group_id": instance.SecurityGroupId,
				"public_ip_id":      instance.PublicIpId,
				"public_ip_address": instance.PublicIpAddress,
				"spec_code":         instance.SpecCode,
				"ssl":               strconv.FormatBool(instance.EnableSsl),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Commands() []pluginsdk.Command {
	return []pluginsdk.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-n", Description: "Network"},
		{Key: "shift-g", Description: "Security Group"},
		{Key: "shift-p", Description: "Parameters"},
		{Key: "shift-w", Description: "Whitelist"},
		{Key: "ctrl+t", Description: "Toggle SSL", StatusLabel: "Updating SSL"},
		{Key: "shift-s", Description: "Show SSL"},
		{Key: "ctrl+k", Description: "Get Certificate", StatusLabel: "Downloading Certificate"},
	}
}

func (r *Resource) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	switch command.Key {
	case "s":
		return r.show(ctx, row)
	case "shift-n":
		return r.network(ctx, row)
	case "shift-g":
		return r.securityGroup(ctx, row)
	case "shift-p":
		return r.parameters(ctx, row)
	case "shift-w":
		return r.whitelist(ctx, row)
	case "ctrl+t":
		return r.toggleSSL(ctx, row)
	case "ctrl+k":
		return r.downloadCertificate(ctx, row)
	case "shift-s":
		return r.showSSL(ctx, row)
	}

	return pluginsdk.Result{}, nil
}
