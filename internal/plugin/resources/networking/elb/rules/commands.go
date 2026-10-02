package rules

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akyriako/o7k/pluginsdk"
	elbrules "github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/rules"
)

func (r *Resource) show(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	policyID := row.Fields["policy_id"]
	if policyID == "" {
		return pluginsdk.Result{}, fmt.Errorf("ELB policy ID is required")
	}

	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, err
	}

	rule, err := elbrules.Get(client, policyID, row.ID).Extract()
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting ELB rule %q: %w", row.ID, err)
	}

	content, err := json.Marshal(rule)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("marshalling ELB rule %q: %w", row.ID, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      row.ID,
			Content: content,
		},
	}, nil
}
