##@ Test
GOOS := $(if $(GOOS),$(GOOS),linux)
GOARCH := $(if $(GOARCH),$(GOARCH),amd64)
GO=GO15VENDOREXPERIMENT="1" CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) GO111MODULE=on go
GOTEST=GO15VENDOREXPERIMENT="1" CGO_ENABLED=1 GO111MODULE=on go test # go race detector requires c
FILES := $$(find . -name "*.go" | grep -vE "vendor")
GOFILTER := grep -vE 'vendor|render.Delims|bindata_assetfs|testutil|\.pb\.go'
GOCHECKER := $(GOFILTER) | awk '{ print } END { if (NR > 0) { exit 1 } }'

default: test

test: check
	$(GO) build ./...

check:
	GO111MODULE=off go get golang.org/x/lint/golint
	@echo "golint"
	@ $(GOLINT)

	@echo "gofmt"
	@gofmt -s -l -w $(FILES) 2>&1 | $(GOCHECKER)
