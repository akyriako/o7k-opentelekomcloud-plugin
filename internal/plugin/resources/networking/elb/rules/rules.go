package rules

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	elbrules "github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/rules"
)

type listOpts struct{}

func (listOpts) ToRuleListQuery() (string, error) {
	return "", nil
}

type Resource struct {
	plugin *plugin.Plugin
}

func NewElbRules(p *plugin.Plugin) *Resource {
	return &Resource{plugin: p}
}

func (r *Resource) Service() string {
	return "elb"
}

func (r *Resource) Kind() string {
	return "elb-rules"
}

func (r *Resource) Title() string {
	return "ELB Rules"
}

func (r *Resource) Aliases() []string {
	return nil
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "type", Title: "TYPE", MinWidth: 18},
		{Key: "compare_type", Title: "COMPARE", MinWidth: 16},
		{Key: "value", Title: "VALUE", MinWidth: 24, Flex: 1},
	}
}

func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	scope := pluginsdk.Scope(ctx)
	policyID := scope["policy_id"]
	if policyID == "" {
		return nil, fmt.Errorf("ELB policy ID is required")
	}

	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return nil, err
	}

	pages, err := elbrules.List(client, policyID, listOpts{}).AllPages()
	if err != nil {
		return nil, fmt.Errorf("listing ELB rules for policy %q: %w", policyID, err)
	}

	result, err := elbrules.ExtractRules(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting ELB rules for policy %q: %w", policyID, err)
	}

	rows := make([]pluginsdk.Row, 0, len(result))

	for _, rule := range result {
		rows = append(rows, pluginsdk.Row{
			ID: rule.ID,
			Fields: map[string]string{
				"id":           rule.ID,
				"type":         string(rule.Type),
				"compare_type": string(rule.CompareType),
				"value":        rule.Value,
				"policy_id":    policyID,
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
