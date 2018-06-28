CMD = dev-watcher
GIT_TAG := $(shell (git describe --abbrev=0 --tags 2> /dev/null || echo v0.0.0) | head -n1)
GIT_HASH := $(shell (git show-ref --head --hash=8 2> /dev/null || echo 00000000) | head -n1)
SRC_DIR := $(shell ls -d */|grep -vE 'vendor')

FMT_ARG ?= -w

.PHONY: all
all: clean $(CMD)

.PHONY: deps
deps:
	# install deps
	@hash dep > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/golang/dep/cmd/dep; \
	fi
	@dep ensure -v

.PHONY: fmt
fmt:
	# gofmt code
	@gofmt -s -l $(FMT_ARG) $(SRC_DIR)
	@go tool vet $(SRC_DIR)

.PHONY: $(CMD)
$(CMD):
	# go build
	go build \
		-ldflags='-X "main.Build=$(GIT_TAG)-$(GIT_HASH)" -X "main.BuildTime=$(BUILD_TIME)"' \
		./cmd/dev-watcher

.PHONY: install
install:
	go install \
			-ldflags='-X "main.Build=$(GIT_TAG)-$(GIT_HASH)" -X "main.BuildTime=$(BUILD_TIME)"' \
			./cmd/dev-watcher

.PHONY: test
test:
	go test -v -coverprofile .cover.out ./...
	@go tool cover -func=.cover.out
	@go tool cover -html=.cover.out -o .cover.html

.PHONY: test/%
test/%:
	go test -v -coverprofile ./$*/.cover.out ./$*
	go tool cover -func=./$*/.cover.out
	go tool cover -html=./$*/.cover.out -o ./$*/.cover.html

.PHONY: clean
clean:
	@rm -rf ./dev-watcher

