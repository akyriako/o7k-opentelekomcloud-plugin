package buckets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/obs"
)

func (r *Resource) show(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	client, err := r.plugin.OBS(ctx)
	if err != nil {
		return pluginsdk.Result{}, err
	}

	metadata, err := client.GetBucketMetadata(&obs.GetBucketMetadataInput{
		Bucket: row.ID,
	})
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting OBS bucket %q metadata: %w", row.ID, err)
	}

	content, err := json.Marshal(metadata)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("marshalling OBS bucket %q metadata: %w", row.ID, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      row.ID,
			Content: content,
		},
	}, nil
}
