SHELL:=/bin/bash -O globstar
PORT=8000
ENV=local

build:
	go mod tidy
	protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative  api/**/*.proto
	go build -v -tags=e2e

run: build
	sudo -E ./multy serve --port=$(PORT) --env=${ENV}

clean:
	find api/proto -name '*.pb.go' -delete
	find . -name '*.lock.hcl' -delete
	find . -name '*.tfstate*' -delete
	rm -rf ./test/**/.terraform
