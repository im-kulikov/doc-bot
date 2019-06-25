NAME=imkulikov/docbot
REPO = $(shell go list -m)
VERSION ?= $(shell git rev-parse --short HEAD)

.PHONY: image publish
image:
	@go mod tidy -v
	@go mod vendor
	docker build --build-arg VERSION=$(VERSION) --build-arg REPO=$(REPO) -t $(NAME):latest .
	docker build --build-arg VERSION=$(VERSION) --build-arg REPO=$(REPO) -t $(NAME):v-$(VERSION) .

publish:
	docker push $(NAME):latest
	docker push $(NAME):v-$(VERSION)
