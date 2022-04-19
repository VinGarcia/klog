
GOBIN=$(shell go env GOPATH)/bin

run-example:
	go run example/main.go

lint: setup
	@$(GOBIN)/golint -set_exit_status -min_confidence 0.9 $(path) $(args)
	@go vet $(path) $(args)
	@echo "Golint & Go Vet found no problems on your code!"

test: setup
	$(GOBIN)/richgo test $(path) $(args)

setup: $(GOBIN)/richgo $(GOBIN)/golint

$(GOBIN)/richgo:
	go get github.com/kyoh86/richgo

$(GOBIN)/golint:
	go get golang.org/x/lint
