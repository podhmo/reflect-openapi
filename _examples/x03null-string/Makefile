export ADDR ?= :44444

openapi.json: main.go
	go run $^ --doc > $@
