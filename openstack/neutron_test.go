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
	s.Token = "cefb1b0ba45744488e6ed702db699327"
	s.NetworkServiceEndpoint = th.Endpoint()
	th.Mux.HandleFunc("/v2.0/tokens", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
				{
					"access": {
						"metadata": {
							"is_admin": 0,
							"roles": [
								"dc6ed0c1bfb847c9b087e1d62068766b"
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
							"expires": "2017-01-19T15:33:11Z",
							"id": "%s",
							"issued_at": "2017-01-19T14:33:11.541197Z",
							"tenant": {
								"description": "",
								"enabled": true,
								"id": "444f244be5e34ce0816e8beccfd332ef",
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
						"description": "demo tenat",
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
			      "admin_state_up": true,
			      "availability_zone_hints": [],
			      "availability_zones": [
				"nova"
			      ],
			      "created_at": "2016-09-08T12:01:34",
			      "description": "",
			      "id": "28dd974d-0ec0-43cc-86ac-06773acb126f",
			      "ipv4_address_scope": null,
			      "ipv6_address_scope": null,
			      "mtu": 1450,
			      "name": "private",
			      "port_security_enabled": true,
			      "provider:network_type": "vxlan",
			      "provider:physical_network": null,
			      "provider:segmentation_id": 19,
			      "revision": 6,
			      "router:external": false,
			      "shared": false,
			      "status": "ACTIVE",
			      "subnets": [
				"ef512c07-9203-4121-9df6-e692a4bc84c5",
				"64c8fbe0-cb8a-41d7-9e65-56f33f9674cb"
			      ],
			      "tags": [],
			      "tenant_id": "111111",
			      "updated_at": "2016-09-08T12:01:38"
			    },
			    {
			      "admin_state_up": true,
			      "availability_zone_hints": [],
			      "availability_zones": [
				"nova"
			      ],
			      "created_at": "2016-09-08T12:01:51",
			      "description": "",
			      "id": "f3722668-e9e7-41dd-8086-5e1b9f5d8209",
			      "ipv4_address_scope": null,
			      "ipv6_address_scope": null,
			      "is_default": true,
			      "mtu": 1500,
			      "name": "public",
			      "port_security_enabled": true,
			      "provider:network_type": "flat",
			      "provider:physical_network": "public",
			      "provider:segmentation_id": null,
			      "revision": 6,
			      "router:external": true,
			      "shared": false,
			      "status": "ACTIVE",
			      "subnets": [
				"4582d819-7ded-4ec5-aa92-27f73781f625",
				"94daf3aa-6faf-43c0-a21c-9656110b3d11"
			      ],
			      "tags": [],
			      "tenant_id": "222222",
			      "updated_at": "2016-09-08T12:02:05"
			    },
			    {
			      "admin_state_up": true,
			      "availability_zone_hints": [],
			      "availability_zones": [
				"nova"
			      ],
			      "created_at": "2016-09-08T12:01:51",
			      "description": "",
			      "id": "f3722668-e9e7-41dd-8086-5e1b9f5d8209",
			      "ipv4_address_scope": null,
			      "ipv6_address_scope": null,
			      "is_default": true,
			      "mtu": 1500,
			      "name": "public",
			      "port_security_enabled": true,
			      "provider:network_type": "flat",
			      "provider:physical_network": "public",
			      "provider:segmentation_id": null,
			      "revision": 6,
			      "router:external": true,
			      "shared": false,
			      "status": "ACTIVE",
			      "subnets": [
				"4582d819-7ded-4ec5-aa92-27f73781f625",
				"94daf3aa-6faf-43c0-a21c-9656110b3d11"
			      ],
			      "tags": [],
			      "tenant_id": "222222",
			      "updated_at": "2016-09-08T12:02:05"
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
			      "allocation_pools": [
				{
				  "end": "2001:db8::1",
				  "start": "2001:db8::1"
				},
				{
				  "end": "2001:db8::ffff:ffff:ffff:ffff",
				  "start": "2001:db8::3"
				}
			      ],
			      "cidr": "2001:db8::/64",
			      "created_at": "2016-09-08T12:02:05",
			      "description": "",
			      "dns_nameservers": [],
			      "enable_dhcp": false,
			      "gateway_ip": "2001:db8::2",
			      "host_routes": [],
			      "id": "4582d819-7ded-4ec5-aa92-27f73781f625",
			      "ip_version": 6,
			      "ipv6_address_mode": null,
			      "ipv6_ra_mode": null,
			      "name": "ipv6-public-subnet",
			      "network_id": "f3722668-e9e7-41dd-8086-5e1b9f5d8209",
			      "revision": 2,
			      "service_types": [],
			      "subnetpool_id": null,
			      "tenant_id": "222222",
			      "updated_at": "2016-09-08T12:02:05"
			    },
			    {
			      "allocation_pools": [
				{
				  "end": "10.0.0.254",
				  "start": "10.0.0.2"
				}
			      ],
			      "cidr": "10.0.0.0/24",
			      "created_at": "2016-09-08T12:01:36",
			      "description": "",
			      "dns_nameservers": [],
			      "enable_dhcp": true,
			      "gateway_ip": "10.0.0.1",
			      "host_routes": [],
			      "id": "64c8fbe0-cb8a-41d7-9e65-56f33f9674cb",
			      "ip_version": 4,
			      "ipv6_address_mode": null,
			      "ipv6_ra_mode": null,
			      "name": "private-subnet",
			      "network_id": "28dd974d-0ec0-43cc-86ac-06773acb126f",
			      "revision": 2,
			      "service_types": [],
			      "subnetpool_id": null,
			      "tenant_id": "222222",
			      "updated_at": "2016-09-08T12:01:36"
			    },
			    {
			      "allocation_pools": [
				{
				  "end": "172.24.4.254",
				  "start": "172.24.4.2"
				}
			      ],
			      "cidr": "172.24.4.0/24",
			      "created_at": "2016-09-08T12:01:56",
			      "description": "",
			      "dns_nameservers": [],
			      "enable_dhcp": false,
			      "gateway_ip": "172.24.4.1",
			      "host_routes": [],
			      "id": "94daf3aa-6faf-43c0-a21c-9656110b3d11",
			      "ip_version": 4,
			      "ipv6_address_mode": null,
			      "ipv6_ra_mode": null,
			      "name": "public-subnet",
			      "network_id": "f3722668-e9e7-41dd-8086-5e1b9f5d8209",
			      "revision": 2,
			      "service_types": [],
			      "subnetpool_id": null,
			      "tenant_id": "222222",
			      "updated_at": "2016-09-08T12:01:56"
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
			      "admin_state_up": true,
			      "availability_zone_hints": [],
			      "availability_zones": [
				"nova"
			      ],
			      "description": "",
			      "distributed": false,
			      "external_gateway_info": {
				"enable_snat": true,
				"external_fixed_ips": [
				  {
				    "ip_address": "172.24.4.3",
				    "subnet_id": "94daf3aa-6faf-43c0-a21c-9656110b3d11"
				  },
				  {
				    "ip_address": "2001:db8::1",
				    "subnet_id": "4582d819-7ded-4ec5-aa92-27f73781f625"
				  }
				],
				"network_id": "f3722668-e9e7-41dd-8086-5e1b9f5d8209"
			      },
			      "flavor_id": null,
			      "ha": false,
			      "id": "a75c645a-6dcc-418c-9371-9be7054c395e",
			      "name": "router1",
			      "revision": 8,
			      "routes": [],
			      "status": "ACTIVE",
			      "tenant_id": "222222"
			    },
			    {
			      "admin_state_up": true,
			      "availability_zone_hints": [],
			      "availability_zones": [
				"nova"
			      ],
			      "description": "",
			      "distributed": false,
			      "external_gateway_info": {
				"enable_snat": true,
				"external_fixed_ips": [
				  {
				    "ip_address": "172.24.4.3",
				    "subnet_id": "94daf3aa-6faf-43c0-a21c-9656110b3d11"
				  },
				  {
				    "ip_address": "2001:db8::1",
				    "subnet_id": "4582d819-7ded-4ec5-aa92-27f73781f625"
				  }
				],
				"network_id": "f3722668-e9e7-41dd-8086-5e1b9f5d8209"
			      },
			      "flavor_id": null,
			      "ha": false,
			      "id": "a75c645a-6dcc-418c-9371-9be7054c395e",
			      "name": "router1",
			      "revision": 8,
			      "routes": [],
			      "status": "ACTIVE",
			      "tenant_id": "222222"
			    },
			    {
			      "admin_state_up": true,
			      "availability_zone_hints": [],
			      "availability_zones": [
				"nova"
			      ],
			      "description": "",
			      "distributed": false,
			      "external_gateway_info": {
				"enable_snat": true,
				"external_fixed_ips": [
				  {
				    "ip_address": "172.24.4.3",
				    "subnet_id": "94daf3aa-6faf-43c0-a21c-9656110b3d11"
				  },
				  {
				    "ip_address": "2001:db8::1",
				    "subnet_id": "4582d819-7ded-4ec5-aa92-27f73781f625"
				  }
				],
				"network_id": "f3722668-e9e7-41dd-8086-5e1b9f5d8209"
			      },
			      "flavor_id": null,
			      "ha": false,
			      "id": "a75c645a-6dcc-418c-9371-9be7054c395e",
			      "name": "router1",
			      "revision": 8,
			      "routes": [],
			      "status": "ACTIVE",
			      "tenant_id": "222222"
			    },
			    {
			      "admin_state_up": true,
			      "availability_zone_hints": [],
			      "availability_zones": [
				"nova"
			      ],
			      "description": "",
			      "distributed": false,
			      "external_gateway_info": {
				"enable_snat": true,
				"external_fixed_ips": [
				  {
				    "ip_address": "172.24.4.3",
				    "subnet_id": "94daf3aa-6faf-43c0-a21c-9656110b3d11"
				  },
				  {
				    "ip_address": "2001:db8::1",
				    "subnet_id": "4582d819-7ded-4ec5-aa92-27f73781f625"
				  }
				],
				"network_id": "f3722668-e9e7-41dd-8086-5e1b9f5d8209"
			      },
			      "flavor_id": null,
			      "ha": false,
			      "id": "a75c645a-6dcc-418c-9371-9be7054c395e",
			      "name": "router1",
			      "revision": 8,
			      "routes": [],
			      "status": "ACTIVE",
			      "tenant_id": "222222"
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
			      "admin_state_up": true,
			      "allowed_address_pairs": [],
			      "binding:host_id": "es-051",
			      "binding:profile": {},
			      "binding:vif_details": {
				"ovs_hybrid_plug": true,
				"port_filter": true
			      },
			      "binding:vif_type": "ovs",
			      "binding:vnic_type": "normal",
			      "created_at": "2016-10-20T13:07:02",
			      "description": "",
			      "device_id": "dd5dfb24-0e5d-41b1-a4a8-4693906d5f5b",
			      "device_owner": "compute:nova",
			      "extra_dhcp_opts": [],
			      "fixed_ips": [
				{
				  "ip_address": "fdaf:f360:d434:0:f816:3eff:fe3b:fe08",
				  "subnet_id": "ef512c07-9203-4121-9df6-e692a4bc84c5"
				},
				{
				  "ip_address": "10.0.0.5",
				  "subnet_id": "64c8fbe0-cb8a-41d7-9e65-56f33f9674cb"
				}
			      ],
			      "id": "004e9c25-de09-4d4c-a2c3-50b05defaac9",
			      "mac_address": "fa:16:3e:3b:fe:08",
			      "name": "",
			      "network_id": "28dd974d-0ec0-43cc-86ac-06773acb126f",
			      "port_security_enabled": true,
			      "revision": 10,
			      "security_groups": [
				"f47fd611-39d9-4999-9b10-41b19e03d40a"
			      ],
			      "status": "ACTIVE",
			      "tenant_id": "222222",
			      "updated_at": "2016-10-20T13:07:25"
			    },
			    {
			      "admin_state_up": true,
			      "allowed_address_pairs": [],
			      "binding:host_id": "es-051",
			      "binding:profile": {},
			      "binding:vif_details": {
				"ovs_hybrid_plug": true,
				"port_filter": true
			      },
			      "binding:vif_type": "ovs",
			      "binding:vnic_type": "normal",
			      "created_at": "2016-10-20T13:07:01",
			      "description": "",
			      "device_id": "92fcd563-3605-4cfa-9241-520566d79e68",
			      "device_owner": "compute:nova",
			      "extra_dhcp_opts": [],
			      "fixed_ips": [
				{
				  "ip_address": "fdaf:f360:d434:0:f816:3eff:fe17:68b3",
				  "subnet_id": "ef512c07-9203-4121-9df6-e692a4bc84c5"
				},
				{
				  "ip_address": "10.0.0.7",
				  "subnet_id": "64c8fbe0-cb8a-41d7-9e65-56f33f9674cb"
				}
			      ],
			      "id": "0a3bdc80-5b3e-4fca-baca-9716d94f56b5",
			      "mac_address": "fa:16:3e:17:68:b3",
			      "name": "",
			      "network_id": "28dd974d-0ec0-43cc-86ac-06773acb126f",
			      "port_security_enabled": true,
			      "revision": 8,
			      "security_groups": [
				"f47fd611-39d9-4999-9b10-41b19e03d40a"
			      ],
			      "status": "ACTIVE",
			      "tenant_id": "222222",
			      "updated_at": "2016-10-20T13:07:23"
			    },
			    {
			      "admin_state_up": true,
			      "allowed_address_pairs": [],
			      "binding:host_id": "es-051",
			      "binding:profile": {},
			      "binding:vif_details": {
				"ovs_hybrid_plug": true,
				"port_filter": true
			      },
			      "binding:vif_type": "ovs",
			      "binding:vnic_type": "normal",
			      "created_at": "2016-09-08T12:01:54",
			      "description": "",
			      "device_id": "a75c645a-6dcc-418c-9371-9be7054c395e",
			      "device_owner": "network:router_interface",
			      "extra_dhcp_opts": [],
			      "fixed_ips": [
				{
				  "ip_address": "10.0.0.1",
				  "subnet_id": "64c8fbe0-cb8a-41d7-9e65-56f33f9674cb"
				}
			      ],
			      "id": "11bc164c-c2dd-4809-9b04-0ef4aaefd8a2",
			      "mac_address": "fa:16:3e:da:6a:9a",
			      "name": "",
			      "network_id": "28dd974d-0ec0-43cc-86ac-06773acb126f",
			      "port_security_enabled": false,
			      "revision": 8,
			      "security_groups": [],
			      "status": "ACTIVE",
			      "tenant_id": "222222",
			      "updated_at": "2016-10-20T13:07:23"
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
					"floating_network_id": "28dc974d-0ec0-43cc-86ac-06773acb126f",
					"router_id": null,
					"fixed_ip_address": null,
					"floating_ip_address": "192.0.0.4",
					"tenant_id": "222222",
					"status": "DOWN",
					"port_id": "0a3bdc81-5b3e-4fca-baca-9716d94f56b5",
					"id": "a75c645a-6dcd-418c-9371-9be7054c395e"
				},
				{
					"floating_network_id": "ef51217-9203-4121-9df6-e692a4bc84c5",
					"router_id": "0a24cb83-faf5-4d7f-b723-3144ed8a2167",
					"fixed_ip_address": "192.0.0.2",
					"floating_ip_address": "10.0.0.3",
					"tenant_id": "222222",
					"status": "DOWN",
					"port_id": "004e9c25-de19-4d4c-a2c3-50b05defaac9",
					"id": "bfdd31f8-ccda-4722-9319-13a7138e226c"
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
