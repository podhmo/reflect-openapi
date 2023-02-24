export ADDR ?= :44444

openapi.json: main.go
	go run ./ --doc --docfile $@ --mdfile README.md

run:
	go run ./
.PHONY: run


# http :44444/users
# http :44444/users/1
