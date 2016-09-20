// +build small

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
	"net/http"
	"testing"

	openstackgophercloud "github.com/rackspace/gophercloud/openstack"
	th "github.com/rackspace/gophercloud/testhelper"

	"github.com/intelsdi-x/snap-plugin-collector-neutron/types"
	"github.com/rackspace/gophercloud"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	Token                  string
	NetworkServiceEndpoint string
}

func (s *TestSuite) SetupSuite() {
	th.SetupHTTP()
	registerRoot()
	registerAuthentication(s)
	registerTenants(s)
	registerNetworks(s)
	registerSubnets(s)
	registerRouters(s)
	registerPorts(s)
	registerFloatingIPs(s)
	registerQuotas(s)
}

func (suite *TestSuite) TearDownSuite() {
	th.TeardownHTTP()
}

func TestRunSuite(t *testing.T) {
	cinderTestSuite := new(TestSuite)
	suite.Run(t, cinderTestSuite)
}

func (s *TestSuite) TestGetAllTenants() {
	Convey("Given list of OpenStack tenants is requested", s.T(), func() {

		Convey("When authentication is required", func() {
			provider, err := Authenticate(th.Endpoint(), "me", "secret", "admin")
			th.AssertNoErr(s.T(), err)
			th.CheckEquals(s.T(), s.Token, provider.TokenID)

			identityClient := openstackgophercloud.NewIdentityV2(provider)

			Convey("and GetAllTenants called", func() {

				tenantList, serr := GetAllTenants(identityClient)

				Convey("Then number of tenants is returned", func() {
					So(len(tenantList), ShouldEqual, 2)
				})

				Convey("and no error reported", func() {
					So(serr, ShouldBeNil)
				})
			})
		})
	})
}

func (s *TestSuite) TestGetNetworkCountPerTenant() {
	Convey("Number of OpenStack networks per tenant is requested", s.T(), func() {

		Convey("When authentication is required", func() {
			provider, serr := Authenticate(th.Endpoint(), "me", "secret", "admin")
			th.AssertNoErr(s.T(), serr)
			th.CheckEquals(s.T(), s.Token, provider.TokenID)
			identityClient := openstackgophercloud.NewIdentityV2(provider)
			tenantList, _ := GetAllTenants(identityClient)

			networkClient, err := openstackgophercloud.NewNetworkV2(provider, gophercloud.EndpointOpts{})
			So(err, ShouldBeNil)

			Convey("and GetNetworkCountPerTenant called", func() {

				networkList, serr := GetNetworkCountPerTenant(networkClient, tenantList)

				Convey("Then number of networks is returned", func() {
					So(len(networkList), ShouldEqual, 2)
					So(networkList["admin"], ShouldEqual, 2)
					So(networkList["demo"], ShouldEqual, 1)
					So(networkList["test"], ShouldEqual, 0)
				})

				Convey("and no error reported", func() {
					So(serr, ShouldBeNil)
				})

			})
		})
	})
}

func (s *TestSuite) TestGetSubnetsCountPerTenant() {
	Convey("Number of OpenStack subnets per tenant is requested", s.T(), func() {

		Convey("When authentication is required", func() {
			provider, serr := Authenticate(th.Endpoint(), "me", "secret", "admin")
			th.AssertNoErr(s.T(), serr)
			th.CheckEquals(s.T(), s.Token, provider.TokenID)
			identityClient := openstackgophercloud.NewIdentityV2(provider)
			tenantList, _ := GetAllTenants(identityClient)

			networkClient, err := openstackgophercloud.NewNetworkV2(provider, gophercloud.EndpointOpts{})
			So(err, ShouldBeNil)

			Convey("and GetSubnetsCountPerTenant called", func() {

				subnetList, serr := GetSubnetsCountPerTenant(networkClient, tenantList)

				Convey("Then number of subnets is returned", func() {
					So(len(subnetList), ShouldEqual, 2)
					So(subnetList["admin"], ShouldEqual, 3)
					So(subnetList["demo"], ShouldEqual, 0)
					So(subnetList["test"], ShouldEqual, 0)
				})

				Convey("and no error reported", func() {
					So(serr, ShouldBeNil)
				})
			})
		})
	})
}

