package instances

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akyriako/o7k/pluginsdk"
	rdsbackups "github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/backups"
	rdsinstances "github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"
)

func (r *Resource) show(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	client, err := r.plugin.RDSV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, err
	}

	result, err := rdsinstances.List(client, rdsinstances.ListOpts{
		Id: row.ID,
	})
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting RDS instance %q: %w", row.ID, err)
	}

	if len(result.Instances) == 0 {
		return pluginsdk.Result{}, fmt.Errorf("RDS instance %q not found", row.ID)
	}

	content, err := json.Marshal(result.Instances[0])
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("marshalling RDS instance %q: %w", row.ID, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      row.ID,
			Content: content,
		},
	}, nil
}

func (r *Resource) securityGroup(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	securityGroupID := row.Fields["security_group_id"]
	if securityGroupID == "" {
		return pluginsdk.Result{}, fmt.Errorf("RDS instance has no security group")
	}

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "securitygroups",
			Field:    "id",
			Value:    securityGroupID,
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

func (r *Resource) flavor(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	scope := make(map[string]string)

	for key, value := range pluginsdk.Scope(ctx) {
		scope[key] = value
	}

	scope["engine"] = row.Fields["engine"]
	scope["version"] = row.Fields["version"]
	scope["spec_code"] = row.Fields["flavor"]

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "rds-flavors",
			Scope:    scope,
		},
	}, nil
}

func (r *Resource) nodes(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	scope := make(map[string]string)

	for key, value := range pluginsdk.Scope(ctx) {
		scope[key] = value
	}

	scope["instance_id"] = row.ID

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "rds-nodes",
			Scope:    scope,
		},
	}, nil
}

func (r *Resource) backups(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	scope := make(map[string]string)

	for key, value := range pluginsdk.Scope(ctx) {
		scope[key] = value
	}

	scope["instance_id"] = row.ID

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "rds-backups",
			Scope:    scope,
		},
	}, nil
}

func (r *Resource) backupPolicy(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	client, err := r.plugin.RDSV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, err
	}

	policy, err := rdsbackups.ShowBackupPolicy(client, row.ID)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting backup policy for RDS instance %q: %w", row.ID, err)
	}

	content, err := json.Marshal(policy)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("marshalling backup policy for RDS instance %q: %w", row.ID, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      row.ID,
			Content: content,
		},
	}, nil
}

func (r *Resource) parameters(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	scope := make(map[string]string)

	for key, value := range pluginsdk.Scope(ctx) {
		scope[key] = value
	}

	scope["instance_id"] = row.ID

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "rds-parameters",
			Scope:    scope,
		},
	}, nil
}
