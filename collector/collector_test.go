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

package collector

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/intelsdi-x/snap-plugin-utilities/str"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
	th "github.com/rackspace/gophercloud/testhelper"
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
	registerEndpoints(s)
	registerTenants(s)
	registerNetworks(s)
	registerSubnets(s)
	registerRouters(s)
	registerPorts(s)
	registerFloatingIPs(s)
	registerQuotas(s)
}

func (s *TestSuite) TearDownSuite() {
	th.TeardownHTTP()
}

func TestRunSuite(t *testing.T) {
	keystoneTestSuite := new(TestSuite)
	suite.Run(t, keystoneTestSuite)
}

func (s *TestSuite) TestGetMetricTypes() {
	Convey("Given config with enpoint, user and password defined", s.T(), func() {
		cfg := setupCfg(th.Endpoint(), "me", "secret", "admin")
		collector := New()
		So(collector, ShouldNotBeNil)
		mts, err := collector.GetMetricTypes(cfg)

		Convey("Then no error should be reported", func() {
			So(err, ShouldBeNil)
		})

		Convey("and proper metric types are returned", func() {
			metricNames := []string{}
			for _, m := range mts {
				metricNames = append(metricNames, m.Namespace().String())
			}

			So(len(mts), ShouldEqual, 28)

			ns := core.NewNamespace(vendor, openstack, pluginName, "admin", networksCountMetric)
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "admin", subnetsCountMetric)
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "admin", routersCountMetric)
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "admin", portsCountMetric)
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "admin", floatingipsCountMetric)
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"subnet")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"network")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"floatingip")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"subnetpool")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"security_group_rule")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"security_group")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"router")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"rbac_policy")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"port")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)

			ns = core.NewNamespace(vendor, openstack, pluginName, "demo", networksCountMetric)
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "demo", subnetsCountMetric)
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "demo", routersCountMetric)
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "demo", portsCountMetric)
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "demo", floatingipsCountMetric)
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"subnet")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"network")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"floatingip")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"subnetpool")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"security_group_rule")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"security_group")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"router")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"rbac_policy")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
			ns = core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"port")
			So(str.Contains(metricNames, ns.String()), ShouldBeTrue)
		})
	})
}

