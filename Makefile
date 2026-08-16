SHELL := bash.exe

# Regenerates the gRPC Go code from proto/*.proto into both backend/ and
# posts/ — each module gets its own checked-in copy, since Docker builds
# each with a context scoped to its own directory (no shared package
# reachable across ./backend and ./posts at build time).
proto-gen:
	protoc -I proto \
		--go_out=backend/internal/genproto/authpb --go_opt=paths=source_relative \
		--go-grpc_out=backend/internal/genproto/authpb --go-grpc_opt=paths=source_relative \
		proto/auth.proto
	protoc -I proto \
		--go_out=posts/internal/genproto/authpb --go_opt=paths=source_relative \
		--go-grpc_out=posts/internal/genproto/authpb --go-grpc_opt=paths=source_relative \
		proto/auth.proto
