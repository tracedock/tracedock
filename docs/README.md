# TraceDock

## Introduction

> :warning: This project is under active development, don't use it in production

**TraceDock** is an extensible, plugin-friendly alternative to the OpenTelemetry Collector.

Its main goal is to provide a flexible telemetry pipeline where you can extend its behavior via external `.so` modules, enabling easy extension, isolation, and rapid experimentation.

TraceDock works as a "dock" where telemetry data arrive, are inspected, processed, and routed to the desired destination.

## Getting started

### Requirements

TraceDock is developed with Goland and requires the version 1.23+ of the language

### Installing

You can run it using docker.

```shell
docker run -ti -v ./config/tracedock.yaml:/etc/tracedock/tracedock.yaml tracedock/tracedock server start --grpc-port=50051 --http-port=8080 -c /etc/tracedock/tracedock.yaml
```

Alternativelly, you can compile it on your own machine.

```shell
git clone https://github.com/tracedock/tracedock.git
cd tracedock
make build

./tracedock server start --grpc-port=50051 --http-port=8080 -c /etc/tracedock/tracedock.yaml
```

### Configure

> :warning: Sample configuration coming soon

TraceDock uses YAML file to configure the pipelines and routing rules.
