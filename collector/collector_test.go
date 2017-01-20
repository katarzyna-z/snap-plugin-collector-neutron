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
