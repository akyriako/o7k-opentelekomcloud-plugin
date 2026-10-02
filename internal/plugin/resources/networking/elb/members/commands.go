package members

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akyriako/o7k/pluginsdk"
	elbmembers "github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/members"
)

func (r *Resource) show(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	poolID := row.Fields["pool_id"]
	if poolID == "" {
		return pluginsdk.Result{}, fmt.Errorf("ELB pool ID is required")
	}

	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, err
	}

	member, err := elbmembers.Get(client, poolID, row.ID).Extract()
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting ELB member %q: %w", row.ID, err)
	}

	content, err := json.Marshal(member)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("marshalling ELB member %q: %w", row.ID, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      row.ID,
			Content: content,
		},
	}, nil
}

func (r *Resource) server(row pluginsdk.Row) (pluginsdk.Result, error) {
	serverID := row.Fields["server_id"]
	if serverID == "" {
		return pluginsdk.Result{}, fmt.Errorf("no server found for ELB member %q", row.ID)
	}

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "servers",
			Field:    "id",
			Value:    serverID,
		},
	}, nil
}
