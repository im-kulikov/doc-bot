NAME=imkulikov/docbot
VERSION ?= $(shell git rev-parse --short HEAD)

.PHONY: image publish
image:
	docker build -t $(NAME):latest .
	docker build -t $(NAME):v-$(VERSION) .

publish:
	docker push $(NAME):latest
	docker push $(NAME):v-$(VERSION)