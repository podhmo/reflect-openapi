export ADDR ?= :44444

openapi.json:
	go run main.go --doc > $@

run:
	go run main.go
.PHONY: run


# http :44444/users
# http :44444/users/1
