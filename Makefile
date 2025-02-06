OUT:=clamav-prometheus-exporter
OWNER:=rekzi
VERSION:=latest
IMAGE:=$(OWNER)/$(OUT):$(VERSION)
IMAGE_EXTRA_ARGS?=

build:
	CGO_ENABLED=0 && go build -installsuffix 'static' -o ${OUT} .
build-version:
	CGO_ENABLED=0 && go build -installsuffix 'static' -o ${OUT} -ldflags="-X main.version=${VERSION}" .
image:
	docker build -t ${OUT} -t ${IMAGE} ${IMAGE_EXTRA_ARGS} .
push:
	docker push ${IMAGE}
clean:
	rm -rf $(OUT)
