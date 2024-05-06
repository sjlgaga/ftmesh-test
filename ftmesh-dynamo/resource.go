// Copyright 2020 Envoyproxy Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package example

import (
	"time"

	v3 "github.com/cncf/xds/go/xds/core/v3"
	v32 "github.com/cncf/xds/go/xds/type/matcher/v3"
	udp_proxyv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/udp/udp_proxy/v3"
	networkv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/matching/common_inputs/network/v3"
	any1 "github.com/golang/protobuf/ptypes/any"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	router "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
)

const (
	ClusterName  = "svc-a"
	RouteName    = "local_route"
	ListenerName = "listener_1"
	ListenerPort = 10980
	UpstreamHost = "10.214.96.108"
	UpstreamPort = 10730
)

func makeCluster(clusterName string) *cluster.Cluster {
	return &cluster.Cluster{
		Name:                 clusterName,
		ConnectTimeout:       durationpb.New(5 * time.Second),
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_EDS},
		LbPolicy:             cluster.Cluster_ROUND_ROBIN,
		//LoadAssignment:       makeEndpoint(clusterName),
		DnsLookupFamily: cluster.Cluster_V4_ONLY,
		EdsClusterConfig: &cluster.Cluster_EdsClusterConfig{
			EdsConfig: &core.ConfigSource{
				ConfigSourceSpecifier: &core.ConfigSource_Ads{
					//	ApiConfigSource: &core.ApiConfigSource{
					//		ApiType: core.ApiConfigSource_GRPC,
					//		GrpcServices: []*core.GrpcService{{
					//			TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					//				EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "xds_cluster"},
					//			},
					//		}},
					//		TransportApiVersion: core.ApiVersion_V3,
					//	},
				},
				InitialFetchTimeout: &durationpb.Duration{Seconds: 10},
			},
			ServiceName: clusterName,
		},
	}
}

func makeEndpoint(clusterName string, upHost string, upPort uint32) *endpoint.ClusterLoadAssignment {
	return &endpoint.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			LbEndpoints: []*endpoint.LbEndpoint{{
				HostIdentifier: &endpoint.LbEndpoint_Endpoint{
					Endpoint: &endpoint.Endpoint{
						Address: &core.Address{
							Address: &core.Address_SocketAddress{
								SocketAddress: &core.SocketAddress{
									Protocol: core.SocketAddress_TCP,
									Address:  upHost,
									PortSpecifier: &core.SocketAddress_PortValue{
										PortValue: upPort,
									},
								},
							},
						},
					},
				},
			}},
		}},
	}
}

func makeRoute(routeName string, clusterName string) *route.RouteConfiguration {
	return &route.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []*route.VirtualHost{{
			Name:    "local_service",
			Domains: []string{"*"},
			Routes: []*route.Route{{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: "/",
					},
				},
				Action: &route.Route_Route{
					Route: &route.RouteAction{
						ClusterSpecifier: &route.RouteAction_Cluster{
							Cluster: clusterName,
						},
						HostRewriteSpecifier: &route.RouteAction_HostRewriteLiteral{
							HostRewriteLiteral: UpstreamHost,
						},
					},
				},
			}},
		}},
	}
}

func makeHTTPListener(listenerName string, routeClusterA string, routeClusterB string, routeClusterDB string, rewritePrefix string, listenPort uint32) *listener.Listener {
	routerConfig, _ := anypb.New(&router.Router{})
	// HTTP filter configuration
	stateConfig := &any1.Any{TypeUrl: "type.googleapis.com/envoy.extensions.filters.http.states_replication.v3.StatesReplication"}
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "ingress_http",
		RouteSpecifier: &hcm.HttpConnectionManager_RouteConfig{
			RouteConfig: &route.RouteConfiguration{
				Name: "local_route",
				VirtualHosts: []*route.VirtualHost{
					{
						Name:    "local_service",
						Domains: []string{"*"},
						Routes: []*route.Route{
							{
								Match: &route.RouteMatch{
									PathSpecifier: &route.RouteMatch_Prefix{
										Prefix: "/svc-a",
									},
								},
								Action: &route.Route_Route{
									Route: &route.RouteAction{
										ClusterSpecifier: &route.RouteAction_Cluster{
											Cluster: routeClusterA,
										},
									},
								},
							},
							{
								Match: &route.RouteMatch{
									PathSpecifier: &route.RouteMatch_Prefix{
										Prefix: "/svc-b",
									},
								},
								Action: &route.Route_Route{
									Route: &route.RouteAction{
										ClusterSpecifier: &route.RouteAction_Cluster{
											Cluster: routeClusterB,
										},
									},
								},
							},
							{
								Match: &route.RouteMatch{
									PathSpecifier: &route.RouteMatch_Prefix{Prefix: "/db"},
								},
								Action: &route.Route_Route{
									Route: &route.RouteAction{
										ClusterSpecifier: &route.RouteAction_Cluster{
											Cluster: routeClusterDB,
										},
										PrefixRewrite: rewritePrefix,
									},
								},
							},
						},
					},
				},
			},
		},
		HttpFilters: []*hcm.HttpFilter{
			{
				Name:       wellknown.StatesReplication,
				ConfigType: &hcm.HttpFilter_TypedConfig{TypedConfig: stateConfig},
			},
			{
				Name:       wellknown.Router,
				ConfigType: &hcm.HttpFilter_TypedConfig{TypedConfig: routerConfig},
			},
		},
	}
	pbst, err := anypb.New(manager)
	if err != nil {
		panic(err)
	}

	return &listener.Listener{
		Name: listenerName,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: listenPort,
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{{
			Filters: []*listener.Filter{{
				Name: wellknown.HTTPConnectionManager,
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		}},
	}
}

func makeConfigSource() *core.ConfigSource {
	source := &core.ConfigSource{}
	source.ResourceApiVersion = resource.DefaultAPIVersion
	source.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &core.ApiConfigSource{
			TransportApiVersion:       resource.DefaultAPIVersion,
			ApiType:                   core.ApiConfigSource_GRPC,
			SetNodeOnFirstMessageOnly: true,
			GrpcServices: []*core.GrpcService{{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "xds_cluster"},
				},
			}},
		},
	}
	return source
}

