main_package_path = ./cmd/buildon
binary_name = buildon

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## audit: run quality control checks
.PHONY: audit
audit: test
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)"
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## tidy: tidy modfiles and modernize and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fix ./...
	go fmt ./...

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs -p 1 ./...

## build: build the applications
.PHONY: build
build:
	go build -o=./bin/${binary_name} ${main_package_path}

## run: run the application
.PHONY: run run/mock
run: build
	./bin/${binary_name}

## run/live : run the application with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/air-verse/air@latest \
		--build.cmd "make build" \
		--build.delay "100" \
		--build.entrypoint "./bin/${binary_name}" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"

## clean: remove unnecessary files
.PHONY: clean
clean:
	rm -rf ./bin
