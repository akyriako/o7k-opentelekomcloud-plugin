package ipgroups

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	elbipgroups "github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/ipgroups"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewElbIpGroups(p *plugin.Plugin) *Resource {
	return &Resource{plugin: p}
}

func (r *Resource) Service() string {
	return "elb"
}

func (r *Resource) Kind() string {
	return "elb-ipgroups"
}

func (r *Resource) Title() string {
	return "ELB IP Groups"
}

func (r *Resource) Aliases() []string {
	return nil
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "ips", Title: "IP ADDRESSES", MinWidth: 32, Flex: 1},
		{Key: "listeners", Title: "LISTENERS", MinWidth: 10},
		{Key: "created_at", Title: "CREATED", MinWidth: 20},
	}
}

func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return nil, err
	}

	scope := pluginsdk.Scope(ctx)
	ipGroupID, scoped := scope["ipgroup_id"]

	if scoped && ipGroupID == "" {
		return []pluginsdk.Row{}, nil
	}

	opts := elbipgroups.ListOpts{}

	if scoped {
		opts.ID = []string{ipGroupID}
	}

	result, err := elbipgroups.List(client, elbipgroups.ListOpts{})
	if err != nil {
		return nil, fmt.Errorf("listing ELB IP groups: %w", err)
	}

	rows := make([]pluginsdk.Row, 0, len(result))

	for _, group := range result {
		ips := make([]string, 0, len(group.IpList))

		for _, ip := range group.IpList {
			ips = append(ips, ip.Ip)
		}

		rows = append(rows, pluginsdk.Row{
			ID: group.ID,
			Fields: map[string]string{
				"id":          group.ID,
				"name":        group.Name,
				"description": group.Description,
				"ips":         strings.Join(ips, ", "),
				"listeners":   strconv.Itoa(len(group.Listeners)),
				"created_at":  group.CreatedAt,
				"updated_at":  group.UpdatedAt,
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

func (r *Resource) show(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, err
	}

	group, err := elbipgroups.Get(client, row.ID)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting ELB IP group %q: %w", row.ID, err)
	}

	content, err := json.Marshal(group)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("marshalling ELB IP group %q: %w", row.ID, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      row.ID,
			Content: content,
		},
	}, nil
}
