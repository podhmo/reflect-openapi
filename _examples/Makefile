doc:
	$(foreach x,$(shell find . -mindepth 2 -name Makefile), $(call NewLineT,make -C $(dir $(x))  openapi.json))

define NewLineT
	$(1)

endef
