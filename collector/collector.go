/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package collector

import (
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	openstackintel "github.com/intelsdi-x/snap-plugin-collector-neutron/openstack"
	openstackgophercloud "github.com/rackspace/gophercloud/openstack"

	"github.com/intelsdi-x/snap-plugin-utilities/config"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/serror"
	"github.com/rackspace/gophercloud"
)

const (
	// version of neutron plugin
	version = 3

	//vendor namespace part
	vendor = "intel"

	//openstack namespace part
	openstack = "openstack"

	//pluginName namespace part
	pluginName = "neutron"

	//pluginType type of plugin
	pluginType = plugin.CollectorPluginType

	//nsLength length of namespace
	nsLength = 5

	//quotaNameIdx position of quota prefix in metric name
	quotaNameIdx = 1

	//metricNameNSPartNumber position of metric name in namespace
	metricNameNSPartNumber = 4

	//tenantNameNSPartNumber position of tenant name in namespace
	tenantNameNSPartNumber = 3

	//networksCountMetric name of metric which indicates  number of tenant networks
	networksCountMetric = "networks_count"

	//subnetsCountMetric name of metric which indicates  number of tenant subnets
	subnetsCountMetric = "subnets_count"

	//routersCountMetric name of metric which indicates  number of tenant routers
	routersCountMetric = "routers_count"

	//portsCountMetric name of metric which indicates  number of tenant ports
	portsCountMetric = "ports_count"

	//floatingipsCountMetric name of metric which indicates  number of tenant  floating IPs
	floatingipsCountMetric = "floatingips_count"

	//quotas prefix for quota metrics
	quotas = "quotas_"

	//cfgUrl name of configuration variable for url  for OpenStack Identity endpoint
	cfgURL = "openstack_auth_url"

	//cfgUser user name used to authenticate
	cfgUser = "openstack_user"

	//cfgPassword password used to authenticate
	cfgPassword = "openstack_password"

	//cfgTenant tenant name used to authenticate
	cfgTenant = "openstack_tenant"
)

//neutronConstMetrics slice of constant metric names
var neutronConstMetrics = []string{
	networksCountMetric,
	subnetsCountMetric,
	routersCountMetric,
	portsCountMetric,
	floatingipsCountMetric,
}

//neutronInfoFields contains information (description and unit) about metrics
var neutronInfoFields = map[string]infoFields{
	networksCountMetric: infoFields{
		description: "number of tenant networks",
		unit:        "",
	},
	subnetsCountMetric: infoFields{
		description: "number of tenant subnets",
		unit:        "",
	},
	routersCountMetric: infoFields{
		description: "number of tenant routers",
		unit:        "",
	},
	portsCountMetric: infoFields{
		description: "number of tenant ports",
		unit:        "",
	},
	floatingipsCountMetric: infoFields{
		description: "number of tenant floating IPs",
		unit:        "",
	},
	quotas + "floatingip": infoFields{
		description: "number of floating IP addresses allowed for a tenant ( -1 means no limit)",
		unit:        "",
	},
	quotas + "ikepolicy": infoFields{
		description: "number of IKE policies allowed for a tenant",
		unit:        "",
	},
	quotas + "ipsec_site_connection": infoFields{
		description: "number of IPSec connections allowed for a tenant",
		unit:        "",
	},
	quotas + "ipsecpolicy": infoFields{
		description: "number of IPSec policies allowed for a tenant",
		unit:        "",
	},
	quotas + "network": infoFields{
		description: "number of networks allowed for a tenant",
		unit:        "",
	},
	quotas + "port": infoFields{
		description: "number of ports allowed for a tenant",
		unit:        "",
	},
	quotas + "rbac_policy": infoFields{
		description: "number of role-based access control (RBAC) policies for a tenant",
		unit:        "",
	},
	quotas + "router": infoFields{
		description: "number of routers allowed for a tenant",
		unit:        "",
	},
	quotas + "security_group": infoFields{
		description: "number of security groups allowed for a tenant",
		unit:        "",
	},
	quotas + "security_group_rule": infoFields{
		description: "number of security group rules allowed for a tenant",
		unit:        "",
	},
	quotas + "subnet": infoFields{
		description: "number of subnets allowed for a tenant",
		unit:        "",
	},
	quotas + "subnetpool": infoFields{
		description: "number of subnet pools allowed for a tenant",
		unit:        "",
	},
}

//Collector neutron plugin struct
type Collector struct {
	provider *gophercloud.ProviderClient
}

//Meta returns meta data for plugin
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		pluginName,
		version,
		pluginType,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType},
		plugin.RoutingStrategy(plugin.StickyRouting),
	)
}

// New creates initialized instance of Glance collector
func New() *Collector {
	return &Collector{}
}

