create_build_dir:
	mkdir -p ./_build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./_build/hgraber-agent-arm64 ./cmd/agent
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./_build/hgraber-agent-amd64 ./cmd/agent

run: create_build_dir
	go build -trimpath -o ./_build/hgraber-agent  ./cmd/agent
	APP_API_ADDR=127.0.0.1:8081 \
	./_build/hgraber-agent --config config-example.yaml

runtrace: create_build_dir
	go build -trimpath -o ./_build/hgraber-agent  ./cmd/agent
	APP_API_ADDR=127.0.0.1:8082 \
	APP_TRACE_ENDPOINT=http://localhost:4318/v1/traces \
	./_build/hgraber-agent --config config-example.yaml

.PHONY: generate
generate:
	go run github.com/ogen-go/ogen/cmd/ogen@v1.2.1 --target internal/controller/api/internal/server -package server --clean agent.yaml

docker: build	
	docker build -f Dockerfile \
		--build-arg "BINARY_PATH=./_build/hgraber-agent-arm64" \
		-t hgraber-next-agent:arm64 .
	docker save hgraber-next-agent:arm64 -o _build/hgraber-next-agent_arm64.tar

	docker build -f Dockerfile \
		--build-arg "BINARY_PATH=./_build/hgraber-agent-amd64" \
		-t hgraber-next-agent:amd64 .
	docker save hgraber-next-agent:amd64 -o _build/hgraber-next-agent_amd64.tar