# snap plugin collector - neutron

snap plugin for collecting metrics from OpenStack Neutron module. 

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
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
 * OpenStack deployment available
 
### Operating systems
All OSs currently supported by snap:
* Linux/amd64

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

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Create Global Config, see description in [snap's Global Config] (https://github.com/intelsdi-x/snap-plugin-collector-neutron/blob/master/README.md#snaps-global-config).
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-collector-neutron/blob/master/README.md#examples).

#### Suggestions
* It is not recommended to set interval for task less than 20 seconds. This may lead to overloading Neutron API with requests.


## Documentation

### Collected Metrics
List of collected metrics is described in [METRICS.md](https://github.com/intelsdi-x/snap-plugin-collector-neutron/blob/master/METRICS.md).

### snap's Global Config
Global configuration files are described in [snap's documentation](https://github.com/intelsdi-x/snap/blob/master/docs/SNAPD_CONFIGURATION.md). You have to add section "neutron" in "collector" section and then specify following options:
- `"openstack_auth_url"` - URL for OpenStack Identity endpoint (ex. `"http://127.0.0.1:5000/v2.0/"`)
- `"openstack_user"` - user name used to authenticate (ex. `"admin"`)
- `"openstack_password"`- password used to authenticate (ex. `"admin"`)
- `"openstack_tenant"` - tenant name used to authenticate (ex. `"admin"`)

Example global configuration file for snap-plugin-collector-neutron plugin (exemplary file in [examples/cfg/] (https://github.com/intelsdi-x/snap-plugin-collector-neutron/blob/master/examples/cfg/):

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
Example running snap-plugin-collector-neutron plugin and writing data to a file.

Make sure that your `$SNAP_PATH` is set, if not:
```
$ export SNAP_PATH=<snapDirectoryPath>/build
```
Other paths to files should be set according to your configuration, using a file you should indicate where it is located.

Create Global Config, see example in [examples/cfg/] (https://github.com/intelsdi-x/snap-plugin-collector-neutron/blob/master/examples/cfg/).

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

Create a task manifest file to use snap-plugin-collector-neutron plugin (exemplary file in [examples/tasks/] (https://github.com/intelsdi-x/snap-plugin-collector-neutron/blob/master/examples/tasks/)):
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

Create a task:
```
$ $SNAP_PATH/bin/snapctl task create -t examples/tasks/task.json
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. The full project is at http://github.com/intelsdi-x/snap.
To reach out on other use cases, visit:
* [snap Gitter channel](https://gitter.im/intelsdi-x/snap)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [Katarzyna Zabrocka](https://github.com/katarzyna-z)
