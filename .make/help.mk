.DEFAULT_GOAL := help

help::
	@egrep -h '\s##\s' $(MAKEFILE_LIST) \
		| awk -F':.*?## | \\| ' '{printf "\033[36m%-22s \033[37m %-40s \033[35m%s \n", $$1, $$2, $$3}'