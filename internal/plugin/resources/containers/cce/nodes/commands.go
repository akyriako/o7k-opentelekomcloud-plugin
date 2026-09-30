package nodes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akyriako/o7k/pluginsdk"
	ccenodes "github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/nodes"
)

func (r *Resource) show(ctx context.Context, nodeID string, clusterID string) (pluginsdk.Result, error) {
	if clusterID == "" {
		return pluginsdk.Result{}, fmt.Errorf("cluster_id is required")
	}

	client, err := r.plugin.CCEV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting CCE client: %w", err)
	}

	node, err := ccenodes.Get(client, clusterID, nodeID)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf(
			"getting CCE node %q from cluster %q: %w",
			nodeID,
			clusterID,
			err,
		)
	}

	content, err := json.Marshal(node)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf(
			"encoding CCE node %q: %w",
			nodeID,
			err,
		)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      node.Metadata.Id,
			Content: content,
		},
	}, nil
}

func (r *Resource) flavor(row pluginsdk.Row) (pluginsdk.Result, error) {
	flavor := row.Fields["flavor"]

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "flavors",
			Field:    "name",
			Value:    flavor,
		},
	}, nil
}

func (r *Resource) image(row pluginsdk.Row) (pluginsdk.Result, error) {
	imageID := row.Fields["image_id"]

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "images",
			Field:    "id",
			Value:    imageID,
		},
	}, nil
}

func (r *Resource) network(row pluginsdk.Row) (pluginsdk.Result, error) {
	subnetID := row.Fields["subnet_id"]

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "networks",
			Field:    "id",
			Value:    subnetID,
		},
	}, nil
}

func (r *Resource) server(row pluginsdk.Row) (pluginsdk.Result, error) {
	serverID := row.Fields["server_id"]

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "servers",
			Field:    "id",
			Value:    serverID,
		},
	}, nil
}
