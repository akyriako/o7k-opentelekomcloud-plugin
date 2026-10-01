package plugin

import (
	"github.com/akyriako/o7k/pluginsdk"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
)

type Plugin struct {
	host      pluginsdk.Host
	provider  *pluginsdk.ClientProvider[*golangsdk.ProviderClient]
	resources []pluginsdk.Resource
}

func New() *Plugin {
	return &Plugin{}
}

func (p *Plugin) SetHost(host pluginsdk.Host) {
	p.host = host
	p.provider = newClient(host)
}

func (p *Plugin) Register(resources ...pluginsdk.Resource) {
	p.resources = append(p.resources, resources...)
}

func (p *Plugin) Metadata() pluginsdk.Metadata {
	return pluginsdk.Metadata{
		Name:    "T Cloud Public",
		Version: "0.2.0-dev.24",
		Color:   "#E20074",
	}
}

func (p *Plugin) Resources() []pluginsdk.Resource {
	return p.resources
}

func (p *Plugin) Host() pluginsdk.Host {
	return p.host
}

func (p *Plugin) Provider() *pluginsdk.ClientProvider[*golangsdk.ProviderClient] {
	return p.provider
}
