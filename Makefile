image:
	docker build -t clamav-prometheus-exporter -t rekzi/clamav-prometheus-exporter:latest .
push:
	docker push rekzi/clamav-prometheus-exporter:latest