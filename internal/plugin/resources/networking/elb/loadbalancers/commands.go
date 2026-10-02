package loadbalancers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/loadbalancers"
)

func (r *Resource) show(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, err
	}

	lb, err := loadbalancers.Get(client, row.ID).Extract()
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting ELB load balancer %q: %w", row.ID, err)
	}

	content, err := json.Marshal(lb)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("marshalling ELB load balancer %q: %w", row.ID, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      row.ID,
			Content: content,
		},
	}, nil
}

func (r *Resource) listeners(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	scope := make(map[string]string)

	for key, value := range pluginsdk.Scope(ctx) {
		scope[key] = value
	}

	scope["loadbalancer_id"] = row.ID

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "elb-listeners",
			Scope:    scope,
		},
	}, nil
}

func (r *Resource) pools(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	scope := make(map[string]string)

	for key, value := range pluginsdk.Scope(ctx) {
		scope[key] = value
	}

	scope["loadbalancer_id"] = row.ID

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "elb-pools",
			Scope:    scope,
		},
	}, nil
}

func (r *Resource) flavors(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	scope := make(map[string]string)

	for key, value := range pluginsdk.Scope(ctx) {
		scope[key] = value
	}

	scope["l4_flavor_id"] = row.Fields["l4_flavor_id"]
	scope["l4_scale_flavor_id"] = row.Fields["l4_scale_flavor_id"]
	scope["l7_flavor_id"] = row.Fields["l7_flavor_id"]
	scope["l7_scale_flavor_id"] = row.Fields["l7_scale_flavor_id"]

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "elb-flavors",
			Scope:    scope,
		},
	}, nil
}
