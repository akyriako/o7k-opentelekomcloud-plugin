package clusters

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"

	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/clusters"
	"gopkg.in/yaml.v3"
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

func (r *Resource) kubeconfig(ctx context.Context, clusterID string) (pluginsdk.Result, error) {
	client, err := r.plugin.CCEV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting CCE client: %w", err)
	}

	cert, err := clusters.GetCert(client, clusterID)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting kubeconfig for CCE cluster %q: %w", clusterID, err)
	}

	content, err := json.Marshal(cert)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("encoding kubeconfig for CCE cluster %q: %w", clusterID, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      clusterID,
			Content: content,
		},
	}, nil
}

func (r *Resource) getKubeconfig(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	current, err := r.plugin.Host().Context(ctx)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting current context: %w", err)
	}

	cloud := current.Cloud
	clusterName := row.Fields["name"]

	if cloud == "" {
		return pluginsdk.Result{}, fmt.Errorf("cloud is required")
	}

	if clusterName == "" {
		return pluginsdk.Result{}, fmt.Errorf("cluster name is required")
	}

	client, err := r.plugin.CCEV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting CCE client: %w", err)
	}

	cert, err := clusters.GetCert(client, row.ID)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting kubeconfig for CCE cluster %q: %w", row.ID, err)
	}

	jsonData, err := json.Marshal(cert)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("encoding kubeconfig for CCE cluster %q: %w", row.ID, err)
	}

	var kubeconfig any
	if err := yaml.Unmarshal(jsonData, &kubeconfig); err != nil {
		return pluginsdk.Result{}, fmt.Errorf("converting kubeconfig for CCE cluster %q: %w", row.ID, err)
	}

	data, err := yaml.Marshal(kubeconfig)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("encoding kubeconfig for CCE cluster %q: %w", row.ID, err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting user home directory: %w", err)
	}

	kubeDir := filepath.Join(home, ".kube")
	if err := os.MkdirAll(kubeDir, 0700); err != nil {
		return pluginsdk.Result{}, fmt.Errorf("creating kubeconfig directory: %w", err)
	}

	filename := fmt.Sprintf("o7k-generated_%s_%s.yaml", cloud, clusterName)
	path := filepath.Join(kubeDir, filename)

	if err := os.WriteFile(path, data, 0600); err != nil {
		return pluginsdk.Result{}, fmt.Errorf("writing kubeconfig %q: %w", path, err)
	}

	return pluginsdk.Result{}, nil
}
