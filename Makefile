build:
	docker build -t neurosnap/nvimsh .
.PHONY: build

push:
	docker push neurosnap/nvimsh:latest
.PHONY: push

upload: build push
.PHONY: upload