func (s *TestSuite) TestGetRoutersCountPerTenant() {
	Convey("Number of OpenStack routers per tenant is requested", s.T(), func() {

		Convey("When authentication is required", func() {
			provider, serr := Authenticate(th.Endpoint(), "me", "secret", "admin")
			th.AssertNoErr(s.T(), serr)
			th.CheckEquals(s.T(), s.Token, provider.TokenID)
			identityClient := openstackgophercloud.NewIdentityV2(provider)
			tenantList, _ := GetAllTenants(identityClient)

			networkClient, err := openstackgophercloud.NewNetworkV2(provider, gophercloud.EndpointOpts{})
			So(err, ShouldBeNil)

			Convey("and GetRoutersCountPerTenant called", func() {

				routerList, serr := GetRoutersCountPerTenant(networkClient, tenantList)

				Convey("Then number of routers is returned", func() {
					So(len(routerList), ShouldEqual, 2)
					So(routerList["admin"], ShouldEqual, 4)
					So(routerList["demo"], ShouldEqual, 0)
					So(routerList["test"], ShouldEqual, 0)
				})

				Convey("and no error reported", func() {
					So(serr, ShouldBeNil)
				})
			})
		})
	})
}

func (s *TestSuite) TestGetPortsCountPerTenant() {
	Convey("Number of OpenStack ports per tenant is requested", s.T(), func() {

		Convey("When authentication is required", func() {
			provider, serr := Authenticate(th.Endpoint(), "me", "secret", "admin")
			th.AssertNoErr(s.T(), serr)
			th.CheckEquals(s.T(), s.Token, provider.TokenID)
			identityClient := openstackgophercloud.NewIdentityV2(provider)
			tenantList, _ := GetAllTenants(identityClient)

			networkClient, err := openstackgophercloud.NewNetworkV2(provider, gophercloud.EndpointOpts{})
			So(err, ShouldBeNil)

			Convey("and GetPortsCountPerTenant called", func() {

				portList, serr := GetPortsCountPerTenant(networkClient, tenantList)

				Convey("Then number of tenants is returned", func() {
					So(len(portList), ShouldEqual, 2)
					So(portList["admin"], ShouldEqual, 3)
					So(portList["demo"], ShouldEqual, 0)
					So(portList["test"], ShouldEqual, 0)
				})

				Convey("and no error reported", func() {
					So(serr, ShouldBeNil)
				})
			})
		})
	})
}

func (s *TestSuite) TestGetFloatingIPsCountPerTenant() {
	Convey("Number of OpenStack floating IPs per tenant is requested", s.T(), func() {

		Convey("When authentication is required", func() {
			provider, serr := Authenticate(th.Endpoint(), "me", "secret", "admin")
			th.AssertNoErr(s.T(), serr)
			th.CheckEquals(s.T(), s.Token, provider.TokenID)
			identityClient := openstackgophercloud.NewIdentityV2(provider)
			tenantList, _ := GetAllTenants(identityClient)

			networkClient, err := openstackgophercloud.NewNetworkV2(provider, gophercloud.EndpointOpts{})
			So(err, ShouldBeNil)

			Convey("and GetFloatingIPsCountPerTenant called", func() {

				floatingipList, serr := GetFloatingIPsCountPerTenant(networkClient, tenantList)

				Convey("Then number of floating IPs for tenants is returned", func() {
					So(len(floatingipList), ShouldEqual, 2)
					So(floatingipList["admin"], ShouldEqual, 2)
					So(floatingipList["demo"], ShouldEqual, 0)
					So(floatingipList["test"], ShouldEqual, 0)
				})
				Convey("and no error reported", func() {
					So(serr, ShouldBeNil)
				})
			})
		})
	})
}

