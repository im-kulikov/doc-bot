NAME=imkulikov/docbot
VERSION ?= $(shell git rev-parse --short HEAD)

.PHONY: image publish
image:
	docker build -t imkulikov/bl-resolver:latest .
	docker build -t imkulikov/bl-resolver:v-$(VERSION) .

publish:
	docker push imkulikov/bl-resolver:latest
	docker push imkulikov/bl-resolver:v-$(VERSION)