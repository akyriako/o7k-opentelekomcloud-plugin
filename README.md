# o7k T Cloud Public  Plugin

T Cloud Public provider plugin for [o7k](https://github.com/akyriako/o7k).

The plugin adds support for browsing and interacting with T Cloud Public resources directly from o7k.

## Installation

Get the URL of the appropriate binary for your platform from the
[GitHub Releases](https://github.com/akyriako/o7k-opentelekomcloud-plugin/releases)
page and install it in the o7k plugin directory, e.g.:

```bash
o7k plugin install https://github.com/akyriako/o7k-opentelekomcloud-plugin/releases/download/v0.1.0/o7k-opentelekomcloud-plugin_0.1.0_linux_amd64
```

Restart o7k after installing or updating the plugin.

## Supported Services

| Service                     | Resources       | Aliases            | Status |
|-----------------------------|-----------------|--------------------|:------:|
| Cloud Container Engine      | Clusters        | cce, cce-clusters  |   ✅   |
| Cloud Container Engine      | Nodes           | cce-nodes          |   ✅   |
| Cloud Container Engine      | Node Pools      | cce-nodepools      |   ✅   |
| Object Storage Service      | Buckets         | buckets, bucket    |   ✅   |
| Relational Database Service | Instances       | rds, rds-instances |   ✅   |
| Relational Database Service | Flavors         | rds-flavors        |   ✅   |
| Relational Database Service | Nodes           | rds-nodes          |   ✅   |
| Relational Database Service | Backups         | rds-backups        |   ✅   |
| Relational Database Service | Parameters      | rds-parameters     |   ✅   |
| Distributed Cache Service   | Instances       | dcs-instances      |   ✅   |
| Distributed Cache Service   | Parameters      | dcs-parameters     |   ✅   |
| Distributed Cache Service   | Whitelists      | dcs-whitelists     |   ✅   |
| Elastic Load Balancing      | Load Balancers  | elb                |   ✅   |
| Elastic Load Balancing      | Listeners       | elb-listeners      |   ✅   |
| Elastic Load Balancing      | Pools           | elb-pools          |   ✅   |
| Elastic Load Balancing      | Members         | elb-members        |   ✅   |
| Elastic Load Balancing      | Health Monitors | elb-healthmonitors |   ✅   |
| Elastic Load Balancing      | Policies        | elb-policies       |   ✅   |
| Elastic Load Balancing      | Rules           | elb-rules          |   ✅   |
| Elastic Load Balancing      | Flavors         | elb-flavors        |   ✅   |
| Elastic Load Balancing      | IP Groups       | elb-ipgroups       |   ✅   |

## Development
Build the plugin locally:

```
make build
```

Build and install it into your local o7k plugin directory:

```
make install
```

The development build includes version, Git commit and build-date metadata.
To test the release build locally:

```
goreleaser release --snapshot --clean
```

## Release

Releases are built with GoReleaser and published as standalone binaries for
Linux and macOS on AMD64 and ARM64. Create a release with:

```
make release-tag TAG=0.1.0
```

This creates and pushes the corresponding `v0.1.0` Git tag, which triggers the
GitHub release workflow.
