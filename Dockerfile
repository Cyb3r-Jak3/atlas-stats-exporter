FROM library/alpine:3.23@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11 AS certs
RUN apk --no-cache add ca-certificates

FROM library/busybox:1.37.0@sha256:1487d0af5f52b4ba31c7e465126ee2123fe3f2305d638e7827681e7cf6c83d5e
ARG TARGETPLATFORM
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY $TARGETPLATFORM/atlas-stats-exporter /usr/bin/atlas-stats-exporter
ENTRYPOINT ["/usr/bin/atlas-stats-exporter"]
HEALTHCHECK CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1