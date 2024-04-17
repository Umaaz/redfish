
.PHONY: cli build-go
cli: build-go

build-go:
	goreleaser build --single-target --snapshot --clean

.PHONY: gen-pkl
gen-pkl:
	pkl-gen-go pkg/config/pkl/AppCopnfig.pkl --generator-settings pkg/config/pkl/generator-settings.pkl

.PHONY: docker
docker:
	 goreleaser --pl

# More exclusions can be added similar with: -not -path './testbed/*'
ALL_SRC := $(shell find . -name '*.go' -type f | sort)

# ALL_PKGS is used with 'go cover'
ALL_PKGS := $(shell go list $(sort $(dir $(ALL_SRC))))
GO_JUNIT_REPORT=2>&1 | go-junit-report -parser gojson -iocopy -out report.xml

.PHONY: test-with-cover-report
test-with-cover-report:
	go test -race -timeout 20m -count=1 -v -cover -json $(ALL_PKGS) $(GO_JUNIT_REPORT)

FILES_TO_FMT=$(shell find . -type d \( -path ./vendor \) -prune -o -name '*.go' -not -name "*.pb.go" -not -name '*.y.go' -print)

.PHONY: fmt check-fmt
fmt:
	@gofumpt -l -w .
	@goimports -w $(FILES_TO_FMT)

check-fmt: fmt
	@git diff --exit-code -- $(FILES_TO_FMT)

.PHONY: vendor-check
vendor-check: gen-pkl
	git diff --exit-code -- **/go.sum **/go.mod vendor/ pkg/deeppb/ pkg/deepql/ modules/querier/stats modules/frontend/v1/frontendv1pb
