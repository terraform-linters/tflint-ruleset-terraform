default: build

test:
	go test ./...

build:
	go build

install: build
	mkdir -p ~/.tflint.d/plugins
	mv ./tflint-ruleset-opentofu ~/.tflint.d/plugins

release:
	cd tools/release; go run main.go
