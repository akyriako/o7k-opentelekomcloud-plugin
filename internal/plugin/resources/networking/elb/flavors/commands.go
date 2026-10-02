package flavors

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akyriako/o7k/pluginsdk"
	elbflavors "github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/flavors"
)

func (r *Resource) show(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, err
	}

	flavor, err := elbflavors.Get(client, row.ID).Extract()
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting ELB flavor %q: %w", row.ID, err)
	}

	content, err := json.Marshal(flavor)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("marshalling ELB flavor %q: %w", row.ID, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      row.ID,
			Content: content,
		},
	}, nil
}
