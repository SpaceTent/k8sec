PROJECT := k8sec
VERSION := $(shell git describe --tag --abbrev=0)
NEXT_VERSION:=$(shell git describe --tags --abbrev=0 | awk -F . '{OFS="."; $$NF+=1; print}')
SHA1 := $(shell git rev-parse HEAD)
NOW := $(shell date -u +'%Y%m%d-%H%M%S')

build:
	go build -o bin/$(PROJECT) -ldflags "-X main.VERSION=$(VERSION) -X main.BUILDTIME=$(NOW)" main.go item.go model.go

release: fmt
	@git tag -a $(NEXT_VERSION) -m "Release $(NEXT_VERSION)"
	@git push --all
	@git push --tags

fmt:
	@go mod tidy
	@goimports -w .
	@gofmt -w -s .
	@go clean ./...


commit: fmt
	@govulncheck ./...
	@git add .
	@git commit -a -m "$(filter-out $@,$(MAKECMDGOALS))"
	@git pull
	@git push

%:
	@:



