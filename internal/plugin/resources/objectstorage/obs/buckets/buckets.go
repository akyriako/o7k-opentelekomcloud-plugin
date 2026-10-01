package buckets

import (
	"context"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/obs"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewObsBuckets(p *plugin.Plugin) *Resource {
	return &Resource{plugin: p}
}

func (r *Resource) Service() string {
	return "obs"
}

func (r *Resource) Kind() string {
	return "obs-buckets"
}

func (r *Resource) Title() string {
	return "OBS Buckets"
}

func (r *Resource) Aliases() []string {
	return []string{"buckets", "bucket"}
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "location", Title: "LOCATION", MinWidth: 14},
		{Key: "type", Title: "TYPE", MinWidth: 12},
		{Key: "created", Title: "CREATED", MinWidth: 20},
	}
}

func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	client, err := r.plugin.OBS(ctx)
	if err != nil {
		return nil, err
	}

	result, err := client.ListBuckets(&obs.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	rows := make([]pluginsdk.Row, 0, len(result.Buckets))

	for _, bucket := range result.Buckets {
		rows = append(rows, pluginsdk.Row{
			ID: bucket.Name,
			Fields: map[string]string{
				"name":     bucket.Name,
				"location": bucket.Location,
				"type":     bucket.BucketType,
				"created":  bucket.CreationDate.Format("2006-01-02 15:04:05"),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Commands() []pluginsdk.Command {
	return []pluginsdk.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}

func (r *Resource) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	switch command.Key {
	case "s":
		return r.show(ctx, row)
	}

	return pluginsdk.Result{}, nil
}