func (s *TestSuite) TestGetQuotasPerTenant() {
	Convey("Given list of OpenStack quotas per tenant is requested", s.T(), func() {

		Convey("When authentication is required", func() {
			provider, serr := Authenticate(th.Endpoint(), "me", "secret", "admin")
			th.AssertNoErr(s.T(), serr)
			th.CheckEquals(s.T(), s.Token, provider.TokenID)
			identityClient := openstackgophercloud.NewIdentityV2(provider)
			tenantList, _ := GetAllTenants(identityClient)

			networkClient, err := openstackgophercloud.NewNetworkV2(provider, gophercloud.EndpointOpts{})
			So(err, ShouldBeNil)

			Convey("and GetQuotasPerTenant called", func() {

				quotaList, serr := GetQuotasPerTenant(networkClient, tenantList)

				Convey("Then list of quotas per tenants is returned", func() {
					So(len(quotaList), ShouldEqual, 2)
					So(len(quotaList["admin"]), ShouldEqual, 9)
					So(len(quotaList["demo"]), ShouldEqual, 9)
					So(len(quotaList["test"]), ShouldEqual, 0)
				})

				Convey("and no error reported", func() {
					So(serr, ShouldBeNil)
				})
			})

			Convey("and GetQuotasPerTenant called for non-existing tenant", func() {
				nonexistingTenantList := []types.Tenant{types.Tenant{ID: "333333", Name: "test"}}
				quotaList, serr := GetQuotasPerTenant(networkClient, nonexistingTenantList)

				Convey("Then list of quotas per tenants is returned", func() {
					So(len(quotaList), ShouldEqual, 0)
				})

				Convey("and no error reported", func() {
					So(serr, ShouldNotBeNil)
				})
			})

		})
	})
}
func (s *TestSuite) TestGetQuotasForTenant() {
	Convey("Given list of OpenStack quotas for particular is requested", s.T(), func() {

		Convey("When authentication is required", func() {
			provider, serr := Authenticate(th.Endpoint(), "me", "secret", "admin")
			th.AssertNoErr(s.T(), serr)
			th.CheckEquals(s.T(), s.Token, provider.TokenID)

			networkClient, err := openstackgophercloud.NewNetworkV2(provider, gophercloud.EndpointOpts{})
			So(err, ShouldBeNil)

			Convey("and GetQuotasForTenant called", func() {

				quotaList, err := GetQuotasForTenant(networkClient, "222222")

				Convey("Then number of quota is returned", func() {
					So(len(quotaList), ShouldEqual, 9)
				})

				Convey("and no error reported", func() {
					So(err, ShouldBeNil)
				})
			})

			Convey("and GetQuotasForTenant called for non-existing tenant", func() {
				quotaList, serr := GetQuotasForTenant(networkClient, "333333")

				Convey("Then number of quota is returned", func() {
					So(len(quotaList), ShouldEqual, 0)
				})

				Convey("and no error reported", func() {
					So(serr, ShouldNotBeNil)
				})
			})
		})
	})
}

func registerRoot() {
	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
				{
					"versions": {
						"values": [
							{
								"status": "experimental",
								"id": "v3.0",
								"links": [
									{ "href": "%s", "rel": "self" }
								]
							},
							{
								"status": "stable",
								"id": "v2.0",
								"links": [
									{ "href": "%s", "rel": "self" }
								]
							}
						]
					}
				}
				`, th.Endpoint()+"v3/", th.Endpoint()+"v2.0/")
	})
}

func registerAuthentication(s *TestSuite) {
	s.Token = "2ed210f132564f21b178afb197ee99e3"
	s.NetworkServiceEndpoint = th.Endpoint()
	th.Mux.HandleFunc("/v2.0/tokens", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
				{
					"access": {
						"metadata": {
							"is_admin": 0,
							"roles": [
								"3083d61996d648ca88d6ff420542f324"
							]
						},
						"serviceCatalog": [
						{
								"endpoints": [
									{
										"adminURL": "%s",
										"id": "3ffe125aa59547029ed774c10b932349",
										"internalURL": "%s",
										"publicURL": "%s",
										"region": "RegionOne"
									}
								],
								"endpoints_links": [],
								"name": "neutron",
								"type": "network"
							}
						],
						"token": {
							"expires": "2016-02-21T14:28:30Z",
							"id": "%s",
							"issued_at": "2016-02-21T13:28:30.656527",
							"tenant": {
								"description": null,
								"enabled": true,
								"id": "97ea299c37bb4e04b3779039ea8aba44",
								"name": "tenant"
							}
						}
					}
				}
			`,
			s.NetworkServiceEndpoint,
			s.NetworkServiceEndpoint,
			s.NetworkServiceEndpoint,
			s.Token)
	})
}

