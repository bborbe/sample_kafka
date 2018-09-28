

.PHONY: format
format:
	go get golang.org/x/tools/cmd/goimports
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +

.PHONY: generate
generate:
	go get github.com/actgardner/gogen-avro/gogen-avro
	go generate ./...

