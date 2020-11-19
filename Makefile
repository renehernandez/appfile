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

generate:
	go generate ${PKGS}
.PHONY: generate

pristine: generate fmt
	git diff | cat
	git ls-files --exclude-standard --modified --deleted --others -x vendor  | grep -v '^go.' | diff /dev/null -
.PHONY: pristine