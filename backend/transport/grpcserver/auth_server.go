// Package grpcserver exposes the monolith's session/friendship data to
// other internal services (currently just posts) over gRPC. It never
// runs behind nginx and is never published in docker-compose — reachable
// only as backend:50051 on the flashat-net bridge network.
package grpcserver

import (
	"context"

	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/db"
	"github.com/leeryan2000/flashat/internal/genproto/authpb"
	"github.com/leeryan2000/flashat/repo"
)

type AuthServer struct {
	authpb.UnimplementedAuthInternalServer
	RedisClient    *db.RedisClient
	FriendshipRepo repo.FriendshipRepo
}

// ValidateSession mirrors middleware.Authenticate's Redis lookup. An
// invalid/expired session is a normal, expected outcome here — it's
// reported via Valid: false, not a gRPC error.
func (s *AuthServer) ValidateSession(ctx context.Context, req *authpb.ValidateSessionRequest) (*authpb.ValidateSessionResponse, error) {
	uid, err := s.RedisClient.GetSession(ctx, req.GetSessionId())
	if err != nil {
		return &authpb.ValidateSessionResponse{Valid: false}, nil
	}
	return &authpb.ValidateSessionResponse{Uid: uid, Valid: true}, nil
}

func (s *AuthServer) GetFriendUIDs(ctx context.Context, req *authpb.GetFriendUIDsRequest) (*authpb.GetFriendUIDsResponse, error) {
	uid, err := uuid.Parse(req.GetUid())
	if err != nil {
		return nil, err
	}

	friendUIDs, err := s.FriendshipRepo.ListAcceptedFriendUIDs(ctx, uid)
	if err != nil {
		return nil, err
	}

	out := make([]string, len(friendUIDs))
	for i, u := range friendUIDs {
		out[i] = u.String()
	}
	return &authpb.GetFriendUIDsResponse{FriendUids: out}, nil
}
