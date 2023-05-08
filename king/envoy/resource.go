package envoy

import (
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
)

func makeCluster() *cluster.Cluster {
	transportsocket := anypb.Any{
		TypeUrl: "type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext",
	}

	toReturn := &cluster.Cluster{
		Name: "RickandMortyAPI",
		ConnectTimeout: &duration.Duration{
			Seconds: 5,
		},
		ClusterDiscoveryType: &cluster.Cluster_Type{
			Type: cluster.Cluster_STRICT_DNS,
		},
		LbPolicy: cluster.Cluster_ROUND_ROBIN,
		LoadAssignment: &endpoint.ClusterLoadAssignment{
			ClusterName: "RickandMortyAPI",
			Endpoints: []*endpoint.LocalityLbEndpoints{{
				LbEndpoints: []*endpoint.LbEndpoint{{
					HostIdentifier: &endpoint.LbEndpoint_Endpoint{
						Endpoint: &endpoint.Endpoint{
							Address: &core.Address{
								Address: &core.Address_SocketAddress{
									SocketAddress: &core.SocketAddress{
										Protocol: core.SocketAddress_TCP,
										Address:  "rickandmortyapi.com",
										PortSpecifier: &core.SocketAddress_PortValue{
											PortValue: 443,
										},
									},
								},
							},
						},
					},
				}},
			}},
		},
		TransportSocket: &core.TransportSocket{
			Name: wellknown.TransportSocketTls,
			ConfigType: &core.TransportSocket_TypedConfig{
				TypedConfig: &transportsocket,
			},
		},
	}

	return toReturn
}

func makeRoute() *route.RouteConfiguration {
	toReturn := &route.RouteConfiguration{
		Name: "RickandMortyAPI",
		VirtualHosts: []*route.VirtualHost{{
			Name:    "RickandMortyAPI",
			Domains: []string{"rickandmortyapi.com"},
			Routes: []*route.Route{{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: "/api",
					},
				},
				Action: &route.Route_Route{
					Route: &route.RouteAction{
						ClusterSpecifier: &route.RouteAction_Cluster{
							Cluster: "RickandMortyAPI",
						},
						Timeout: &duration.Duration{
							Seconds: 5,
						},
					},
				},
			}},
		}},
	}

	return toReturn
}

func makeHTTPListener() *listener.Listener {
	rte := &route.RouteConfiguration{
		Name: "RickandMortyAPI",
		VirtualHosts: []*route.VirtualHost{{
			Name:    "RickandMortyAPI",
			Domains: []string{"*"},
			Routes: []*route.Route{{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: "/api/rnk",
					},
				},
				Action: &route.Route_Route{
					Route: &route.RouteAction{
						ClusterSpecifier: &route.RouteAction_Cluster{
							Cluster: "RickandMortyAPI",
						},
						Timeout: &duration.Duration{
							Seconds: 5,
						},
						PrefixRewrite: "/api",
						HostRewriteSpecifier: &route.RouteAction_HostRewriteLiteral{
							HostRewriteLiteral: "rickandmortyapi.com",
						},
					},
				},
			}},
			RequestHeadersToRemove: []string{"x-api-key"},
		}},
	}

	httpfc := anypb.Any{
		TypeUrl: "type.googleapis.com/envoy.extensions.filters.http.router.v3.Router",
	}

	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "ingress_http",
		RouteSpecifier: &hcm.HttpConnectionManager_RouteConfig{
			RouteConfig: rte,
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name: "envoy.filters.http.router",
			ConfigType: &hcm.HttpFilter_TypedConfig{
				TypedConfig: &httpfc,
			},
		}},
	}

	pbst, err := ptypes.MarshalAny(manager)
	if err != nil {
		log.Fatal(err)
	}

	toReturn := &listener.Listener{
		Name: "RickandMortyAPI",
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: 10000,
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{{
			Filters: []*listener.Filter{{
				Name: "envoy.filters.network.http_connection_manager",
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		}},
	}

	return toReturn
}

func GenerateSnapshot() *cache.Snapshot {
	snap, _ := cache.NewSnapshot("1",
		map[resource.Type][]types.Resource{
			resource.ClusterType:  {makeCluster()},
			resource.RouteType:    {makeRoute()},
			resource.ListenerType: {makeHTTPListener()},
		},
	)
	return snap
}
