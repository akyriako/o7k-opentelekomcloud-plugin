package members

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k-opentelekomcloud-plugin/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	elbmembers "github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/members"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/ports"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewElbMembers(p *plugin.Plugin) *Resource {
	return &Resource{plugin: p}
}

func (r *Resource) Service() string {
	return "elb"
}

func (r *Resource) Kind() string {
	return "elb-members"
}

func (r *Resource) Title() string {
	return "ELB Members"
}

func (r *Resource) Aliases() []string {
	return nil
}

func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "address", Title: "ADDRESS", MinWidth: 16},
		{Key: "port", Title: "PORT", MinWidth: 8},
		{Key: "status", Title: "STATUS", MinWidth: 12},
		{Key: "weight", Title: "WEIGHT", MinWidth: 8},
		{Key: "admin_state_up", Title: "ADMIN UP", MinWidth: 10},
		{Key: "ip_version", Title: "IP VERSION", MinWidth: 10},
	}
}

// List intentionally performs a per-member neutron port lookup to resolve the
// backend server ID. This is a one-off exception because ELB pools typically
// contain a limited number of members and server navigation provides substantial
// value to the user.
//
// Do NOT use this implementation as a general pattern for o7k resources or
// plugins. N+1 API requests should be actively avoided; prefer bulk APIs,
// server-side filtering, or data already returned by the primary list operation
// whenever possible.
func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	scope := pluginsdk.Scope(ctx)
	poolID := scope["pool_id"]
	if poolID == "" {
		return nil, fmt.Errorf("ELB pool ID is required")
	}

	client, err := r.plugin.ELBV3(ctx)
	if err != nil {
		return nil, err
	}

	networkingClient, err := r.plugin.NetworkV2(ctx)
	if err != nil {
		return nil, err
	}

	pages, err := elbmembers.List(client, poolID, elbmembers.ListOpts{}).AllPages()
	if err != nil {
		return nil, err
	}

	result, err := elbmembers.ExtractMembers(pages)
	if err != nil {
		return nil, err
	}

	rows := make([]pluginsdk.Row, 0, len(result))

	for _, member := range result {
		serverID := ""

		portPages, err := ports.List(networkingClient, ports.ListOpts{
			FixedIps: []string{"ip_address=" + member.Address},
		}).AllPages()
		if err != nil {
			return nil, fmt.Errorf("listing ports for ELB member %q: %w", member.ID, err)
		}

		memberPorts, err := ports.ExtractPorts(portPages)
		if err != nil {
			return nil, fmt.Errorf("extracting ports for ELB member %q: %w", member.ID, err)
		}

		for _, port := range memberPorts {
			if port.DeviceID != "" {
				serverID = port.DeviceID
				break
			}
		}

		rows = append(rows, pluginsdk.Row{
			ID: member.ID,
			Fields: map[string]string{
				"id":             member.ID,
				"name":           member.Name,
				"address":        member.Address,
				"port":           strconv.Itoa(member.ProtocolPort),
				"status":         member.OperatingStatus,
				"weight":         strconv.Itoa(member.Weight),
				"admin_state_up": strconv.FormatBool(member.AdminStateUp),
				"ip_version":     member.IpVersion,
				"subnet_id":      member.SubnetID,
				"pool_id":        poolID,
				"server_id":      serverID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Commands() []pluginsdk.Command {
	return []pluginsdk.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-s", Description: "Server"},
	}
}

func (r *Resource) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	switch command.Key {
	case "s":
		return r.show(ctx, row)
	case "shift-s":
		return r.server(row)
	}

	return pluginsdk.Result{}, nil
}
