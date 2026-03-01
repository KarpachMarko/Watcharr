.PHONY: dev version

dev:
	./node_modules/.bin/vite dev

dev_host:
	./node_modules/.bin/vite dev --host

# Don't judge me.
version:
	@$(eval VER := $(filter-out $@,$(MAKECMDGOALS)))
	@echo "Updating version to $(VER)"
	echo -n "$(VER)" > ./server/VERSION
	npm version $(VER) --git-tag-version false
%:
	@:
