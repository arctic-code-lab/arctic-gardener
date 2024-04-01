# Install tools using asdf
tools:
	@scripts/install-asdf-plugins.sh
	@scripts/install-asdf-versions.sh

build:
	@mkdir -p build
	@GOOS=linux GOARCH=arm64 go build -o build/arctic-gardener ./cmd/arctic-gardener

publish:
	@scp build/arctic-gardener homebridge.local:~

clean:
	@rm -rf build

.PHONY: tools build publish clean