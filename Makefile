GO_MODULE := GO111MODULE=on
BUILD_TAGS := -tags netgo
BUILD_FLAGS := -ldflags '-extldflags "-static"'

build:
	$(GO_MODULE) go build -o ./bin/ $(BUILD_TAGS) $(BUILD_FLAGS)

test:
	$(GO_MODULE) go test -cover ./...

.PHONY: banuser
banuser:
	$(GO_MODULE) go build -o ./bin/ $(BUILD_TAGS) $(BUILD_FLAGS) ./cmd/banuser

.PHONY: followback
followback:
	$(GO_MODULE) go build -o ./bin/ $(BUILD_TAGS) $(BUILD_FLAGS) ./cmd/followback
