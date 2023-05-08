package main

import (
	"context"
	clusterService "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoveryGrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointService "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerService "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeService "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimeService "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretService "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	callbacks "king/envoy"
	myResource "king/envoy"
	"king/routes"
	"log"
	"net"
)

func main() {
	// Create cache
	snapshotCache := cache.NewSnapshotCache(false, cache.IDHash{}, nil)

	// Add snapshot to cache
	nodeID := "my_node_id"
	ctx := context.Background()
	if err := snapshotCache.SetSnapshot(ctx, nodeID, myResource.GenerateSnapshot()); err != nil {
		log.Fatalf("failed to set snapshot: %v", err)
	}

	// Initialize gRPC server
	grpcServer := grpc.NewServer()

	// Configure callbacks on xDS server
	xdsServer := server.NewServer(context.Background(), snapshotCache, &callbacks.Callbacks{})

	// Registry xDS server with gRPC server
	discoveryGrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, xdsServer)
	endpointService.RegisterEndpointDiscoveryServiceServer(grpcServer, xdsServer)
	clusterService.RegisterClusterDiscoveryServiceServer(grpcServer, xdsServer)
	routeService.RegisterRouteDiscoveryServiceServer(grpcServer, xdsServer)
	listenerService.RegisterListenerDiscoveryServiceServer(grpcServer, xdsServer)
	secretService.RegisterSecretDiscoveryServiceServer(grpcServer, xdsServer)
	runtimeService.RegisterRuntimeDiscoveryServiceServer(grpcServer, xdsServer)

	// Initialize gin server
	router := gin.Default()
	// Setup routes
	routes.SetupRoutes(router)

	// Start gin server in a separate goroutine
	go func() {
		if err := router.Run(":8081"); err != nil {
			log.Fatalf("failed to serve gin server: %v", err)
		}
	}()

	// Start gRPC server in a separate goroutine
	go func() {
		lis, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("gRPC server listening on %s\n", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()

	// Use a channel to block main and keep it running
	forever := make(chan bool)
	<-forever
}
