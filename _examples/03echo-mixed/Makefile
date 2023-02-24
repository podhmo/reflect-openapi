export ADDR ?= :44444

openapi.json: main.go
	go generate

run:
	go run ./
.PHONY: run


# http :44444/users
# http :44444/users/1
