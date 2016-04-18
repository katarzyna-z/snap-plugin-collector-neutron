# snap plugin collector - neutron

## Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Data Type | Description
----------------|:-------------------------|:-----------------------
/intel/openstack/neutron/\<tenant_name\>/networks_count | int64 | number of tenant networks
/intel/openstack/neutron/\<tenant_name\>/subnets_count  | int64 | number of tenant subnets
/intel/openstack/neutron/\<tenant_name\>/routers_count | int64 | number of tenant routers
/intel/openstack/neutron/\<tenant_name\>/ports_count | int64 | number of tenant ports
/intel/openstack/neutron/\<tenant_name\>/floatingips_count | int64 | number of tenant floating IPs
/intel/openstack/neutron/\<tenant_name\>/quotas_floatingip | int64 | number of floating IP addresses allowed for a tenant ( -1 means no limit)
/intel/openstack/neutron/\<tenant_name\>/quotas_ikepolicy | int64 | number of IKE policies allowed for a tenant
/intel/openstack/neutron/\<tenant_name\>/quotas_ipsec_site_connection | int64 | number of  IPSec connections allowed for a tenant
/intel/openstack/neutron/\<tenant_name\>/quotas_ipsecpolicy | int64 | number of IPSec policies  allowed for a tenant
/intel/openstack/neutron/\<tenant_name\>/quotas_network | int64 | number of networks allowed for a tenant
/intel/openstack/neutron/\<tenant_name\>/quotas_port | int64 |  number of ports allowed for a tenant
/intel/openstack/neutron/\<tenant_name\>/quotas_rbac_policy | int64 | number of role-based access control (RBAC) policies for a tenant
/intel/openstack/neutron/\<tenant_name\>/quotas_router | int64 | number of routers allowed for a tenant
/intel/openstack/neutron/\<tenant_name\>/quotas_security_group | int64 | number of security groups allowed for a tenant
/intel/openstack/neutron/\<tenant_name\>/quotas_security_group_rule | int64 | number of security group rules allowed for a tenant
/intel/openstack/neutron/\<tenant_name\>/quotas_subnet | int64 | number of subnets allowed for a tenant
/intel/openstack/neutron/\<tenant_name\>/quotas_subnetpool | int64 | number of subnet pools allowed for a tenant