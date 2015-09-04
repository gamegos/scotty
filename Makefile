APP_NAME = scotty
REPO = github.com/gamegos/scotty

WORKDIR = $(shell pwd)
OUTPUT_DIR = $(WORKDIR)/bin
IMAGE_NAME = gamegos/$(APP_NAME)

build:
	CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' -o $(OUTPUT_DIR)/$(APP_NAME)

install-deps:
	go get ./...

build-in-container:
	docker run --rm -v $(WORKDIR):/gopath/src/$(REPO) -i -t google/golang make -C /gopath/src/$(REPO) install-deps build

bin/$(APP_NAME): build-in-container

# docker image
build-image: bin/$(APP_NAME)
	docker build -t $(IMAGE_NAME) .

export-image: build-image
	docker save -o $(OUTPUT_DIR)/$(APP_NAME)-image $(IMAGE_NAME)
	docker rmi $(IMAGE_NAME)

delete-image:
	docker rmi $(IMAGE_NAME)

.PHONY: build install-deps build-in-container build-image export-image delete-image
