NAME	:= api-rback
SRCS	:= $(shell find . -type d -name archive -prune -o -type f -name '*.go')
LDFLAGS	:= -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

bin/$(NAME)/static: $(SRCS)
	CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME)

bin/$(NAME): $(SRCS)
	go build -o bin/$(NAME)

bin/$(NAME)/embed/static: $(SRCS)
	CGO_ENABLED=0 go build -a -tags netgo -tags embed -installsuffix netgo $(LDFLAGS) -o bin/$(NAME)

bin/$(NAME)/embed: $(SRCS)
	go build -tags embed -o bin/$(NAME)

.PHONY: deps
deps:
	go get -v

.PHONY: cross-build
cross-build: deps
	for os in darwin linux windows; do \
		for arch in amd64 386; do \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$$os-$$arch/$(NAME); \
		done; \
	done

.PHONY: cross-build
cross-build/embed: deps
	for os in darwin linux windows; do \
		for arch in amd64 386; do \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -a -tags netgo -tags embed -installsuffix netgo $(LDFLAGS) -o dist/$$os-$$arch/$(NAME); \
		done; \
	done