func registerTenants(s *TestSuite) {
	th.Mux.HandleFunc("/v2.0/tenants", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
			{
				"tenants": [
					{
						"description": "Test tenat",
						"enabled": true,
						"id": "111111",
						"name": "demo"
					},
					{
						"description": "admin tenant",
						"enabled": true,
						"id": "222222",
						"name": "admin"
					}
				],
				"tenants_links": []
			}
		`)
	})
}

func registerNetworks(s *TestSuite) {
	th.Mux.HandleFunc("/v2.0/networks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
			{
				"networks": [
					{
						"status": "ACTIVE",
						"subnets": [
							"54d6f61d-db07-451c-9ab3-b9609b6b6f0b"
						],
						"name": "private-network",
						"admin_state_up": true,
						"tenant_id": "111111",
						"shared": true,
						"id": "d32019d3-bc6e-4319-9c1d-6722fc136a21"
					},
					{
						"status": "ACTIVE",
						"subnets": [
							"54d6f61d-db07-451c-9ab3-b9609b6b6f0b"
						],
						"name": "private-network",
						"admin_state_up": true,
						"tenant_id": "222222",
						"shared": true,
						"id": "d32019d3-bc6e-4319-9c1d-6722fc136a22"
					},
					{
						"status": "ACTIVE",
						"subnets": [
							"08eae331-0402-425a-923c-34f7cfe39c1b"
						],
						"name": "private",
						"admin_state_up": true,
						"tenant_id": "222222",
						"shared": true,
						"id": "db193ab3-96e3-4cb3-8fc5-05f4296d0324"
					}
				]
			}
			`)
	})
}

func registerSubnets(s *TestSuite) {
	th.Mux.HandleFunc("/v2.0/subnets", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"subnets": [
				{
					"name": "private-subnet",
					"enable_dhcp": true,
					"network_id": "db193ab3-96e3-4cb3-8fc5-05f4296d0324",
					"tenant_id": "222222",
					"dns_nameservers": [],
					"allocation_pools": [
						{
							"start": "10.0.0.2",
							"end": "10.0.0.254"
						}
					],
					"host_routes": [],
					"ip_version": 4,
					"gateway_ip": "10.0.0.1",
					"cidr": "10.0.0.0/24",
					"id": "08eae331-0402-425a-923c-34f7cfe39c1b"
				},
				{
					"name": "my_subnet",
					"enable_dhcp": true,
					"network_id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
					"tenant_id": "222222",
					"dns_nameservers": [],
					"allocation_pools": [
						{
							"start": "192.0.0.2",
							"end": "192.255.255.254"
						}
					],
					"host_routes": [],
					"ip_version": 4,
					"gateway_ip": "192.0.0.1",
					"cidr": "192.0.0.0/8",
					"id": "54d6f61d-db07-451c-9ab3-b9609b6b6f0b"
				},
				{
					"name": "my_subnet",
					"enable_dhcp": true,
					"network_id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
					"tenant_id": "222222",
					"dns_nameservers": [],
					"allocation_pools": [
						{
							"start": "192.0.0.2",
							"end": "192.255.255.254"
						}
					],
					"host_routes": [],
					"ip_version": 4,
					"gateway_ip": "192.0.0.1",
					"cidr": "192.0.0.0/8",
					"id": "54d6f61d-db07-451c-9ab3-b9609b6b6f0c"
				}
			]
		}
        `)
	})
}

func registerRouters(s *TestSuite) {
	th.Mux.HandleFunc("/v2.0/routers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"routers": [
				{
					"status": "ACTIVE",
					"external_gateway_info": null,
					"name": "second_routers",
					"admin_state_up": true,
					"tenant_id": "222222",
					"distributed": false,
					"id": "7177abc4-5ae9-4bb7-b0d4-89e94a4abf3b"
				},
				{
					"status": "ACTIVE",
					"external_gateway_info": {
						"network_id": "3c5bcddd-6af9-4e6b-9c3e-c153e521cab8"
					},
					"name": "router1",
					"admin_state_up": true,
					"tenant_id": "222222",
					"distributed": false,
					"id": "a9254bdb-2613-4a13-ac4c-adc581fba50d"
				},
				{
					"status": "ACTIVE",
					"external_gateway_info": {
						"network_id": "3c5bcddd-6af9-4e6b-9c3e-c153e521cab8"
					},
					"name": "router1",
					"admin_state_up": true,
					"tenant_id": "222222",
					"distributed": false,
					"id": "a9254bdb-2613-4a13-ac4c-adc581fba50e"
				},
				{
					"status": "ACTIVE",
					"external_gateway_info": {
						"network_id": "3c5bcddd-6af9-4e6b-9c3e-c153e521cabg"
					},
					"name": "router1",
					"admin_state_up": true,
					"tenant_id": "222222",
					"distributed": false,
					"id": "a9254bdb-2613-4a13-ac4c-adc581fba50f"
				}
			]
		}
		`)
	})
}

