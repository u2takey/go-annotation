.PHONY: install build example


install:
	go install ./cmd/annotation-gen/.

build:
	go build -o bin/annotation-gen ./cmd/annotation-gen/.

example: install
	go generate ./examples/. && go test ./examples/.