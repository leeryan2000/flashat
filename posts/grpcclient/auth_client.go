// Package grpcclient dials the backend monolith's internal AuthInternal
// gRPC service (see backend/transport/grpcserver).
package grpcclient

import (
	"github.com/leeryan2000/flashat-posts/internal/genproto/authpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewAuthClient dials addr (e.g. "backend:50051") once at startup and
// returns a long-lived client — grpc.NewClient connects lazily and
// reconnects on its own, so this is safe to call before the backend's
// listener is necessarily up yet (docker-compose depends_on only orders
// container start, not readiness).
//
// Traffic is unencrypted (insecure.NewCredentials()) — acceptable here
// since it never leaves the single-host flashat-net bridge network, but
// not what a real multi-host service mesh would do (that would need
// mTLS).
func NewAuthClient(addr string) (authpb.AuthInternalClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	return authpb.NewAuthInternalClient(conn), conn, nil
}
