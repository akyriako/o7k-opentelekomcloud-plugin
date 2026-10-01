package rds

import (
	"context"
	"fmt"
	"strings"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	rdsflavors "github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/flavors"
)

type FlavorsResource struct {
	plugin *plugin.Plugin
}

func NewRdsFlavors(p *plugin.Plugin) *FlavorsResource {
	return &FlavorsResource{plugin: p}
}

func (r *FlavorsResource) Service() string {
	return "rds"
}

func (r *FlavorsResource) Kind() string {
	return "rds-flavors"
}

func (r *FlavorsResource) Title() string {
	return "RDS Flavors"
}

func (r *FlavorsResource) Aliases() []string {
	return nil
}

func (r *FlavorsResource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "spec_code", Title: "SPEC CODE", MinWidth: 28, Flex: 1},
		{Key: "vcpus", Title: "VCPUS", MinWidth: 8},
		{Key: "ram", Title: "RAM", MinWidth: 8},
		{Key: "mode", Title: "MODE", MinWidth: 10},
		{Key: "versions", Title: "VERSIONS", MinWidth: 16},
	}
}

func (r *FlavorsResource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	scope := pluginsdk.Scope(ctx)

	client, err := r.plugin.RDSV3(ctx)
	if err != nil {
		return nil, err
	}

	result, err := rdsflavors.ListFlavors(client, rdsflavors.ListOpts{
		DatabaseName: scope["engine"],
		VersionName:  scope["version"],
		SpecCode:     scope["spec_code"],
	})
	if err != nil {
		return nil, err
	}

	rows := make([]pluginsdk.Row, 0, len(result))

	for _, flavor := range result {
		rows = append(rows, pluginsdk.Row{
			ID: flavor.Id,
			Fields: map[string]string{
				"id":        flavor.Id,
				"spec_code": flavor.SpecCode,
				"vcpus":     flavor.VCPUs,
				"ram":       formatRAM(flavor.RAM),
				"mode":      flavor.InstanceMode,
				"versions":  strings.Join(flavor.VersionName, ", "),
			},
		})
	}

	return rows, nil
}

func (r *FlavorsResource) Commands() []pluginsdk.Command {
	return nil
}

func (r *FlavorsResource) Execute(context.Context, pluginsdk.Command, pluginsdk.Row) (pluginsdk.Result, error) {
	return pluginsdk.Result{}, nil
}

func formatRAM(ram int) string {
	return fmt.Sprintf("%d GB", ram)
}
