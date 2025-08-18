
GOBIN=$(shell go env GOPATH)/bin

path=./...

run-example:
	go run example/main.go

lint: setup
	@$(GOBIN)/staticcheck $(path) $(args)
	@$(GOBIN)/errcheck $(path) $(args)
	@go vet $(path) $(args)
	@echo "StaticCheck & ErrCheck & Go Vet found no problems on your code!"

test: setup
	$(GOBIN)/richgo test $(path) $(args)

setup: $(GOBIN)/richgo $(GOBIN)/staticcheck $(GOBIN)/errcheck

$(GOBIN)/richgo:
	go get github.com/kyoh86/richgo

$(GOBIN)/staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest

$(GOBIN)/errcheck:
	go install github.com/kisielk/errcheck@latest
