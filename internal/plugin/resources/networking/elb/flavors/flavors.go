package flavors

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	elbflavors "github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/flavors"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewElbFlavors(p *plugin.Plugin) *Resource {
	return &Resource{plugin: p}
}

func (r *Resource) Service() string {
	return "elb"
}

func (r *Resource) Kind() string {
	return "elb-flavors"
}

func (r *Resource) Title() string {
	return "ELB Flavors"
}

func (r *Resource) Aliases() []string {
	return nil
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "relationship", Title: "RELATIONSHIP", MinWidth: 14},
		{Key: "type", Title: "TYPE", MinWidth: 10},
		{Key: "connection", Title: "CONNECTIONS", MinWidth: 12},
		{Key: "cps", Title: "CPS", MinWidth: 8},
		{Key: "qps", Title: "QPS", MinWidth: 8},
		{Key: "https_cps", Title: "HTTPS CPS", MinWidth: 10},
		{Key: "bandwidth", Title: "BANDWIDTH", MinWidth: 12},
		{Key: "lcu", Title: "LCU", MinWidth: 8},
		{Key: "shared", Title: "SHARED", MinWidth: 8},
		{Key: "sold_out", Title: "SOLD OUT", MinWidth: 10},
	}
}

func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return nil, err
	}

	pages, err := elbflavors.List(client, elbflavors.ListOpts{}).AllPages()
	if err != nil {
		return nil, fmt.Errorf("listing ELB flavors: %w", err)
	}

	result, err := elbflavors.ExtractFlavors(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting ELB flavors: %w", err)
	}

	rows := make([]pluginsdk.Row, 0, len(result))

	scope := pluginsdk.Scope(ctx)

	relationships := map[string]string{}

	if id := scope["l4_flavor_id"]; id != "" {
		relationships[id] = "L4"
	}

	if id := scope["l4_scale_flavor_id"]; id != "" {
		relationships[id] = "L4 SCALE"
	}

	if id := scope["l7_flavor_id"]; id != "" {
		relationships[id] = "L7"
	}

	if id := scope["l7_scale_flavor_id"]; id != "" {
		relationships[id] = "L7 SCALE"
	}

	scoped := len(relationships) > 0

	for _, flavor := range result {
		relationship := relationships[flavor.ID]

		if scoped && relationship == "" {
			continue
		}
		rows = append(rows, pluginsdk.Row{
			ID: flavor.ID,
			Fields: map[string]string{
				"id":         flavor.ID,
				"name":       flavor.Name,
				"type":       flavor.Type,
				"connection": strconv.Itoa(flavor.Info.Connection),
				"cps":        strconv.Itoa(flavor.Info.Cps),
				"qps":        strconv.Itoa(flavor.Info.Qps),
				"https_cps":  strconv.Itoa(flavor.Info.HttpsCps),
				"bandwidth":  strconv.Itoa(flavor.Info.Bandwidth),
				"lcu":        strconv.Itoa(flavor.Info.Lcu),
				"shared":     strconv.FormatBool(flavor.Shared),
				"sold_out":   strconv.FormatBool(flavor.SoldOut),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Commands() []pluginsdk.Command {
	return []pluginsdk.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}

func (r *Resource) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	switch command.Key {
	case "s":
		return r.show(ctx, row)
	}

	return pluginsdk.Result{}, nil
}
