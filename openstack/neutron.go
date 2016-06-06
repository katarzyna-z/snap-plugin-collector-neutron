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

package openstack

import (
	"fmt"

	"github.com/intelsdi-x/snap-plugin-collector-neutron/openstack/tenantquotas"
	"github.com/intelsdi-x/snap-plugin-collector-neutron/types"
	"github.com/intelsdi-x/snap/core/serror"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/identity/v2/tenants"
	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/openstack/networking/v2/ports"
	"github.com/rackspace/gophercloud/openstack/networking/v2/subnets"
)

const (
	quotaPath = "quota"
)

// GetAllTenants is used to retrieve list of available tenants
func GetAllTenants(client *gophercloud.ServiceClient) ([]types.Tenant, serror.SnapError) {
	tnts := []types.Tenant{}

	pager := tenants.List(client, nil)
	page, err := pager.AllPages()
	if err != nil {
		return tnts, serror.New(err)
	}

	tenantList, err := tenants.ExtractTenants(page)
	if err != nil {
		return tnts, serror.New(err)
	}

	for _, t := range tenantList {
		tnts = append(tnts, types.Tenant{Name: t.Name, ID: t.ID})
	}
	return tnts, nil
}

// GetNetworkCountPerTenant is used to retrieve number of networks per tenant
func GetNetworkCountPerTenant(client *gophercloud.ServiceClient, tenantList []types.Tenant) (map[string]int64, serror.SnapError) {
	tenantNetworksCount := map[string]int64{}

	pager := networks.List(client, nil)
	page, err := pager.AllPages()
	if err != nil {
		return tenantNetworksCount, serror.New(err)
	}

	networkList, err := networks.ExtractNetworks(page)
	if err != nil {
		return tenantNetworksCount, serror.New(err)
	}

	for _, tnt := range tenantList {
		if _, ok := tenantNetworksCount[tnt.Name]; !ok {
			tenantNetworksCount[tnt.Name] = 0
		}

		for _, net := range networkList {
			if tnt.ID == net.TenantID {
				tenantNetworksCount[tnt.Name]++
			}
		}
	}
	return tenantNetworksCount, nil
}

// GetSubnetsCountPerTenant is used to retrieve number of subnets per tenant
func GetSubnetsCountPerTenant(client *gophercloud.ServiceClient, tenantList []types.Tenant) (map[string]int64, serror.SnapError) {
	tenantSubnetsCount := map[string]int64{}

	pager := subnets.List(client, nil)
	page, err := pager.AllPages()
	if err != nil {
		return tenantSubnetsCount, serror.New(err)
	}

	subnetList, err := subnets.ExtractSubnets(page)
	if err != nil {
		return tenantSubnetsCount, serror.New(err)
	}

	for _, tnt := range tenantList {
		if _, ok := tenantSubnetsCount[tnt.Name]; !ok {
			tenantSubnetsCount[tnt.Name] = 0
		}

		for _, subnet := range subnetList {
			if tnt.ID == subnet.TenantID {
				tenantSubnetsCount[tnt.Name]++
			}
		}
	}
	return tenantSubnetsCount, nil
}

//GetRoutersCountPerTenant  is used to retrieve number of routers per tenant
func GetRoutersCountPerTenant(client *gophercloud.ServiceClient, tenantList []types.Tenant) (map[string]int64, serror.SnapError) {
	tenantRoutersCount := map[string]int64{}

	pager := routers.List(client, routers.ListOpts{})
	page, err := pager.AllPages()
	if err != nil {
		return tenantRoutersCount, serror.New(err)
	}

	routerList, err := routers.ExtractRouters(page)
	if err != nil {
		return tenantRoutersCount, serror.New(err)
	}

	for _, tnt := range tenantList {
		if _, ok := tenantRoutersCount[tnt.Name]; !ok {
			tenantRoutersCount[tnt.Name] = 0
		}

		for _, router := range routerList {
			if tnt.ID == router.TenantID {
				tenantRoutersCount[tnt.Name]++
			}
		}
	}
	return tenantRoutersCount, nil
}

//GetPortsCountPerTenant  is used to retrieve number of ports per tenant
func GetPortsCountPerTenant(client *gophercloud.ServiceClient, tenantList []types.Tenant) (map[string]int64, serror.SnapError) {
	tenantPortsCount := map[string]int64{}

	pager := ports.List(client, ports.ListOpts{})
	page, err := pager.AllPages()
	if err != nil {
		return tenantPortsCount, serror.New(err)
	}

	portList, err := ports.ExtractPorts(page)
	if err != nil {
		return tenantPortsCount, serror.New(err)
	}

	for _, tnt := range tenantList {
		if _, ok := tenantPortsCount[tnt.Name]; !ok {
			tenantPortsCount[tnt.Name] = 0
		}

		for _, port := range portList {
			if tnt.ID == port.TenantID {
				tenantPortsCount[tnt.Name]++
			}
		}
	}
	return tenantPortsCount, nil
}

//GetFloatingIPsCountPerTenant is used to retrieve number of floating IPs per tenant
func GetFloatingIPsCountPerTenant(client *gophercloud.ServiceClient, tenantList []types.Tenant) (map[string]int64, serror.SnapError) {
	tenantFloatingipsCount := map[string]int64{}

	opts := floatingips.ListOpts{}
	pager := floatingips.List(client, opts)

	page, err := pager.AllPages()
	if err != nil {
		return tenantFloatingipsCount, serror.New(err)
	}

	floatingipList, err := floatingips.ExtractFloatingIPs(page)
	if err != nil {
		return tenantFloatingipsCount, serror.New(err)
	}

	for _, tnt := range tenantList {
		if _, ok := tenantFloatingipsCount[tnt.Name]; !ok {
			tenantFloatingipsCount[tnt.Name] = 0
		}

		for _, floatip := range floatingipList {
			if tnt.ID == floatip.TenantID {
				tenantFloatingipsCount[tnt.Name]++
			}
		}
	}
	return tenantFloatingipsCount, nil
}

//GetQuotasPerTenant is used to retrieve quotas per tenants
func GetQuotasPerTenant(client *gophercloud.ServiceClient, tenantList []types.Tenant) (map[string]map[string]int64, serror.SnapError) {
	var tenantQuotas map[string]map[string]int64
	tenantQuotas = make(map[string]map[string]int64)

	for _, tnt := range tenantList {
		quotasMap, serr := GetQuotasForTenant(client, tnt.ID)
		if serr != nil {
			return tenantQuotas, serr
		}
		tenantQuotas[tnt.Name] = quotasMap
	}
	return tenantQuotas, nil
}

//GetQuotasForTenant is used to retrieve quotas for specified tenant
func GetQuotasForTenant(client *gophercloud.ServiceClient, tenantID string) (map[string]int64, serror.SnapError) {
	quotas, err := tenantquotas.Get(client, tenantID).Extract()
	if err != nil {
		return nil, serror.New(err)
	}

	quotasMap, ok := quotas[quotaPath]
	if !ok {
		f := map[string]interface{}{"quotas": quotas, "tenantName": tenantID}
		return nil, serror.New(fmt.Errorf("GetQuotasForTenant: inocorrect response format"), f)
	}
	return quotasMap, nil
}