func makeUDPListener(listenerName string) *listener.Listener {
	sourceIPInput, err := anypb.New(&networkv3.SourceIPInput{})
	if err != nil {
		panic(err)
	}
	routeAction, err := anypb.New(&udp_proxyv3.Route{Cluster: ClusterName})
	if err != nil {
		panic(err)
	}
	udpProxyAny, err := anypb.New(&udp_proxyv3.UdpProxyConfig{
		StatPrefix:                "service",
		UsePerPacketLoadBalancing: true,
		RouteSpecifier: &udp_proxyv3.UdpProxyConfig_Matcher{
			Matcher: &v32.Matcher{
				MatcherType: &v32.Matcher_MatcherList_{
					MatcherList: &v32.Matcher_MatcherList{
						Matchers: []*v32.Matcher_MatcherList_FieldMatcher{
							{
								Predicate: &v32.Matcher_MatcherList_Predicate{
									MatchType: &v32.Matcher_MatcherList_Predicate_SinglePredicate_{
										SinglePredicate: &v32.Matcher_MatcherList_Predicate_SinglePredicate{
											Input: &v3.TypedExtensionConfig{
												Name:        "envoy.matching.inputs.source_ip",
												TypedConfig: sourceIPInput,
											},
											Matcher: &v32.Matcher_MatcherList_Predicate_SinglePredicate_ValueMatch{
												ValueMatch: &v32.StringMatcher{
													MatchPattern: &v32.StringMatcher_Exact{
														Exact: "127.0.0.1",
													},
												},
											},
										},
									},
								},
								OnMatch: &v32.Matcher_OnMatch{
									OnMatch: &v32.Matcher_OnMatch_Action{
										Action: &v3.TypedExtensionConfig{
											Name:        "envoy.extensions.filters.udp.udp_proxy.v3.Route",
											TypedConfig: routeAction,
										},
									},
								},
							},
						}},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	return &listener.Listener{
		Name: listenerName,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_UDP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: ListenerPort,
					},
				},
			},
		},
		ListenerFilters: []*listener.ListenerFilter{
			&listener.ListenerFilter{
				Name: "envoy.filters.udp_listener.udp_proxy",
				ConfigType: &listener.ListenerFilter_TypedConfig{
					TypedConfig: udpProxyAny,
				},
			},
		},
	}
}

func GenerateSnapshot() *cache.Snapshot {
	snap, _ := cache.NewSnapshot("1",
		map[resource.Type][]types.Resource{
			resource.ClusterType: {
				makeCluster("svc-a-110"), makeCluster("svc-a-108"), makeCluster("svc-a-107"),
				makeCluster("svc-b-107"), makeCluster("svc-b-110"), makeCluster("svc-b-108"),
				makeCluster("db-110"), makeCluster("db-108"), makeCluster("db-107")},
			// resource.RouteType:    {makeRoute(RouteName, ClusterName)},
			resource.ListenerType: {
				makeHTTPListener("listener-110", "svc-a-110", "svc-b-110", "db-110", "/db", 10729),
				makeHTTPListener("listener-108", "svc-a-108", "svc-b-108", "db-108", "/", 10728),
				makeHTTPListener("listener-107", "svc-a-107", "svc-b-107", "db-107", "/db", 10727)},
			resource.EndpointType: {
				makeEndpoint("svc-a-110", "127.0.0.1", 10730),
				makeEndpoint("svc-a-108", "128.110.219.68", 10728),
				makeEndpoint("svc-a-107", "128.110.219.68", 10729),
				makeEndpoint("svc-b-110", "128.110.219.70", 10728),
				makeEndpoint("svc-b-108", "127.0.0.1", 20730),
				makeEndpoint("svc-b-107", "128.110.219.70", 10728),
				makeEndpoint("db-108", "127.0.0.1", 8000),
				makeEndpoint("db-110", "128.110.219.68", 10728),
				makeEndpoint("db-107", "128.110.219.68", 10728),
			},
		},
	)
	return snap
}
