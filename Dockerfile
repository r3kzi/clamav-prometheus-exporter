# First stage: build the executable.
FROM golang:1.21 as builder
WORKDIR /go/src/github.com/rekzi/clamav-prometheus-exporter/
COPY . .
# CGO_ENABLED=0 to build a statically-linked executable

ENV CGO_ENABLED=0
RUN go build -installsuffix 'static' -o clamav-prometheus-exporter .

# Final stage: the running container.
# Use a minimal image for running the application
FROM alpine:3.18 AS final

# Install necessary certificates
RUN apk add --no-cache ca-certificates

# Set the working directory for the runtime
WORKDIR /bin/

# Import the compiled executable from the first stage.
COPY --from=builder /go/src/github.com/rekzi/clamav-prometheus-exporter/clamav-prometheus-exporter .

RUN addgroup prometheus
RUN adduser -S -u 1000 prometheus \
    && chown -R prometheus:prometheus /bin

USER 1000
# Default port for metrics exporters
EXPOSE 9810
ENTRYPOINT [ "/bin/clamav-prometheus-exporter" ]
