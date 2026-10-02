package listeners

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/listeners"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewElbListeners(p *plugin.Plugin) *Resource {
	return &Resource{plugin: p}
}

func (r *Resource) Service() string {
	return "elb"
}

func (r *Resource) Kind() string {
	return "elb-listeners"
}

func (r *Resource) Title() string {
	return "ELB Listeners"
}

func (r *Resource) Aliases() []string {
	return nil
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "protocol", Title: "PROTOCOL", MinWidth: 10},
		{Key: "port", Title: "PORT", MinWidth: 8},
		{Key: "admin_state_up", Title: "ADMIN UP", MinWidth: 10},
		{Key: "default_pool_id", Title: "DEFAULT POOL", MinWidth: 36},
		{Key: "http2", Title: "HTTP/2", MinWidth: 8},
	}
}

func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	scope := pluginsdk.Scope(ctx)
	loadbalancerID := scope["loadbalancer_id"]
	if loadbalancerID == "" {
		return nil, fmt.Errorf("ELB load balancer ID is required")
	}

	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return nil, err
	}

	pages, err := listeners.List(client, listeners.ListOpts{
		LoadBalancerID: []string{loadbalancerID},
	}).AllPages()
	if err != nil {
		return nil, err
	}

	result, err := listeners.ExtractListeners(pages)
	if err != nil {
		return nil, err
	}

	rows := make([]pluginsdk.Row, 0, len(result))

	for _, listener := range result {
		rows = append(rows, pluginsdk.Row{
			ID: listener.ID,
			Fields: map[string]string{
				"id":                 listener.ID,
				"name":               listener.Name,
				"protocol":           listener.Protocol,
				"port":               strconv.Itoa(listener.ProtocolPort),
				"admin_state_up":     strconv.FormatBool(listener.AdminStateUp),
				"default_pool_id":    listener.DefaultPoolID,
				"http2":              strconv.FormatBool(listener.Http2Enable),
				"loadbalancer_id":    loadbalancerID,
				"tls_container_ref":  listener.DefaultTlsContainerRef,
				"ca_container_ref":   listener.CAContainerRef,
				"tls_ciphers_policy": listener.TlsCiphersPolicy,
				"security_policy_id": listener.SecurityPolicy,
				"ipgroup_id":         listener.IpGroup.IpGroupID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Commands() []pluginsdk.Command {
	return []pluginsdk.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-p", Description: "Policies"},
		{Key: "shift-g", Description: "IP Group"},
	}
}

func (r *Resource) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	switch command.Key {
	case "s":
		return r.show(ctx, row)
	case "shift-p":
		return r.policies(ctx, row)
	case "shift-g":
		return r.ipGroup(ctx, row)
	}

	return pluginsdk.Result{}, nil
}
