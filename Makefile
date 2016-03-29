GO15VENDOREXPERIMENT := 1
GO ?= go
BINARY = kitauth
COVERAGEDIR = coverage

all: fmt build test cover
install-deps:
	glide install
fmt:
	find . -not -path "./vendor/*" -name '*.go' -type f | sed 's#\(.*\)/.*#\1#' | sort -u | xargs -n1 -I {} bash -c "cd {} && goimports -w *.go && gofmt -w -s -l *.go"
clean:
	$(GO) clean
	rm -f $(BINARY)
	rm -f bin/$(BINARY)
	rm -rf $(COVERAGEDIR)
build:
	if [ ! -d bin ]; then mkdir bin; fi
	$(GO) build -v -o bin/$(BINARY)
test:
	#if [ ! -d $(COVERAGEDIR) ]; then mkdir $(COVERAGEDIR); fi
	#$(GO) test -v ./bindings -cover -coverprofile=$(COVERAGEDIR)/bindings.coverprofile
cover:
	#$(GO) tool cover -html=$(COVERAGEDIR)/bindings.coverprofile -o $(COVERAGEDIR)/bindings.html
