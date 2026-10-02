package policies

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akyriako/o7k/pluginsdk"
	elbpolicies "github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/policies"
)

func (r *Resource) show(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, err
	}

	policy, err := elbpolicies.Get(client, row.ID).Extract()
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting ELB policy %q: %w", row.ID, err)
	}

	content, err := json.Marshal(policy)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("marshalling ELB policy %q: %w", row.ID, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      row.ID,
			Content: content,
		},
	}, nil
}

func (r *Resource) rules(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	scope := make(map[string]string)

	for key, value := range pluginsdk.Scope(ctx) {
		scope[key] = value
	}

	scope["policy_id"] = row.ID

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "elb-rules",
			Scope:    scope,
		},
	}, nil
}
