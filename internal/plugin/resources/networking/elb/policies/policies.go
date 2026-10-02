package policies

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	elbpolicies "github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/policies"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewElbPolicies(p *plugin.Plugin) *Resource {
	return &Resource{plugin: p}
}

func (r *Resource) Service() string {
	return "elb"
}

func (r *Resource) Kind() string {
	return "elb-policies"
}

func (r *Resource) Title() string {
	return "ELB Policies"
}

func (r *Resource) Aliases() []string {
	return nil
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "action", Title: "ACTION", MinWidth: 20},
		{Key: "priority", Title: "PRIORITY", MinWidth: 10},
		{Key: "position", Title: "POSITION", MinWidth: 10},
		{Key: "status", Title: "PROVISIONING", MinWidth: 14},
	}
}

func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	scope := pluginsdk.Scope(ctx)
	listenerID := scope["listener_id"]
	if listenerID == "" {
		return nil, fmt.Errorf("ELB listener ID is required")
	}

	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return nil, err
	}

	pages, err := elbpolicies.List(client, elbpolicies.ListOpts{
		ListenerID: []string{listenerID},
	}).AllPages()
	if err != nil {
		return nil, fmt.Errorf("listing ELB policies: %w", err)
	}

	result, err := elbpolicies.ExtractPolicies(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting ELB policies: %w", err)
	}

	rows := make([]pluginsdk.Row, 0, len(result))

	for _, policy := range result {
		rows = append(rows, pluginsdk.Row{
			ID: policy.ID,
			Fields: map[string]string{
				"id":                   policy.ID,
				"name":                 policy.Name,
				"action":               string(policy.Action),
				"priority":             strconv.Itoa(policy.Priority),
				"position":             strconv.Itoa(policy.Position),
				"status":               policy.Status,
				"listener_id":          policy.ListenerID,
				"redirect_pool_id":     policy.RedirectPoolID,
				"redirect_listener_id": policy.RedirectListenerID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Commands() []pluginsdk.Command {
	return []pluginsdk.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-r", Description: "Rules"},
	}
}

func (r *Resource) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	switch command.Key {
	case "s":
		return r.show(ctx, row)
	case "shift-r":
		return r.rules(ctx, row)
	}

	return pluginsdk.Result{}, nil
}
