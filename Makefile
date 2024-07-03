create_build_dir:
	mkdir -p ./_build

run: create_build_dir
	go build -trimpath -o ./_build/hgraber-agent  ./cmd/agent
	./_build/hgraber-agent --token agent-token --addr 127.0.0.1:8081 --debug --export-path .hidden/exported

.PHONY: generate
generate:
	go run github.com/ogen-go/ogen/cmd/ogen --target internal/controller/api/internal/server -package server --clean agent.yaml