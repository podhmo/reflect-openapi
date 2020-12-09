doc:
	$(foreach x,$(shell find . -name Makefile -d 2), $(call NewLineT,make -C $(dir $(x))  openapi.json))

define NewLineT
	$(1)

endef
