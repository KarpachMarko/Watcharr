# Don't judge me.
version:
	@$(eval VER := $(filter-out $@,$(MAKECMDGOALS)))
	@echo "Updating version to $(VER)"
	echo "$(VER)" > ./server/VERSION
	npm version $(VER) --git-tag-version false
%:
	@:
