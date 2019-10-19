BIN         := ./bin/cni-plugin-logger-cloudwatch

export GO111MODULE=on

build: 
	go build -mod=vendor -o $(BIN) 

test:
	golint ./pkg/...
	go test -cover ./pkg/... 

