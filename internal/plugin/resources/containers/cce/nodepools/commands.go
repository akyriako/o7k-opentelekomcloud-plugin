package nodepools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/nodepools"
)

func (r *Resource) show(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	clusterID := row.Fields["cluster_id"]
	if clusterID == "" {
		return pluginsdk.Result{}, fmt.Errorf("cluster_id is required")
	}

	client, err := r.plugin.CCEV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting CCE client: %w", err)
	}

	nodePool, err := nodepools.Get(client, clusterID, row.ID)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting CCE node pool %q: %w", row.ID, err)
	}

	content, err := json.Marshal(nodePool)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("encoding CCE node pool %q: %w", row.ID, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      row.ID,
			Content: content,
		},
	}, nil
}

func (r *Resource) flavor(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "flavors",
			Field:    "name",
			Value:    row.Fields["flavor"],
		},
	}, nil
}

func (r *Resource) securityGroup(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	securityGroups := strings.Split(row.Fields["security_groups"], ",")
	if len(securityGroups) == 0 || securityGroups[0] == "" {
		return pluginsdk.Result{}, fmt.Errorf("node pool has no custom security groups")
	}

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "securitygroups",
			Field:    "id",
			Value:    securityGroups[0],
		},
	}, nil
}

func (r *Resource) scale(ctx context.Context, row pluginsdk.Row, delta int) (pluginsdk.Result, error) {
	clusterID := row.Fields["cluster_id"]
	if clusterID == "" {
		return pluginsdk.Result{}, fmt.Errorf("cluster_id is required")
	}

	client, err := r.plugin.CCEV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting CCE client: %w", err)
	}

	nodePool, err := nodepools.Get(client, clusterID, row.ID)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting CCE node pool %q: %w", row.ID, err)
	}

	nodeCount := nodePool.Spec.InitialNodeCount + delta
	if nodeCount < 0 {
		return pluginsdk.Result{}, fmt.Errorf("node pool cannot be scaled below 0 nodes")
	}

	opts := nodepools.UpdateOpts{
		Metadata: nodepools.UpdateMetaData{
			Name: nodePool.Metadata.Name,
		},
		Spec: nodepools.UpdateSpec{
			NodeTemplate:     nodePool.Spec.NodeTemplate,
			InitialNodeCount: nodeCount,
		},
	}

	_, err = nodepools.Update(client, clusterID, row.ID, opts)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("scaling CCE node pool %q to %d nodes: %w", row.ID, nodeCount, err)
	}

	return pluginsdk.Result{}, nil
}
