FROM gcr.io/distroless/static-debian12:nonroot@sha256:a9329520abc449e3b14d5bc3a6ffae065bdde0f02667fa10880c49b35c109fd1
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/atlas-stats-exporter /usr/bin/atlas-stats-exporter
ENTRYPOINT ["/usr/bin/atlas-stats-exporter"]