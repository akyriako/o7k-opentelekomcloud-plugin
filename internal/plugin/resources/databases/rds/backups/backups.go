package backups

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	rdsbackups "github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/backups"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewRdsBackups(p *plugin.Plugin) *Resource {
	return &Resource{plugin: p}
}

func (r *Resource) Service() string {
	return "rds"
}

func (r *Resource) Kind() string {
	return "rds-backups"
}

func (r *Resource) Title() string {
	return "RDS Backups"
}

func (r *Resource) Aliases() []string {
	return nil
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 12},
		{Key: "type", Title: "TYPE", MinWidth: 12},
		{Key: "size", Title: "SIZE", MinWidth: 12},
		{Key: "engine", Title: "ENGINE", MinWidth: 12},
		{Key: "version", Title: "VERSION", MinWidth: 10},
		{Key: "begin_time", Title: "BEGIN", MinWidth: 22},
		{Key: "end_time", Title: "END", MinWidth: 22},
	}
}

func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	scope := pluginsdk.Scope(ctx)
	instanceID := scope["instance_id"]

	if instanceID == "" {
		return nil, fmt.Errorf("instance_id is required")
	}

	client, err := r.plugin.RDSV3(ctx)
	if err != nil {
		return nil, err
	}

	result, err := rdsbackups.List(client, rdsbackups.ListOpts{
		InstanceID: instanceID,
	})
	if err != nil {
		return nil, fmt.Errorf("listing backups for RDS instance %q: %w", instanceID, err)
	}

	rows := make([]pluginsdk.Row, 0, len(result))

	for _, backup := range result {
		rows = append(rows, pluginsdk.Row{
			ID: backup.ID,
			Fields: map[string]string{
				"id":          backup.ID,
				"name":        backup.Name,
				"status":      string(backup.Status),
				"type":        backup.Type,
				"size":        formatSize(backup.Size),
				"engine":      backup.Datastore.Type,
				"version":     backup.Datastore.Version,
				"begin_time":  backup.BeginTime,
				"end_time":    backup.EndTime,
				"instance_id": backup.InstanceID,
				"description": backup.Description,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Commands() []pluginsdk.Command {
	return nil
}

func (r *Resource) Execute(context.Context, pluginsdk.Command, pluginsdk.Row) (pluginsdk.Result, error) {
	return pluginsdk.Result{}, nil
}

func formatSize(size int) string {
	if size < 1024 {
		return fmt.Sprintf("%d kB", size)
	}

	if size < 1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(size)/1024)
	}

	return fmt.Sprintf("%.1f GB", float64(size)/(1024*1024))
}
