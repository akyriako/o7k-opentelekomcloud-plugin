package healthmonitors

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/monitors"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewElbHealthMonitors(p *plugin.Plugin) *Resource {
	return &Resource{plugin: p}
}

func (r *Resource) Service() string {
	return "elb"
}

func (r *Resource) Kind() string {
	return "elb-healthmonitors"
}

func (r *Resource) Title() string {
	return "ELB Health Monitors"
}

func (r *Resource) Aliases() []string {
	return nil
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "type", Title: "TYPE", MinWidth: 10},
		{Key: "delay", Title: "DELAY", MinWidth: 8},
		{Key: "timeout", Title: "TIMEOUT", MinWidth: 8},
		{Key: "max_retries", Title: "MAX RETRIES", MinWidth: 12},
		{Key: "max_retries_down", Title: "MAX RETRIES DOWN", MinWidth: 16},
		{Key: "admin_state_up", Title: "ADMIN UP", MinWidth: 10},
	}
}

func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	scope := pluginsdk.Scope(ctx)
	monitorID := scope["monitor_id"]

	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return nil, err
	}

	if monitorID != "" {
		monitor, err := monitors.Get(client, monitorID).Extract()
		if err != nil {
			return nil, fmt.Errorf("getting ELB health monitor %q: %w", monitorID, err)
		}

		return []pluginsdk.Row{r.row(*monitor)}, nil
	}

	pages, err := monitors.List(client, monitors.ListOpts{}).AllPages()
	if err != nil {
		return nil, fmt.Errorf("listing ELB health monitors: %w", err)
	}

	result, err := monitors.ExtractMonitors(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting ELB health monitors: %w", err)
	}

	rows := make([]pluginsdk.Row, 0, len(result))

	for _, monitor := range result {
		rows = append(rows, r.row(monitor))
	}

	return rows, nil
}

func (r *Resource) row(monitor monitors.Monitor) pluginsdk.Row {
	poolID := ""
	if len(monitor.Pools) > 0 {
		poolID = monitor.Pools[0].ID
	}

	return pluginsdk.Row{
		ID: monitor.ID,
		Fields: map[string]string{
			"id":               monitor.ID,
			"name":             monitor.Name,
			"type":             string(monitor.Type),
			"delay":            strconv.Itoa(monitor.Delay),
			"timeout":          strconv.Itoa(monitor.Timeout),
			"max_retries":      strconv.Itoa(monitor.MaxRetries),
			"max_retries_down": strconv.Itoa(monitor.MaxRetriesDown),
			"admin_state_up":   strconv.FormatBool(monitor.AdminStateUp),
			"monitor_port":     strconv.Itoa(monitor.MonitorPort),
			"http_method":      monitor.HTTPMethod,
			"url_path":         monitor.URLPath,
			"domain_name":      monitor.DomainName,
			"expected_codes":   monitor.ExpectedCodes,
			"pool_id":          poolID,
		},
	}
}

func (r *Resource) Commands() []pluginsdk.Command {
	return []pluginsdk.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-p", Description: "Pool"},
	}
}

func (r *Resource) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	switch command.Key {
	case "s":
		return r.show(ctx, row)
	case "shift-p":
		return r.pool(ctx, row)
	}

	return pluginsdk.Result{}, nil
}