// GetMetricTypes returns list of available metric types
// It returns error in case retrieval was not successful
func (c *Collector) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	mts := []plugin.MetricType{}
	items, err := config.GetConfigItems(cfg, cfgURL, cfgUser, cfgPassword, cfgTenant)
	if err != nil {
		return nil, err
	}
	domain_name := ""
	domain_id := ""

	endpoint := items[cfgURL].(string)
	user := items[cfgUser].(string)
	password := items[cfgPassword].(string)
	tenant := items[cfgTenant].(string)
	dom_name, _ := config.GetConfigItem(cfg, "domain_name")
	dom_id, _ := config.GetConfigItem(cfg, "domain_id")
	if dom_name != nil {
		domain_name = dom_name.(string)
	}
	if dom_id != nil {
		domain_id = dom_id.(string)
	}
	if c.provider == nil {
		provider, serr := openstackintel.Authenticate(endpoint, user, password, tenant, domain_name, domain_id)
		if serr != nil {
			log.WithFields(serr.Fields()).Warn(serr.Error())
			return nil, serr
		}
		c.provider = provider
	}

	// Retrieve list of all available tenants for provided endpoint, user and password

	identityClient := openstackgophercloud.NewIdentityV2(c.provider)
	allTenants, serr := openstackintel.GetAllTenants(identityClient)
	if serr != nil {
		log.WithFields(serr.Fields()).Warn(serr.Error())
		return nil, serr
	}

	networkClient, err := openstackgophercloud.NewNetworkV2(c.provider, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, err
	}

	// Generate available namespace from tenants (user counts per tenant)
	for _, tenant := range allTenants {

		q, serr := openstackintel.GetQuotasForTenant(networkClient, tenant.ID)
		if serr != nil {
			log.WithFields(serr.Fields()).Warn(serr.Error())
			return nil, serr
		}

		for k := range q {
			info := getInfoFields(quotas + k)
			mts = append(mts, plugin.MetricType{

				Namespace_:   core.NewNamespace(vendor, openstack, pluginName, tenant.Name, quotas+k),
				Config_:      cfg.ConfigDataNode,
				Description_: info.description,
				Unit_:        info.unit,
			})
		}

		for m := range neutronConstMetrics {
			info := getInfoFields(neutronConstMetrics[m])
			mts = append(mts, plugin.MetricType{
				Namespace_:   core.NewNamespace(vendor, openstack, pluginName, tenant.Name, neutronConstMetrics[m]),
				Config_:      cfg.ConfigDataNode,
				Description_: info.description,
				Unit_:        info.unit,
			})
		}
	}
	return mts, nil
}

