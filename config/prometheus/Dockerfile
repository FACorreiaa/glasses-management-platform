# Dockerfile for Prometheus Fly App

# Use the official Prometheus image as the base
FROM prom/prometheus:v2.51.2

# Copy your local prometheus.yml configuration file into the image
# at the location where Prometheus expects to find it.
COPY prometheus.yml /etc/prometheus/prometheus.yml