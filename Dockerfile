FROM library/alpine:3.24@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6 AS certs
RUN apk --no-cache add ca-certificates

FROM library/busybox:1.38.0@sha256:fd8d9aa63ba2f0982b5304e1ee8d3b90a210bc1ffb5314d980eb6962f1a9715d
ARG TARGETPLATFORM
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY $TARGETPLATFORM/atlas-stats-exporter /usr/bin/atlas-stats-exporter
ENTRYPOINT ["/usr/bin/atlas-stats-exporter"]
HEALTHCHECK CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1