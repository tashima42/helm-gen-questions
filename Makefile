# The plugin version, helm names the package after it.
VERSION ?= 0.0.0+snapshot-$(git rev-parse --short HEAD)
LDFLAGS := -s -w -X github.com/tashima42/helm-questions/cmd.version=$(VERSION)
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

.PHONY: build build-all package test clean

# Builds the binary of the current platform, so the repository can be used with `helm plugin install .`
build:
	@mkdir -p bin
	CGO_ENABLED=0 go build -trimpath -ldflags '$(LDFLAGS)' \
		-o bin/helm-questions_$$(go env GOOS)_$$(go env GOARCH)$$(go env GOEXE) .

# Builds the binaries of all the platforms listed in plugin.yaml.
build-all:
	@mkdir -p bin
	@for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; ext=; [ $$os = windows ] && ext=.exe; \
		echo "building $$os/$$arch"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -trimpath -ldflags '$(LDFLAGS)' \
			-o bin/helm-questions_$${os}_$${arch}$$ext . || exit 1; \
	done

# Packages the plugin with the binaries of all the platforms into dist/, with its checksum. Not signed.
package: build-all
	./package.sh

test:
	go test ./...

clean:
	rm -rf bin dist
