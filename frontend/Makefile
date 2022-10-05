VERSION_HASH_VALUE = $(shell git rev-parse --short HEAD)

.PHONY: minify-static-development
minify-static-development:
	MINIFY=true VERSION_HASH=_version_ npm run build
	@echo "./static/_version_ created"

.PHONY: minify-static
minify-static:
	MINIFY=true VERSION_HASH=$(VERSION_HASH_VALUE) npm run build
	@echo "./static/$(VERSION_HASH_VALUE) created"

.PHONY: print-version-hash
print-version-hash:
	@echo $(VERSION_HASH_VALUE)
