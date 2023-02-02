SHELL := bash

test:
	go test ./... |& grep -v '^# '
.PHONY: test

ci: test
.PHONY: ci
