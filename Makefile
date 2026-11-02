.PHONY: test lint demo logs

RUN = go run ./plugin/scripts

test:
	test -z "$$(gofmt -l plugin tests examples)"
	go vet ./...
	for d in before after comments; do (cd examples/shop/$$d && go vet ./...) || exit 1; done
	go test ./...

lint:
	claude plugin validate ./plugin --strict
	claude plugin validate . --strict

demo:
	$(RUN) compare examples/logs/shop-before.log examples/logs/shop-after.log
	$(RUN) compare examples/logs/shop-before.log examples/logs/shop-comments.log

# Rebuilds the shop logs from examples/shop. Compiler output depends on the
# Go version, so the saved logs say which version made them in examples/logs/README.md.
logs:
	for d in before after comments; do (cd examples/shop/$$d && go build -gcflags=-m=2 -o /dev/null ./... 2> ../../logs/shop-$$d.log) || exit 1; done
	cd examples/shop/before && go build -gcflags=-m -o /dev/null ./... 2> ../../logs/shop-before-m1.log
	cd examples/shop/after && go build -gcflags=-m -o /dev/null ./... 2> ../../logs/shop-after-m1.log
