IMAGE_TAG=latest
NAME=clamav-prometheus-exporter
OWNER=rekzi

build:
	CGO_ENABLED=0 && go build -installsuffix 'static' -o ${NAME} .
image:
	docker build -t ${NAME} -t ${OWNER}/${NAME}:${IMAGE_TAG} .
push:
	docker push ${OWNER}/${NAME}:${IMAGE_TAG}
