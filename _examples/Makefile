doc:
	$(foreach x,$(shell find . -name Makefile -d 2), $(call NewLineT,make -C $(dir $(x))  openapi.json))

update-deps:
	$(foreach x,$(shell find . -name Makefile -d 2), $(call NewLineT,cd $(dir $(x)) && go get -u -v && go mod tidy))

define NewLineT
	$(1)

endef
