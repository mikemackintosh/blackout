bin=blackout
OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))

build:
	go build -o bin/$(bin)-$(OS)-$(ARCH) ./cmd/$(bin)/...
