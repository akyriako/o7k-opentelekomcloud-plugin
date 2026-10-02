package healthmonitors

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/monitors"
)

func (r *Resource) show(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, err
	}

	monitor, err := monitors.Get(client, row.ID).Extract()
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting ELB health monitor %q: %w", row.ID, err)
	}

	content, err := json.Marshal(monitor)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("marshalling ELB health monitor %q: %w", row.ID, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      row.ID,
			Content: content,
		},
	}, nil
}

func (r *Resource) pool(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	poolID := row.Fields["pool_id"]
	if poolID == "" {
		return pluginsdk.Result{}, fmt.Errorf("ELB health monitor %q has no associated pool", row.ID)
	}

	scope := make(map[string]string)

	for key, value := range pluginsdk.Scope(ctx) {
		scope[key] = value
	}

	delete(scope, "monitor_id")
	scope["pool_id"] = poolID

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "elb-pools",
			Field:    "id",
			Value:    poolID,
			Scope:    scope,
		},
	}, nil
}
