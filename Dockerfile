FROM scratch
COPY atlas-stats-exporter /
ENTRYPOINT ["/atlas-stats-exporter"]
