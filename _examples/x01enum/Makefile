export ADDR ?= :44444

openapi.json: main.go
	go run $^ --doc > $@
openapi-without-ref.json: main.go
	WITHOUT_REF=1 go run $^ --doc > $@

diff: openapi-without-ref.json openapi.json
	diff -u $^ || :
.PHONY: diff
