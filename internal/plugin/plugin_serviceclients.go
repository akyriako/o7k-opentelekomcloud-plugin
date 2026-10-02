package plugin

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/pluginsdk"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/credentials"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/obs"
)

func (p *Plugin) OBS(ctx context.Context) (*obs.ObsClient, error) {
	current, err := p.host.Context(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting current context: %w", err)
	}

	cloud, err := loadCloud(current)
	if err != nil {
		return nil, err
	}

	if current.Region != "" {
		cloud.RegionName = current.Region
	}

	if cloud.RegionName == "" {
		cloud.RegionName = "eu-de"
	}

	accessKey := cloud.AuthInfo.AccessKey
	secretKey := cloud.AuthInfo.SecretKey
	securityToken := cloud.AuthInfo.SecurityToken

	if accessKey == "" || secretKey == "" {
		provider, err := p.provider.Client(ctx)
		if err != nil {
			return nil, fmt.Errorf("getting provider client: %w", err)
		}

		identity, err := openstack.NewIdentityV3(provider, golangsdk.EndpointOpts{})
		if err != nil {
			return nil, fmt.Errorf("creating identity client: %w", err)
		}

		credential, err := credentials.CreateTemporary(identity, credentials.CreateTemporaryOpts{
			Methods: []string{"token"},
			Token:   provider.TokenID,
		}).Extract()
		if err != nil {
			return nil, fmt.Errorf("creating temporary OBS credentials: %w", err)
		}

		accessKey = credential.AccessKey
		secretKey = credential.SecretKey
		securityToken = credential.SecurityToken
	}

	endpoint := fmt.Sprintf("https://obs.%s.otc.t-systems.com", cloud.RegionName)

	options := []obs.Configurer{
		obs.WithRegion(cloud.RegionName),
		obs.WithRequestContext(ctx),
	}

	if cloud.RegionName != "" {
		endpoint = fmt.Sprintf("https://obs.%s.otc.t-systems.com", cloud.RegionName)
		options = append(options, obs.WithRegion(cloud.RegionName))
	}

	if securityToken != "" {
		options = append(options, obs.WithSecurityToken(securityToken))
	}

	client, err := obs.New(accessKey, secretKey, endpoint, options...)
	if err != nil {
		return nil, fmt.Errorf("creating OBS client: %w", err)
	}

	return client, nil
}

func (p *Plugin) CCEV3(ctx context.Context) (*golangsdk.ServiceClient, error) {
	return pluginsdk.GetServiceClient(ctx, p.provider, "cce", func(provider *golangsdk.ProviderClient, current pluginsdk.Context) (*golangsdk.ServiceClient, error) {
		return openstack.NewCCE(provider, golangsdk.EndpointOpts{Region: current.Region})
	})
}

func (p *Plugin) NetworkV2(ctx context.Context) (*golangsdk.ServiceClient, error) {
	return pluginsdk.GetServiceClient(ctx, p.provider, "network", func(provider *golangsdk.ProviderClient, current pluginsdk.Context) (*golangsdk.ServiceClient, error) {
		return openstack.NewNetworkV2(provider, golangsdk.EndpointOpts{Region: current.Region})
	})
}

func (p *Plugin) RDSV3(ctx context.Context) (*golangsdk.ServiceClient, error) {
	return pluginsdk.GetServiceClient(ctx, p.provider, "rds-v3", func(provider *golangsdk.ProviderClient, current pluginsdk.Context) (*golangsdk.ServiceClient, error) {
		return openstack.NewRDSV3(provider, golangsdk.EndpointOpts{Region: current.Region})
	})
}

func (p *Plugin) ELBV3(ctx context.Context) (*golangsdk.ServiceClient, error) {
	return pluginsdk.GetServiceClient(ctx, p.provider, "elb-v3", func(provider *golangsdk.ProviderClient, current pluginsdk.Context) (*golangsdk.ServiceClient, error) {
		return openstack.NewELBV3(provider, golangsdk.EndpointOpts{Region: current.Region})
	})
}
