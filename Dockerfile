FROM library/alpine:3.23@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11 AS certs
RUN apk --no-cache add ca-certificates

FROM library/busybox:1.38.0@sha256:fd8d9aa63ba2f0982b5304e1ee8d3b90a210bc1ffb5314d980eb6962f1a9715d
ARG TARGETPLATFORM
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY $TARGETPLATFORM/atlas-stats-exporter /usr/bin/atlas-stats-exporter
ENTRYPOINT ["/usr/bin/atlas-stats-exporter"]
HEALTHCHECK CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1