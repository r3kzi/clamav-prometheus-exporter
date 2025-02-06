OUT:=clamav-prometheus-exporter
OWNER:=rekzi
VERSION:=$(file < VERSION)
IMAGE:=$(OWNER)/$(OUT):$(VERSION)
IMAGE_EXTRA_ARGS?=

build:
	CGO_ENABLED=0 && go build -installsuffix 'static' -o ${OUT} .
build-version:
	CGO_ENABLED=0 && go build -installsuffix 'static' -o ${OUT} -ldflags="-X main.version=${VERSION}" .
clean:
	rm -rf $(OUT)
image:
	docker build -t ${OUT}:latest -t ${IMAGE} ${IMAGE_EXTRA_ARGS} .
push:
	docker push ${OUT}:latest && docker push ${IMAGE}
run:
	go run .