// CollectMetrics returns list of requested metric values
// It returns error in case retrieval was not successful
func (c *Collector) CollectMetrics(metricTypes []plugin.MetricType) ([]plugin.MetricType, error) {
	items, err := config.GetConfigItems(metricTypes[0], cfgURL, cfgUser, cfgPassword, cfgTenant)
	if err != nil {
		return nil, err
	}

	domain_name := ""
	domain_id := ""

	endpoint := items[cfgURL].(string)
	user := items[cfgUser].(string)
	password := items[cfgPassword].(string)
	tenant := items[cfgTenant].(string)
	dom_name, _ := config.GetConfigItem(metricTypes[0], "domain_name")
	dom_id, _ := config.GetConfigItem(metricTypes[0], "domain_id")
	if dom_name != nil {
		domain_name = dom_name.(string)
	}
	if dom_id != nil {
		domain_id = dom_id.(string)
	}

	if c.provider == nil {
		provider, serr := openstackintel.Authenticate(endpoint, user, password, tenant, domain_name, domain_id)
		if serr != nil {
			log.WithFields(serr.Fields()).Warn(serr.Error())
			return nil, serr
		}
		c.provider = provider
	}

	identityClient := openstackgophercloud.NewIdentityV2(c.provider)
	tenantList, serr := openstackintel.GetAllTenants(identityClient)
	if serr != nil {
		log.WithFields(serr.Fields()).Warn(serr.Error())
		return nil, serr
	}

	networkClient, err := openstackgophercloud.NewNetworkV2(c.provider, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, err
	}

	var done sync.WaitGroup
	done.Add(6)

	var tenantNetworks map[string]int64
	go func() {
		var serr serror.SnapError
		tenantNetworks, serr = openstackintel.GetNetworkCountPerTenant(networkClient, tenantList)
		if serr != nil {
			log.WithFields(serr.Fields()).Warn(serr.Error())
			panic(serr)
		}
		done.Done()
	}()

	var tenantSubnets map[string]int64
	go func() {
		var serr serror.SnapError
		tenantSubnets, serr = openstackintel.GetSubnetsCountPerTenant(networkClient, tenantList)
		if serr != nil {
			log.WithFields(serr.Fields()).Warn(serr.Error())
			panic(serr)
		}
		done.Done()
	}()

	var tenantRouters map[string]int64
	go func() {
		var serr serror.SnapError
		tenantRouters, serr = openstackintel.GetRoutersCountPerTenant(networkClient, tenantList)
		if serr != nil {
			log.WithFields(serr.Fields()).Warn(serr.Error())
			panic(serr)
		}
		done.Done()
	}()

	var tenantPorts map[string]int64
	go func() {
		var serr serror.SnapError
		tenantPorts, serr = openstackintel.GetPortsCountPerTenant(networkClient, tenantList)
		if serr != nil {
			log.WithFields(serr.Fields()).Warn(serr.Error())
			panic(serr)
		}
		done.Done()
	}()

	var tenantFloatingips map[string]int64
	go func() {
		var serr serror.SnapError
		tenantFloatingips, serr = openstackintel.GetFloatingIPsCountPerTenant(networkClient, tenantList)
		if serr != nil {
			log.WithFields(serr.Fields()).Warn(serr.Error())
			panic(serr)
		}
		done.Done()
	}()

	var tenantQuotasList map[string]map[string]int64
	go func() {
		var serr serror.SnapError
		tenantQuotasList, serr = openstackintel.GetQuotasPerTenant(networkClient, tenantList)
		if serr != nil {
			log.WithFields(serr.Fields()).Warn(serr.Error())
			panic(err)
		}
		done.Done()
	}()

	done.Wait()

	metrics := []plugin.MetricType{}
	for _, metricType := range metricTypes {

		namespace := metricType.Namespace()
		if len(namespace) != nsLength {
			f := map[string]interface{}{"namespace": metricType.Namespace().String()}
			serr := serror.New(fmt.Errorf("Incorrect namespace length"), f)
			log.WithFields(serr.Fields()).Warn(serr.String())
			continue
		}

		metric := plugin.MetricType{
			Timestamp_: time.Now(),
			Namespace_: namespace,
		}

		tenantName := namespace[tenantNameNSPartNumber].Value
		var val int64
		var ok bool
		switch namespace[metricNameNSPartNumber].Value {
		case networksCountMetric:
			val, ok = tenantNetworks[tenantName]
		case subnetsCountMetric:
			val, ok = tenantSubnets[tenantName]
		case routersCountMetric:
			val, ok = tenantRouters[tenantName]
		case portsCountMetric:
			val, ok = tenantPorts[tenantName]
		case floatingipsCountMetric:
			val, ok = tenantFloatingips[tenantName]
		default:

			if !strings.HasPrefix(namespace[metricNameNSPartNumber].Value, quotas) {
				f := map[string]interface{}{"namespace": "/" + metricType.Namespace().String()}
				serr := serror.New(fmt.Errorf("Incorrect namespace, prefix '%s' is desired", quotas), f)
				log.WithFields(serr.Fields()).Warn(serr.String())
				continue
			}

			quotaName := namespace[metricNameNSPartNumber].Value[len(quotas):]
			val, ok = tenantQuotasList[tenantName][quotaName]
		}

		if !ok {
			f := map[string]interface{}{"namespace": metricType.Namespace().String(), "tenantName": tenantName}
			serr := serror.New(fmt.Errorf("Incorrect namespace, metric with specified namespace does not exist"), f)
			log.WithFields(serr.Fields()).Warn(serr.String())
			continue
		}
		metric.Data_ = val
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

// GetConfigPolicy returns config policy
// It returns error in case retrieval was not successful
func (c *Collector) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cp := cpolicy.New()
	config := cpolicy.NewPolicyNode()

	r1, err := cpolicy.NewStringRule(cfgURL, true)
	if err != nil {
		return cp, err
	}
	r1.Description = "URL for OpenStack Identity endpoint"
	config.Add(r1)

	r2, err := cpolicy.NewStringRule(cfgUser, true)
	if err != nil {
		return cp, err
	}
	r2.Description = "user name used to authenticate"
	config.Add(r2)

	r3, err := cpolicy.NewStringRule(cfgPassword, true)
	if err != nil {
		return cp, err
	}
	r3.Description = "password used to authenticate"
	config.Add(r3)

	r4, err := cpolicy.NewStringRule(cfgTenant, true)
	if err != nil {
		return cp, err
	}
	r4.Description = " tenant name used to authenticate"
	config.Add(r4)

	cp.Add([]string{""}, config)
	return cp, nil
}

func getInfoFields(metric string) infoFields {
	info, ok := neutronInfoFields[metric]
	if !ok {
		info = infoFields{description: "", unit: ""}
	}
	return info
}

type infoFields struct {
	description string
	unit        string
}
