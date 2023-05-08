package envoy

import (
	"context"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discoveryGrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"log"
)

type Callbacks struct{}

func (c Callbacks) OnFetchRequest(context.Context, *discoveryGrpc.DiscoveryRequest) error {
	return nil
}

func (c Callbacks) OnFetchResponse(*discoveryGrpc.DiscoveryRequest, *discoveryGrpc.DiscoveryResponse) {
}

func (c Callbacks) OnStreamOpen(_ context.Context, id int64, typ string) error {
	log.Printf("stream open for %s with id %d", typ, id)
	return nil
}

func (c Callbacks) OnStreamClosed(id int64, _ *core.Node) {
	log.Printf("stream closed for id %d", id)
}

func (c Callbacks) OnStreamRequest(int64, *discoveryGrpc.DiscoveryRequest) error {
	return nil
}

func (c Callbacks) OnStreamResponse(context.Context, int64, *discoveryGrpc.DiscoveryRequest, *discoveryGrpc.DiscoveryResponse) {
}

func (c Callbacks) OnDeltaStreamOpen(context.Context, int64, string) error {
	return nil
}

func (c Callbacks) OnDeltaStreamClosed(int64, *core.Node) {}

func (c Callbacks) OnStreamDeltaRequest(int64, *discoveryGrpc.DeltaDiscoveryRequest) error {
	return nil
}

func (c Callbacks) OnStreamDeltaResponse(int64, *discoveryGrpc.DeltaDiscoveryRequest, *discoveryGrpc.DeltaDiscoveryResponse) {
}
