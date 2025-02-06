OUT:=clamav-prometheus-exporter
OWNER:=rekzi
VERSION:=latest
IMAGE:=$(OWNER)/$(OUT):$(VERSION)
IMAGE_EXTRA_ARGS?=

build:
	CGO_ENABLED=0 && go build -installsuffix 'static' -o ${OUT} -ldflags="-X main.version=${VERSION}" .
clean:
	rm -rf $(OUT)
image:
	docker build -t ${IMAGE} --build-arg=VERSION=${VERSION} ${IMAGE_EXTRA_ARGS} .
push:
	docker push ${IMAGE}
run:
	go run .
