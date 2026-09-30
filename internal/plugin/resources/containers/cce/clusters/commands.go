package clusters

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"

	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/clusters"
)

func (r *Resource) show(ctx context.Context, id string) (pluginsdk.Result, error) {
	client, err := r.plugin.CCEV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting CCE client: %w", err)
	}

	cluster, err := clusters.Get(client, id)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting CCE cluster %q: %w", id, err)
	}

	content, err := json.Marshal(cluster)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("encoding CCE cluster %q: %w", id, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      cluster.Metadata.Id,
			Content: content,
		},
	}, nil
}

func (r *Resource) nodes(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	scope := make(map[string]string)
	maps.Copy(scope, pluginsdk.Scope(ctx))
	scope["cluster_id"] = row.ID

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "cce-nodes",
			Scope:    scope,
		},
	}, nil
}
