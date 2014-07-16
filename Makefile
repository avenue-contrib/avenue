PROG=avenue
OUT_DIR=bin


binary: deps
	@go build -o $(OUT_DIR)/$(PROG)

deps: env
	@command -v depman >/dev/null 2>&1 || { go get github.com/vube/depman; }
	@depman install

env:
	@mkdir -p $(OUT_DIR)

clean:
	rm -rf $(OUT_DIR)

examples:
	find ./examples -depth 1 -type d -exec go build -o $(OUT_DIR)/{} ./{} \;

.PHONY: examples