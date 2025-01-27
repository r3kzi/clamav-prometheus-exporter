build:
	CGO_ENABLED=0 && go build -installsuffix 'static' -o clamav-prometheus-exporter .
image:
	docker build -t clamav-prometheus-exporter -t rekzi/clamav-prometheus-exporter:latest .
push:
	docker push rekzi/clamav-prometheus-exporter:latest