func (s *TestSuite) TestCollectMetrics() {
	Convey("Given set of metric types", s.T(), func() {
		cfg := setupCfg(th.Endpoint(), "admin", "secret", "admin")

		ns1 := core.NewNamespace(vendor, openstack, pluginName, "admin", networksCountMetric)
		ns2 := core.NewNamespace(vendor, openstack, pluginName, "admin", subnetsCountMetric)
		ns3 := core.NewNamespace(vendor, openstack, pluginName, "admin", routersCountMetric)
		ns4 := core.NewNamespace(vendor, openstack, pluginName, "admin", portsCountMetric)
		ns5 := core.NewNamespace(vendor, openstack, pluginName, "admin", floatingipsCountMetric)
		ns6 := core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"subnet")
		ns7 := core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"network")
		ns8 := core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"floatingip")
		ns9 := core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"subnetpool")
		ns10 := core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"security_group_rule")
		ns11 := core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"security_group")
		ns12 := core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"router")
		ns13 := core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"rbac_policy")
		ns14 := core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"port")

		ns15 := core.NewNamespace(vendor, openstack, pluginName, "demo", networksCountMetric)
		ns16 := core.NewNamespace(vendor, openstack, pluginName, "demo", subnetsCountMetric)
		ns17 := core.NewNamespace(vendor, openstack, pluginName, "demo", routersCountMetric)
		ns18 := core.NewNamespace(vendor, openstack, pluginName, "demo", portsCountMetric)
		ns19 := core.NewNamespace(vendor, openstack, pluginName, "demo", floatingipsCountMetric)
		ns20 := core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"subnet")
		ns21 := core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"network")
		ns22 := core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"floatingip")
		ns23 := core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"subnetpool")
		ns24 := core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"security_group_rule")
		ns25 := core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"security_group")
		ns26 := core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"router")
		ns27 := core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"rbac_policy")
		ns28 := core.NewNamespace(vendor, openstack, pluginName, "demo", quotas+"port")

		mTypes := []plugin.MetricType{
			plugin.MetricType{Namespace_: ns1, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns2, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns3, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns4, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns5, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns6, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns7, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns8, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns9, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns10, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns11, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns12, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns13, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns14, Config_: cfg.ConfigDataNode},

			plugin.MetricType{Namespace_: ns15, Config_: cfg.ConfigDataNode},

			plugin.MetricType{Namespace_: ns16, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns17, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns18, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns19, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns20, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns21, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns22, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns23, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns24, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns25, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns26, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns27, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns28, Config_: cfg.ConfigDataNode},
		}

		Convey("When ColelctMetrics() is called", func() {
			collector := New()

			Convey("Then instance of plugin should be created", func() {
				So(collector, ShouldNotBeNil)
			})

			mts, err := collector.CollectMetrics(mTypes)

			Convey("Then no error should be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("and proper metric types are returned", func() {
				So(mts, ShouldNotBeNil)

				metricNames := map[string]interface{}{}
				for _, m := range mts {
					ns := m.Namespace().String()
					metricNames[ns] = m.Data()
				}

				So(len(mts), ShouldEqual, 28)

				//networks_count
				val, ok := metricNames[ns1.String()]
				So(ok, ShouldBeTrue)
				v, ok := val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 2)

				val, ok = metricNames[ns15.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 1)

				//subnets_counts
				val, ok = metricNames[ns2.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 3)

				val, ok = metricNames[ns16.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 0)

				//routers_count
				val, ok = metricNames[ns3.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 4)

				val, ok = metricNames[ns17.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 0)

				//ports_count
				val, ok = metricNames[ns4.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 3)

				val, ok = metricNames[ns18.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 0)

				//floatingips_count
				val, ok = metricNames[ns5.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 2)

				val, ok = metricNames[ns19.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 0)

				//quotas_subnet
				val, ok = metricNames[ns6.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 10)

				val, ok = metricNames[ns20.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 11)

				//quotas_network
				val, ok = metricNames[ns7.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 13)

				val, ok = metricNames[ns21.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 12)

				//quotas_floatingip
				val, ok = metricNames[ns8.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 50)

				val, ok = metricNames[ns22.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 51)

				//quotas_subnetpool
				val, ok = metricNames[ns9.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, -1)

				val, ok = metricNames[ns23.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 0)

				//quotas_security_group_rule
				val, ok = metricNames[ns10.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 100)

				val, ok = metricNames[ns24.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 101)

				//quotas_security_group
				val, ok = metricNames[ns11.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 10)

				val, ok = metricNames[ns25.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 11)

				//quotas_router
				val, ok = metricNames[ns12.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 15)

				val, ok = metricNames[ns26.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 16)

				//quotas_rbac_policy
				val, ok = metricNames[ns13.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, -1)

				val, ok = metricNames[ns27.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 0)

				//quotas_port
				val, ok = metricNames[ns14.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 50)

				val, ok = metricNames[ns28.String()]
				So(ok, ShouldBeTrue)
				v, ok = val.(int64)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 51)
			})
		})
	})

	Convey("Given set of incorrec metric types", s.T(), func() {
		cfg := setupCfg(th.Endpoint(), "admin", "secret", "admin")
		ns1 := core.NewNamespace(vendor, openstack, pluginName, "admin123", networksCountMetric)
		ns2 := core.NewNamespace(vendor, openstack, pluginName, "admin123", networksCountMetric, "test")
		ns3 := core.NewNamespace(vendor, openstack, pluginName, "admin", "test")
		ns4 := core.NewNamespace(vendor, openstack, pluginName, "admin", quotas+"test")
		mTypes := []plugin.MetricType{
			plugin.MetricType{Namespace_: ns1, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns2, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns3, Config_: cfg.ConfigDataNode},
			plugin.MetricType{Namespace_: ns4, Config_: cfg.ConfigDataNode},
		}

		Convey("When ColelctMetrics() is called", func() {
			collector := New()

			Convey("Then instance of plugin should be created", func() {
				So(collector, ShouldNotBeNil)
			})

			mts, err := collector.CollectMetrics(mTypes)

			Convey("Then no error should be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then metric list should have correct size", func() {
				So(len(mts), ShouldEqual, 0)
			})
		})
	})
}

