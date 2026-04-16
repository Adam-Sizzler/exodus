.PHONY: download-monaco-deps clean bump-patch bump-minor bump-major tag-release

MONACO_FILES = \
	public/assets/wasm_exec.js \
	public/assets/xray.schema.json \
	public/assets/xray.schema.cn.json \
	public/assets/main.wasm

MONACO_BASE_URL ?= https://teamdominant.github.io/xray-monaco-editor

download-monaco-deps: $(MONACO_FILES)

public/assets/wasm_exec.js:
	@mkdir -p public/assets
	curl -o $@ $(MONACO_BASE_URL)/wasm_exec.js

public/assets/xray.schema.json:
	@mkdir -p public/assets
	curl -o $@ $(MONACO_BASE_URL)/xray.schema.json

public/assets/xray.schema.cn.json:
	@mkdir -p public/assets
	curl -o $@ $(MONACO_BASE_URL)/xray.schema.cn.json

public/assets/main.wasm:
	@mkdir -p public/assets
	curl -o $@ $(MONACO_BASE_URL)/main.wasm

clean:
	rm -f $(MONACO_FILES)

bump-patch:
	npm version patch --no-git-tag-version
	npm install

bump-minor:
	npm version minor --no-git-tag-version
	npm install

bump-major:
	npm version major --no-git-tag-version
	npm install

tag-release:
	@VERSION=$$(node -p "require('./package.json').version") && \
	echo "Creating signed tag for version $$VERSION..." && \
	git tag -s "$$VERSION" -m "Release $$VERSION" && \
	git push origin --follow-tags && \
	echo "Signed tag $$VERSION created and pushed"
