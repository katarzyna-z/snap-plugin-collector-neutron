# Snap plugin collector - neutron

Snap plugin for collecting metrics from OpenStack Neutron.

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
You can get the pre-built binaries for your OS and architecture from plugin's [Github Releases](releases) page.

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
This builds the plugin in `./build`

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Create Global Config, see description in [Snap's Global Config] (https://github.com/intelsdi-x/snap-plugin-collector-neutron/blob/master/README.md#snaps-global-config).
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-collector-neutron/blob/master/README.md#examples).

#### Suggestions
* It is not recommended to set interval for task less than 20 seconds. This may lead to overloading Neutron API with requests.

## Documentation

### Collected Metrics
List of collected metrics is described in [METRICS.md](https://github.com/intelsdi-x/snap-plugin-collector-neutron/blob/master/METRICS.md).

### snap's Global Config
Global configuration files are described in [Snap's documentation](https://github.com/intelsdi-x/snap/blob/master/docs/SNAPTELD_CONFIGURATION.md). You have to add section "neutron" in "collector" section and then specify following options:
- `"openstack_auth_url"` - URL for OpenStack Identity endpoint (ex. `"http://127.0.0.1:5000/v2.0/"`)
- `"openstack_user"` - user name used to authenticate (ex. `"admin"`)
- `"openstack_password"`- password used to authenticate (ex. `"admin"`)
- `"openstack_tenant"` - tenant name used to authenticate (ex. `"admin"`)

Example global configuration file for snap-plugin-collector-neutron plugin (exemplary file in [examples/cfg/] (https://github.com/intelsdi-x/snap-plugin-collector-neutron/blob/master/examples/cfg/):

```
{
  "control": {
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
}
```


### Examples
Example running snap-plugin-collector-neutron plugin and writing data to a file using [snap-plugin-publisher-file](https://github.com/intelsdi-x/snap-plugin-publisher-file).

Create Global Config, see example in [examples/cfg/] (https://github.com/intelsdi-x/snap-plugin-collector-neutron/blob/master/examples/cfg/).

Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started),
in one terminal window, run `snapteld` (in this case with logging set to 1, trust disabled and global configuration saved in cfg.json):
```
$ $SNAP_PATH/bin/snapteld -l 1 -t 0 --config cfg.json
```

In another terminal window:

Download and load Snap plugins:
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-neutron/latest/linux/x86_64/snap-plugin-collector-neutron
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
$ snaptel plugin load snap-plugin-collector-neutron
$ snaptel plugin load snap-plugin-publisher-file
```

See available metrics for your system

```
$ snaptel metric list
```

Create a [Task Manifest](https://github.com/intelsdi-x/snap/blob/master/docs/TASKS.md) file to use snap-plugin-collector-neutron plugin (exemplary file in [examples/tasks/] (https://github.com/intelsdi-x/snap-plugin-collector-neutron/blob/master/examples/tasks/)):
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
$ snaptel task create -t examples/tasks/task.json
```

And watch the metrics populate:
```
$ snaptel task watch <task_id>
```

Stop previously created task:
```
$ snaptel task stop <task_id>
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. The full project is at http://github.com/intelsdi-x/snap.
To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[Snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [Katarzyna Zabrocka](https://github.com/katarzyna-z)