func (s *TestSuite) TestGetConfigPolicy() {
	Convey("Meta should return metadata for the plugin", s.T(), func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, pluginName)
		So(meta.Version, ShouldResemble, version)
		So(meta.Type, ShouldResemble, plugin.CollectorPluginType)
	})

	Convey("Given config with enpoint, user and password defined", s.T(), func() {
		collector := New()

		Convey("So collector should not be nil", func() {
			So(collector, ShouldNotBeNil)
		})

		configPolicy, err := collector.GetConfigPolicy()

		Convey("collector.GetConfigPolicy() should return a config policy", func() {
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})

			Convey("So we should not get an err retreiving the config policy", func() {
				So(err, ShouldBeNil)
			})

			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})

			correctConfig := make(map[string]ctypes.ConfigValue)
			correctConfig[cfgURL] = ctypes.ConfigValueStr{Value: th.Endpoint()}
			correctConfig[cfgUser] = ctypes.ConfigValueStr{Value: "admin"}
			correctConfig[cfgPassword] = ctypes.ConfigValueStr{Value: "pass"}
			correctConfig[cfgTenant] = ctypes.ConfigValueStr{Value: "admin"}

			cfg, errs := configPolicy.Get([]string{""}).Process(correctConfig)

			Convey("So config policy should process correctConfig and return a config", func() {
				So(cfg, ShouldNotBeNil)
			})

			Convey("So correctConfig processing should return no errors", func() {
				So(errs.HasErrors(), ShouldBeFalse)
			})

			wrongConfig1 := make(map[string]ctypes.ConfigValue)

			cfg1, errs1 := configPolicy.Get([]string{""}).Process(wrongConfig1)

			Convey("So config policy should not return a config after processing wrongConfig1", func() {
				So(cfg1, ShouldBeNil)
			})

			Convey("So wrongConfig1 processing should return no errors", func() {
				So(errs1.HasErrors(), ShouldBeTrue)
			})

			wrongConfig2 := make(map[string]ctypes.ConfigValue)

			cfg2, errs2 := configPolicy.Get([]string{""}).Process(wrongConfig2)

			correctConfig[cfgURL] = ctypes.ConfigValueStr{Value: th.Endpoint()}
			correctConfig[cfgUser] = ctypes.ConfigValueStr{Value: "admin"}

			Convey("So config policy should not return a config after processing wrongConfig2", func() {
				So(cfg2, ShouldBeNil)
			})
			Convey("So wrongConfig2 processing should return errors", func() {
				So(errs2.HasErrors(), ShouldBeTrue)
			})
		})
	})
}

func setupCfg(endpoint, user, password, tenant string) plugin.ConfigType {
	node := cdata.NewNode()
	node.AddItem(cfgURL, ctypes.ConfigValueStr{Value: endpoint})
	node.AddItem(cfgUser, ctypes.ConfigValueStr{Value: user})
	node.AddItem(cfgPassword, ctypes.ConfigValueStr{Value: password})
	node.AddItem(cfgTenant, ctypes.ConfigValueStr{Value: tenant})
	return plugin.ConfigType{ConfigDataNode: node}
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
			`, s.NetworkServiceEndpoint,
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

func registerEndpoints(s *TestSuite) {
	th.Mux.HandleFunc("/v2/endpoints", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
			{
				"endpoints": [
					{
						"enabled": true,
						"id": "035d8ea06d88496e929310a7cda173a6",
						"interface": "public",
						"links": {
							"self": "https://public.fuel.local:5000/v3/endpoints/035d8ea06d88496e929310a7cda173a6"
						},
						"region": "RegionOne",
						"region_id": "RegionOne",
						"service_id": "dc52eef88c2d470fb68912a2641eaab4",
						"url": "https://public.fuel.local:8773/services/Cloud"
					},
					{
						"enabled": true,
						"id": "0a587d5392a54d4e8eb8d7328a7acbf1",
						"interface": "public",
						"links": {
							"self": "https://public.fuel.local:5000/v3/endpoints/0a587d5392a54d4e8eb8d7328a7acbf1"
						},
						"region": "RegionOne",
						"region_id": "RegionOne",
						"service_id": "615e06490105462cbbbab919bbe1c725",
						"url": "https://public.fuel.local:8777"
					},
					{
						"enabled": true,
						"id": "0b74729c70814c5cb65be5e7b56d56ed",
						"interface": "admin",
						"links": {
							"self": "https://public.fuel.local:5000/v3/endpoints/0b74729c70814c5cb65be5e7b56d56ed"
						},
						"region": "RegionOne",
						"region_id": "RegionOne",
						"service_id": "efbf568dd1234f52a73869c8cab10d93",
						"url": "http://192.168.20.2:9696"
					},
					{
						"enabled": true,
						"id": "159572c2ce0d480db916fc4986812d0b",
						"interface": "internal",
						"links": {
							"self": "https://public.fuel.local:5000/v3/endpoints/159572c2ce0d480db916fc4986812d0b"
						},
						"region": "RegionOne",
						"region_id": "RegionOne",
						"service_id": "efbf568dd1234f52a73869c8cab10d93",
						"url": "http://192.168.20.2:9696"
					}
				]
			}
		`)
	})
}
