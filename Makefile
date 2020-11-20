VERSION ?= $(shell git describe --abbrev=0 --tags)
PKGS := $(shell go list ./... | grep -v "/vendor/\|/examples")

fmt:
	go fmt ${PKGS}
.PHONY: fmt

check:
	go vet ${PKGS}
.PHONY: check

test:
	go test -v ${PKGS} -cover -race -p=1
.PHONY: test

build:
	go build -ldflags '-X github.com/renehernandez/appfile/internal/version.Version=${VERSION}'
.PHONY: build

generate:
	go generate ${PKGS}
.PHONY: generate

pristine: generate fmt
	git diff | cat
	git ls-files --exclude-standard --modified --deleted --others -x vendor  | grep -v '^go.' | diff /dev/null -
.PHONY: pristine

tools:
	go get -u github.com/mitchellh/gox
.PHONY: tools

release:
	env CGO_ENABLED=0 gox -osarch '!darwin/386' -os '!openbsd !freebsd !netbsd' -arch '!mips !mipsle !mips64 !mips64le !s390x' -output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}" -ldflags '-X github.com/renehernandez/appfile/internal/version.Version=${VERSION}'
.PHONY: release

clean:
	rm dist/appfile_*
.PHONY: clean

