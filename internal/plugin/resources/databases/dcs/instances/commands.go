package instances

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/akyriako/o7k/pluginsdk"
	dcsinstances "github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v2/instance"
	dcsssl "github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v2/ssl"
)

func (r *Resource) show(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	client, err := r.plugin.DCSV2(ctx)
	if err != nil {
		return pluginsdk.Result{}, err
	}

	instance, err := dcsinstances.Get(client, row.ID)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting DCS instance %q: %w", row.ID, err)
	}

	content, err := json.Marshal(instance)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("marshalling DCS instance %q: %w", row.ID, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      row.ID,
			Content: content,
		},
	}, nil
}

func (r *Resource) network(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "networks",
			Field:    "id",
			Value:    row.Fields["subnet_id"],
		},
	}, nil
}

func (r *Resource) securityGroup(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "securitygroups",
			Field:    "id",
			Value:    row.Fields["security_group_id"],
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
			Resource: "dcs-parameters",
			Scope:    scope,
		},
	}, nil
}

func (r *Resource) whitelist(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	scope := make(map[string]string)

	for key, value := range pluginsdk.Scope(ctx) {
		scope[key] = value
	}

	scope["instance_id"] = row.ID

	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "dcs-whitelists",
			Scope:    scope,
		},
	}, nil
}

//func (r *Resource) toggleSSL(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
//	client, err := r.plugin.DCSV2(ctx)
//	if err != nil {
//		return pluginsdk.Result{}, err
//	}
//
//	enabled, err := strconv.ParseBool(row.Fields["ssl"])
//	if err != nil {
//		return pluginsdk.Result{}, fmt.Errorf("invalid SSL status %q: %w", row.Fields["ssl"], err)
//	}
//
//	target := !enabled
//
//	_, err = dcsssl.Update(client, dcsssl.SslOpts{
//		InstanceId: row.ID,
//		Enabled:    &target,
//	})
//	if err != nil {
//		return pluginsdk.Result{}, fmt.Errorf("updating SSL for DCS instance %q: %w", row.ID, err)
//	}
//
//	return pluginsdk.Result{}, nil
//}

func (r *Resource) toggleSSL(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	client, err := r.plugin.DCSV2(ctx)
	if err != nil {
		return pluginsdk.Result{}, err
	}

	info, err := dcsssl.Get(client, row.ID)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting SSL status for DCS instance %q: %w", row.ID, err)
	}

	target := !info.Enabled

	_, err = dcsssl.Update(client, dcsssl.SslOpts{
		InstanceId: row.ID,
		Enabled:    &target,
	})
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("updating SSL for DCS instance %q: %w", row.ID, err)
	}

	return pluginsdk.Result{}, nil
}

func (r *Resource) downloadCertificate(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	client, err := r.plugin.DCSV2(ctx)
	if err != nil {
		return pluginsdk.Result{}, err
	}

	cert, err := dcsssl.DownloadCert(client, row.ID)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting SSL certificate download link for DCS instance %q: %w", row.ID, err)
	}

	if cert.Link == "" {
		return pluginsdk.Result{}, fmt.Errorf("empty SSL certificate download link for DCS instance %q", row.ID)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cert.Link, nil)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("creating certificate download request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("downloading SSL certificate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return pluginsdk.Result{}, fmt.Errorf("downloading SSL certificate: HTTP %s", resp.Status)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting home directory: %w", err)
	}

	dir := filepath.Join(home, "Downloads")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return pluginsdk.Result{}, fmt.Errorf("creating Downloads directory: %w", err)
	}

	fileName := filepath.Base(cert.FileName)
	if fileName == "" || fileName == "." {
		fileName = "dcs-" + row.ID + "-certificate.crt"
	}

	path := filepath.Join(dir, fileName)

	file, err := os.Create(path)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("creating certificate file %q: %w", path, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return pluginsdk.Result{}, fmt.Errorf("writing certificate file %q: %w", path, err)
	}

	return pluginsdk.Result{}, nil
}

func (r *Resource) showSSL(ctx context.Context, row pluginsdk.Row) (pluginsdk.Result, error) {
	client, err := r.plugin.DCSV2(ctx)
	if err != nil {
		return pluginsdk.Result{}, err
	}

	info, err := dcsssl.Get(client, row.ID)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting SSL information for DCS instance %q: %w", row.ID, err)
	}

	content, err := json.Marshal(info)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("marshalling SSL information for DCS instance %q: %w", row.ID, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      row.ID,
			Content: content,
		},
	}, nil
}
