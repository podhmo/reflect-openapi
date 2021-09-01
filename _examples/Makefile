doc:
	$(foreach x,$(shell find . -mindepth 2 -name Makefile), $(call NewLineT,make -C $(dir $(x))  openapi.json))

update-deps:
	$(foreach x,$(shell find . -mindepth 2 -name Makefile), $(call NewLineT,cd $(dir $(x)) && go get -u -v && go mod tidy))

define NewLineT
	$(1)

endef
