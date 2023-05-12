fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

run-ui:
	cd ./ui && npm install
	cd ./ui && npm run format
	cd ./ui && npm run dev

run-server: fmt vet
	go run ./main.go

build-ui:
	cd ./ui && npm install
	cd ./ui && npm run build

build-server: fmt vet
	go build -o ./dist/aeto-web ./main.go

build: build-ui build-server