func registerPorts(s *TestSuite) {
	th.Mux.HandleFunc("/v2.0/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"ports": [
				{
					"status": "ACTIVE",
					"binding:host_id": "devstack",
					"name": "",
					"admin_state_up": true,
					"network_id": "70c1db1f-b701-45bd-96e0-a313ee3430b3",
					"tenant_id": "222222",
					"device_owner": "network:router_gateway",
					"mac_address": "fa:16:3e:58:42:ed",
					"fixed_ips": [
						{
							"subnet_id": "008ba151-0b8c-4a67-98b5-0d2b87666062",
							"ip_address": "172.24.4.2"
						}
					],
					"id": "d80b1a3b-4fc1-49f3-952e-1e2ab7081d8b",
					"security_groups": [],
					"device_id": "9ae135f4-b6e0-4dad-9e91-3c223e385824"
				},
				{
					"status": "ACTIVE",
					"binding:host_id": "devstack",
					"name": "",
					"admin_state_up": true,
					"network_id": "70c1db1f-b701-45bd-96e0-a313ee3430b3",
					"tenant_id": "222222",
					"device_owner": "network:router_gateway",
					"mac_address": "fa:16:3e:58:42:ed",
					"fixed_ips": [
						{
							"subnet_id": "008ba151-0b8c-4a67-98b5-0d2b87666062",
							"ip_address": "172.24.4.2"
						}
					],
					"id": "d80b1a3b-4fc1-49f3-952e-1e2ab7081d8c",
					"security_groups": [],
					"device_id": "9ae135f4-b6e0-4dad-9e91-3c223e385824"
				},
				{
					"status": "ACTIVE",
					"binding:host_id": "devstack",
					"name": "",
					"admin_state_up": true,
					"network_id": "70c1db1f-b701-45bd-96e0-a313ee3430b3",
					"tenant_id": "222222",
					"device_owner": "network:router_gateway",
					"mac_address": "fa:16:3e:58:42:ed",
					"fixed_ips": [
						{
							"subnet_id": "008ba151-0b8c-4a67-98b5-0d2b87666062",
							"ip_address": "172.24.4.2"
						}
					],
					"id": "d80b1a3b-4fc1-49f3-952e-1e2ab7081d8d",
					"security_groups": [],
					"device_id": "9ae135f4-b6e0-4dad-9e91-3c223e385824"
				}
			]
		}
      `)
	})
}

func registerFloatingIPs(s *TestSuite) {
	th.Mux.HandleFunc("/v2.0/floatingips", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"floatingips": [
				{
					"floating_network_id": "6d67c30a-ddb4-49a1-bec3-a65b286b4170",
					"router_id": null,
					"fixed_ip_address": null,
					"floating_ip_address": "192.0.0.4",
					"tenant_id": "222222",
					"status": "DOWN",
					"port_id": null,
					"id": "2f95fd2b-9f6a-4e8e-9e9a-2cbe286cbf9e"
				},
				{
					"floating_network_id": "90f742b1-6d17-487b-ba95-71881dbc0b64",
					"router_id": "0a24cb83-faf5-4d7f-b723-3144ed8a2167",
					"fixed_ip_address": "192.0.0.2",
					"floating_ip_address": "10.0.0.3",
					"tenant_id": "222222",
					"status": "DOWN",
					"port_id": "74a342ce-8e07-4e91-880c-9f834b68fa25",
					"id": "ada25a95-f321-4f59-b0e0-f3a970dd3d63"
				}
			]
		}
		`)
	})
}

func registerQuotas(s *TestSuite) {
	th.Mux.HandleFunc("/v2.0/quotas/222222", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
				{
					"quota": {
						"subnet": 10,
						"network": 13,
						"floatingip": 50,
						"subnetpool": -1,
						"security_group_rule": 100,
						"security_group": 10,
						"router": 15,
						"rbac_policy": -1,
						"port": 50
					}
				}
			`)
	})
	th.Mux.HandleFunc("/v2.0/quotas/111111", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
			{
				"quota": {
					"subnet": 11,
					"network": 12,
					"floatingip": 51,
					"subnetpool": 0,
					"security_group_rule": 101,
					"security_group": 11,
					"router": 16,
					"rbac_policy": 0,
					"port": 51
				}
			}
		`)
	})
}
