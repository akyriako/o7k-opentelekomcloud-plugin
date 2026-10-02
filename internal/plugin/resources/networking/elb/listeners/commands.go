package listeners

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/listeners"
)

func (r *Resource) show(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, err
	}

	listener, err := listeners.Get(client, row.ID).Extract()
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting ELB listener %q: %w", row.ID, err)
	}

	content, err := json.Marshal(listener)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("marshalling ELB listener %q: %w", row.ID, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      row.ID,
			Content: content,
		},
	}, nil
}

func (r *Resource) policies(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	scope := make(map[string]string)

	for key, value := range pluginsdk.Scope(ctx) {
		scope[key] = value
	}

	scope["listener_id"] = row.ID

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "elb-policies",
			Scope:    scope,
		},
	}, nil
}

func (r *Resource) ipGroup(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	scope := make(map[string]string)

	for key, value := range pluginsdk.Scope(ctx) {
		scope[key] = value
	}

	scope["ipgroup_id"] = row.Fields["ipgroup_id"]

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "elb-ipgroups",
			Scope:    scope,
		},
	}, nil
}
