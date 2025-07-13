FROM gcr.io/distroless/static-debian12
COPY atlas-stats-exporter /
ENTRYPOINT ["/atlas-stats-exporter"]
