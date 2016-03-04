# snap-plugin-collector-neutron

snap plugin for collecting metrics from OpenStack Neutron module. 

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started
Plugin collects metrics by communicating with OpenStack by REST API using Networking API v2.0.
It can be used locally on the OpenStack host or remotely communicating with the OpenStack host via HTTP(S).

### System Requirements
 - Linux
 - OpenStack deployment available

### Installation
#### Download neutron plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [Github Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-neutron
Clone repo into `$GOPATH/src/github/intelsdi-x/`:
```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-neutron
```
Build the plugin by running make in repo:
```
$ make
```
This builds the plugin in `/build/rootfs`

## Documentation

### Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Data Type | Description
----------|-----------|-----------------------
intel/openstack/neutron/\<tenant_name\>/networks_count | int64 | number of tenant networks
intel/openstack/neutron/\<tenant_name\>/subnets_count  | int64 | number of tenant subnets
intel/openstack/neutron/\<tenant_name\>/routers_count | int64 | number of tenant routers
intel/openstack/neutron/\<tenant_name\>/ports_count | int64 | number of tenant ports
intel/openstack/neutron/\<tenant_name\>/floatingips_count | int64 | number of tenant floating IPs
intel/openstack-neutron/\<tenant_name\>/quotas_{floatingip,ikepolicy,ipsec_site_connection,ipsecpolicy,network, port,router,security_group,security_group_rule,subnet}  | int64 |  quotas for a tenant. 

### snap's Global Config
Global configuration files are described in snap's documentation. You have to add section "neutron" in "collector" section and then specify following options:
- `"openstack_auth_url"` - URL for OpenStack Identity endpoint (ex. `"http://127.0.0.1:5000/v2.0/"`)
- `"openstack_user"` - user name used to authenticate (ex. `"admin"`)
- `"openstack_password"`- password used to authenticate (ex. `"admin"`)
- `"openstack_tenant"` - tenant name used to authenticate (ex. `"admin"`)

Example global configuration file for snap-plugin-collector-neutron plugin (exemplary file in examples/cfg/cfg.json):

```
{
  "control": {
    "cache_ttl": "5s"
  },
  "scheduler": {
    "default_deadline": "5s",
    "worker_pool_size": 5
  },
  "plugins": {
    "collector": {
      "neutron": {
        "all": {
          "openstack_auth_url": "http://localhost:5000/v2.0/",
          "openstack_user": "admin",
          "openstack_password": "admin",
          "openstack_tenant": "admin"
        }
      }
    },
    "publisher": {},
    "processor": {}
  }
}
```


### Examples
It is not suggested to set interval below 20 seconds. This may lead to overloading OpenStack Identity service with authentication requests.


Example running snap-plugin-collector-neutron plugin and writing data to a file.

In one terminal window, open the snap daemon (in this case with logging set to 1,  trust disabled and global configuration saved in cfg.json ):
```
$ $SNAP_PATH/bin/snapd -l 1 -t 0 --config cfg.json
```
In another terminal window:

Load snap-plugin-collector-neutron plugin
```
$ $SNAP_PATH/bin/snapctl plugin load snap-plugin-collector-neutron
```
Load file plugin for publishing:
```
$ $SNAP_PATH/bin/snapctl plugin load $SNAP_PATH/plugin/snap-publisher-file
```
See available metrics for your system

```
$ $SNAP_PATH/bin/snapctl metric list
```

Create task manifest file to use snap-plugin-collector-neutron plugin (exemplary file in examples/tasks/task.json):
```
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "60s"
    },
    "workflow": {
        "collect": {
            "metrics": {
		        "/intel/openstack/neutron/admin/networks_count": {},
		        "/intel/openstack/neutron/admin/subnets_count": {},
		        "/intel/openstack/neutron/admin/routers_count": {},
		        "/intel/openstack/neutron/admin/ports_count": {},
                "/intel/openstack/neutron/admin/floatingips_count": {},
                "/intel/openstack/neutron/admin/quotas_subnet": {},
		        "/intel/openstack/neutron/admin/quotas_network": {},
		        "/intel/openstack/neutron/admin/quotas_floatingip": {},
                "/intel/openstack/neutron/admin/quotas_security_group_rule": {},
		        "/intel/openstack/neutron/admin/quotas_security_group": {},
		        "/intel/openstack/neutron/admin/quotas_router": {},
                "/intel/openstack/neutron/admin/quotas_port": {}
           },
            "config": {
            },
            "process": null,
            "publish": [
                {
                    "plugin_name": "file",
                    "config": {
                        "file": "/tmp/published_neutron"
                    }
                }
            ]
        }
    }
}
```

Create task:
```
$ $SNAP_PATH/bin/snapctl task create -t examples/tasks/task.json
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [Katarzyna Zabrocka](https://github.com/katarzyna-z)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
