package pools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akyriako/o7k/pluginsdk"
	elbpools "github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/pools"
)

func (r *Resource) show(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, err
	}

	pool, err := elbpools.Get(client, row.ID).Extract()
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting ELB pool %q: %w", row.ID, err)
	}

	content, err := json.Marshal(pool)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("marshalling ELB pool %q: %w", row.ID, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      row.ID,
			Content: content,
		},
	}, nil
}

func (r *Resource) members(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	scope := make(map[string]string)

	for key, value := range pluginsdk.Scope(ctx) {
		scope[key] = value
	}

	scope["pool_id"] = row.ID

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "elb-members",
			Scope:    scope,
		},
	}, nil
}

func (r *Resource) healthMonitor(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	monitorID := row.Fields["monitor_id"]
	if monitorID == "" {
		return pluginsdk.Result{}, fmt.Errorf("ELB pool %q has no health monitor", row.ID)
	}

	scope := make(map[string]string)

	for key, value := range pluginsdk.Scope(ctx) {
		scope[key] = value
	}

	scope["monitor_id"] = monitorID
	scope["pool_id"] = row.ID

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "elb-healthmonitors",
			Scope:    scope,
		},
	}, nil
}
