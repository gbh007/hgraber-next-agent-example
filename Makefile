create_build_dir:
	mkdir -p ./_build

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./_build/hgraber-agent-arm64 ./cmd/agent
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./_build/hgraber-agent-amd64 ./cmd/agent

.PHONY: generate
generate:
	go run github.com/ogen-go/ogen/cmd/ogen@v1.2.1 --target open_api/agentAPI -package agentAPI --clean open_api/agent.yaml
	go run github.com/ogen-go/ogen/cmd/ogen@v1.2.1 --target open_api/serverAPI -package serverAPI --clean open_api/server.yaml

.PHONY: docker
docker: build	
	docker build -f Dockerfile \
		--build-arg "BINARY_PATH=./_build/hgraber-agent-arm64" \
		-t hgraber-next-agent:arm64 .
	docker save hgraber-next-agent:arm64 -o _build/hgraber-next-agent_arm64.tar

	docker build -f Dockerfile \
		--build-arg "BINARY_PATH=./_build/hgraber-agent-amd64" \
		-t hgraber-next-agent:amd64 .
	docker save hgraber-next-agent:amd64 -o _build/hgraber-next-agent_amd64.tar

.PHONY: run-example
run-example: create_build_dir
	CGO_ENABLED=0 go build -trimpath -o ./_build/hgraber-agent ./cmd/agent

	./_build/hgraber-agent --config config-example.yaml

.PHONY: run
run: create_build_dir
	CGO_ENABLED=0 go build -trimpath -o ./_build/hgraber-agent ./cmd/agent

	./_build/hgraber-agent

.PHONY: scan
scan: create_build_dir
	CGO_ENABLED=0 go build -trimpath -o ./_build/hgraber-agent ./cmd/agent

	./_build/hgraber-agent --scan