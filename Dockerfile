FROM gcr.io/distroless/static-debian12:nonroot@sha256:2b7c93f6d6648c11f0e80a48558c8f77885eb0445213b8e69a6a0d7c89fc6ae4
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/atlas-stats-exporter /usr/bin/atlas-stats-exporter
ENTRYPOINT ["/usr/bin/atlas-stats-exporter"]